# Architecture

This document is a map of the k8gb code base for people who want to change it. For what k8gb does and why, start with the [introduction](docs/intro.md). For component diagrams, see [K8GB Internal Components](docs/components.md). For the security view of the same system, with its actors, actions and trust boundaries, see the [self-assessment](self-assessment.md).

## Bird's-eye view

k8gb is a Kubernetes operator that does global load balancing through DNS. Every cluster in a GSLB setup runs its own k8gb, there is no central control cluster, and the clusters share state only through DNS.

In each cluster:

1. The k8gb controller watches `Gslb` resources. Each `Gslb` points at a networking resource through `spec.resourceRef`: an Ingress, a LoadBalancer Service, a Gateway API route or an Istio VirtualService.
2. The controller checks whether the Services behind that resource have ready endpoints, and writes the addresses the resource is exposed on into a `DNSEndpoint` resource.
3. CoreDNS, running k8gb's [CRD plugin](https://github.com/k8gb-io/coredns-crd-plugin), answers queries for the load-balanced zone from those `DNSEndpoint` resources.
4. To learn what is healthy in the other clusters, the controller asks their CoreDNS for `localtargets` records, then merges the answers according to the `Gslb` strategy: round robin, weighted round robin, failover or GeoIP. See [strategies](docs/strategy.md).
5. The parent zone in the edge DNS provider delegates the load-balanced zone to the CoreDNS instances. The controller keeps those NS and glue records up to date through ExternalDNS, or through the Infoblox API when Infoblox is the edge DNS. With `ZoneDelegation` resources, a cluster can start and stop serving zones at runtime. See [dynamic zones](docs/dynamic_zones.md).

## Code map

| Path | What lives there |
|---|---|
| `main.go` | Starts the controller manager and registers the reconcilers below. |
| `api/v1beta1io/` | The `k8gb.io/v1beta1` API types: `Gslb` and `ZoneDelegation`. |
| `api/v1beta1/` | The legacy `k8gb.absa.oss/v1beta1` `Gslb` type, kept so existing resources can be migrated. |
| `controllers/gslb_controller_*.go` | The main `Gslb` reconciler. |
| `controllers/gslb_legacy_controller.go`, `controllers/gslb_migration_controller.go`, `controllers/conversion_k8gbio.go` | Handling of legacy `Gslb` resources and their migration to `k8gb.io`. |
| `controllers/bootstrap_controller_*.go` | Watches the CoreDNS Service and any Ingress labeled `k8gb.io/ip-source=true`, and refreshes zone delegation when their exposed IPs change. |
| `controllers/zone_delegation_*.go`, `controllers/zones/` | The `ZoneDelegation` reconciler and the zone service behind it. |
| `controllers/refresolver/` | Resolves `spec.resourceRef` for each supported kind and finds the IPs it is exposed on. |
| `controllers/resolver/` | Reads and validates the controller configuration, which the Helm chart passes in as environment variables. |
| `controllers/providers/k8gbendpoint/` | Builds the application `DNSEndpoint`, including the `localtargets` lookups in peer clusters and the strategy logic. |
| `controllers/providers/dns/` | Zone delegation in the edge DNS: the ExternalDNS provider and the Infoblox client. |
| `controllers/providers/metrics/` | Prometheus metrics. |
| `controllers/geotags/`, `controllers/ipresolver/`, `controllers/tracing/` | Cluster geo tags, discovery of the IPs CoreDNS is exposed on, and OpenTelemetry tracing. |
| `chart/k8gb/` | The Helm chart: CRDs, RBAC, the controller Deployment, and the CoreDNS and ExternalDNS subcharts. |
| `terratest/`, `chainsaw/` | End-to-end tests that run against local k3d clusters. The cluster configs are in `k3d/`. |
| `deploy/` | Manifests and Helm values used by the local setup and demos in the `Makefile`. |
| `docs/`, `mkdocs.yml`, `overrides/` | The [k8gb.io](https://www.k8gb.io) website. |
| `adr/` | Architecture decision records. |
| `grafana/` | Grafana dashboards for the controller metrics. |
| `olm/` | Templates and a script that generate the Operator Lifecycle Manager bundle. |
