# SPDX-FileCopyrightText: Copyright (c) 2026 The k8gb Contributors.
# SPDX-License-Identifier: Apache-2.0

"""Experimental text-only HTTP router over Dynamo Relay CKF hints."""

import argparse
import asyncio
import logging
import time
from dataclasses import dataclass
from urllib.parse import urlsplit

import aiohttp
import yaml
from aiohttp import web
from ckf import sequence_hashes
from relay_client import Region, RelayConsumer

LOG = logging.getLogger(__name__)
ROUTER = web.AppKey("router", object)


@dataclass(frozen=True)
class Config:
    model: str
    tokenizer_path: str
    local_region: str
    regions: tuple
    block_size: int = 32
    freshness_seconds: float = 5
    remote_penalty_tokens: int = 32
    request_timeout_seconds: float = 180

    @classmethod
    def load(cls, path):
        with open(path) as source:
            raw = yaml.safe_load(source)
        config = cls(**{**raw, "regions": tuple(raw["regions"])})
        if config.model != "Qwen/Qwen3-0.6B" or not config.tokenizer_path:
            raise ValueError("prototype is pinned to Qwen/Qwen3-0.6B and its local tokenizer")
        names = [region["name"] for region in config.regions]
        if not 2 <= len(names) <= 16 or len(set(names)) != len(names) or config.local_region not in names:
            raise ValueError("configure distinct regions including local_region")
        if config.block_size != 32 or config.freshness_seconds <= 0 or config.remote_penalty_tokens < 0:
            raise ValueError("invalid prototype routing limits")
        for region in config.regions:
            url = urlsplit(region["frontend"])
            if url.scheme not in ("http", "https") or not url.hostname or url.username or url.query or url.fragment:
                raise ValueError("frontend must be a fixed regional HTTP origin")
            if url.path not in ("", "/"):
                raise ValueError("frontend origin must not include a path")
            if not region["relay"]:
                raise ValueError("relay address is required")
        return config


class Router:
    def __init__(self, config, tokenizer):
        self.config = config
        self.tokenizer = tokenizer
        self.regions = {
            item["name"]: Region(
                item["name"],
                item["frontend"].rstrip("/"),
                item["relay"],
                config.model,
                config.block_size,
                config.freshness_seconds,
            )
            for item in config.regions
        }
        self.domain = None
        self.session = None

    def choose(self, tokens, excluded=()):
        config = self.config
        local = self.regions[config.local_region]
        domains = {pool.domain for pool in local.pools.values()}
        if domains:
            self.domain = next(iter(domains)) if len(domains) == 1 else None
        hashes = sequence_hashes(tokens, config.block_size)
        candidates = []
        for region in self.regions.values():
            if region.name in excluded or not region.http_healthy():
                continue
            pools = [pool for pool in region.eligible_pools() if self.domain is not None and pool.domain == self.domain]
            is_local = region.name == config.local_region
            if not is_local and not pools:
                continue
            matched = max((pool.replica.prefix_blocks(hashes) * config.block_size for pool in pools), default=0)
            penalty = 0 if is_local else config.remote_penalty_tokens
            score = matched - penalty - region.inflight * config.block_size
            candidates.append(
                {"region": region.name, "matched_tokens": matched, "score": score, "inflight": region.inflight}
            )
        if not candidates:
            raise web.HTTPServiceUnavailable(text="No eligible regional inference endpoint")
        selected = max(
            candidates, key=lambda item: (item["score"], item["region"] == config.local_region, item["region"])
        )
        return {**selected, "entry_region": config.local_region, "candidates": candidates}

    async def poll_health(self, region):
        url = region.frontend + "/v1/models/" + self.config.model + "/ready"
        while True:
            try:
                async with self.session.get(url, timeout=aiohttp.ClientTimeout(total=2)) as response:
                    region.healthy = response.status == 200
                    await response.read()
            except (aiohttp.ClientError, TimeoutError):
                region.healthy = False
            region.healthy_updated = time.monotonic()
            await asyncio.sleep(1)

    def state(self):
        now = time.monotonic()
        return {
            "entry_region": self.config.local_region,
            "regions": {
                name: {
                    "serving_healthy": region.http_healthy(),
                    "relay_identity": region.identity,
                    "ready_pool_count": len(region.eligible_pools()),
                    "inflight": region.inflight,
                    "last_relay_error": region.error,
                    "pools": [
                        {
                            "sequence": pool.replica.sequence,
                            "complete": pool.replica.words is not None,
                            "age_seconds": round(now - pool.replica.updated, 3) if pool.replica.updated else None,
                            "namespace": pool.descriptor.serving_endpoint.namespace,
                        }
                        for pool in region.pools.values()
                    ],
                }
                for name, region in self.regions.items()
            },
        }


