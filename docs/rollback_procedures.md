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

You may also come across issues with CRD conflicts during rollback. If so, follow these steps:

1. **Delete existing CRDs:**

```bash
kubectl delete crd gslbs.k8gb.io
kubectl delete crd healthchecks.k8gb.io
```
2. **Reinstall previous version:**

```bash
helm upgrade k8gb k8gb/k8gb --version v0.14.0 -n k8gb -f new-values.yaml
```
3. **Verify installation:**

```bash
helm list -n k8gb
kubectl get pods -n k8gb
kubectl get gslb -A
```

### Alternative: Manual CRD Annotation and Labeling
```bash
kubectl annotate crd dnsendpoints.externaldns.k8s.io meta.helm.sh/release-name=k8gb --overwrite

kubectl annotate crd dnsendpoints.externaldns.k8s.io meta.helm.sh/release-namespace=k8gb --overwrite

kubectl label crd dnsendpoints.externaldns.k8s.io app.kubernetes.io/managed-by=Helm --overwrite

helm upgrade k8gb k8gb/k8gb --version v0.14.0 -n k8gb -f new-values.yaml
```
This method manually sets the necessary Helm annotations and labels on the CRDs to allow Helm to manage them during the rollback.

You can also fix this by using the `--force` flag with Helm upgrade:
```bash
helm upgrade k8gb k8gb/k8gb --version v0.14.0 -n k8gb -f new-values.yaml --force
```