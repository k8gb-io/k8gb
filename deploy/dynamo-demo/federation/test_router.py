# SPDX-FileCopyrightText: Copyright (c) 2026 The k8gb Contributors.
# SPDX-License-Identifier: Apache-2.0

"""Protocol, fail-closed routing and real HTTP streaming regression tests."""

import asyncio
import struct
import time
from contextlib import asynccontextmanager

import aiohttp
import grpc
import pytest
import relay_pb2 as pb
import relay_pb2_grpc as rpc
import xxhash
from aiohttp import web
from aiohttp.test_utils import TestServer
from ckf import ALT_DOMAIN, HEADER, Format, Replica, sequence_hashes
from relay_client import MARKER, Pool, Region, RelayConsumer, dc_identity
from router import ROUTER, Config, Router, create_app, validate_payload

TOKENS = list(range(128))
MODEL = "Qwen/Qwen3-0.6B"
LAYOUT = Format(1, 1234, 64)


def bucket_words(tokens=TOKENS):
    words = [0] * LAYOUT.buckets
    for value in sequence_hashes(tokens, 32):
        mixed = xxhash.xxh3_64_intdigest(struct.pack("<Q", value), seed=LAYOUT.seed)
        fp = (mixed & 0xFFFF) or 1
        first = (mixed >> 16) & (LAYOUT.buckets - 1)
        delta = xxhash.xxh3_64_intdigest(struct.pack("<H", fp), seed=LAYOUT.seed ^ ALT_DOMAIN) & (LAYOUT.buckets - 1)
        for bucket in (first, first ^ (delta or 1)):
            for slot in range(4):
                if (words[bucket] >> (slot * 16)) & 0xFFFF == 0:
                    words[bucket] |= fp << (slot * 16)
                    break
            else:
                continue
            break
        else:
            raise AssertionError("test filter is full")
    return words


def frame(kind, body, sequence, region="us"):
    return (
        HEADER.pack(
            b"CBI1",
            1,
            kind,
            16,
            4,
            1,
            LAYOUT.seed,
            LAYOUT.buckets,
            dc_identity(region),
            sequence,
            xxhash.xxh3_64_intdigest(body) & 0xFFFFFFFF,
        )
        + body
    )


def snapshot(words=None, sequence=1, region="us"):
    words = bucket_words() if words is None else words
    return frame(1, struct.pack("<IIQ", 0, 1, 0) + struct.pack(f"<{len(words)}Q", *words), sequence, region)


def descriptor(region="us", incarnation=1):
    return pb.KvPoolDescriptor(
        producer=pb.ProducerIdentity(
            pool_id=pb.KvPoolId(
                identity_version=1,
                dc_id=dc_identity(region),
                indexer_domain=pb.IndexerDomainId(
                    cache_semantics=pb.DigestIdentity(digest=b"c" * 16, source=1),
                    routing_scope=pb.DigestIdentity(digest=b"r" * 16, source=1),
                ),
            ),
            producer_incarnation=incarnation,
            layout_generation=1,
            ckf_format=pb.CkfFormat(
                format_version=1, seed=LAYOUT.seed, bucket_count=LAYOUT.buckets, fingerprint_bits=16, slots_per_bucket=4
            ),
        ),
        serving_endpoint=pb.DynamoEndpointId(namespace="dynamo", component="backend", endpoint="generate"),
        registrations=[
            pb.ModelRegistration(
                canonical_model_id=MODEL, target=pb.ModelTarget(base=pb.BaseModelTarget(base_model="/opt/model"))
            )
        ],
        query_semantics=pb.KvQuerySemantics(kv_block_size=32, hash_format=1),
        pool_roles=[4],
    )


def add_pool(region, warm=True):
    pool = Pool.from_descriptor(descriptor(region.name), region.name, MODEL, 32)
    pool.replica.apply(1, snapshot(bucket_words() if warm else [0] * 64, region=region.name), 0, 1)
    region.pools[pool.key] = pool
    region.ready_models = {(pool.key, MODEL)}
    region.readiness_updated = region.healthy_updated = time.monotonic()
    region.healthy = True
    return pool


