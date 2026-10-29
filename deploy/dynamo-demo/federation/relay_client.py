# SPDX-FileCopyrightText: Copyright (c) 2026 The k8gb Contributors.
# SPDX-License-Identifier: Apache-2.0

"""Independent, fail-closed catalog, readiness and CKF stream consumers."""

import asyncio
import logging
import struct
import time
from dataclasses import dataclass, field

import blake3
import grpc
import relay_pb2 as pb
import relay_pb2_grpc as rpc
from ckf import Format, InvalidPublication, Replica

MARKER = 0x4B565231
LOG = logging.getLogger(__name__)


def dc_identity(name):
    raw = name.encode()
    material = b"dynamo/indexer-dc/v1" + struct.pack("<I", len(raw)) + raw
    return int.from_bytes(blake3.blake3(material).digest()[:8], "little")


def wire_key(value):
    return value.SerializeToString(deterministic=True)


def relay_key(value):
    return value.drt_instance_id, value.relay_incarnation


def envelope(message, identity):
    if message.protocol_version != 1 or message.contract_marker != MARKER or relay_key(message.relay) != identity:
        raise InvalidPublication("relay contract or incarnation changed")


async def next_message(stream, timeout):
    value = await asyncio.wait_for(stream.read(), timeout)
    if value is grpc.aio.EOF:
        raise ConnectionError("relay stream closed")
    return value


@dataclass
class Pool:
    descriptor: object
    replica: Replica
    canonical_models: frozenset

    @property
    def key(self):
        return wire_key(self.descriptor.producer.pool_id)

    @property
    def domain(self):
        return wire_key(self.descriptor.producer.pool_id.indexer_domain)

    @classmethod
    def from_descriptor(cls, descriptor, region, model, block_size):
        producer = descriptor.producer
        identity = producer.pool_id
        if identity.identity_version != 1 or identity.dc_id != dc_identity(region):
            raise InvalidPublication("unexpected pool identity")
        for digest in (identity.indexer_domain.cache_semantics, identity.indexer_domain.routing_scope):
            if len(digest.digest) != 16 or digest.source not in (1, 2):
                raise InvalidPublication("unsupported indexer domain")
        if not producer.producer_incarnation or not producer.layout_generation:
            raise InvalidPublication("missing producer generation")
        semantics = descriptor.query_semantics
        if semantics.hash_format != 1 or semantics.kv_block_size != block_size:
            raise InvalidPublication("prototype requires DYNAMO_STANDARD_V1 and the configured block size")
        if set(descriptor.pool_roles) != {4}:
            raise InvalidPublication("prototype supports aggregated serving pools only")
        targets = set()
        names = set()
        for registration in descriptor.registrations:
            if registration.target.WhichOneof("target") != "base":
                continue
            if model == registration.canonical_model_id or model in registration.aliases:
                targets.add(registration.target.base.base_model)
                names.add(registration.canonical_model_id)
        if len(targets) != 1 or not names:
            raise InvalidPublication("model binding absent or ambiguous")
        layout = producer.ckf_format
        ckf = Format(
            layout.format_version, layout.seed, layout.bucket_count, layout.fingerprint_bits, layout.slots_per_bucket
        )
        return cls(descriptor, Replica(ckf, identity.dc_id), frozenset(names))


@dataclass
class Region:
    name: str
    frontend: str
    relay: str
    model: str
    block_size: int
    freshness: float
    pools: dict = field(default_factory=dict)
    ready_models: set = field(default_factory=set)
    readiness_updated: float = 0
    healthy_updated: float = 0
    healthy: bool = False
    identity: tuple = None
    error: str = None
    inflight: int = 0
    unavailable_until: float = 0

    def http_healthy(self):
        now = time.monotonic()
        return self.healthy and now - self.healthy_updated < self.freshness and now >= self.unavailable_until

    def eligible_pools(self):
        now = time.monotonic()
        if now - self.readiness_updated >= self.freshness:
            return []
        return [
            pool
            for pool in self.pools.values()
            if pool.replica.words is not None
            and now - pool.replica.updated < self.freshness
            and any((pool.key, model) in self.ready_models for model in pool.canonical_models)
        ]

    def reset(self):
        self.pools.clear()
        self.ready_models.clear()
        self.readiness_updated = 0
        self.identity = None


