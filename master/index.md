# K8GB - Kubernetes Global Balancer

[CNCF Incubating Project](https://landscape.cncf.io/?item=orchestration-management--coordination-service-discovery--k8gb) | [Roadmap](https://github.com/orgs/k8gb-io/projects/2/views/2) | [Join #k8gb on CNCF Slack](https://cloud-native.slack.com/archives/C021P656HGB) | [Join our Community Meetings](#community-meetings)

> **Container images:** the Helm chart default is `registry.k8gb.io/k8gb-io/k8gb` (see `chart/k8gb/values.yaml`). Signed/OCI artifacts are also published to `ghcr.io/k8gb-io/*` (see [CONTRIBUTING.md](https://www.k8gb.io/latest/CONTRIBUTING/index.md)). The Docker Hub `absaoss/k8gb` badge above reflects historical pulls.

A Global Service Load Balancing solution with a focus on having cloud native qualities and work natively in a Kubernetes context.

Just a single Gslb CRD to enable the Global Load Balancing:

```
apiVersion: k8gb.io/v1beta1
kind: Gslb
metadata:
  name: test-gslb-failover
  namespace: test-gslb
spec:
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    name: test-gslb-failover
  strategy:
    type: failover # Global load balancing strategy
    primaryGeoTag: eu-west-1 # Primary cluster geo tag
```

Legacy `k8gb.absa.oss/v1beta1` resources are still accepted and automatically migrated to `k8gb.io/v1beta1`. Legacy objects are labeled `k8gb.io/migrated-to-k8gb-io=true` and emit a warning event telling users to edit the `k8gb.io` object going forward. Legacy example: `docs/examples/legacy/gslb-legacy.yaml`.

[Global load balancing](https://cloud.redhat.com/blog/global-load-balancer-approaches), commonly referred to as GSLB (Global Server Load Balancing) solutions, has been typically the domain of proprietary network software and hardware vendors and installed and managed by siloed network teams.

k8gb is a completely open source, cloud native, global load balancing solution for Kubernetes.

k8gb focuses on load balancing traffic across geographically dispersed Kubernetes clusters using multiple load balancing [strategies](https://www.k8gb.io/latest/strategy/index.md) to meet requirements such as region failover for high availability.

Global load balancing for any Kubernetes Service can now be enabled and managed by any operations or development teams in the same Kubernetes native way as any other custom resource.

## Key Differentiators

- Load balancing is based on timeproof DNS protocol which is perfect for global scope and extremely reliable
- No dedicated management cluster and no single point of failure
- Kubernetes native application health checks utilizing status of Liveness and Readiness probes for load balancing decisions
- Configuration with a single Kubernetes CRD of Gslb kind

## Quick Start

Simply run

```
make deploy-full-local-setup
```

It will deploy two local [k3s](https://k3s.io/) clusters via [k3d](https://k3d.io/), [expose associated CoreDNS service for UDP DNS traffic](https://www.k8gb.io/latest/exposing_dns/index.md), and install k8gb with test applications and two sample Gslb resources on top.

This setup is adapted for local scenarios and works without external DNS provider dependency.

Consult with [local playground](https://www.k8gb.io/latest/local/index.md) documentation to learn all the details of experimenting with local setup.

To demonstrate k8gb as a global resilience layer for AI inference endpoints, run:

```
make deploy-full-local-setup FULL_LOCAL_SETUP_WITH_AI_DEMO=true
make ai-inference-demo
```

This deploys a lightweight Ollama model in both clusters, exposes an OpenAI-compatible endpoint, and shows failover through the same global hostname.

See the [AI inference resilience demo](https://www.k8gb.io/latest/ai-inference-demo/index.md) for local and real-cluster usage.

Optionally, you can run `make deploy-prometheus` and check the metrics on the test clusters (http://localhost:9090, http://localhost:9091).

## Motivation and Architecture

k8gb was born out of the need for an open source, cloud native GSLB solution at Absa Group in South Africa.

As part of the bank's wider container adoption running multiple, geographically dispersed Kubernetes clusters, the need for a global load balancer that was driven from the health of Kubernetes Services was required and for which there did not seem to be an existing solution.

Yes, there are proprietary network software and hardware vendors with GSLB solutions and products, however, these were costly, heavyweight in terms of complexity and adoption, and were not Kubernetes native in most cases, requiring dedicated hardware or software to be run outside of Kubernetes.

This was the problem we set out to solve with k8gb.

Born as a completely open source project and following the popular Kubernetes operator pattern, k8gb can be installed in a Kubernetes cluster and via a Gslb custom resource, can provide independent GSLB capability to any Ingress or Service in the cluster, without the need for handoffs and coordination between dedicated network teams.

k8gb commoditizes GSLB for Kubernetes, putting teams in complete control of exposing Services across geographically dispersed Kubernetes clusters across public and private clouds.

k8gb requires no specialized software or hardware, relying completely on other OSS/CNCF projects, has no single point of failure, and fits in with any existing Kubernetes deployment workflow (e.g. GitOps, Kustomize, Helm, etc.) or tools.

Please see the extended architecture documentation [here](https://www.k8gb.io/latest/intro/index.md)

Internal k8gb architecture and its components are described [here](https://www.k8gb.io/latest/components/index.md)

## Installation and Configuration Tutorials

For comprehensive installation and configuration guides, see the [Installation and Configuration Tutorials](https://www.k8gb.io/latest/tutorials/index.md) documentation.

## Adopters

A list of publicly known users of the K8GB project can be found in [ADOPTERS.md](https://www.k8gb.io/latest/ADOPTERS/index.md). We encourage all users of K8GB to add themselves to this list!

## Case Studies

- Featured end-user story: [Millennium bcp CNCF case study](https://www.cncf.io/case-studies/millennium-bcp/)

## Production Readiness

You can use k8gb in on-prem, cloud and hybrid environments.

k8gb is tested with the following environment options and integration scenarios.

| Type                      | Implementation                                                                                                  |
| ------------------------- | --------------------------------------------------------------------------------------------------------------- |
| Kubernetes Version        | >= `1.32` (CI matrix: `v1.32` / `v1.33` / `v1.34`)                                                              |
| Environment               | Any conformant Kubernetes cluster on-prem or in cloud                                                           |
| Ingress / Gateway         | Gateway API (recommended), Istio, AWS Load Balancer Controller; Ingress NGINX only as legacy [\*](#clarify)     |
| EdgeDNS                   | Infoblox, Route53, NS1, CloudFlare, AzureDNS, GCP Cloud DNS                                                     |
| Ingress Integration Modes | Gateway API (`HTTPRoute` / related), Kubernetes Ingress (embedded/referenced), Istio (VirtualService + Gateway) |
| Automated E2E Validation  | Multi-cluster automated E2E testing with Terratest and Chainsaw                                                 |

\* Every entry in this table reflects a setup that has been tested and verified with a working k8gb installation. * This matrix is intentionally conservative, not exhaustive. k8gb is Kubernetes-native and designed for standards-compliant clusters and ingress implementations; if your stack is not listed, validate it in staging before production rollout.

## Presentations and Talks

|                                         | **KubeCon NA 2025**                     |
| --------------------------------------- | --------------------------------------- |
| **Container Days Conference 2025**      | **Container Days Conference 2025**      |
| **KCD Czech&Slovak 2025**               | **Cloud Native Rejekts 2025**           |
| **KubeCon China 2025**                  | **KubeCon EU 2025**                     |
| **ChatLoopBackOff - Episode 42 (K8gb)** | **KubeCon NA 2024**                     |
| **Open Source Summit EU 2024**          | **KubeCon EU 2024**                     |
| **KubeCon NA 2023**                     | **KubeCon EU 2023**                     |
| **FOSDEM 2022**                         | **KCDBengaluru 2023**                   |
| **Crossplane Community Day**            | **KubeCon NA 2021**                     |
| **#29 DoK Community**                   | **AWS Containers from the Couch show**  |
| **OpenShift Commons Briefings**         | **Demo at Kubernetes SIG Multicluster** |

You can also find recordings from our community meetings at [k8gb youtube channel](https://www.youtube.com/channel/UCwvtktvdZu_pg-t-INvuW5g). The meeting runs every other Wednesday from 13:00 CET to 13:30 CET [calendar](https://zoom-lfx.platform.linuxfoundation.org/meetings/k8gb?view=month). Join the [Zoom Meeting](https://zoom-lfx.platform.linuxfoundation.org/meeting/92572060749?password=645f8346-1952-44fa-bd9b-45208260fc10).

## Online Articles and Publications

- https://oilbeater.com/en/2024/04/18/k8gb-best-cloudnative-gslb/
- https://medium.com/@kaanara/designing-a-high-availability-kubernetes-cluster-infrastructure-for-a-major-digital-travel-agency-7663465871fd
- https://www.redhat.com/en/blog/global-load-balancing-red-hat-openshift-k8gb
- https://blog.neilfren.ch/global-service-loadbalancers-for-your-homelab/
- https://andrewbaker.ninja/2021/01/22/external-k8gb-presentation-to-kubernetes-sig-multicluster/

## Books Covering k8gb

| **Kubernetes - An Enterprise Guide - Second Edition** | **Kubernetes – An Enterprise Guide - Third Edition** |
| ----------------------------------------------------- | ---------------------------------------------------- |
|                                                       |                                                      |

## Contributing

See [CONTRIBUTING](https://www.k8gb.io/latest/CONTRIBUTING/index.md)
