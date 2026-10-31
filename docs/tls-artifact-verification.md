# TLS artifact verification and upgrades

The [eMudhra incident](https://go.dev/blog/emudhra-incident) requires rebuilding Go
clients with a patched compiler or updating their trust stores. Updating `go.mod`
alone does not change an already published container. Track the remaining release
work in [#2592](https://github.com/k8gb-io/k8gb/issues/2592).

## What this repository checks

The operator uses Go 1.27.3 or later and a digest-pinned distroless base. Renovate
can update that digest through its Dockerfile manager. Both the amd64 and arm64
images must pass `hack/verify-image.sh`:

- The binary build metadata must report a released Go version >= 1.27.3 and the
  requested operating system and architecture.
- The final image's CA bundle must parse successfully and contain certificates.
- Neither `AAA Certificate Services` nor `emSign ECC Root CA - G3` may be present,
  identified by their public-key hashes, following the
  [distroless check](https://github.com/GoogleContainerTools/distroless/blob/29b6e28900d76ffcde3cdc69af252be9c76c8635/common/cacerts_test.go).

The check copies files from a stopped container and never executes the inspected
binary. It prints the Docker image identity, available registry digests, compiler
version, binary and bundle hashes, certificate count, and affected roots. It is
specific to this remediation; it is not a general vulnerability scanner or a
replacement for signature verification.

The optional embedded ExternalDNS deployment explicitly uses v0.24.0, independently
of its Helm dependency (chart 1.23.0 defaults to v0.23.0). The nested
`extdns.image.repository` and `extdns.image.tag` values let Renovate's `helm-values`
manager track the image even before another upstream chart is published.
Updating the imported `sigs.k8s.io/external-dns` Go module does not update this
separate deployment.

ExternalDNS v0.24.0 is built with Go 1.27.3 on both supported architectures. Its
CA bundle still contains the two affected root keys. The Go verifier blocks the
compromised intermediates; refreshing ExternalDNS's base trust store remains an
upstream follow-up. The `external-dns` check profile therefore requires patched
Go and reports the roots without requiring their removal. This exception does
not apply to the operator.

## Upgrading installations

Install the new k8gb release once it is published. Inspect any existing
`k8gb.imageTag` and `extdns.image.tag` overrides: explicit values, including values
retained by `helm upgrade --reuse-values`, can keep an old image in use. Render the
upgrade and check the running deployments' images after rollout. Release v1.0.1
predates the operator rebuild; changing only the ExternalDNS tag does not update
the operator binary.

For embedded ExternalDNS, the values are:

```yaml
extdns:
  enabled: true
  image:
    repository: registry.k8s.io/external-dns/external-dns
    tag: v0.24.0
```

For independently managed ExternalDNS, upgrade its deployment separately to
v0.24.0 or a newer verified image. Preserve the provider credentials, domain and
label filters, TXT ownership settings, and managed record types. A k8gb chart
upgrade cannot update a deployment managed by another release or controller.

The image-only upgrade retains the existing `externaldns.k8s.io/v1alpha1`
DNSEndpoint CRDs. ExternalDNS v0.24.0 still reads that API, including k8gb's A and
NS records; a CRD replacement is not required for the Go security fix. Its
[upstream release](https://github.com/kubernetes-sigs/external-dns/releases/tag/v0.24.0)
also introduces stricter CRD validation. If separately adopting those schemas,
audit existing records against the new DNS-name, record-type, target, TTL and
size constraints before applying them. Helm does not automatically upgrade CRDs
from a dependency's `crds/` directory.

## Before publishing a release

Chainsaw CI verifies the snapshot operator images and the ExternalDNS image
selected by the deployed chart on both architectures, then runs the DNS
reconciliation suite. The existing upgrade-testing workflow exercises upgrades
from the stable release.

The release workflow checks both operator architectures after GoReleaser pushes
the images, while the GitHub release is still a draft. A failed check must be
resolved before publishing that draft. To repeat the content checks manually
(requires Docker and the repository's Go toolchain):

```sh
operator=ghcr.io/k8gb-io/k8gb:<release-tag>
externaldns=registry.k8s.io/external-dns/external-dns:v0.24.0
for arch in amd64 arm64; do
  docker pull --platform "linux/$arch" "$operator"
  bash hack/verify-image.sh "$operator" "linux/$arch"
  docker pull --platform "linux/$arch" "$externaldns"
  bash hack/verify-image.sh "$externaldns" "linux/$arch" external-dns
done
```

Before closing #2592:

1. Require successful artifact, DNS reconciliation, and upgrade checks for the
   release candidate. Check other DNS providers used by the installation as well;
   an RFC2136 test does not exercise cloud-provider authentication.
2. Verify the signatures and provenance using the
   [release verification instructions](../CONTRIBUTING.md#signed-releases).
3. Publish the operator release and Helm chart; render/install the published
   chart and confirm that `registry.k8gb.io` can pull the new operator image.
4. Record the release version, immutable image index/platform digests, compiler
   versions, CA-check output, and verification results on the issue. Use digests
   when repeating checks so a subsequently moved tag cannot change the evidence.

Publishing, checking the published chart/mirror and signatures/provenance, and
attaching the final release evidence remain release-time work.
