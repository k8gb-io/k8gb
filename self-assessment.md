# K8gb security self-assessment

2023-10-31

K8gb is a Global Service Load Balancing solution with a focus on having cloud native qualities and work natively in a Kubernetes context.

## Metadata

Quick reference information, later used for indexing.

|   |  |
| -- | -- |
| Software | https://github.com/k8gb-io/k8gb  |
| Security Provider | No |
| Languages | Golang |
| SBOM | https://github.com/k8gb-io/k8gb/releases/download/v0.11.5/k8gb_0.11.5_linux_amd64.tar.gz.sbom.json |
| | |

### Security links

| Doc | url |
| -- | -- |
| Security file | https://github.com/k8gb-io/k8gb/blob/master/SECURITY.md |
| Security insights | https://github.com/k8gb-io/k8gb/blob/master/SECURITY-INSIGHTS.yml |
| Cosign pub-key | https://github.com/k8gb-io/k8gb/blob/master/cosign.pub |

### Intended Use

To increase the software supply chain security, we encourage our users to consume k8gb container images with Kyverno's admission webhook 
([/policy](https://kyverno.io/docs/writing-policies/verify-images/sigstore/#verifying-image-signatures)) that will ensure that
images are signed and nobody had tempered with them. Our public key that can be used to verify this is in the root or our repository.
