# Experimental cross-region KV routing

This opt-in extension to the [Dynamo demo](dynamo-demo.md) adds an HTTP entry
router in each region. k8gb selects a reachable entry through DNS. The entry
uses Dynamo KV Relay cache hints to choose a regional Dynamo frontend, which
then selects a local worker. A request can enter EU and execute in US when US
has a longer cached prefix. Model KV tensors stay in the selected region.

```text
client -> k8gb DNS -> EU HTTP entry -> US Dynamo frontend -> US Qwen worker
                          |                    ^
                          +-- US Relay (gRPC)--+  cache hints / readiness
                          +-- EU Relay (gRPC)      cache hints / readiness
```

The implementation lives in `deploy/dynamo-demo/federation`. It is a standalone
prototype, separate from `python -m dynamo.global_router`. It does not change
the k8gb controller or require shared Dynamo discovery between regions.

## Versions and scope

The CPU runtime is pinned to `ai-dynamo==1.7.0.dev20261026` and vLLM CPU 0.31.0.
The Relay WAN API required here is newer than the baseline 1.6 demo. The copied
Apache-2.0 `relay.proto`, CBI1 parser, CKF addressing and text hash pipeline
follow [Dynamo commit 6b5b6459e21d](https://github.com/ai-dynamo/dynamo/tree/6b5b6459e21df07c1a5f5dfd9e5981bd66e783a5/lib/llm/src/kv_dc_relay).
Generated Python stubs are built from that schema. Protocol version, marker,
producer generation, hash format, block size and indexer domain are checked.
Retest this extension when changing the runtime or protocol revision.

The operator can remain at 1.6.0: the renderer uses existing v1beta1 DGD fields
and sets `runtimeVersionOverride: 1.7.0` for the nightly image tag. Qwen remains
the real `Qwen/Qwen3-0.6B` model at revision
`c1899de289a04d12100db370d81485cdf75e47ca`, with 32-token blocks and prefix
caching. Frontends, workers and relays all use the same nightly runtime.
The frontend also receives a writable `/.cache` volume for the nightly's
Model Deployment Card artifact cache when running as UID 1000.

Supported requests are base-model, plain-text `/v1/chat/completions`, including
SSE streaming, with one aggregated model pool per region. Tools, multimodal
input, LoRA, cache salts, worker pinning and
other OpenAI endpoints are outside this prototype. Requests are limited to
4096 input tokens, 256 output tokens and 1 MiB JSON. The entry applies the same
pinned Qwen chat template as the regional frontend before hashing tokens.

## Build

Run from the k8gb repository root. Use Python 3.12 or 3.13 for the helper tools
(the entry image uses Python 3.12):

```sh
python3.13 -m venv /tmp/k8gb-federation-venv
. /tmp/k8gb-federation-venv/bin/activate
pip install -r deploy/dynamo-demo/federation/requirements.txt
python -m grpc_tools.protoc -I deploy/dynamo-demo/federation \
  --python_out=deploy/dynamo-demo/federation \
  --grpc_python_out=deploy/dynamo-demo/federation \
  deploy/dynamo-demo/federation/relay.proto

docker build -t k8gb-dynamo-cpu:1.7.0.dev20261026 \
  --build-arg DYNAMO_VERSION=1.7.0.dev20261026 \
  --build-arg CPU_ARM64_IMAGE=vllm/vllm-openai-cpu:v0.31.0-arm64@sha256:045fb648a61cfcd32d4b9ca5d0c391013afadae0178d2e23ac292597f2f64036 \
  --build-arg CPU_AMD64_IMAGE=vllm/vllm-openai-cpu:v0.31.0-x86_64@sha256:8024248339dc6878daa5349344ed29d49c8a6732f6bdf7400fda32ce33e4b30b \
  -f deploy/dynamo-demo/Dockerfile deploy/dynamo-demo
docker build -t k8gb-dynamo-global-router:prototype deploy/dynamo-demo/federation
```

The second image packages only the HTTP entry, tokenizer and gRPC consumer.
Import both images into both playground clusters:

```sh
k3d image import k8gb-dynamo-cpu:1.7.0.dev20261026 k8gb-dynamo-global-router:prototype -c test-gslb1
k3d image import k8gb-dynamo-cpu:1.7.0.dev20261026 k8gb-dynamo-global-router:prototype -c test-gslb2
```

For other clusters, publish them to your registry and adjust `CPU_IMAGE` /
`ROUTER_IMAGE` in `federation/render.py`.

## Two clusters with k8gb

Start with the two k3d playground clusters and Dynamo operator installation
from the baseline demo. Allow space for CPU images on every k3d node; a 120 GiB
Docker VM disk provides room for the builds, imports and model files. Disk
pressure can garbage-collect unused imported images and prevent replacement
pods from starting.

Save one router configuration per region, following
`federation/config.example.yaml`. That example's service names work for two
namespaces in one cluster. For two clusters, use addresses reachable from pods
in both clusters:

```yaml
model: Qwen/Qwen3-0.6B
tokenizer_path: /opt/model
local_region: eu
regions:
- name: eu
  frontend: http://172.20.0.2:30800
  relay: 172.20.0.2:30551
- name: us
  frontend: http://172.20.0.3:30800
  relay: 172.20.0.3:30551
```

Replace the example IPs with the two actual node addresses. Save this as
`/tmp/federation-eu.yaml`; make a copy with `local_region: us` as
`/tmp/federation-us.yaml`. Origins must point directly to regional frontends.
Using the global hostname or another federation entry would create a loop.
The `--with-gslb` profile exposes Relay on NodePort 30551 and the direct regional
frontend on NodePort 30800; it routes the existing Ingress to the new HTTP entry.
These ports must be free in the demo clusters.

Find the nodes' `INTERNAL-IP` addresses with:

```sh
kubectl --context k3d-test-gslb1 get nodes -o wide
kubectl --context k3d-test-gslb2 get nodes -o wide
```

When upgrading a running baseline demo, stop its frontend and workers before
applying the new profile to avoid loading old and new models simultaneously:

```sh
for context in k3d-test-gslb1 k3d-test-gslb2; do
  kubectl --context "$context" -n dynamo-demo patch dgd dynamo-demo --type=json \
    -p '[{"op":"test","path":"/spec/components/0/name","value":"Frontend"},{"op":"test","path":"/spec/components/1/name","value":"worker"},{"op":"replace","path":"/spec/components/0/replicas","value":0},{"op":"replace","path":"/spec/components/1/replicas","value":0}]'
  kubectl --context "$context" -n dynamo-demo wait --for=delete pod \
    -l nvidia.com/dynamo-graph-deployment-name=dynamo-demo --timeout=180s
done
```

Render and apply the regional profiles:

```sh
python deploy/dynamo-demo/federation/render.py \
  --config /tmp/federation-eu.yaml --with-gslb > /tmp/federation-eu-resources.yaml
python deploy/dynamo-demo/federation/render.py \
  --config /tmp/federation-us.yaml --with-gslb > /tmp/federation-us-resources.yaml
kubectl --context k3d-test-gslb1 apply -f /tmp/federation-eu-resources.yaml
kubectl --context k3d-test-gslb2 apply -f /tmp/federation-us-resources.yaml

for context in k3d-test-gslb1 k3d-test-gslb2; do
  kubectl --context "$context" -n dynamo-demo wait \
    dgd/dynamo-demo --for=condition=Ready --timeout=1200s
  kubectl --context "$context" -n dynamo-demo rollout status \
    deployment/dynamo-global-router --timeout=300s
done
```

This intentionally upgrades the demo DGD's runtime and changes its Ingress
backend. Use dedicated demo namespaces. The renderer sets
`globalDynamoNamespace: true` on frontend and worker to give both regional
indexers the same logical namespace prefix (`dynamo`, with an operator-generated
worker-spec hash suffix). Relay watches all logical namespaces inside its own
Kubernetes namespace. Kubernetes discovery remains
isolated to each Kubernetes namespace. Keep worker specs, model paths and
tokenizer revisions identical so cache domains compare equal. Do not put
unrelated graphs into these dedicated namespaces.

The default is one worker per region. Use `--worker-replicas 2` on both renders
to retain the baseline's four-worker topology; allow at least 32 GiB for that
playground. Wait for both DGDs and HTTP entries to become ready. Relay uses the
DGD's operator-created discovery service account. No GPU, etcd, NATS or Grove
is required.

This profile uses plaintext HTTP and gRPC on a private demo network. It is not
a public deployment configuration. Authentication, TLS and tenant isolation
must be added before using it beyond the playground. Debug routes expose
routing state and should remain private.

## Prove forwarding and cache reuse

Port-forward the EU entry and the direct US frontend in two terminals:

```sh
kubectl --context k3d-test-gslb1 -n dynamo-demo \
  port-forward svc/dynamo-global-router 18080:8080
kubectl --context k3d-test-gslb2 -n dynamo-demo \
  port-forward svc/dynamo-demo-frontend 18001:8000
```

Then run:

```sh
python deploy/dynamo-demo/federation/check.py \
  --warm-url http://localhost:18001 --entry-url http://localhost:18080 \
  --expect-region us
curl -s http://localhost:18080/debug/state
```

The check warms a long prefix directly in US, waits for a fresh remote CKF
match of at least 512 tokens, sends the same request through EU, and asserts
`X-Dynamo-Selected-Region: us` plus positive
`usage.prompt_tokens_details.cached_tokens`. It also checks SSE reaches
`data: [DONE]`. Each run generates a fresh prefix so previously warmed EU caches
do not affect the demonstration.

For a frontend failure, set the US Frontend's `replicas` to zero in its DGD.
The entry health poll should remove US and serve new requests locally. For a
metadata outage, scale the US `dynamo-relay` Deployment to zero: its hints must
become ineligible immediately on stream closure, or within five seconds for
a silent connection loss. Restore the replica counts after testing.

## Demonstrate DNS entry failover

Create the client pod and use the global-hostname curl command from
[step 2 of the baseline guide](dynamo-demo.md#2-send-a-request-through-the-global-hostname).
Add `-D -` to curl to see `X-Dynamo-Entry-Region` and
`X-Dynamo-Selected-Region`: the entry region follows DNS, while the selected
region follows the request's cache hints. Initially the entry should be `eu`.

Stop the EU entry:

```sh
kubectl --context k3d-test-gslb1 -n dynamo-demo scale \
  deployment/dynamo-global-router --replicas=0
dig @127.0.0.1 -p 1053 dynamo.cloud.example.com +short +tcp
```

Repeat curl and dig until the response has `X-Dynamo-Entry-Region: us` and
BIND9 returns secondary ingress addresses. Temporary HTTP errors are possible
during convergence. Repeat curl with `"stream":true` to verify streaming too.

Restore the entry and wait for rollout:

```sh
kubectl --context k3d-test-gslb1 -n dynamo-demo scale \
  deployment/dynamo-global-router --replicas=1
kubectl --context k3d-test-gslb1 -n dynamo-demo rollout status \
  deployment/dynamo-global-router --timeout=300s
```

Repeat curl and dig until the entry header and DNS addresses return to EU.
Restore the entry even if you stop the walkthrough early. This tests entry
failure; existing streams are not migrated and DNS caches affect recovery time.

## Routing and failure behavior

Each entry consumes a complete pool catalog, independent serving-readiness
snapshots, and per-pool CKF snapshots/deltas/heartbeats. Partial snapshots,
checksum failures, sequence gaps, producer replacement and relay restarts
invalidate the affected replica. Reconnection starts with a new snapshot.
Unknown or stale remote readiness excludes the remote region. A healthy local
frontend can continue serving when Relay metadata is unavailable.

The score is matched prefix tokens minus 32 tokens for a remote region and
32 tokens for each request currently in flight through this entry. Equal
scores prefer the local region. This is a demonstration heuristic, without
WAN latency prediction or global backend load accounting; the Relay load RPCs
are not consumed. CKFs are approximate hints and may produce false positives.
The selected regional KV router still chooses the actual worker.

HTTP bodies and SSE chunks are forwarded without rewriting. Client disconnects
cancel the upstream request. Only connection-establishment failures permit
trying another region; ambiguous failures or an already-started response are
never replayed. `/ready` reflects whether the entry has an eligible inference
origin; `/live` only reports process liveness. k8gb continues to handle DNS-level
entry failover, independently of per-request cache selection.
With this profile, stopping a local frontend can leave its entry ready when a
remote origin remains eligible. To demonstrate DNS-level entry failover, stop
the entry Deployment itself rather than only its local frontend.

## Regression tests

```sh
pip install pytest pytest-asyncio
python -m pytest -q deploy/dynamo-demo/federation/test_router.py
```

Tests include CBI1 corruption and partial snapshots, stale hints/readiness,
incompatible domains, local autonomy, real gRPC producer replacement, and real
HTTP SSE forwarding, connection failure/timeout fallback, unsupported model
targets and readiness dependencies, and no replay after partial responses.

## Measured prototype validation

### Two k3d clusters

On 2026-10-29, the NodePort/Ingress/Gslb profile passed acceptance in the existing
two-cluster playground on ARM64 Colima: k3s v1.36.2-k3s1, k8gb v1.0.0,
Dynamo operator 1.6.0, runtime 1.7.0.dev20261026 and vLLM CPU 0.31.0.
The VM had 8 CPUs, 24 GiB RAM and a 120 GiB disk, with one real Qwen worker per
region. Both clusters share that VM; this does not measure WAN behavior.

| Check | Observed result |
| --- | --- |
| Warm US, enter EU | US selected from a 992-token CKF match; EU match was zero |
| Forwarded inference and SSE | Same US worker as warm-up; generated `hello`, with 992 cached tokens out of 1014 prompt tokens |
| Stop and restart US Relay | US hints withdrawn, EU served locally, then a fresh snapshot restored remote routing |
| Stop and restore US frontend | US excluded while EU continued inference; remote routing recovered afterwards |
| Global hostname | DNS selected the EU entry, which forwarded the warmed request to US |
| Stop EU entry | HTTP entry and delegated BIND9 addresses moved to US; regular and streaming inference succeeded |
| Restore EU entry | HTTP entry and BIND9 addresses returned to EU; the warmed request still executed in US |
| Worker continuity | Both worker pod identities and restart counts remained unchanged during the successful run |

The regression suite passes 25 cases, including connection-establishment timeout
fallback and rejection of unknown model targets and readiness dependencies.

### Initial single-cluster validation

On 2026-10-27, this profile was exercised in an isolated Linux ARM64 Kind cluster
(Kubernetes 1.30.13), with two separate Kubernetes namespaces representing EU
and US and one real CPU Qwen worker in each. Docker had 8 CPUs and 16 GiB RAM;
the test cluster used about 10.6 GiB after startup. This validates the regional
protocol and HTTP path; it is not a two-datacenter deployment or a throughput
benchmark. The existing 1.6 two-cluster DNS acceptance run remains documented
in the baseline guide.

| Check | Observed result |
| --- | --- |
| Warm US, enter through EU | US selected from a 960-token CKF match; EU match was zero |
| Non-streaming inference | Same US worker as warm-up; 960 cached tokens out of 983 prompt tokens |
| SSE inference | US selected; generated `hello`, usage with 960 cached tokens, and `[DONE]` |
| Stop US Relay | US hints invalidated; EU continued real local inference |
| Silent EU-to-US partition | DROP rules in the EU entry pod blocked both US Relay and frontend; remote eligibility disappeared after about 5.3 seconds in this run |
| During that partition | US origin stayed healthy via a separate path; EU answered locally |
| Restore connectivity | A fresh subscription restored the 960-token US match and remote selection |
| Stop US frontend in DGD | US excluded despite fresh Relay hints; EU answered locally |
| Regression suite | 19 tests passed, including gRPC producer replacement, client disconnect and partial-response handling |

The partition affected only the EU entry pod's traffic to the two US service
addresses. It did not simulate host loss, inter-datacenter latency, asymmetric
routing in both directions, or a whole-region failure. Four-worker execution
of this nightly profile, x86-64 inference and production WAN behavior remain
unvalidated. The two-cluster run above adds entry DNS failover and failback;
it does not establish a recovery-time SLO.
