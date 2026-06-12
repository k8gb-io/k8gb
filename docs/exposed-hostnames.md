# Exposed Hostnames Annotation

The `k8gb.io/exposed-hostnames` annotation allows GSLB resources to dynamically resolve IP addresses from hostnames instead of using static IPs.

## Problem

In environments where nodes are behind NAT or use dynamic IPs (e.g., DDNS), hardcoding IP addresses with `k8gb.io/exposed-ip-addresses` doesn't work because the IPs can change. The `k8gb.io/exposed-hostnames` annotation solves this by resolving hostnames to IPs on each reconciliation loop.

## Usage

Add the annotation to your GSLB resource with a comma-separated list of hostnames:

```yaml
apiVersion: k8gb.io/v1beta1
kind: Gslb
metadata:
  name: my-app-gslb
  namespace: default
  annotations:
    k8gb.io/exposed-hostnames: "node1.ddns.example.com,node2.ddns.example.com"
spec:
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    matchLabels:
      app: my-app
  strategy:
    type: roundRobin
    dnsTtlSeconds: 30
  splitBrainThresholdSeconds: 300
```

## How It Works

1. On each reconciliation, k8gb reads the `k8gb.io/exposed-hostnames` annotation
2. Each hostname is resolved to its current IP address(es) using the configured DNS servers
3.  The resolved IPs are written to the generated `DNSEndpoint` records responsible for publishing the `localtargets-*` and final `Gslb` DNS A records.

This works across all supported resource types: Ingress, LoadBalancer Service, Istio VirtualService, and Gateway API.

## Precedence

When determining exposed IPs, k8gb checks annotations in this order:

1. `k8gb.io/exposed-hostnames` — resolve hostnames dynamically
2. `k8gb.io/exposed-ip-addresses` — use static IP addresses
3. Resource status — fall back to IPs from the resource's `.status.loadBalancer`

If `exposed-hostnames` is set, it takes precedence over `exposed-ip-addresses`.

## Validation

Hostnames are validated at admission time using a `ValidatingAdmissionPolicy` (when enabled). Each hostname must be a valid RFC 1123 hostname.

Invalid example:
```yaml
k8gb.io/exposed-hostnames: "not a valid hostname!"  # rejected
```

Valid examples:
```yaml
k8gb.io/exposed-hostnames: "node1.example.com"
k8gb.io/exposed-hostnames: "node1.example.com,node2.example.com"
```