class RelayConsumer:
    def __init__(self, region):
        self.region = region
        self.pool_tasks = {}

    async def run(self):
        region = self.region
        options = [("grpc.max_receive_message_length", 8 * 1024 * 1024)]
        async with grpc.aio.insecure_channel(region.relay, options=options) as channel:
            stub = rpc.KvEventRelayStub(channel)
            while True:
                readiness = None
                catalog = None
                try:
                    info = await stub.GetRelayInfo(
                        pb.RelayInfoRequest(contract_marker=MARKER), timeout=region.freshness
                    )
                    region.identity = relay_key(info.relay)
                    envelope(info, region.identity)
                    if not all(region.identity):
                        raise InvalidPublication("missing relay identity")
                    readiness = asyncio.create_task(self.watch_readiness(stub, region.identity))
                    catalog = stub.WatchKvPoolCatalog(
                        pb.WatchKvPoolCatalogRequest(
                            subscriber_id=f"demo-{region.name}-catalog", contract_marker=MARKER
                        )
                    )
                    revision = -1
                    async for update in catalog:
                        envelope(update, region.identity)
                        if update.revision < revision:
                            raise InvalidPublication("catalog revision regressed")
                        revision = update.revision
                        desired = {}
                        for descriptor in update.snapshot.pools:
                            try:
                                pool = Pool.from_descriptor(descriptor, region.name, region.model, region.block_size)
                            except InvalidPublication as error:
                                LOG.warning("region=%s unsupported pool: %s", region.name, error)
                                continue
                            if pool.key in desired:
                                raise InvalidPublication("duplicate pool in catalog")
                            desired[pool.key] = pool
                        if len(desired) > 1:
                            raise InvalidPublication("prototype requires one aggregated model pool per region")
                        for key in set(region.pools) - set(desired):
                            await self.retire(key)
                        for key, pool in desired.items():
                            current = region.pools.get(key)
                            if current and wire_key(current.descriptor.producer) == wire_key(pool.descriptor.producer):
                                current.descriptor = pool.descriptor
                                current.canonical_models = pool.canonical_models
                                continue
                            await self.retire(key)
                            region.pools[key] = pool
                            self.pool_tasks[key] = asyncio.create_task(self.watch_pool(stub, pool, region.identity))
                        region.error = None
                    raise ConnectionError("catalog stream closed")
                except asyncio.CancelledError:
                    raise
                except (grpc.RpcError, InvalidPublication, TimeoutError, ConnectionError) as error:
                    if asyncio.current_task().cancelling():
                        raise asyncio.CancelledError
                    region.error = str(error)
                    LOG.warning("region=%s catalog unavailable: %s", region.name, error)
                finally:
                    region.reset()
                    if catalog is not None:
                        catalog.cancel()
                    if readiness is not None:
                        readiness.cancel()
                        await asyncio.gather(readiness, return_exceptions=True)
                    for key in list(self.pool_tasks):
                        await self.retire(key)
                await asyncio.sleep(1)

    async def retire(self, key):
        pool = self.region.pools.pop(key, None)
        if pool is not None:
            pool.replica.clear()
        task = self.pool_tasks.pop(key, None)
        if task is not None:
            task.cancel()
            await asyncio.gather(task, return_exceptions=True)

    async def watch_pool(self, stub, pool, identity):
        region = self.region
        while True:
            stream = stub.SubscribeKvPool(
                pb.SubscribeKvPoolRequest(
                    subscriber_id=f"demo-{region.name}-pool",
                    expected_producer=pool.descriptor.producer,
                    contract_marker=MARKER,
                )
            )
            try:
                while True:
                    update = await next_message(stream, region.freshness)
                    envelope(update, identity)
                    if wire_key(update.producer) != wire_key(pool.descriptor.producer):
                        raise InvalidPublication("pool producer changed")
                    pool.replica.apply(update.kind, update.payload, update.base_sequence, update.sequence)
            except asyncio.CancelledError:
                raise
            except (grpc.RpcError, InvalidPublication, TimeoutError, ConnectionError) as error:
                if asyncio.current_task().cancelling():
                    raise asyncio.CancelledError
                region.error = str(error)
                LOG.warning("region=%s pool stream invalidated: %s", region.name, error)
            finally:
                pool.replica.clear()
                stream.cancel()
            await asyncio.sleep(1)

    async def watch_readiness(self, stub, identity):
        region = self.region
        while True:
            stream = stub.SubscribeServingReadiness(
                pb.SubscribeServingReadinessRequest(
                    subscriber_id=f"demo-{region.name}-readiness", contract_marker=MARKER
                )
            )
            revision = -1
            try:
                while True:
                    update = await next_message(stream, region.freshness)
                    envelope(update, identity)
                    if update.revision < revision:
                        raise InvalidPublication("readiness revision regressed")
                    revision = update.revision
                    ready = set()
                    for entry in update.entries:
                        if entry.state != 2 or set(entry.present_roles) != {4} or entry.missing_roles:
                            continue
                        for member in entry.members:
                            if member.HasField("pool_id") and set(member.roles) == {4}:
                                ready.add((wire_key(member.pool_id), entry.canonical_model_id))
                    region.ready_models = ready
                    region.readiness_updated = time.monotonic()
            except asyncio.CancelledError:
                raise
            except (grpc.RpcError, InvalidPublication, TimeoutError, ConnectionError) as error:
                if asyncio.current_task().cancelling():
                    raise asyncio.CancelledError
                region.error = str(error)
                LOG.warning("region=%s readiness unavailable: %s", region.name, error)
            finally:
                region.ready_models.clear()
                region.readiness_updated = 0
                stream.cancel()
            await asyncio.sleep(1)
