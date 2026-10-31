# Security Policy

We're always looking for ways to improve our security posture and taking security vulnerabilities seriously.

## Reporting a Vulnerability

CRA stewardship: This project is supported under the Linux Foundation CRA stewardship framework, as described at https://www.linuxfoundation.org/security. Security vulnerabilities should be reported through the mechanisms described below, which we will coordinate with our CRA steward. For actively exploited vulnerabilities and severe incidents that may require CRA escalation, please use the project's emergency security reporting mechanisms as appropriate.

To report a vulnerability, send an email to the official k8gb maintainers list [cncf-k8gb-maintainers@lists.cncf.io](mailto:cncf-k8gb-maintainers@lists.cncf.io), or report it privately on GitHub with [Report a vulnerability](https://github.com/k8gb-io/k8gb/security/advisories/new) on the repository's Security tab.

If you have evidence that the vulnerability is being actively exploited, say so at the top of your report. Maintainers prioritize those reports and escalate them to the Linux Foundation CRA steward as described in [CONTRIBUTING.md](CONTRIBUTING.md#escalating-to-the-cra-steward).

Please make sure to provide the following information:
- vulnerability description
- version of the k8gb package
- version of the k8s cluster
- steps to reproduce the vulnerability

## Review Process

The standard response time from k8gb project maintainers for a vulnerability report is less than 14 days.

Once k8gb project maintainers have confirmed the relevance of the report, a [draft security advisory](https://github.com/k8gb-io/k8gb/security/advisories) will be created on Github. Please provide Github username to k8gb project maintainers if you wish to be invited to participate in advisory discussions. You can also request to be updated about advisory discussions directly via email.

If the vulnerability is accepted, you'll be informed about the determined timeline for developing a patch, public disclosure, and patch release. You are expected to participate in the discussion of the timeline and respect the agreed dates for public disclosure.

## Published Advisories

Once a vulnerability is publicly disclosed, its advisory is published on the repository's [security advisories page](https://github.com/k8gb-io/k8gb/security/advisories).

## OSPS Baseline review evidence

The [published Level 2 results](https://grc.store/targets/k8gb-io/k8gb) include
controls that need human review. The following evidence accompanies the automated
results; it does not constitute an independent security audit or a badge award.

### Privileged release workflows (OSPS-BR-01.03)

The privileged-context finding in the 2026-10-31 scan concerns
[helm_sign_oci.yaml](https://github.com/k8gb-io/k8gb/blob/master/.github/workflows/helm_sign_oci.yaml). Its automatic trigger
follows a successful `Helm Publish` workflow; it also permits manual dispatch.
[Helm Publish](https://github.com/k8gb-io/k8gb/blob/master/.github/workflows/helm_publish.yaml) runs on published releases or
manual dispatch, with no pull-request trigger. The signing workflow uses a GitHub
hosted runner, checks out `inputs.tag || github.ref`, and does not select a pull
request head or the triggering workflow's head SHA. Manual tag selection and
release publishing are trusted maintainer actions.

The signing job has `contents: read`, `packages: write`, and `id-token: write`.
It fetches chart archives from this repository's `gh-pages` branch, pushes them to
GHCR, and signs them with GitHub OIDC. It does not restore an Actions cache or
download artifacts from a pull-request workflow. Writers to `gh-pages`, release
publishing credentials, and the pinned third-party actions are part of this trust
boundary. The job uses Harden Runner with an outbound endpoint allowlist.

Changes to the workflow triggers, checkout references, chart sources, caching, or
runner type require this evidence to be revisited. Release signing improvements
are tracked separately in [#2568](https://github.com/k8gb-io/k8gb/pull/2568).

### External interfaces (OSPS-SA-02.01)

| Interface | Documentation and contract |
| --- | --- |
| Kubernetes configuration and watched resources | [API types](https://github.com/k8gb-io/k8gb/tree/master/api), [CRD schemas](https://github.com/k8gb-io/k8gb/tree/master/chart/k8gb/crd), [resource references](https://github.com/k8gb-io/k8gb/blob/master/docs/resource_ref.md), and [dynamic zones](https://github.com/k8gb-io/k8gb/blob/master/docs/dynamic_zones.md) |
| DNS clients and inter-cluster health discovery | [Component data flows](https://github.com/k8gb-io/k8gb/blob/master/docs/components.md), [DNS exposure](https://github.com/k8gb-io/k8gb/blob/master/docs/exposing_dns.md), [strategies](https://github.com/k8gb-io/k8gb/blob/master/docs/strategy.md), and [trust boundaries](https://github.com/k8gb-io/k8gb/blob/master/self-assessment.md#actors) |
| Edge DNS providers and credentials | [Provider guides](https://github.com/k8gb-io/k8gb/blob/master/docs/tutorials.md), including [Route53](https://github.com/k8gb-io/k8gb/blob/master/docs/deploy_route53.md), [Google Cloud DNS](https://github.com/k8gb-io/k8gb/blob/master/docs/provider_gcp.md), and [Infoblox](https://github.com/k8gb-io/k8gb/blob/master/docs/deploy_infoblox.md) |
| Metrics, health endpoints, and optional tracing | [Metrics](https://github.com/k8gb-io/k8gb/blob/master/docs/metrics.md), [tracing](https://github.com/k8gb-io/k8gb/blob/master/docs/traces.md), and [ports and access assumptions](https://github.com/k8gb-io/k8gb/blob/master/self-assessment.md#minimal-attack-surface) |
| Installation and runtime configuration | [Helm chart documentation](https://github.com/k8gb-io/k8gb/blob/master/chart/k8gb/README.md) and [values schema](https://github.com/k8gb-io/k8gb/blob/master/chart/k8gb/values.schema.json) |

The [architecture map](https://github.com/k8gb-io/k8gb/blob/master/ARCHITECTURE.md) connects these interfaces to the code.
k8gb configures DNS and does not proxy application traffic; the referenced
Ingress, Service, Gateway API route, or Istio resource handles application traffic.

### Security self-assessment (OSPS-SA-03.01)

The [self-assessment](https://github.com/k8gb-io/k8gb/blob/master/self-assessment.md), reviewed in
[#2574](https://github.com/k8gb-io/k8gb/pull/2574), covers intended use, components,
actors and actions, trust boundaries, security goals and non-goals, development
practices, and vulnerability response. It records deployment and RBAC hardening
gaps, tracked in [#2587](https://github.com/k8gb-io/k8gb/issues/2587). The
[Security Insights assessment entry](https://github.com/k8gb-io/k8gb/blob/master/SECURITY-INSIGHTS.yml) links this evidence
and records its review date.

### Vulnerability publication (OSPS-VM-04.01)

The public disclosure destination is the
[GitHub security advisories page](https://github.com/k8gb-io/k8gb/security/advisories),
as specified in [Published Advisories](#published-advisories). The 2026-10-31 scan
found no published advisories, so it could not infer this practice from advisory
history. Private reporting and coordinated disclosure follow the process above.
