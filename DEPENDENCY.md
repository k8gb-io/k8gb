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

## Archive/Deprecation

When a third-party package is discontinued, the K8gb maintainers must fensure to replace the package with a suitable alternative.

## Enforcement

This policy is enforced by the K8gb maintainers.
Maintainers are expected to review each other's code changes to ensure that they comply with this policy.

## Exceptions

Exceptions to this policy may be granted by the K8gb project lead on a case-by-case basis.

## Credits

This policy was adapted from the [Kubescape Community](https://github.com/kubescape/kubescape/blob/master/docs/environment-dependencies-policy.md) & [Project Capsule](https://github.com/projectcapsule/capsule/blob/main/DEPENDENCY.md)