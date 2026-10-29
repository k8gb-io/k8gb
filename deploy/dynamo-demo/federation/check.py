# SPDX-FileCopyrightText: Copyright (c) 2026 The k8gb Contributors.
# SPDX-License-Identifier: Apache-2.0

"""Warm a direct regional origin, then prove cache-directed HTTP forwarding."""

import argparse
import json
import time
from urllib.request import Request, urlopen

MODEL = "Qwen/Qwen3-0.6B"


def payload():
    return {
        "model": MODEL,
        "messages": [
            {
                "role": "system",
                "content": ("This is a routing demonstration. Preserve this shared prefix in the KV cache. " * 64),
            },
            {"role": "user", "content": "Reply with exactly: hello"},
        ],
        "max_tokens": 8,
        "temperature": 0,
        "chat_template_kwargs": {"enable_thinking": False},
        "nvext": {"extra_fields": ["worker_id"]},
    }


def request(url, data=None):
    body = None if data is None else json.dumps(data).encode()
    with urlopen(Request(url, data=body, headers={"Content-Type": "application/json"}), timeout=180) as response:
        return dict(response.headers), response.read()


def check(warm_url, entry_url, expected_region):
    body = payload()
    _, warm = request(warm_url.rstrip("/") + "/v1/chat/completions", body)
    print(json.dumps({"phase": "warm", "response": json.loads(warm)}), flush=True)
    deadline = time.monotonic() + 30
    while True:
        _, raw = request(entry_url.rstrip("/") + "/debug/route", body)
        decision = json.loads(raw)
        if decision["region"] == expected_region and decision["matched_tokens"] >= 512:
            break
        if time.monotonic() >= deadline:
            raise AssertionError(f"No fresh remote cache hint: {decision}")
        time.sleep(0.5)
    print(json.dumps({"phase": "decision", **decision}), flush=True)
    headers, raw = request(entry_url.rstrip("/") + "/v1/chat/completions", body)
    result = json.loads(raw)
    selected = next(value for key, value in headers.items() if key.lower() == "x-dynamo-selected-region")
    assert selected == expected_region, headers
    cached = result["usage"]["prompt_tokens_details"]["cached_tokens"]
    assert cached > 0, result
    print(json.dumps({"phase": "forwarded", "headers": headers, "response": result}), flush=True)
    stream = {**body, "stream": True, "stream_options": {"include_usage": True}}
    headers, raw = request(entry_url.rstrip("/") + "/v1/chat/completions", stream)
    selected = next(value for key, value in headers.items() if key.lower() == "x-dynamo-selected-region")
    assert selected == expected_region and b"data: [DONE]" in raw, (headers, raw)
    print(json.dumps({"phase": "stream", "headers": headers, "sse": raw.decode()}), flush=True)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--warm-url", required=True)
    parser.add_argument("--entry-url", required=True)
    parser.add_argument("--expect-region", default="us")
    args = parser.parse_args()
    check(args.warm_url, args.entry_url, args.expect_region)
