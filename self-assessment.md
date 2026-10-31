# k8gb Self-assessment

Security reviewers: Yury Tsarev, Jirka Kremser

Last reviewed: 2026-10-31. This revision follows the [OpenSSF security self-assessment template](https://security-assessments.openssf.org/self-assessment/template/self-assessment/) and corrects details that have changed since the 2025 version. Remaining hardening opportunities and documentation gaps are described in the relevant sections.

This document provides a self-assessment of the k8gb project following the guidelines outlined by the CNCF TAG Security and Compliance group and the OpenSSF. The purpose is to evaluate k8gb’s current security posture and alignment with best practices, ensuring that it is suitable for adoption at a CNCF incubation level.

## Table of Contents

* [Metadata](#metadata)
  * [Security links](#security-links)
  * [Intended use](#intended-use)
* [Overview](#overview)
  * [Actors](#actors)
  * [Actions](#actions)
  * [Background](#background)
  * [Goals](#goals)
  * [Non-goals](#non-goals)
* [Self-assessment use](#self-assessment-use)
* [Security functions and features](#security-functions-and-features)
* [Project compliance](#project-compliance)
  * [Future state](#future-state)
* [Secure development practices](#secure-development-practices)
  * [Development pipeline](#development-pipeline)
  * [Communication channels](#communication-channels)
  * [Ecosystem](#ecosystem)
* [Security issue resolution](#security-issue-resolution)
  * [Responsible disclosure practice](#responsible-disclosure-practice)
  * [Vulnerability response process](#vulnerability-response-process)
  * [Incident response](#incident-response)
* [Appendix](#appendix)

## Metadata

### Security links
|   |  |
| - | - |
| Assessment Stage | Complete. |
| Software | https://github.com/k8gb-io/k8gb  |
| Security Provider | No, k8gb is not a security provider. It relies on Kubernetes' built-in security mechanisms (e.g., RBAC, network policies) and does not implement standalone security features. |
| Languages | Golang |
| SBOM | SPDX SBOMs are attached to every release, for example https://github.com/k8gb-io/k8gb/releases/download/v1.0.1/k8gb_1.0.1_linux_amd64.tar.gz.sbom.json |
| Security Insights | https://github.com/k8gb-io/k8gb/blob/master/SECURITY-INSIGHTS.yml |
| Security File | https://github.com/k8gb-io/k8gb/blob/master/SECURITY.md |
| Cosign pub-key | https://github.com/k8gb-io/k8gb/blob/master/cosign.pub |
| | |

### Intended use

To increase the software supply chain security, we encourage our users to consume k8gb container images with Kyverno's admission webhook
([/policy](https://kyverno.io/docs/policy-types/cluster-policy/verify-images/sigstore/)) that will ensure that
images are signed and nobody has tampered with them. The public key that can be used to verify this is in the root of our repository.

## Overview

k8gb is implemented using the Kubernetes operator pattern, with a `Gslb` CRD (Custom Resource Definition) to enable Global Load Balancing and a `ZoneDelegation` CRD for dynamic zones. k8gb provides independent GSLB (Global Server Load Balancing) capability to any Ingress or Service without a dedicated management cluster, instead relying on timeproof DNS (Domain Name System). k8gb has no single point of failure, and uses Kubernetes native application health checks (such as liveness and readiness probes) to update DNS records to aid in load balancing decisions.

### Background

Global enterprises moving to the cloud need a global load balancer to make intelligent routing decisions based on the health and availability of Kubernetes services across multiple clusters and geographic regions. Traditional proprietary software and vendor solutions exist, but these are often expensive, complex, and not cloud native. Many require dedicated hardware appliances and operate outside the Kubernetes ecosystem, creating operational overhead and vendor lock-in.

k8gb addresses these challenges as a vendor-neutral, CNCF (Cloud Native Computing Foundation) Incubating project. It is designed as a cloud-native Kubernetes Global Load Balancer that operates entirely within the Kubernetes ecosystem. Unlike traditional solutions, k8gb does not require any special software, dedicated hardware, or external management systems - it relies exclusively on open source software (OSS) and CNCF projects, integrating seamlessly with existing Kubernetes workflows such as GitOps, Kustomize, and Helm package management.

### Actors

The individual parts of k8gb that interact to provide the desired functionality.

- CoreDNS
  - Role: Embedded [custom CoreDNS](https://github.com/k8gb-io/coredns-crd-plugin) to serve DNS (Domain Name System) requests.
  - Isolation: runs as its own Pod
- ExternalDNS
  - Role: Integrated ExternalDNS to automate zone delegation configuration.  
  - Isolation: runs as its own Pod
- k8gb Controller
  - Role: Coordinates logic according to the GSLB (Global Server Load Balancing) strategy.
  - Isolation: runs as its own Pod, under a ServiceAccount bound to a ClusterRole that is defined in `chart/k8gb/templates/cluster_role.yaml`
- Peer k8gb clusters
  - Role: Each cluster in a GSLB setup runs its own k8gb, and the clusters learn each other's healthy endpoints over DNS.
  - Isolation: separate Kubernetes clusters; k8gb controllers discover peer endpoints by querying peer CoreDNS instances and edge DNS resolvers over DNS.

### Actions

The steps that k8gb performs in order to provide the desired functionality. See https://www.k8gb.io/images/k8gb-components.svg.

- CoreDNS
  - Functionality:
    - Serves external DNS requests with dynamically constructed DNS responses.
- ExternalDNS
  - Functionality:
    - Reads data from DNSEndpoint CR to update an external DNS provider (e.g., Route53).
- k8gb Controller
  - Functionality:
    - Watches all namespaces for `Gslb` custom resources, and resolves the networking resource each one references through `spec.resourceRef` (Ingress, Service, Gateway API routes or Istio VirtualService).
    - Reads local endpoint health from ready EndpointSlices.
    - Creates DNSEndpoint CR (Custom Resource) populated with information from `Gslb` Ingress status (application FQDN - Fully Qualified Domain Name and active IP addresses, which are used for dynamic A record composition).
    - Queries the CoreDNS instances of the peer clusters for their `localtargets` records to learn which endpoints are healthy there, and merges them into the records it serves according to the GSLB strategy.
    - Creates DNSEndpoint to configure DNS (Domain Name System) zone delegation in an external DNS provider.

### Goals

The intended goals of k8gb, including the security guarantees it provides.

#### Secure and Verified Builds

Ensure all k8gb releases are signed and verified to guarantee authenticity and integrity, protecting users from tampered or malicious builds.

#### Minimal Attack Surface

k8gb is designed with a minimal attack surface. The only port meant for traffic from outside the cluster is DNS (53/tcp and 53/udp) on the CoreDNS Service, which is `ClusterIP` by default and exposed by the operator as needed (see [Exposing DNS](https://www.k8gb.io/exposing_dns/)). k8gb has no management interface or administrative API, and it does not proxy application traffic.

Inside the cluster, the components also serve plain HTTP endpoints without authentication:
- the k8gb controller serves Prometheus metrics on port 8080
- CoreDNS serves Prometheus metrics on port 9153, and health and readiness checks on ports 8080 and 8181
- optional OpenTelemetry tracing, off by default, adds an OTLP collector

The chart does not ship a NetworkPolicy, so restricting access to these ports is up to the operator.

#### Secure Deployment Practices

k8gb implements comprehensive secure deployment practices through multiple layers:

**Secure Defaults**: The k8gb controller container runs as a non-root user (UID 1000) with a read-only root filesystem, privilege escalation disabled, and resource requests and limits set. Its image is based on `gcr.io/distroless/static:nonroot`. CoreDNS and ExternalDNS use their upstream Helm charts' security defaults together with overrides in the k8gb chart and operator-supplied values. For example, the k8gb chart overrides CoreDNS's added capabilities with an empty list.

**RBAC**: The controller uses a ClusterRole for cross-namespace reconciliation and cluster-scoped `ZoneDelegation` resources. It reads Services, EndpointSlices, ConfigMaps, Gateway API routes and Istio resources. The current rules grant wildcard resources and verbs for `k8gb.io` and `k8gb.absa.oss`, and wildcard verbs for Ingresses and DNSEndpoints. A separate Role permits ConfigMap creation in the installation namespace and get/update/patch/delete on `k8gb-zone-delegation`. Replacing wildcard permissions with explicit resources and verbs is a candidate for further hardening, with validation of legacy reconciliation, migration, status updates and finalizers. Reducing cluster-wide ConfigMap list/watch access also requires changing the controller's cache configuration. These improvements are tracked in [#2587](https://github.com/k8gb-io/k8gb/issues/2587).

**Documentation and Guidance**: The DNS provider guides explain how to give ExternalDNS or k8gb its provider credentials, for example IRSA for Route53, Workload Identity for Google Cloud DNS, and Kubernetes Secrets for most other providers. k8gb does not yet have a dedicated production hardening guide. A guide covering NetworkPolicy, Pod Security Standards, DNS-authoring permissions, credential handling and artifact verification is tracked in [#2587](https://github.com/k8gb-io/k8gb/issues/2587).

**Kubernetes-Native Security**: k8gb leverages Kubernetes built-in security mechanisms rather than implementing custom security layers, ensuring compatibility with existing security policies and reducing complexity.

These goals aim to make k8gb a reliable and secure solution for global load balancing while minimizing risks and ensuring trust in the project's artifacts.

See also [Intended use](#intended-use).

### Non-goals

The k8gb project is focused on providing global load balancing and failover capabilities for Kubernetes applications. However, there are certain security-related features and responsibilities that are intentionally out of scope. Below are the non-goals, along with explanations for why they are not within the project's scope:

#### Traffic Handling and TLS Termination:

k8gb does not pass application traffic through itself or handle TLS/HTTP connections directly. These responsibilities are offloaded to referenced networking resources, such as Kubernetes Ingress controllers or service meshes.

**Reason**: k8gb is designed to manage DNS-based global load balancing and failover, not to act as a proxy or gateway for application traffic.

#### Security of Referenced Networking Resources:

k8gb does not enforce or manage security configurations (e.g., TLS settings) for the networking resources it references, such as Ingress controllers.

**Reason**: Security configurations for these resources are the responsibility of the user or the specific networking solution being used.

## Self-assessment use

This self-assessment is created by the k8gb team to perform an internal analysis of the project's security. It is not intended to provide a security audit of k8gb, or function as an independent assessment or attestation of k8gb's security health.

This document serves to provide k8gb users with an initial understanding of k8gb's security, where to find existing security documentation, k8gb plans for security, and general overview of k8gb security practices, both for development of k8gb as well as security of k8gb.

This document provides the CNCF TAG Security and Compliance with an initial understanding of k8gb to assist in a joint-assessment, necessary for projects under incubation. Taken together, this document and the joint-assessment serve as a cornerstone for if and when k8gb seeks graduation and is preparing for a security audit.

## Security functions and features

| Component | Applicability | Description of Importance |
| --------- | ------------- | ------------------------- |
| DNS-Based Traffic Management| `Critical`| k8gb uses DNS for global load balancing and failover, ensuring that traffic is routed to healthy clusters without passing through k8gb itself. This design minimizes the attack surface compared to traditional proxy-based load balancers and reduces the risk of traffic interception or manipulation versus solutions that handle application traffic directly. |
| Minimal Port Exposure | `Critical` | Only DNS (53/tcp and 53/udp) is meant to be reachable from outside the cluster, and there is no management interface or administrative API. In-cluster metrics and health endpoints are plain HTTP, so operators should restrict them (see [Minimal Attack Surface](#minimal-attack-surface)). |
| Integration with Kubernetes RBAC | `Critical` | k8gb relies on Kubernetes Role-Based Access Control (RBAC) to enforce authorization. Creating or modifying `Gslb`, `ZoneDelegation` or k8gb-labelled `DNSEndpoint` resources can influence DNS answers. The deprecated Ingress annotation interface also allows Ingresses with `k8gb.io/strategy` to create or update Gslbs through the controller. Changes to networking resources referenced by an existing Gslb can affect its DNS records too. Operators should account for these paths when granting write permissions; restricting direct CRD writes alone does not cover all DNS-authoring paths. |
| Cross-cluster health over DNS | `Critical` | Each k8gb cluster learns which endpoints are healthy in its peers by querying their CoreDNS over plain DNS. The DNS path between clusters and the edge resolvers is therefore part of the trusted environment. |
| Kubernetes Secrets for Sensitive Data | `Security Relevant` | DNS provider credentials are kept in Kubernetes Secrets or supplied through workload identity, and access is controlled by Kubernetes RBAC. Encryption at rest depends on how the cluster is configured. k8gb does not configure it. |
| Release signing and provenance | `Security Relevant` | Release binaries, archives, checksums, SBOMs and container images are signed with cosign using the project key in [`cosign.pub`](https://github.com/k8gb-io/k8gb/blob/master/cosign.pub). Helm charts in the OCI registry are signed keyless, and SLSA provenance is generated for binaries and images. See [Signed releases](https://github.com/k8gb-io/k8gb/blob/master/CONTRIBUTING.md#signed-releases). |
| Secure Default Configurations | `Security Relevant` | k8gb provides secure default configurations to help users deploy the project in a way that aligns with Kubernetes security best practices, reducing the risk of misconfiguration compared to complex traditional load balancing solutions that often require extensive manual security hardening. |

## Project compliance

List of what standards k8gb is compliant with, and how that compliance has been validated, or Future State

### Future state

k8gb maintains strong compliance with several industry standards and frameworks:

**Supply Chain Levels for Software Artifacts (SLSA)**:
- k8gb currently implements SLSA Level 3 compliance using the official `slsa-framework/slsa-github-generator`
- All release artifacts include signed provenance attestations with cryptographic verification
- Build process isolation and non-falsifiable provenance are implemented
- Both container images and release binaries have complete SLSA provenance chains

**OpenSSF Scorecard**:
- The Scorecard workflow runs on pushes to `master` and weekly, and publishes its results to [scorecard.dev](https://scorecard.dev/viewer/?uri=github.com/k8gb-io/k8gb)

**OpenSSF Best Practices**:
- k8gb holds the [OpenSSF Best Practices](https://www.bestpractices.dev/projects/4866) passing badge

**OpenSSF OSPS Baseline**:
- k8gb is working toward OSPS Baseline Level 2 as part of the Fall 2026 Security Slam, tracked in [#2563](https://github.com/k8gb-io/k8gb/issues/2563)

**Container Security Standards**:
- Container images are signed with cosign using the project key in `cosign.pub`. Helm charts in the OCI registry are signed keyless through GitHub OIDC.
- Software Bill of Materials (SBOM) generation for all releases
- Multi-architecture container builds with an SBOM attestation and SLSA provenance

**Current Compliance Status**:
k8gb does not hold formal certifications. The items above are self-reported, and each one links to results that can be checked independently.

## Secure development practices

k8gb is committed to maintaining a secure software development lifecycle (SDLC) by implementing robust practices and automation. Below are the key measures in place to ensure the security and integrity of the project.

### Development pipeline

In order to secure the SDLC from development to deployment, the following measures are in place. Please consult the roadmap for information about how this list is growing.

- Branch protection on the default (`master`) branch, verified on 2026-10-31:
  - Require a pull request with at least one approval
  - Dismiss stale approvals when new commits are pushed
  - Enforce the configured protection on administrators
  - Require linear history
  - Disable force pushes and branch deletion
  - Signed commits, Code Owner reviews, approval of the most recent reviewable push, and conversation resolution are not required
  - The [active master ruleset](https://github.com/k8gb-io/k8gb/rules/26274974) requires the DCO app's check and an up-to-date branch, with no bypass actors
  - Required CI checks and their activation procedure are documented in [Required pull-request checks](.github/required-checks.md); the live ruleset records which checks are enforced
  - The classic branch protection rule remains in place alongside the ruleset during rollout

- Contributions:
  - Every commit needs a Developer Certificate of Origin (DCO) sign-off, enforced by the DCO GitHub app
  - Pull requests need an approval from a maintainer listed in `CODEOWNERS` (see [CONTRIBUTING.md](https://github.com/k8gb-io/k8gb/blob/master/CONTRIBUTING.md#contribution-flow))

- CI/CD Pipeline:
  - Use GitHub Actions for continuous integration and deployment (CI/CD)
  - Selected CI/CD jobs run `step-security/harden-runner` with egress blocking. Coverage is not universal: the `go-inspect` job in `build.yml` is explicitly exempt, and `scorecard.yml` does not configure Harden Runner. Most third-party actions are pinned to commit SHAs; the SLSA generator workflows are referenced by version tags.
  - Include comprehensive testing and quality assurance to catch issues early:
    - **Code Quality**: golangci-lint pipeline for code quality and potential security issues
    - **Quality Monitoring**: Go Report Card integration (https://goreportcard.com/report/github.com/k8gb-io/k8gb)
    - **Functional Testing**: Terratest and Chainsaw end-to-end testing pipelines for functionality validation
  - **Security-Focused Scanning and Analysis**:
    - **Static Security Analysis**: CodeQL pipeline for vulnerability detection and security issue identification
    - **Security Posture Assessment**: OpenSSF Scorecard pipeline for continuous security best practices monitoring

- Release Process:
  - **Automated Release Pipeline**: A version bump pull request triggers tagging and the release workflow, and a maintainer publishes the resulting draft release (see [Release process](https://github.com/k8gb-io/k8gb/blob/master/CONTRIBUTING.md#release-process))
  - **Cryptographic Signing**: Release binaries, archives, checksums, SBOMs and container images are signed with cosign using the project key in `cosign.pub`. Helm charts in the OCI registry are signed keyless.
  - **SLSA Level 3 Compliance**: Generate provenance attestations using official SLSA framework for supply chain security
  - **Software Bill of Materials (SBOM)**: Comprehensive SBOM generation for transparency and vulnerability tracking

- Software Composition Analysis
  - Integration of dependency management tool (Mend Renovate) to monitor and secure third-party dependencies.
  - FOSSA scan pipeline

### Communication channels

- Internal: maintainers coordinate through the private `#k8gb-maintainers` channel in CNCF Slack and the maintainers mailing list (cncf-k8gb-maintainers@lists.cncf.io). General project discussions also take place in the public `#k8gb` CNCF Slack channel and community meetings.
- Inbound: users reach the team through GitHub issues and discussions, the `#k8gb` Slack channel, and the community meetings. Security reports go through the private channels in [SECURITY.md](https://github.com/k8gb-io/k8gb/blob/master/SECURITY.md).
- Outbound: releases are announced through GitHub releases and their notes, the [k8gb.io blog](https://www.k8gb.io/blog/), LinkedIn and X.

|   |  |
| - | - |
| Contributing | https://github.com/k8gb-io/k8gb/blob/master/CONTRIBUTING.md |
| Mailing list | cncf-k8gb-maintainers@lists.cncf.io |
| Slack | https://cloud-native.slack.com/archives/C021P656HGB |
| LinkedIn | https://www.linkedin.com/company/k8gb/ |
| Twitter / X | https://x.com/k8gb_io |
| Community meetings | https://zoom-lfx.platform.linuxfoundation.org/meetings/k8gb?view=week |
| | |

### Ecosystem

k8gb works with the networking resources teams already run. A `Gslb` can point at a Kubernetes Ingress, a Service, a Gateway API route or an Istio VirtualService. k8gb serves DNS through CoreDNS and its own CRD plugin ([coredns-crd-plugin](https://github.com/k8gb-io/coredns-crd-plugin)), and it delegates zones through ExternalDNS or directly to Infoblox. The supported edge DNS providers are Route53, Azure DNS, Google Cloud DNS, Infoblox, NS1, Windows DNS, Cloudflare and RFC2136 servers such as BIND. Organizations using k8gb are listed in [ADOPTERS.md](https://github.com/k8gb-io/k8gb/blob/master/ADOPTERS.md).

## Security issue resolution

The k8gb security policy is maintained in the https://github.com/k8gb-io/k8gb/blob/master/SECURITY.md file.

### Responsible disclosure practice

The k8gb project accepts vulnerability reports as outlined in the security policy https://github.com/k8gb-io/k8gb/blob/master/SECURITY.md#reporting-a-vulnerability. Reports can be sent privately by email to the maintainers list or through GitHub private vulnerability reporting on the repository's Security tab. k8gb is supported under the Linux Foundation CRA stewardship framework, and reports of actively exploited vulnerabilities are escalated to the LF steward.

### Vulnerability response process

The k8gb maintainer group is responsible for triaging and responding to vulnerability reports, with a standard response time of less than 14 days. The repository administrators are @ytsarev, @k0da and @kuritka, as recorded in `SECURITY-INSIGHTS.yml`. Once a report is confirmed, maintainers open a draft GitHub security advisory, agree a timeline for the fix and disclosure with the reporter, and publish the advisory on the [security advisories page](https://github.com/k8gb-io/k8gb/security/advisories) once the fix is released.

### Incident response

The k8gb incident response process is outlined in [SECURITY.md](https://github.com/k8gb-io/k8gb/blob/master/SECURITY.md#review-process). Actively exploited vulnerabilities and security compromises of the project's own infrastructure are escalated to the Linux Foundation steward at steward@linuxfoundation.org, as described in [CONTRIBUTING.md](https://github.com/k8gb-io/k8gb/blob/master/CONTRIBUTING.md#escalating-to-the-cra-steward). Remediation does not wait for the escalation decision.

[Signed releases](https://github.com/k8gb-io/k8gb/blob/master/CONTRIBUTING.md#signed-releases) documents generating a replacement cosign key pair and updating `COSIGN_PRIVATE_KEY`, `COSIGN_PUBLIC_KEY`, `COSIGN_PASSWORD` and `cosign.pub`. A complete compromised-release or signing-key response runbook is not yet documented. It should cover suspending affected publishing, investigating affected credentials and artifacts, restoring a trusted build environment, rotating compromised credentials, publishing verified replacement artifacts, and notifying users which artifacts and keys must no longer be trusted. Suspect artifacts must be investigated or rebuilt from trusted inputs before signing them again. Keyless OCI chart signing requires separate consideration of the GitHub workflow and OIDC identity.

## Appendix
- Known Issues Over Time
  - Known issues are currently tracked in the project roadmap.
  - No security advisories have been published to date. Published advisories are listed on the [security advisories page](https://github.com/k8gb-io/k8gb/security/advisories), and other known issues and bugs are tracked in GitHub Issues.
- OpenSSF Best Practices
  - [![OpenSSF Best Practices](https://www.bestpractices.dev/projects/4866/badge)](https://www.bestpractices.dev/projects/4866)
  - [![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/k8gb-io/k8gb/badge)](https://scorecard.dev/viewer/?uri=github.com/k8gb-io/k8gb)
- Case Studies
  - Multicloud load balancing at Millennium bcp, described in the [CNCF case study](https://www.cncf.io/case-studies/millennium-bcp/).
  - Disaster Recovery for Financial Services
    A financial services company leverages k8gb for disaster recovery, ensuring that critical applications remain available even during data center failures. The k8gb failover mechanism automatically redirects traffic to backup clusters, maintaining business continuity. More in this KubeCon presentation: https://www.youtube.com/watch?v=U46hlF0Z3xs
- Related Projects / Vendors
  - Kubernetes Ingress Controllers:
    k8gb is often compared to Kubernetes Ingress controllers (e.g., NGINX Ingress, Traefik). While Ingress controllers handle traffic routing within a single cluster, k8gb focuses on global load balancing and failover across multiple clusters, making it complementary to Ingress solutions.
  - External DNS Providers:
    k8gb integrates with external DNS providers (e.g., Route 53, Cloudflare) to manage DNS-based traffic routing. Unlike standalone DNS solutions, k8gb provides Kubernetes-native automation for global load balancing and failover.
  - Service Meshes:
    Service meshes (e.g., Istio, Linkerd) provide advanced traffic management and security features within a cluster. k8gb, on the other hand, operates at the DNS level, enabling cross-cluster traffic management without requiring changes to application code.