def config(eu="http://eu:8000", us="http://us:8000", freshness=5):
    return Config(
        MODEL,
        "/unused",
        "eu",
        (
            {"name": "eu", "frontend": eu, "relay": "localhost:1"},
            {"name": "us", "frontend": us, "relay": "localhost:2"},
        ),
        freshness_seconds=freshness,
    )


def test_snapshot_delta_and_prefix_queries():
    replica = Replica(LAYOUT, dc_identity("us"))
    replica.apply(1, snapshot(), 0, 1)
    assert replica.prefix_blocks(sequence_hashes(TOKENS, 32)) == 4
    assert replica.prefix_blocks(sequence_hashes([999] + TOKENS[1:], 32)) == 0
    body = struct.pack("<QI", 1, 64) + b"".join(struct.pack("<IQ", i, 0) for i in range(64))
    replica.apply(2, frame(2, body, 2), 1, 2)
    assert replica.prefix_blocks(sequence_hashes(TOKENS, 32)) == 0
    replica.apply(3, b"", 2, 2)
    assert replica.sequence == 2


@pytest.mark.parametrize("failure", ["checksum", "gap", "wrong_dc", "unsupported", "bounds", "heartbeat"])
def test_bad_publication_discards_existing_cache(failure):
    replica = Replica(LAYOUT, dc_identity("us"))
    replica.apply(1, snapshot(), 0, 1)
    kind, payload, base, sequence = 1, snapshot(sequence=2), 0, 2
    if failure == "checksum":
        payload = payload[:-1] + bytes([payload[-1] ^ 1])
    elif failure == "gap":
        kind, base, sequence = 2, 7, 8
        payload = frame(2, struct.pack("<QI", 7, 0), sequence)
    elif failure == "wrong_dc":
        payload = snapshot(sequence=2, region="eu")
    elif failure == "unsupported":
        payload = payload[:8] + b"\x08" + payload[9:]
    elif failure == "bounds":
        kind, base = 2, 1
        payload = frame(2, struct.pack("<QI", 1, 1) + struct.pack("<IQ", 64, 9), sequence)
    elif failure == "heartbeat":
        kind, payload = 3, b""
    with pytest.raises(ValueError):
        replica.apply(kind, payload, base, sequence)
    assert replica.words is None and replica.updated == 0


def test_partial_snapshot_is_never_queryable():
    replica = Replica(LAYOUT, dc_identity("us"))
    words = bucket_words()
    first = frame(1, struct.pack("<IIQ", 0, 2, 0) + struct.pack("<32Q", *words[:32]), 1)
    second = frame(1, struct.pack("<IIQ", 1, 2, 32) + struct.pack("<32Q", *words[32:]), 1)
    replica.apply(1, first, 0, 1)
    assert replica.words is None
    replica.apply(1, second, 0, 1)
    assert replica.prefix_blocks(sequence_hashes(TOKENS, 32)) == 4


def test_routing_stale_readiness_domain_and_local_autonomy():
    router = Router(config(), None)
    local = add_pool(router.regions["eu"], warm=False)
    remote = add_pool(router.regions["us"])
    assert router.choose(TOKENS)["region"] == "us"
    router.regions["us"].readiness_updated -= 10
    assert router.choose(TOKENS)["region"] == "eu"
    router.regions["us"].readiness_updated = time.monotonic()
    remote.descriptor.producer.pool_id.indexer_domain.routing_scope.digest = b"x" * 16
    assert router.choose(TOKENS)["region"] == "eu"
    remote.descriptor.producer.pool_id.indexer_domain.CopyFrom(local.descriptor.producer.pool_id.indexer_domain)
    assert router.choose(TOKENS)["region"] == "us"
    remote.replica.updated -= 10
    assert router.choose(TOKENS)["region"] == "eu"
    router.regions["eu"].reset()
    assert router.choose(TOKENS)["region"] == "eu"


