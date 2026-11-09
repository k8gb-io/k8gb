# AI Agent Instructions for k8gb

This file provides repository-specific guidance for AI coding agents working on `k8gb` (Kubernetes Global Balancer). Read [`CONTRIBUTING.md`](CONTRIBUTING.md) and [`AI_POLICY.md`](AI_POLICY.md) for complete project policies.

## Quick Reference

| Task | Command | Notes |
| :--- | :--- | :--- |
| Full pre-PR verification | `make check` | Runs license header check, linter, and unit tests. Always run before opening a PR. |
| Lint | `make lint` | Uses rules defined in `.golangci.yaml`. |
| Add license headers | `make license` | Adds Apache-2.0 headers (`.golic.yaml`) to new Go files. |
| Unit tests & codegen | `make test` | Runs unit tests and regenerates API deepcopy code and CRD YAML manifests. |
| Regenerate mocks | `make mocks` | Regenerates Go mocks under `controllers/mocks/` using `mockery`. |
| Build binary | `make build` | Builds the `k8gb` operator binary. |
| Local multi-cluster setup | `make deploy-full-local-setup` | Deploys two local `k3d` clusters with `k8gb`, CoreDNS, and sample apps. |
| End-to-end tests | `make terratest` | Runs the Go-based Terratest E2E suite (or use Kyverno Chainsaw in `chainsaw/`). |

## AI Policy & Commit Sign-Off Rules

All AI-assisted contributions are governed by [`AI_POLICY.md`](AI_POLICY.md). Agents and contributors must follow these rules:

- **Human Accountability:** The named human contributor is solely responsible for the correctness, security, licensing, quality, and reviewability of all submitted changes, and must be able to explain and defend the submission without relying on an AI tool.
- **DCO Sign-Off (`git commit -s`):** All commits must be signed off (`git commit -s`), which adds `Signed-off-by: Your Name <your.email@example.com>` using the committer's own Git identity. Never copy another person's sign-off, and never sign off as an AI.
- **No AI Attribution Trailers:** Never list an AI system or agent in `Signed-off-by:`, `Co-authored-by:`, or any similar human-contributor Git trailer.
- **No Autonomous Posting:** Never configure an AI agent to autonomously open, post, or reply to issues, pull requests, or review comments on a contributor's behalf.
- **Verify Against Authoritative Sources:** Do not add an API, dependency, behavior claim, citation, or configuration solely because an AI tool proposed it; verify against the codebase and upstream documentation.
- **Licensing & Third-Party Code:** All contributions must comply with the Apache-2.0 license. Never include pre-existing third-party copyrighted material unless permissions and attributions are verified.
- **Security Reporting:** Never open a public issue or pull request containing vulnerability details discovered during AI analysis; follow the confidential reporting process in [`SECURITY.md`](SECURITY.md).

## Architecture Overview

`k8gb` is a cloud-native Global Server Load Balancing (GSLB) Kubernetes operator built with `controller-runtime` that load-balances traffic across geographically dispersed Kubernetes clusters using DNS.

- **`Gslb` (`k8gb.io/v1beta1`):** The primary user-facing Custom Resource Definition (CRD). It references a target Kubernetes `Ingress`, `Service`, or Gateway API resource (`spec.resourceRef`) or embedded ingress spec, and defines the global load-balancing strategy (`roundRobin`, `failover`, `geoip`) and `serviceHealthPolicy`. Legacy `k8gb.absa.oss/v1beta1` resources are accepted and automatically migrated at runtime to `k8gb.io/v1beta1`.
- **`ZoneDelegation` (`k8gb.io/v1beta1`):** A cluster-scoped CRD (`loadBalancedZone`, `parentZone`, `dnsZoneNegTTL`, `doFinalize`) that dynamically activates DNS zones only on clusters where applications are deployed (preventing premature `NXDOMAIN` responses in shared or multi-tenant clusters). When reconciled, `k8gb` updates the `k8gb-zone-delegation` `ConfigMap` mounted into CoreDNS pods. Its `.status.dnsServers` acts as the single source of truth for participating cluster DNS server IPs. Static zones defined via `k8gb.dnsZones` in Helm values are converted into `ZoneDelegation` resources internally.
- **`DNSEndpoint` (`externaldns.k8s.io/v1alpha1`):** Managed by the `k8gb` controller for two distinct purposes:
  1. **Application `DNSEndpoint`s:** Created per `Gslb` resource in the application namespace, storing local healthy endpoint targets (`localtargets-*`) and combined global DNS records resolved across clusters. These are watched and served directly by the in-cluster CoreDNS pods.
  2. **Delegation `DNSEndpoint`s:** Created in the `k8gb` namespace when the ExternalDNS provider is enabled so ExternalDNS can configure NS and glue records in the parent edge DNS zone (`parentZone`).
