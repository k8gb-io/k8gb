# K8GB Rollback Procedures

## Rollback Compatibility

| From | To | Method | Status |
|------|----|---------| -------|
| v0.17.0 | v0.16.0 | `helm rollback` | ✅ Works |
| v0.16.0+| v0.15.0 | `helm upgrade` + values file | ⚠️ Needs CoreDNS fix |
| v0.15.0+ | v0.14.0 | `helm upgrade` + schema conversion | ⚠️ Needs dnsZones conversion |

## v0.16.0+ → v0.15.0 Rollback

### Issue
Direct rollback fails due to CoreDNS service configuration changes.

Note: v0.16.0 introduced zones configuration that breaks rollback to v0.15.0

### Solution
Use values file with v0.15.0 compatible CoreDNS config:

```yaml
# v015-rollback-values.yaml
coredns:
# v0.15.0 compatible CoreDNS configuration
  servers:
  - port: 5353
    servicePort: 53
    plugins:
    - name: prometheus
      parameters: 0.0.0.0:9153
```

**Rollback command:**
```bash
helm upgrade k8gb k8gb/k8gb --version v0.15.0 -n k8gb -f v015-rollback-values.yaml
```

## v0.15.0+→ v0.14.0 Rollback

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

You may also run into CRD ownership conflicts when downgrading/rolling back, because the CRD already exists but Helm can’t “adopt” it. Since CRDs are cluster-scoped, Helm relies on specific annotations/labels to recognize them as part of the release.

```bash
kubectl annotate crd dnsendpoints.externaldns.k8s.io meta.helm.sh/release-name=k8gb --overwrite

kubectl annotate crd dnsendpoints.externaldns.k8s.io meta.helm.sh/release-namespace=k8gb --overwrite

kubectl label crd dnsendpoints.externaldns.k8s.io app.kubernetes.io/managed-by=Helm --overwrite

helm upgrade k8gb k8gb/k8gb --version v0.14.0 -n k8gb -f new-values.yaml
```
This method manually sets the necessary Helm annotations and labels on the CRDs to allow Helm to manage them during the rollback.
> NOTE: by rollback, we mean downgrading to the previous version, not specific `helm rollback` functionality 
If the chart version changes the CRD schema (or conversion webhook behavior), downgrading can still fail even after adoption, because existing CR instances may not validate against the older schema. In that case you may need to align CRD versions more explicitly (or avoid downgrading CRDs altogether).
You can use --force, but treat it as a last resort because it may replace resources to apply changes.:
```bash
helm upgrade k8gb k8gb/k8gb --version v0.14.0 -n k8gb -f new-values.yaml --force
```
