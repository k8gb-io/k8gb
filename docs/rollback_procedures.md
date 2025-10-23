# K8GB Rollback Procedures

## v0.15.0 → v0.14.0 Rollback

### Issue

Direct rollback fails due to breaking changes in v0.15.0 Helm values schema.

In v0.15.0, the dnsZones array replaces the dnsZone and edgeDNSZone values.
dnsZone + edgeDNSZone → dnsZones array

### Solution: Convert Values Before Rollback

1. **Extract current values:**

```bash
helm get values k8gb -n k8gb > new-values.yaml
```

2. **Convert values to v0.14.0 schema:**

```bash
# v0.15.0 format
k8gb:
  dnsZones:
    - parentZone: "example.com"
      loadBalancedZone: "cloud.example.com"

# v0.14.0 format
k8gb:
  dnsZone: "cloud.example.com"
  edgeDNSZone: "example.com"
```

3. **Rollback:**

```bash
helm upgrade k8gb k8gb/k8gb --version v0.14.0 -n k8gb -f new-values.yaml
```

4. **Verify rollback:**

```bash
helm list -n k8gb
kubectl get pods -n k8gb
kubectl get gslb -A
```