def validate_payload(payload, model):
    allowed = {
        "model",
        "messages",
        "stream",
        "stream_options",
        "max_tokens",
        "temperature",
        "top_p",
        "seed",
        "stop",
        "nvext",
        "chat_template_kwargs",
    }
    if not isinstance(payload, dict) or set(payload) - allowed or payload.get("model") != model:
        raise web.HTTPBadRequest(text="Prototype supports the configured base model and plain-text chat fields only")
    messages = payload.get("messages")
    if not isinstance(messages, list) or not 1 <= len(messages) <= 128:
        raise web.HTTPBadRequest(text="messages must be a nonempty list")
    for message in messages:
        if (
            not isinstance(message, dict)
            or set(message) != {"role", "content"}
            or message["role"] not in ("system", "user", "assistant")
            or not isinstance(message["content"], str)
        ):
            raise web.HTTPBadRequest(
                text="Only role/content text messages are supported; tools and multimodal input are unsupported"
            )
    if type(payload.get("stream", False)) is not bool:
        raise web.HTTPBadRequest(text="stream must be boolean")
    tokens = payload.get("max_tokens", 32)
    if type(tokens) is not int or not 1 <= tokens <= 256:
        raise web.HTTPBadRequest(text="max_tokens must be 1..256 (default 32)")
    kwargs = payload.get("chat_template_kwargs", {})
    if (
        not isinstance(kwargs, dict)
        or set(kwargs) - {"enable_thinking"}
        or any(type(v) is not bool for v in kwargs.values())
    ):
        raise web.HTTPBadRequest(text="Only the boolean enable_thinking template option is supported")
    nvext = payload.get("nvext", {})
    if (
        not isinstance(nvext, dict)
        or set(nvext) - {"extra_fields"}
        or nvext.get("extra_fields", []) not in ([], ["worker_id"])
    ):
        raise web.HTTPBadRequest(text="Only nvext.extra_fields=[worker_id] is supported")
    options = payload.get("stream_options", {})
    if (
        not isinstance(options, dict)
        or set(options) - {"include_usage"}
        or any(type(v) is not bool for v in options.values())
    ):
        raise web.HTTPBadRequest(text="Unsupported stream_options")
    return {**payload, "max_tokens": tokens}


async def prepare_request(request):
    if request.headers.get("X-Dynamo-Federation-Hop"):
        raise web.HTTPLoopDetected(text="Regional origins must point directly to a Dynamo frontend")
    if request.headers.get("x-dynamo-worker-instance-id") or request.headers.get("x-dynamo-dp-rank"):
        raise web.HTTPBadRequest(text="Worker affinity headers are unsupported at the global entry")
    try:
        payload = await request.json()
    except (ValueError, UnicodeDecodeError) as error:
        raise web.HTTPBadRequest(text="Expected a JSON request") from error
    router = request.app[ROUTER]
    payload = validate_payload(payload, router.config.model)
    tokens = await asyncio.to_thread(
        router.tokenizer.apply_chat_template,
        payload["messages"],
        tokenize=True,
        return_dict=False,
        add_generation_prompt=True,
        **payload.get("chat_template_kwargs", {}),
    )
    if len(tokens) > 4096:
        raise web.HTTPBadRequest(text="Prototype input exceeds 4096 tokens")
    return router, payload, tokens


