# K8GB Rollback Procedures

> NOTE: by rollback, we mean downgrading to a previous chart/app version (typically
> `helm upgrade --version …`), not specifically `helm rollback`, unless noted.

## Rollback Compatibility

| From | To | Method | Status |
|------|----|---------| -------|
| v0.20.x | v0.19.x | `helm upgrade` | ⚠️ Chart OK if no ZoneDelegation / TLSRoute / new features in use; see [API-group caveats](#rolling-back-across-the-k8gbio-api-group-migration-v019) |
| v0.19.x | v0.18.x | `helm upgrade` | ❌ Unsupported once objects were migrated to `k8gb.io` (one-way); restore from backup if you must leave v0.19 |
| v0.18.x | v0.17.x | `helm upgrade` | ✅ Usually works (no API-group change); verify CRDs / values |
| v0.17.0 | v0.16.0 | `helm rollback` | ✅ Works |
| v0.16.0+| v0.15.0 | `helm upgrade` + values file | ⚠️ Needs CoreDNS fix |
| v0.15.0+ | v0.14.0 | `helm upgrade` + schema conversion | ⚠️ Needs dnsZones conversion |

## v0.20.x → v0.19.x Rollback

### When this is safe

`helm upgrade` back to `v0.19.x` is reasonable when the cluster does **not** depend on features introduced in v0.20 (notably `ZoneDelegation`, TLSRoute support, and related status/finalizer behavior). Keep your existing Helm values (including `dnsZones` / `edgeDNSServers` / `extdns`).

```bash
helm upgrade k8gb k8gb/k8gb --version v0.19.0 -n k8gb -f current-values.yaml
kubectl get pods -n k8gb
kubectl get gslb -A
kubectl get gslb.k8gb.io -A
```

### When this is not safe

If you created `ZoneDelegation` objects (or other v0.20-only CRs) and the older chart/CRDs cannot represent them, delete or migrate those resources **before** downgrading, or restore the control plane from a pre-upgrade backup. Downgrading CRDs under live incompatible instances can leave objects that fail validation.

## Rolling back across the `k8gb.io` API-group migration (v0.19+)

v0.19 introduced the vendor-neutral canonical API group `k8gb.io/v1beta1` alongside legacy
`k8gb.absa.oss/v1beta1`. Migration is **one-way** and **label-triggered**. See
[Controlled legacy migration model](service_upgrade.md#controlled-legacy-migration-model) and
[ADR-0002](../adr/0002-migrate-gslb-api-group-to-vendor-neutral-k8gb-io.md).

### Labels and finalizer

| Label / finalizer | Meaning |
|---|---|
| `k8gb.io/migration-requested=true` | Operator requested migration for this legacy object; legacy runtime reconcile pauses for that object |
| `k8gb.io/migrated-to-k8gb-io=true` | Migration finished; edit the canonical `k8gb.io` Gslb going forward |
| `k8gb.io/legacy-migration-protection` | Finalizer on the legacy object during/after migration for safe ownerReference cleanup (`controllers/gslb_migration_controller.go`) |

Trigger (for reference — do **not** use this as a rollback step):

```bash
kubectl label gslb.k8gb.absa.oss -n <ns> <name> k8gb.io/migration-requested=true --overwrite
```

### What is unsupported

| Situation | Rollback guidance |
|---|---|
| Still on legacy only (`k8gb.absa.oss`, no migration labels) | Downgrade v0.19 → v0.18 may work; verify pods and `kubectl get gslb.k8gb.absa.oss -A` |
| Migration requested or completed (`migration-requested` / `migrated-to-k8gb-io`) | **Do not** expect a supported in-place rollback to a pre-migration world. Canonical `k8gb.io` objects and ownership changes are not automatically reversed |
| Need pre-migration state | Restore Gslb/Ingress/DNSEndpoint objects (and CRDs if needed) from a backup taken **before** setting `k8gb.io/migration-requested=true` |

There is no supported “un-migrate” controller path that recreates a writable legacy object from a canonical `k8gb.io` Gslb.

### If you must leave v0.19+ after migration

1. Prefer staying on a v0.19+ release and fixing forward.
2. If a downgrade is unavoidable, restore from backup rather than inventing a reverse conversion.
3. After restore, confirm which API group your automation applies (`k8gb.io` vs `k8gb.absa.oss`) before starting the operator again.

## v0.18.x → v0.17.x Rollback

No API-group migration in this step. Typical path:

```bash
helm get values k8gb -n k8gb > current-values.yaml
helm upgrade k8gb k8gb/k8gb --version v0.17.0 -n k8gb -f current-values.yaml
kubectl get pods -n k8gb
kubectl get gslb -A
```

If Helm refuses to adopt existing CRDs, see the CRD ownership notes under
[v0.15.0+ → v0.14.0](#v0150-v0140-rollback).

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
If the chart version changes the CRD schema (or conversion webhook behavior), downgrading can still fail even after adoption, because existing CR instances may not validate against the older schema. In that case you may need to align CRD versions more explicitly (or avoid downgrading CRDs altogether).
You can use --force, but treat it as a last resort because it may replace resources to apply changes.:
```bash
helm upgrade k8gb k8gb/k8gb --version v0.14.0 -n k8gb -f new-values.yaml --force
```