@pytest.mark.parametrize(
    "field,value",
    [
        ("tools", []),
        ("model", "other"),
        ("cache_salt", "tenant"),
        ("max_tokens", True),
        ("stream_options", {"include_usage": "yes"}),
    ],
)
def test_unsupported_request_is_rejected(field, value):
    payload = {"model": MODEL, "messages": [{"role": "user", "content": "hello"}], field: value}
    with pytest.raises(web.HTTPBadRequest):
        validate_payload(payload, MODEL)


class Tokenizer:
    def apply_chat_template(self, *args, **kwargs):
        return TOKENS


@asynccontextmanager
async def origin(handler):
    app = web.Application()

    async def healthy(request):
        return web.json_response({"ready": True})

    app.router.add_get("/v1/models/{model:.*}/ready", healthy)
    app.router.add_post("/v1/chat/completions", handler)
    async with TestServer(app) as server:
        yield str(server.make_url("")).rstrip("/")


async def setup_router(app):
    router = app[ROUTER]
    add_pool(router.regions["eu"], warm=False)
    add_pool(router.regions["us"])


def test_real_chat_template_dependency():
    from tokenizers import Tokenizer as BackendTokenizer
    from tokenizers.models import WordLevel
    from transformers import PreTrainedTokenizerFast

    tokenizer = PreTrainedTokenizerFast(
        tokenizer_object=BackendTokenizer(WordLevel({"[UNK]": 0, "hello": 1}, unk_token="[UNK]")),
        chat_template="{{ messages[0]['content'] }}",
    )
    assert tokenizer.apply_chat_template([{"role": "user", "content": "hello"}], tokenize=True, return_dict=False) == [
        1
    ]


@pytest.mark.asyncio
async def test_client_disconnect_closes_upstream_and_releases_inflight():
    stopped = asyncio.Event()

    async def streaming(request):
        response = web.StreamResponse(headers={"Content-Type": "text/event-stream"})
        await response.prepare(request)
        try:
            while True:
                await response.write(b"data: progress\n\n")
                await asyncio.sleep(0.01)
        except ConnectionResetError:
            pass
        finally:
            stopped.set()
        return response

    async with origin(streaming) as eu, origin(streaming) as us:
        app = create_app(config(eu, us), Tokenizer(), consume_relays=False)
        async with TestServer(app) as server, aiohttp.ClientSession() as client:
            await setup_router(app)
            async with client.post(
                server.make_url("/v1/chat/completions"),
                json={
                    "model": MODEL,
                    "messages": [{"role": "user", "content": "hi"}],
                    "stream": True,
                },
            ) as response:
                assert await response.content.read(16) == b"data: progress\n\n"
            await asyncio.wait_for(stopped.wait(), 2)
            assert app[ROUTER].regions["us"].inflight == 0


@pytest.mark.asyncio
async def test_real_http_sse_passthrough_and_connector_fallback():
    calls = []

    async def handler(request):
        calls.append(request.headers["X-Dynamo-Federation-Hop"])
        response = web.StreamResponse(headers={"Content-Type": "text/event-stream"})
        await response.prepare(request)
        await response.write(b'data: {"text":"hello"}\n\n')
        await asyncio.sleep(0.01)
        await response.write(b"data: [DONE]\n\n")
        return response

    async with origin(handler) as eu, origin(handler) as us:
        app = create_app(config(eu, us), Tokenizer(), consume_relays=False)
        async with TestServer(app) as server, aiohttp.ClientSession() as client:
            await setup_router(app)
            payload = {"model": MODEL, "messages": [{"role": "user", "content": "hi"}], "stream": True}
            async with client.post(server.make_url("/v1/chat/completions"), json=payload) as response:
                assert response.headers["X-Dynamo-Selected-Region"] == "us"
                assert await response.read() == b'data: {"text":"hello"}\n\ndata: [DONE]\n\n'
            app[ROUTER].regions["us"].frontend = "http://127.0.0.1:1"
            async with client.post(server.make_url("/v1/chat/completions"), json=payload) as response:
                assert response.headers["X-Dynamo-Selected-Region"] == "eu"
                assert b"[DONE]" in await response.read()
            assert calls == ["1", "1"]


