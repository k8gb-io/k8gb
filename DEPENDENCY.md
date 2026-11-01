# Environment Dependencies Policy

## Purpose

This policy describes how K8gb maintainers consume third-party packages.

## Scope

This policy applies to all K8gb maintainers and all third-party packages used in the K8gb project.

## Policy

K8gb maintainers must follow these guidelines when consuming third-party packages:

- Only use third-party packages that are necessary for the functionality of K8gb.
- Use the latest version of all third-party packages whenever possible.
- Avoid using third-party packages that are known to have security vulnerabilities.
- Pin all third-party packages to specific versions in the K8gb codebase.
- Use a dependency management tool, such as Go modules, to manage third-party dependencies.
- Dependencies must pass all automated tests before being merged into the K8gb codebase.

## Procedure

When adding a new third-party package to K8gb, maintainers must follow these steps:

1. Evaluate the need for the package. Is it necessary for the functionality of K8gb? 
2. Research the package. Is it well-maintained? Does it have a good reputation? 
3. Choose a version of the package. Use the latest version whenever possible. 
4. Pin the package to the specific version in the K8gb codebase. 
5. Update the K8gb documentation to reflect the new dependency.

## Obtaining Dependencies

- Go dependencies are declared in [`go.mod`](go.mod) and locked by checksum in [`go.sum`](go.sum). The Go toolchain downloads them through the default Go module proxy and checks them against the Go checksum database; the repository does not override `GOPROXY` or `GOSUMDB`.
- The Helm chart depends on the CoreDNS and ExternalDNS charts, pinned to exact versions in [`chart/k8gb/Chart.yaml`](chart/k8gb/Chart.yaml).
- Build, lint and test tools are pinned to exact versions in the [`Makefile`](Makefile), either as `*_VERSION` variables or in `go install <module>@<version>`.
- Most GitHub Actions in [`.github/workflows`](.github/workflows) are pinned to full commit SHAs.

## Tracking and Updating Dependencies

- [Renovate](https://docs.renovatebot.com/) tracks Go modules, Helm chart dependencies, GitHub Actions, and the tool and test-environment versions in the `Makefile` and `k3d/` configs. It opens pull requests for available updates, using the rules in [`renovate.json5`](renovate.json5). Most minor and patch updates are grouped into one pull request; dependencies still on a 0.x version, the k3s version and the Go version get their own.
- Pending and open updates are listed on the [Dependency Dashboard](https://github.com/k8gb-io/k8gb/issues/1048) issue.
- Some paths and dependencies are excluded from automatic updates. They are listed under `ignorePaths` and `ignoreDeps` in [`renovate.json5`](renovate.json5).
- Renovate pull requests go through the same CI and maintainer review as any other change. Renovate does not merge them on its own.
- The [FOSSA workflow](.github/workflows/fossa.yml) scans dependency licenses on pushes and on pull requests to `master`.
- Each release includes an SPDX SBOM for every release archive, generated with Syft, that records the dependencies shipped in that release.

## Archive/Deprecation

When a third-party package is discontinued, the K8gb maintainers must ensure to replace the package with a suitable alternative.

## Enforcement

This policy is enforced by the K8gb maintainers.
Maintainers are expected to review each other's code changes to ensure that they comply with this policy.

## Exceptions

Exceptions to this policy may be granted by the K8gb project lead on a case-by-case basis.

## Credits

This policy was adapted from the [Kubescape Community](https://github.com/kubescape/kubescape/blob/master/docs/environment-dependencies-policy.md) & [Project Capsule](https://github.com/projectcapsule/capsule/blob/main/DEPENDENCY.md)
