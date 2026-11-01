<!-- omit in toc -->
# Local playground for testing and development

- [Environment prerequisites](#environment-prerequisites)
- [Running project locally](#running-project-locally)
- [Verify installation](#verify-installation)
- [Run integration tests](#run-integration-tests)
- [Cleaning](#cleaning)
- [AI inference resilience demo](#ai-inference-resilience-demo)
- [Sample demo](#sample-demo)
  - [Gateway API walkthrough (recommended)](#gateway-api-walkthrough-recommended)
  - [Multi-Service Ingress](#multi-service-ingress)
  - [Round Robin (Ingress)](#round-robin-ingress)
  - [Failover (Ingress)](#failover-ingress)

 ---
**NOTE**

This tutorial relies on some makefile targets, to be able to fully understand what's happening under the covers, check the Makefile
[source code](https://github.com/k8gb-io/k8gb/blob/master/Makefile).
Or you may want to run a target first with `-n` switch that will print what it is going to do (example: `make -n test-round-robin`).
For more user-centric targets in that makefile consult `make help`.

---

## Environment prerequisites

- [Install **Go**](https://go.dev/dl/) — use the version pinned in [`go.mod`](https://github.com/k8gb-io/k8gb/blob/master/go.mod) (CI uses `go-version-file: ./go.mod`)

- [Install **yq v4**](https://github.com/mikefarah/yq) (Mike Farah's implementation), required by `make test` to generate CRD manifests.

- [Install **Git**](https://git-scm.com/downloads)

- [Install **Docker**](https://docs.docker.com/get-docker/)
  > Ensure you are able to push/pull from your docker registry

  > To run multiple clusters, reserve 8GB of memory*
    ![](images/docker_settings.png)
      <div>
        <sup><i>* above screenshot is provided for <strong>Docker for Mac</strong>, options for other Docker distributions may vary
        </i></sup>
      </div>

 - [install **kubectl**](https://kubernetes.io/docs/tasks/tools/install-kubectl/) to operate k8s clusters

 - [install **helm3**](https://helm.sh/docs/intro/install/) to deploy k8gb and related test workloads

 - [install **k3d**](https://k3d.io/#installation) to run local [k3s](https://k3s.io/) clusters (minimum v5.3.0 version is required)

 - Install **GNU coreutils** for `timeout`, used to bound image imports. It is included in most Linux distributions; on macOS, run `brew install coreutils` to provide `gtimeout`.

 - [install **golangci-lint**](https://golangci-lint.run/docs/welcome/install/) for code quality checks

## Running project locally

To spin-up a local environment using two k3s clusters and deploy a test application to both clusters, execute the command below:
```sh
make deploy-full-local-setup
```

The setup installs the following components on each workload cluster:

| Component | Version | Role |
|---|---|---|
| [Traefik](https://traefik.io/) | `36.3.0` (Helm) | Ingress controller (replaces retired ingress-nginx) |
| [Gateway API CRDs](https://gateway-api.sigs.k8s.io/) | `v1.6.3` | Standard Gateway/HTTPRoute/TCPRoute/UDPRoute/TLSRoute APIs |
| [Istio](https://istio.io/) | `v1.31.1` | GatewayClass `istio`; Ingress Gateway on port `8000` |
| [k8gb](https://k8gb.io/) | current | GSLB operator |

> **ARM64 / macOS**: the setup is validated on Apple Silicon (`arm64`). k3d pulls multi-arch k3s images. Ensure Docker Desktop has at least 8 GB RAM allocated and Rosetta is not required by k3s images (all images are native `arm64`).

> **Linux CI**: the same `make deploy-full-local-setup` target runs in CI. Port mappings use direct host-network binding on the k3d agent nodes for Traefik (ports 80/443) and standard loadbalancer ports for Istio (8000/8080), Prometheus and CoreDNS.

By default, local setup also applies a small set of legacy `k8gb.absa.oss/v1beta1` migration demo resources on each local test cluster. The examples reuse workloads from the standard setup and cover referenced Ingresses, embedded Ingress migration, and annotation-sensitive LoadBalancer Service reconciliation. The setup prints legacy/canonical manifests so you can see migration output immediately.

The Service examples share one LoadBalancer Service backed by the existing `multiservice-blue` pods:

- `legacy-service-runtime-demo` remains in legacy runtime mode.
- `legacy-service-migration-demo` is migrated to `k8gb.io/v1beta1`.

Verify that both paths retain `k8gb.io/hostname` and reconcile a `DNSEndpoint`:

```sh
kubectl get gslb.k8gb.absa.oss -n test-gslb legacy-service-runtime-demo \
  -o jsonpath='{.metadata.annotations.k8gb\.io/hostname}{" / "}{.status.hosts}{"\n"}'
kubectl get dnsendpoint -n test-gslb legacy-service-runtime-demo

kubectl get gslb.k8gb.io -n test-gslb legacy-service-migration-demo \
  -o jsonpath='{.metadata.annotations.k8gb\.io/hostname}{" / "}{.status.hosts}{"\n"}'
kubectl get dnsendpoint -n test-gslb legacy-service-migration-demo
```

If you want to skip this demo output:

```sh
make deploy-full-local-setup SHOW_LEGACY_MIGRATION_DEMO=false
```

## Verify installation

If local setup runs well, check if clusters are correctly installed

```sh
kubectl cluster-info --context k3d-edgedns && kubectl cluster-info --context k3d-test-gslb1 && kubectl cluster-info --context k3d-test-gslb2
```

Cluster [test-gslb1](https://github.com/k8gb-io/k8gb/tree/master/k3d/test-gslb1.yaml) is exposing external DNS on default port `:5053`
while [test-gslb2](https://github.com/k8gb-io/k8gb/tree/master/k3d/test-gslb2.yaml) on port `:5054`.

Cluster [edgedns](https://github.com/k8gb-io/k8gb/tree/master/k3d/edge-dns.yaml) runs BIND and acts as EdgeDNS holding Delegated Zone for our test setup and answers
on port `:1053`.

```sh
dig @localhost -p 1053 roundrobin.cloud.example.com +short +tcp
```
Should return ***two A records*** from both clusters (IP addresses and order may differ):
```
172.20.0.2
172.20.0.5
172.20.0.4
172.20.0.6
```

You can verify that correct IP addresses of all the nodes in both clusters were populated:
```sh
for c in k3d-test-gslb{1,2}; do kubectl get no -ocustom-columns="NAME:.metadata.name,IP:status.addresses[0].address" --context $c; done
```

Returns a result similar to:
```
NAME                      IP
k3d-test-gslb1-agent-0    172.20.0.2
k3d-test-gslb1-server-0   172.20.0.4
NAME                      IP
k3d-test-gslb2-server-0   172.20.0.6
k3d-test-gslb2-agent-0    172.20.0.5
```

Or you can ask specific CoreDNS instance for its local targets:
```sh
dig -p 5053 @localhost localtargets-roundrobin.cloud.example.com && \
dig -p 5054 @localhost localtargets-roundrobin.cloud.example.com
```
As expected result you should see **two A records** divided between both clusters.
```sh
...
...
;; ANSWER SECTION:
localtargets-roundrobin.cloud.example.com. 30 IN A 172.20.0.4
localtargets-roundrobin.cloud.example.com. 30 IN A 172.20.0.2
...
...
localtargets-roundrobin.cloud.example.com. 30 IN A 172.20.0.5
localtargets-roundrobin.cloud.example.com. 30 IN A 172.20.0.6
```

### Verify Traefik Ingress address publication

Traefik runs with `hostNetwork: true` on the k3d agent node, binding directly to ports 80 and 443. Because no LoadBalancer service is used, addresses are derived from the node IP rather than a service `status.loadBalancer`.

Verify Traefik is ready and the Ingress objects have been picked up:
```sh
kubectl -n k8gb get pod -l app.kubernetes.io/name=traefik
kubectl -n k8gb get ingressclass traefik
# Ingress resources in test-gslb should show the agent node address
kubectl -n test-gslb get ingress roundrobin-ingress -o wide
```

Both clusters have [podinfo](https://github.com/stefanprodan/podinfo) installed on top.
Hit the Ingress directly through the host port to confirm reachability:
```sh
curl localhost:80 -H "Host:roundrobin.cloud.example.com" && curl localhost:81 -H "Host:roundrobin.cloud.example.com"
```

The local setup also deploys small multi-service Ingresses backed by two services and canonical `k8gb.io/v1beta1` GSLBs using both service health policies.

```sh
curl localhost:80/blue -H "Host:multiservice-all.cloud.example.com"
curl localhost:80/green -H "Host:multiservice-any.cloud.example.com"
kubectl --context k3d-test-gslb1 -n test-gslb get gslb \
  multiservice-gslb-all multiservice-gslb-any \
  -o custom-columns=NAME:.metadata.name,POLICY:.spec.serviceHealthPolicy,HEALTH:.status.serviceHealth
```

### Verify Gateway status.addresses

The Istio Gateway (`gatewayapi` in namespace `istio-ingress`) is served by the Istio Ingress Gateway on port `8000`. Verify `Programmed`, `Accepted`, and that `status.addresses` is populated on both clusters:

```sh
for c in k3d-test-gslb{1,2}; do
  echo "=== $c ==="
  kubectl --context $c -n istio-ingress get gateway gatewayapi \
    -o jsonpath='Programmed={.status.conditions[?(@.type=="Programmed")].status}{"\n"}Accepted={.status.conditions[?(@.type=="Accepted")].status}{"\n"}Addresses={.status.addresses[*].value}{"\n"}'
done
```

Expected output (addresses will differ per environment):
```
=== k3d-test-gslb1 ===
Programmed=True
Accepted=True
Addresses=172.20.0.3
=== k3d-test-gslb2 ===
Programmed=True
Accepted=True
Addresses=172.20.0.7
```

Verify an HTTPRoute has `ResolvedRefs` and `Accepted`:
```sh
kubectl --context k3d-test-gslb1 -n test-gslb-gatewayapi get httproute roundrobin-httproute \
  -o jsonpath='{range .status.parents[*]}parent={.parentRef.name} Accepted={.conditions[?(@.type=="Accepted")].status} ResolvedRefs={.conditions[?(@.type=="ResolvedRefs")].status}{"\n"}{end}'
```

Confirm reachability through the Istio port:
```sh
curl localhost:8000 -H "Host:roundrobin-httproute.cloud.example.com"
curl localhost:8001 -H "Host:roundrobin-httproute.cloud.example.com"
```

> **Port coexistence**: Traefik binds ports 80/443 on the agent node (hostNetwork). Istio binds port 8000/8080 via the k3d loadbalancer. CoreDNS uses port 5053/5054. BIND9/EdgeDNS uses port 1053. These ranges do not overlap.

## Run integration tests

There is wide range of scenarios which **GSLB** provides and all of them are covered within [tests](https://github.com/k8gb-io/k8gb/tree/master/terratest).
To check whether everything is running properly execute [terratest](https://terratest.gruntwork.io/) :

```sh
make terratest
```

Terratest and Chainsaw test fixtures reference Ingress resources via the `traefik` IngressClass. Resolver unit tests are independent of the controller. TCP/UDP/TLS forwarding tests target the Istio Gateway port (`8000`), not the Traefik port.

## Cleaning

Clean up your local development clusters with
```sh
make destroy-full-local-setup
```

## AI inference resilience demo

The local setup can also deploy a lightweight Ollama inference endpoint to demonstrate k8gb failover for AI traffic:

```sh
make deploy-full-local-setup FULL_LOCAL_SETUP_WITH_AI_DEMO=true
```

If the local setup is already running:

```sh
make ai-inference-demo AI_DEMO_ACTION=deploy
make ai-inference-demo AI_DEMO_ACTION=probe
make ai-inference-demo AI_DEMO_ACTION=failover
make ai-inference-demo AI_DEMO_ACTION=probe
make ai-inference-demo AI_DEMO_ACTION=failback
```

The demo caches the downloaded model in a PVC, so failback does not re-download the model. If model startup is slow, inspect progress with:

```sh
make ai-inference-demo AI_DEMO_ACTION=logs
```

The AI inference demo uses the `traefik` IngressClass. Streaming SSE responses pass through Traefik without modification; the default read/write timeouts are sufficient for the demo workload.

See [AI Inference Resilience Demo](ai-inference-demo.md) for the full local and real-environment flow.

## Sample demo

### Gateway API walkthrough (recommended)

The recommended path for new deployments uses **Gateway API**: a `Gateway` with `gatewayClassName: istio` (already installed by the setup), `HTTPRoute` resources that attach to it, and `Gslb` objects that reference the `HTTPRoute` via `spec.resourceRef`.

The manifests are in [`deploy/gslb/`](https://github.com/k8gb-io/k8gb/tree/master/deploy/gslb/) and are applied automatically by `deploy-test-apps`.

#### Round robin via HTTPRoute

```sh
# Both clusters serve roundrobin-httproute.cloud.example.com via Istio port 8000/8001
curl localhost:8000 -H "Host:roundrobin-httproute.cloud.example.com"
curl localhost:8001 -H "Host:roundrobin-httproute.cloud.example.com"
```

Check GSLB status for the HTTPRoute-backed resource:
```sh
kubectl --context k3d-test-gslb1 -n test-gslb-gatewayapi get gslb roundrobin-httproute -o yaml | \
  yq '{"status": .status}'
```

#### Failover via HTTPRoute

The `failover-httproute` Gslb uses `primaryGeoTag: eu`. Cluster 1 is `eu`; cluster 2 is `us`.

```sh
# Verify primary is serving
curl localhost:8000 -H "Host:failover-httproute.cloud.example.com"

# Scale down eu primary
kubectl --context k3d-test-gslb1 -n test-gslb-gatewayapi scale deploy frontend-podinfo --replicas=0

# Wait for DNS TTL (30s) and confirm traffic shifts to us
curl localhost:8000 -H "Host:failover-httproute.cloud.example.com"

# Restore
kubectl --context k3d-test-gslb1 -n test-gslb-gatewayapi scale deploy frontend-podinfo --replicas=1
```

Alternatively, use the automated failover helper (targets the Ingress-based `failover-ingress` by default; override with env vars to target the HTTPRoute demo):
```sh
GSLB_NAME=failover-httproute NAMESPACE=test-gslb-gatewayapi make demo-failover
```

> **Known limitation — controller outage detection**: k8gb derives health from Service endpoint readiness, not from the gateway or ingress controller process itself. If the Traefik or Istio pod crashes while application pods remain ready, k8gb continues advertising the cluster as healthy. This gap is tracked in [#1754](https://github.com/k8gb-io/k8gb/issues/1754) and is not addressed by this migration.

---

### Multi-Service Ingress

`multiservice-all.cloud.example.com` and `multiservice-any.cloud.example.com` both route `/blue` and `/green` to separate Services via Traefik. The `All` GSLB requires both Services to be healthy, while the `Any` GSLB remains healthy as long as at least one backing Service is healthy.

```sh
kubectl --context k3d-test-gslb1 -n test-gslb scale deploy multiservice-green --replicas=0
kubectl --context k3d-test-gslb1 -n test-gslb get gslb \
  multiservice-gslb-all multiservice-gslb-any \
  -o custom-columns=NAME:.metadata.name,POLICY:.spec.serviceHealthPolicy,HEALTH:.status.serviceHealth
kubectl --context k3d-test-gslb1 -n test-gslb scale deploy multiservice-green --replicas=1
```

### Round Robin (Ingress)

Both clusters have [podinfo](https://github.com/stefanprodan/podinfo) installed on top, where each
cluster has been tagged to serve a different region. In this demo we will hit podinfo by `wget -qO - roundrobin.cloud.example.com` and depending
on the region, podinfo will return **us** or **eu**. In the current round robin implementation IP addresses are randomly picked.
See [Gslb manifest with round robin strategy](https://github.com/k8gb-io/k8gb/tree/master/deploy/gslb/k8gb.io_v1beta1_gslb_cr_roundrobin_ingress.yaml)

Try to run the following command several times and watch the `message` field.
```sh
make test-round-robin
```
As expected result you should see podinfo message changing

```text
{
  "hostname": "frontend-podinfo-856bb46677-8p45m",
  ...
  "message": "us",
  ...
}
```
```text
{
  "hostname": "frontend-podinfo-856bb46677-8p45m",
  ...
  "message": "eu",
  ...
}
```

### Failover (Ingress)

Both clusters have [podinfo](https://github.com/stefanprodan/podinfo) installed on top where each
cluster has been tagged to serve a different region. In this demo we will hit podinfo by `wget -qO - failover.cloud.example.com` and depending
on whether podinfo is running inside the cluster it returns only **eu** or **us**.
See [Gslb manifest with failover strategy](https://github.com/k8gb-io/k8gb/tree/master/deploy/gslb/k8gb.io_v1beta1_gslb_cr_failover_ingress.yaml)

Switch GSLB to failover mode:
```sh
make init-failover
```
Now both clusters are running in failover mode and podinfo is running on both of them.
Run several times command below and watch `message` field.
```sh
make test-failover
```
You will see only **eu** podinfo is responsive:
```text
{
  "hostname": "frontend-podinfo-856bb46677-8p45m",
  ...
  "message": "eu",
  ...
}
```
Stop podinfo on **current (eu)** cluster:
```
make stop-test-app
```
Several times hit application again
```sh
make test-failover
```
As expected result you should see only podinfo from **second cluster (us)** is responding:
```text
{
  "hostname": "frontend-podinfo-856bb46677-v5nll",
  ...
  "message": "us",
  ...
}
```
It might happen that podinfo will be unavailable for a while due to
[DNS sync interval](https://github.com/k8gb-io/k8gb/pull/81) and default k8gb DNS TTL of 30 seconds
```text
wget: server returned error: HTTP/1.1 503 Service Temporarily Unavailable
```
Start podinfo again on **current (eu)** cluster:
```sh
make start-test-app
```
and hit several times hit podinfo:
```sh
make test-failover
```
After DNS sync interval is over **eu** will be back
```text
{
  "hostname": "frontend-podinfo-6945c9ddd7-xksrc",
  ...
  "message": "eu",
  ...
}
```
Optionally you can switch GSLB back to round-robin mode
```sh
make init-round-robin
```
