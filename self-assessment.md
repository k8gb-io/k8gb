# k8gb Self-assessment

Security reviewers: Yury Tsarev, Jirka Kremser

k8gb is a Global Service Load Balancing solution with a focus on having cloud native qualities and working natively in a Kubernetes context.

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
  * [Deployment pipeline](#deployment-pipeline)
  * [Communication channels](#communication-channels)
* [Security issue resolution](#security-issue-resolution)
  * [Responsible disclosure practice](#responsible-disclosure-practice)
  * [Incident response](#incident-response)
* [Appendix](#appendix)

## Metadata

### Security links
|   |  |
| - | - |
| Software | https://github.com/k8gb-io/k8gb  |
| Security Provider | No |
| Languages | Golang |
| SBOM | https://github.com/k8gb-io/k8gb/releases/download/v0.11.5/k8gb_0.11.5_linux_amd64.tar.gz.sbom.json |
| Security Insights | https://github.com/k8gb-io/k8gb/blob/master/SECURITY-INSIGHTS.yml |
| Security File | https://github.com/k8gb-io/k8gb/blob/master/SECURITY.md |
| Cosign pub-key | https://github.com/k8gb-io/k8gb/blob/master/cosign.pub |
| | |

### Intended use

To increase the software supply chain security, we encourage our users to consume k8gb container images with Kyverno's admission webhook 
([/policy](https://kyverno.io/docs/writing-policies/verify-images/sigstore/#verifying-image-signatures)) that will ensure that
images are signed and nobody had tempered with them. Our public key that can be used to verify this is in the root or our repository.

## Overview

k8gb is implemented using the Kubernetes operator pattern with a single CRD to enable Global Load Balancing. k8gb provides independent GSLB capability to any Ingress or Service without a dedicated management cluster, instead relying on timeproof DNS. k8gb has no single point of failure, and uses Kubernetes native application health checks (such as liveness and readiness probes) to update DNS to aid in load balancing decisions.

### Background

Global enterprises moving to the cloud need a global load balancer to make decisions based on Kubernetes services. Proprietary software and vendors exist, but these are expensive, complex, and not cloud native. Some require dedicated hardware and run outside Kubernetes.

k8gb is a vendor-neutral, CNCF Sandbox project. It is the only cloud native Kubernetes Global Load Balancer. k8gb does not require any special software or hardware - it relies only on OSS / CNCF projects, and fits with existing Kubernetes workflows like: GitOps, Kustomize, and Helm.

### Actors

@ytsarev please check this section.

The individual parts of k8gb that interact to provide the desired functionality.

- CoreDNS
  - Role: Embedded [custom CoreDNS](https://github.com/k8gb-io/coredns-crd-plugin) to serve DNS requests.
  - Isolation: runs as its own Pod
- ExternalDNS
  - Role: Integrated ExternalDNS to automate zone delegation configuration.  
  - Isolation: runs as its own Pod
- k8gb Controller
  - Role: Coordinates logic according to the GSLB strategy.
  - Isolation: runs as its own Pod

### Actions

@ytsarev please check this section.

The steps that k8gb performs in order to provide the desired functionality.

- CoreDNS
  - Functionality:
    - Serves external DNS requests with dynamically constructed DNS responses.
- ExternalDNS
  - Functionality:
    - Reads data from DNSEndpoint CR to update an external DNS provider (e.g., Route53).
- k8gb Controller
  - Functionality:
    - Watches all namespaces for `Gslb` custom resources.
    - Creates DNSEndpoint CR populated with information from `Gslb` Ingress status (application FQDN and active IP addresses, which are used for dynamic A record composition).
    - Creates DNSEndpoint to configure DNS zone delegation in an external DNS provider.

### Goals

@ytsarev please check this section.

The intended goals of k8gb, including the security guarantees it provides.

#### Secure and Verified Builds

Ensure all K8GB releases are signed and verified to guarantee authenticity and integrity, protecting users from tampered or malicious builds.

#### Minimal Attack Surface

Expose only the necessary ports for GSLB operations, such as DNS (53/tcp and 53/udp), to reduce the attack surface and enhance security.

#### Secure Deployment Practices

Provide secure default configurations and documentation to help users deploy K8GB in a way that aligns with Kubernetes security best practices.

These goals aim to make K8GB a reliable and secure solution for global load balancing while minimizing risks and ensuring trust in the project's artifacts.

See also [Intended use](#intended-use).

### Non-goals

@ytsarev please check this section.

Non-goals that a reasonable reader of k8gb’s documentation could believe may be in scope.

## Self-assessment use

This self-assessment is created by the k8gb team to perform an internal analysis of the project's security. It is not intended to provide a security audit of k8gb, or function as an independent assessment or attestation of k8gb's security health.

This document serves to provide k8gb users with an initial understanding of k8gb's security, where to find existing security documentation, k8gb plans for security, and general overview of k8gb security practices, both for development of k8gb as well as security of k8gb.

This document provides the CNCF TAG-Security with an initial understanding of k8gb to assist in a joint-assessment, necessary for projects under incubation. Taken together, this document and the joint-assessment serve as a cornerstone for if and when k8gb seeks graduation and is preparing for a security audit.

## Security functions and features

@ytsarev please check this section.

| Component | Applicability | Description of Importance |
| --------- | ------------- | ------------------------- |
| Component Name 1 | `Critical` or `Security Relevant` | why this feature is an important part of k8gb’s design and why it should be part of the threat model |
| Component Name N | `Critical` or `Security Relevant` | why this feature is an important part of k8gb’s design and why it should be part of the threat model |

## Project compliance

@ytsarev please check this section.

List of what standards k8gb is compliant with, and how that compliance has been validated, or Future State

### Future state

@ytsarev please check this section.

If k8gb is not compliant with any standards, note that here. Why is k8gb not compliant with any standards, and why that is the case. Will it need to be compliant in the future?

## Secure development practices

@ytsarev please check this section.

k8gb strives to implement the highest standard of secure development best practices, as noted below.

### Deployment pipeline

In order to secure the SDLC from development to deployment, the following measures are in place. Please consult the roadmap for information about how this list is growing.

- Branch protection on the default (`main`) branch:
  - Require signed commits
  - Require a pull request before merging
    - Require approvals: 1
    - Dismiss stale pull request approvals when new commits are pushed
    - Require review from Code Owners
    - Require approval of the most recent reviewable push
    - Require conversation resolution before merging
  - Require status checks to pass before merging
    - Require branches to be up to date before merging

### Communication channels

|   |  |
| - | - |
| Contributing | https://github.com/k8gb-io/k8gb/blob/master/CONTRIBUTING.md |
| Mailing list | cncf-k8gb-maintainers@lists.cncf.io |
| Slack | https://cloud-native.slack.com/archives/C021P656HGB |
| LinkedIn | https://www.linkedin.com/company/k8gb/ |
| Twitter / X | https://x.com/k8gb_io |
| Community meetings | https://zoom-lfx.platform.linuxfoundation.org/meetings/k8gb?view=week |
| | |

## Security issue resolution

The k8gb security policy is maintained in the https://github.com/k8gb-io/k8gb/blob/master/SECURITY.md file.

### Responsible disclosure practice

The k8gb project accepts vulnerability reports as outlined in the security policy https://github.com/k8gb-io/k8gb/blob/master/SECURITY.md#reporting-a-vulnerability.

### Incident response

The k8gb incident response process is outline in the security policy https://github.com/k8gb-io/k8gb/blob/master/SECURITY.md#review-process.

## Appendix
- Known Issues Over Time
  - Known issues are currently tracked in the project roadmap.
- OpenSSF Best Practices
  - [![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/k8gb-io/k8gb/badge)](https://securityscorecards.dev/viewer/?uri=github.com/k8gb-io/k8gb)
- Case Studies
- Related Projects / Vendors