async def debug_route(request):
    router, _, tokens = await prepare_request(request)
    return web.json_response(router.choose(tokens))


async def chat(request):
    router, payload, tokens = await prepare_request(request)
    excluded = set()
    while True:
        decision = router.choose(tokens, excluded)
        region = router.regions[decision["region"]]
        region.inflight += 1
        upstream = None
        try:
            headers = {"X-Dynamo-Federation-Hop": "1"}
            if "Authorization" in request.headers:
                headers["Authorization"] = request.headers["Authorization"]
            try:
                upstream = await router.session.post(
                    region.frontend + "/v1/chat/completions",
                    json=payload,
                    headers=headers,
                    timeout=aiohttp.ClientTimeout(total=router.config.request_timeout_seconds, sock_connect=2),
                )
            except aiohttp.ClientConnectorError:
                region.unavailable_until = time.monotonic() + router.config.freshness_seconds
                excluded.add(region.name)
                continue
            except (aiohttp.ClientError, TimeoutError) as error:
                raise web.HTTPBadGateway(text="Regional frontend failed before response headers") from error
            if upstream.status >= 500:
                region.unavailable_until = time.monotonic() + router.config.freshness_seconds
            response = web.StreamResponse(
                status=upstream.status,
                headers={
                    "Content-Type": upstream.headers.get("Content-Type", "application/json"),
                    "X-Dynamo-Entry-Region": router.config.local_region,
                    "X-Dynamo-Selected-Region": region.name,
                    "X-Dynamo-KV-Matched-Tokens": str(decision["matched_tokens"]),
                    "Cache-Control": "no-cache",
                },
            )
            LOG.info(
                "entry=%s selected=%s matched_tokens=%s",
                router.config.local_region,
                region.name,
                decision["matched_tokens"],
            )
            await response.prepare(request)
            async for chunk in upstream.content.iter_chunked(65536):
                await response.write(chunk)
            await response.write_eof()
            return response
        finally:
            if upstream is not None:
                upstream.close()
            region.inflight -= 1


async def ready(request):
    try:
        request.app[ROUTER].choose([])
    except web.HTTPServiceUnavailable:
        return web.json_response({"ready": False}, status=503)
    return web.json_response({"ready": True})


async def live(request):
    return web.json_response({"live": True})


async def state(request):
    return web.json_response(request.app[ROUTER].state())


def create_app(config, tokenizer=None, consume_relays=True):
    if tokenizer is None:
        from transformers import AutoTokenizer

        tokenizer = AutoTokenizer.from_pretrained(config.tokenizer_path, local_files_only=True)
        tokenizer.apply_chat_template(
            [{"role": "user", "content": "startup check"}],
            tokenize=True,
            return_dict=False,
            add_generation_prompt=True,
        )
    app = web.Application(client_max_size=1024 * 1024)
    router = Router(config, tokenizer)
    app[ROUTER] = router

    async def lifecycle(app):
        async with aiohttp.ClientSession() as session:
            router.session = session
            tasks = [asyncio.create_task(router.poll_health(region)) for region in router.regions.values()]
            if consume_relays:
                tasks.extend(asyncio.create_task(RelayConsumer(region).run()) for region in router.regions.values())
            try:
                yield
            finally:
                for task in tasks:
                    task.cancel()
                await asyncio.gather(*tasks, return_exceptions=True)

    app.cleanup_ctx.append(lifecycle)
    app.router.add_post("/v1/chat/completions", chat)
    app.router.add_post("/debug/route", debug_route)
    app.router.add_get("/debug/state", state)
    app.router.add_get("/ready", ready)
    app.router.add_get("/live", live)
    return app


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--config", required=True)
    parser.add_argument("--port", type=int, default=8080)
    args = parser.parse_args()
    logging.basicConfig(level=logging.INFO)
    web.run_app(create_app(Config.load(args.config)), port=args.port, handler_cancellation=True)
