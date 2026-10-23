# Agent Instructions for k8gb

## Project Overview
k8gb is a cloud-native Global Server Load Balancing (GSLB) solution for Kubernetes. It enables multi-cluster traffic management and failover using DNS.

## Architecture
k8gb consists of several key components:
- **Gslb Controller**: Manages GSLB custom resources (`Gslb`), orchestrating health checks, geo-IP resolution, and DNS configuration across clusters.
- **DNSEndpoint**: Custom Resource Definition (CRD) managed by external-dns that represents generated DNS records.
- **CoreDNS Plugin**: Custom CoreDNS plugin (`k8gb` plugin) that dynamically resolves DNS queries according to cluster health and GSLB strategy.
- **External-DNS Integration**: Synchronizes DNSEndpoint resources with edge DNS providers and cross-cluster CoreDNS instances.

## Build, Test, and Development Commands
Common commands from `Makefile`:
- `make test`: Run unit tests.
- `make lint`: Run linters to verify code format and quality.
- `make mocks`: Generate mocks for testing interfaces.
- `make deploy-full-local-setup`: Spin up local test environment using multi-cluster setups for local end-to-end testing.

## Coding Conventions
- Standard Go code layout and formatting rules (`gofmt`, `golangci-lint`).
- Maintain modular structure with clear separation between controller logic, DNS resolution, and provider integration.
- Include unit tests and integration tests for new features and bug fixes.

## Developer Certificate of Origin (DCO)
All commits must include a `Signed-off-by` line certifying compliance with the Developer Certificate of Origin (DCO).
Example:
`Signed-off-by: Abhinav Singh Yadav <abhinavsinghyadav17oct@gmail.com>`