@pytest.mark.asyncio
async def test_failed_upstream_is_not_replayed_after_response_starts():
    local_calls = []

    async def local(request):
        local_calls.append(1)
        return web.json_response({"ok": True})

    async def broken(request):
        response = web.StreamResponse(headers={"Content-Type": "text/event-stream"})
        await response.prepare(request)
        await response.write(b"data: started\n\n")
        await asyncio.sleep(0.02)
        request.transport.close()
        return response

    async with origin(local) as eu, origin(broken) as us:
        app = create_app(config(eu, us), Tokenizer(), consume_relays=False)
        async with TestServer(app) as server, aiohttp.ClientSession() as client:
            await setup_router(app)
            async with client.post(
                server.make_url("/v1/chat/completions"),
                json={
                    "model": MODEL,
                    "messages": [{"role": "user", "content": "hi"}],
                    "stream": True,
                },
            ) as response:
                with pytest.raises(aiohttp.ClientPayloadError):
                    await response.read()
            assert local_calls == []


@pytest.mark.asyncio
async def test_real_grpc_producer_replacement_and_disconnect_invalidate_cache():
    class Relay(rpc.KvEventRelayServicer):
        def __init__(self):
            self.changes = asyncio.Queue()
            self.identity = pb.RelayIdentity(drt_instance_id=1, relay_incarnation=2)
            self.subscribed = asyncio.Queue()

        async def GetRelayInfo(self, request, context):
            return pb.RelayInfo(protocol_version=1, contract_marker=MARKER, relay=self.identity)

        async def WatchKvPoolCatalog(self, request, context):
            while True:
                revision, value = await self.changes.get()
                yield pb.KvPoolCatalogUpdate(
                    protocol_version=1,
                    contract_marker=MARKER,
                    relay=self.identity,
                    revision=revision,
                    snapshot=pb.KvPoolCatalogSnapshot(pools=[value]),
                )

        async def SubscribeServingReadiness(self, request, context):
            entry = pb.TopologyEntry(
                canonical_model_id=MODEL,
                state=2,
                present_roles=[4],
                members=[pb.TopologyMember(pool_id=descriptor().producer.pool_id, roles=[4])],
            )
            while True:
                yield pb.ServingReadinessUpdate(
                    protocol_version=1, contract_marker=MARKER, relay=self.identity, revision=1, entries=[entry]
                )
                await asyncio.sleep(0.05)

        async def SubscribeKvPool(self, request, context):
            await self.subscribed.put(request.expected_producer.producer_incarnation)
            if request.expected_producer.producer_incarnation == 1:
                yield pb.FilterUpdate(
                    protocol_version=1,
                    contract_marker=MARKER,
                    relay=self.identity,
                    producer=request.expected_producer,
                    sequence=1,
                    kind=1,
                    payload=snapshot(),
                )
            await asyncio.Event().wait()

    server = grpc.aio.server()
    relay = Relay()
    rpc.add_KvEventRelayServicer_to_server(relay, server)
    port = server.add_insecure_port("127.0.0.1:0")
    await server.start()
    region = Region("us", "http://unused", f"127.0.0.1:{port}", MODEL, 32, 0.5)
    consumer = asyncio.create_task(RelayConsumer(region).run())
    try:
        await relay.changes.put((1, descriptor()))
        assert await asyncio.wait_for(relay.subscribed.get(), 2) == 1
        for _ in range(100):
            if region.eligible_pools():
                break
            await asyncio.sleep(0.01)
        assert len(region.eligible_pools()) == 1
        await relay.changes.put((2, descriptor(incarnation=2)))
        assert await asyncio.wait_for(relay.subscribed.get(), 2) == 2
        assert region.eligible_pools() == []
        await server.stop(0)
        for _ in range(100):
            if not region.pools:
                break
            await asyncio.sleep(0.01)
        assert not region.pools and not region.ready_models
    finally:
        consumer.cancel()
        await asyncio.gather(consumer, return_exceptions=True)
        await server.stop(0)
