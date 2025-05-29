# Automatic CoreDNS exposed IP discovery
k8gb automatically discovers the correct IP addresses for DNS delegation during bootstrap.

## How Address Discovery Works
When k8gb needs to create DNS delegation records (delegated DNSEndpoints sitting in k8gb namespace), it needs to know CoreDNS IP addresses.

**- If you are running CoreDNS service as a LoadBalancer service** (`.Values.coredns.serviceType`: `LoadBalancer`): k8gb automatically uses the external IP addresses of CoreDNS service type `LoadBalancer`.

**- If you are NOT running CoreDNS service as a LoadBalancer service** (common for local development and testing):
There are no external IP addresses. In this case, k8gb will use the address of the first Ingress labeled with `k8gb.io/ip-source="true"`.

This ensures that DNS records created by k8gb always point to the correct addresses, regardless of your environment or infrastructure setup.

> **Note:** k8gb fails on bootstrap if it cannot discover the address of the CoreDNS service or the first Ingress with the `k8gb.io/ip-source="true"` label. This is to prevent misconfigurations that could lead to DNS resolution issues.

## Example for Local Development
If CoreDNS is not exposed as a LoadBalancer, label your Ingress to let k8gb know where to find the correct IP:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: k8gb-init-ingress
  namespace: k8gb
  labels:
    k8gb.io/ip-source: "true"
spec:
...
```