- **CoreDNS CRD Plugin ([`k8gb-io/coredns-crd-plugin`](https://github.com/k8gb-io/coredns-crd-plugin)):** An external CoreDNS plugin (`k8s_crd`) maintained in a separate repository. It runs inside the `k8gb` CoreDNS instance, watches `DNSEndpoint` resources via the Kubernetes API, and answers authoritative DNS queries (including `edns0` client-subnet GeoIP matching and round-robin responses) for the delegated GSLB zones.

## Repository Directory Map & Related Repositories

### Directory Map

- `api/` — Kubernetes CRD Go type definitions (`api/v1beta1/`, `api/v1beta1io/`, and `externaldns` types) and generated deepcopy code.
- `controllers/` — Operator reconciliation loops (`Gslb` and `ZoneDelegation`), DNS providers (`controllers/providers/`), configuration/dependency resolver (`controllers/depresolver/`), logging, Prometheus metrics, OpenTelemetry tracing, and generated mocks (`controllers/mocks/`).
- `chart/k8gb/` — Helm chart templates, `values.yaml`, and auto-generated CRD manifests (`chart/k8gb/crd/`).
- `docs/` — Project documentation built with MkDocs (`mkdocs.yml`) and link rewriting (`docs/link_fixer.py`).
- `terratest/` — Go-based Terratest integration and end-to-end test suite.
- `chainsaw/` — Kyverno Chainsaw declarative end-to-end test scenarios.
- `deploy/` — Local `k3d`/`kind` cluster manifests, Helm value profiles (`stable.yaml`, `next.yaml`), sample `Gslb` resources, and test workloads.
- `dns-provider-test/` — Integration test suite for external DNS providers (Route53, Infoblox, Cloudflare, Azure, GCP, NS1, RFC2136).
- `adr/` — Architecture Decision Records.

### Related Repositories

- [`k8gb-io/coredns-crd-plugin`](https://github.com/k8gb-io/coredns-crd-plugin) — CoreDNS `k8s_crd` plugin that serves DNS records directly from `DNSEndpoint` Custom Resources.
- [`k8gb-io/external-dns`](https://github.com/k8gb-io/external-dns) — ExternalDNS integration images and configurations used with `k8gb`.

## Coding Conventions

- **Linting:** Lint rules live in `.golangci.yaml`; run `make lint`.
- **Pre-PR Checks:** Run `make check` (which verifies license headers, lint rules, and unit tests) before opening a PR.
- **License Headers:** Every new Go file requires the Apache-2.0 license header; run `make license` to add it automatically.
- **Logging & Metrics:** Use the project's `zerolog` wrapper (`logging.Logger()`) with contextual fields. Do not discard errors with `_` outside tests; log reconciliation errors and increment error metrics via `metrics.Metrics()`.
- **Commit & PR Message Format:** Follow the format in [`CONTRIBUTING.md`](CONTRIBUTING.md) ("Commit and Pull Request message"):
  - Subject line (`<what was changed>`) must be 70 characters or fewer.
  - Second line must be blank.
  - Body (`<why this change was made>`) must wrap at 80 characters, followed by an optional footer (e.g., `Fixes #123`).

## Boundaries: What Agents Must Not Touch Without Permission

Never modify or hand-edit the following files without explicit human permission:

- `.github/workflows/*.yaml` — CI/CD and release pipelines (consult the maintainer team first).
- `chart/k8gb/crd/*.yaml` — Auto-generated CRD YAML manifests (never hand-edit).
- `api/*/zz_generated.deepcopy.go` — Auto-generated Kubernetes deepcopy functions (never hand-edit).
- `controllers/mocks/` — Auto-generated test mocks (regenerate via `make mocks`, never hand-edit).
- `chart/k8gb/Chart.yaml` — Helm chart version and metadata (modified only during the release workflow described in `CONTRIBUTING.md`).
- `cosign.pub` and `SECURITY-INSIGHTS.yml` — Cryptographic verification keys and security metadata.

## k8gb-Specific Pitfalls

- **`make test` regenerates API code and CRD YAMLs:** Running `make test` automatically runs `controller-gen` to regenerate deepcopy code and the CRD YAML files in `chart/k8gb/crd/`. Whenever you change API types under `api/`, run `make test`, commit the regenerated files, and **never** hand-edit `chart/k8gb/crd/*.yaml`.
- **Regenerating mocks after interface changes:** If you modify an interface that has a generated mock, run `make mocks` to update `controllers/mocks/` rather than editing mock files manually.
- **End-to-end testing requires `terratest` or `chainsaw`:** Standard `go test ./...` runs unit tests (`make test`). Multi-cluster end-to-end flows are tested using `terratest/` (`make terratest`) or Kyverno `chainsaw/` (`chainsaw/`) against local `k3d` clusters (`make deploy-full-local-setup`).
- **Adding new documentation pages (`mkdocs.yml`):** When adding a new documentation page under `docs/`, you must add it to **both** the `nav` section and the `llmstxt` plugin's `sections` in `mkdocs.yml`. Follow the link conventions in `CONTRIBUTING.md` (root files use `docs/<page>.md`; files inside `docs/` use relative links) so `mkdocs build --strict` passes.
