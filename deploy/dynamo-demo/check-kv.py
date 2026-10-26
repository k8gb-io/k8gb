"""Check local Dynamo KV placement with two simulated workers, via its HTTP API."""

import json
import time
import urllib.request
import uuid

BASE = "http://127.0.0.1:8000"
MODEL = "Qwen/Qwen3-0.6B"


def request(content, worker=None, query_only=False):
    payload = {
        "model": MODEL,
        "messages": [{"role": "user", "content": content}],
        "max_tokens": 1,
        "stream": True,
        "nvext": {
            "extra_fields": ["worker_id"],
            "annotations": ["query_instance_id:"] if query_only else [],
        },
    }
    headers = {"Content-Type": "application/json"}
    if worker is not None:
        headers["x-dynamo-worker-instance-id"] = str(worker)
        headers["x-dynamo-dp-rank"] = "0"
    req = urllib.request.Request(
        BASE + "/v1/chat/completions", json.dumps(payload).encode(), headers
    )
    selected = None
    with urllib.request.urlopen(req, timeout=30) as response:
        for line in response:
            if not line.startswith(b"data:"):
                continue
            data = line[5:].strip()
            if data == b"[DONE]":
                continue
            chunk = json.loads(data)
            if "error" in chunk:
                raise RuntimeError(chunk)
            ids = chunk.get("nvext", {}).get("worker_id", {})
            selected = ids.get("decode_worker_id") or ids.get("prefill_worker_id") or selected
    if selected is None:
        raise RuntimeError("Dynamo returned no worker identity")
    return selected


def main():
    with urllib.request.urlopen(BASE + "/health?instances=true", timeout=10) as response:
        instances = json.load(response)["instances"]
    workers = sorted({
        instance["instance_id"] for instance in instances
        if instance["component"] == "backend" and instance["endpoint"] == "generate"
    })
    if len(workers) != 2:
        raise RuntimeError(f"Expected two local workers, found {workers}")

    # Warm a different, long prefix on each worker. Only warm-up is pinned.
    nonce = uuid.uuid4().hex
    prefixes = [f"{label}-{nonce} " * 128 for label in ("alpha", "bravo")]
    for prefix, worker in zip(prefixes, workers):
        actual = request(prefix, worker=worker)
        if actual != worker:
            raise RuntimeError(f"Warm-up pin ignored: expected {worker}, got {actual}")

    deadline = time.monotonic() + 30
    while True:
        # Query-only requests cannot accidentally warm the wrong worker while
        # waiting for asynchronous KV events. No affinity headers are supplied.
        selected = [request(prefix, query_only=True) for prefix in prefixes]
        if selected == workers:
            break
        if time.monotonic() >= deadline:
            raise RuntimeError(f"KV placement did not converge: expected {workers}, got {selected}")
        time.sleep(0.5)

    for label, prefix, worker in zip(("alpha", "bravo"), prefixes, workers):
        actual = request(prefix)
        if actual != worker:
            raise RuntimeError(f"Unpinned {label} request selected {actual}, expected {worker}")
        print(f"prefix={label} warm_worker={worker} unpinned_worker={actual}", flush=True)
    print("PASS: local KV-aware placement with simulated cache state; no cross-DC routing.")


if __name__ == "__main__":
    main()
