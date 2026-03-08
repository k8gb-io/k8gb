# Legacy GSLB Migration Manual Acceptance Tests

This document defines **manual** acceptance tests for the legacy API migration:

- from `k8gb.absa.oss/v1beta1` to `k8gb.io/v1beta1`
- including embedded Ingress migration to `resourceRef`
- including ownerReference handoff to prevent accidental Ingress garbage collection

These checks are intentionally not automated. They are meant for:

- **Pre-merge confidence** on the feature branch
- **Post-merge upgrade validation** before or during release rollout

## Scope

The migration behavior validated here:

1. Canonical GSLB is created from legacy GSLB (same name/namespace).
2. Legacy GSLB is labeled `k8gb.io/migrated-to-k8gb-io=true`.
3. Embedded legacy ingress is migrated to `spec.resourceRef` and detached from embedded `spec.ingress`.
4. For embedded legacy mode, legacy GSLB ownerReference is removed from the Ingress.
5. Existing canonical GSLB is not overwritten.
6. Already migrated legacy objects are ignored.
7. Legacy deletion path removes migration finalizer after cleanup.
8. Embedded migration does not fail if the Ingress is missing.
9. OwnerReference cleanup also runs for already-migrated embedded legacy objects.
10. Finalizer cleanup does not block deletion when embedded Ingress is missing.

For controlled migration mode, cases that expect canonical object creation require explicit trigger label:

- `k8gb.io/migration-requested=true`

## Prerequisites

- Kubernetes cluster with k8gb installed from the target branch/release candidate
- Legacy CRD installed (`gslbs.k8gb.absa.oss`)
- `kubectl` configured for the test cluster
- `jq` recommended for easier inspection

## Pre-Merge Manual Acceptance

Use a dedicated namespace:

```bash
kubectl create ns migration-e2e || true
```

Migration trigger helper (use in cases that expect conversion):

```bash
kubectl label gslb.k8gb.absa.oss -n migration-e2e <name> k8gb.io/migration-requested=true --overwrite
```

Cases requiring this trigger: 1, 2, 3, 4, 6, and 8.

### Case 1: Legacy referenced GSLB migrates to canonical

Apply legacy referenced GSLB:

```bash
kubectl apply -f - <<'EOF_CASE1'
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: legacy-ref
  namespace: migration-e2e
spec:
  ingressClassName: nginx
  rules:
    - host: legacy-ref.cloud.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: dummy
                port:
                  number: 80
---
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: legacy-ref
  namespace: migration-e2e
spec:
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    name: legacy-ref
  strategy:
    type: roundRobin
EOF_CASE1
```

Validate:

```bash
kubectl get gslb.k8gb.io -n migration-e2e legacy-ref
kubectl get gslb.k8gb.absa.oss -n migration-e2e legacy-ref -o jsonpath='{.metadata.labels.k8gb\.io/migrated-to-k8gb-io}'
```

Expected:

- Canonical object exists.
- Legacy label equals `true`.

### Case 2: Embedded legacy GSLB migrates to `resourceRef`

Apply embedded legacy GSLB:

```bash
kubectl apply -f - <<'EOF_CASE2'
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: legacy-embedded
  namespace: migration-e2e
spec:
  ingress:
    ingressClassName: nginx
    rules:
      - host: legacy-embedded.cloud.example.com
        http:
          paths:
            - path: /
              pathType: Prefix
              backend:
                service:
                  name: dummy
                  port:
                    number: 80
  strategy:
    type: roundRobin
EOF_CASE2
```

Validate canonical spec:

```bash
kubectl get gslb.k8gb.io -n migration-e2e legacy-embedded -o yaml
```

Expected in canonical object:

- `spec.resourceRef.apiVersion = networking.k8s.io/v1`
- `spec.resourceRef.kind = Ingress`
- `spec.resourceRef.name = legacy-embedded`
- `spec.ingress` is empty

### Case 3: Existing canonical object is not overwritten

Apply canonical first, then legacy with same name:

```bash
kubectl apply -f - <<'EOF_CASE3'
apiVersion: k8gb.io/v1beta1
kind: Gslb
metadata:
  name: existing-canonical
  namespace: migration-e2e
spec:
  strategy:
    type: failover
    primaryGeoTag: eu
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    name: existing-canonical
---
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: existing-canonical
  namespace: migration-e2e
spec:
  strategy:
    type: roundRobin
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    name: existing-canonical
EOF_CASE3
```

Validate:

```bash
kubectl get gslb.k8gb.io -n migration-e2e existing-canonical -o jsonpath='{.spec.strategy.type}'
kubectl get gslb.k8gb.absa.oss -n migration-e2e existing-canonical -o jsonpath='{.metadata.labels.k8gb\.io/migrated-to-k8gb-io}'
```

Expected:

- Canonical strategy remains `failover` (unchanged).
- Legacy label equals `true`.

### Case 4: OwnerReference handoff and deletion safety

For an embedded legacy object, verify Ingress survives legacy deletion.

Capture ownerRefs before deletion:

```bash
kubectl get ingress -n migration-e2e legacy-embedded -o json | jq '.metadata.ownerReferences'
kubectl get gslb.k8gb.absa.oss -n migration-e2e legacy-embedded -o jsonpath='{.metadata.finalizers}'
```

Delete legacy object:

```bash
kubectl delete gslb.k8gb.absa.oss -n migration-e2e legacy-embedded
```

Validate:

```bash
kubectl get ingress -n migration-e2e legacy-embedded
kubectl get ingress -n migration-e2e legacy-embedded -o json | jq '.metadata.ownerReferences'
# Expected to fail with NotFound (legacy object should be gone after finalizer cleanup)
kubectl get gslb.k8gb.absa.oss -n migration-e2e legacy-embedded || true
```

Expected:

- Ingress still exists.
- OwnerReferences do not contain the deleted legacy GSLB UID.
- Legacy object is deleted after cleanup.

### Case 5: Already migrated legacy object is ignored

Apply an already-migrated legacy object:

```bash
kubectl apply -f - <<'EOF_CASE5'
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: already-migrated
  namespace: migration-e2e
  labels:
    k8gb.io/migrated-to-k8gb-io: "true"
spec:
  strategy:
    type: roundRobin
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    name: already-migrated
EOF_CASE5
```

Validate:

```bash
# Expected to fail with NotFound (controller should ignore already-migrated legacy objects)
kubectl get gslb.k8gb.io -n migration-e2e already-migrated || true
kubectl get events -n migration-e2e --sort-by=.lastTimestamp | grep -E 'LegacyIgnored|already-migrated'
```

Expected:

- Canonical object is not created automatically for this case.
- `LegacyIgnored` event is emitted.

### Case 6: Embedded legacy migration succeeds even if Ingress is missing

Apply embedded legacy GSLB without creating a matching Ingress resource:

```bash
kubectl apply -f - <<'EOF_CASE6'
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: embedded-no-ingress
  namespace: migration-e2e
spec:
  ingress:
    ingressClassName: nginx
    rules:
      - host: embedded-no-ingress.cloud.example.com
        http:
          paths:
            - path: /
              pathType: Prefix
              backend:
                service:
                  name: dummy
                  port:
                    number: 80
  strategy:
    type: roundRobin
EOF_CASE6
```

Validate:

```bash
kubectl get gslb.k8gb.io -n migration-e2e embedded-no-ingress
kubectl get gslb.k8gb.absa.oss -n migration-e2e embedded-no-ingress -o jsonpath='{.metadata.labels.k8gb\.io/migrated-to-k8gb-io}'
```

Expected:

- Canonical object exists even though Ingress is missing.
- Legacy label equals `true`.

### Case 7: Already-migrated embedded legacy still detaches stale Ingress ownerRef

Apply already-migrated embedded legacy object:

```bash
kubectl apply -f - <<'EOF_CASE7'
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: stale-owner
  namespace: migration-e2e
  labels:
    k8gb.io/migrated-to-k8gb-io: "true"
spec:
  ingress:
    ingressClassName: nginx
    rules:
      - host: stale-owner.cloud.example.com
        http:
          paths:
            - path: /
              pathType: Prefix
              backend:
                service:
                  name: dummy
                  port:
                    number: 80
  strategy:
    type: roundRobin
EOF_CASE7
```

Create an Ingress with a stale ownerReference to the legacy GSLB UID, then retrigger reconcile:

```bash
LEGACY_UID=$(kubectl get gslb.k8gb.absa.oss -n migration-e2e stale-owner -o jsonpath='{.metadata.uid}')

kubectl apply -f - <<EOF_CASE7_ING
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: stale-owner
  namespace: migration-e2e
  ownerReferences:
    - apiVersion: k8gb.absa.oss/v1beta1
      kind: Gslb
      name: stale-owner
      uid: ${LEGACY_UID}
spec:
  ingressClassName: nginx
  rules:
    - host: stale-owner.cloud.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: dummy
                port:
                  number: 80
EOF_CASE7_ING

kubectl annotate gslb.k8gb.absa.oss -n migration-e2e stale-owner migration-recheck-ts="$(date +%s)" --overwrite
```

Validate:

```bash
kubectl get ingress -n migration-e2e stale-owner -o json | jq '.metadata.ownerReferences'
```

Expected:

- Legacy ownerReference is removed from the Ingress.

### Case 8: Existing canonical + embedded legacy still performs owner cleanup

Create canonical object first:

```bash
kubectl apply -f - <<'EOF_CASE8_CANONICAL'
apiVersion: k8gb.io/v1beta1
kind: Gslb
metadata:
  name: canonical-embedded-cleanup
  namespace: migration-e2e
spec:
  strategy:
    type: failover
    primaryGeoTag: eu
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    name: canonical-embedded-cleanup
EOF_CASE8_CANONICAL
```

Apply embedded legacy object with same name:

```bash
kubectl apply -f - <<'EOF_CASE8_LEGACY'
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: canonical-embedded-cleanup
  namespace: migration-e2e
spec:
  ingress:
    ingressClassName: nginx
    rules:
      - host: canonical-embedded-cleanup.cloud.example.com
        http:
          paths:
            - path: /
              pathType: Prefix
              backend:
                service:
                  name: dummy
                  port:
                    number: 80
  strategy:
    type: roundRobin
EOF_CASE8_LEGACY
```

Inject stale ownerReference and retrigger reconcile:

```bash
LEGACY_UID=$(kubectl get gslb.k8gb.absa.oss -n migration-e2e canonical-embedded-cleanup -o jsonpath='{.metadata.uid}')

kubectl apply -f - <<EOF_CASE8_ING
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: canonical-embedded-cleanup
  namespace: migration-e2e
  ownerReferences:
    - apiVersion: k8gb.absa.oss/v1beta1
      kind: Gslb
      name: canonical-embedded-cleanup
      uid: ${LEGACY_UID}
spec:
  ingressClassName: nginx
  rules:
    - host: canonical-embedded-cleanup.cloud.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: dummy
                port:
                  number: 80
EOF_CASE8_ING

kubectl annotate gslb.k8gb.absa.oss -n migration-e2e canonical-embedded-cleanup migration-recheck-ts="$(date +%s)" --overwrite
```

Validate:

```bash
kubectl get gslb.k8gb.io -n migration-e2e canonical-embedded-cleanup -o jsonpath='{.spec.strategy.type}'
kubectl get ingress -n migration-e2e canonical-embedded-cleanup -o json | jq '.metadata.ownerReferences'
```

Expected:

- Canonical strategy stays `failover` (not overwritten).
- Legacy ownerReference is removed from Ingress.

### Case 9: Deletion path does not block when embedded Ingress is absent

Apply embedded legacy object:

```bash
kubectl apply -f - <<'EOF_CASE9'
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: delete-no-ingress
  namespace: migration-e2e
spec:
  ingress:
    ingressClassName: nginx
    rules:
      - host: delete-no-ingress.cloud.example.com
        http:
          paths:
            - path: /
              pathType: Prefix
              backend:
                service:
                  name: dummy
                  port:
                    number: 80
  strategy:
    type: roundRobin
EOF_CASE9
```

Ensure no matching Ingress exists, then delete legacy object:

```bash
kubectl delete ingress -n migration-e2e delete-no-ingress --ignore-not-found
kubectl delete gslb.k8gb.absa.oss -n migration-e2e delete-no-ingress
# Expected to fail with NotFound once deletion completes
kubectl get gslb.k8gb.absa.oss -n migration-e2e delete-no-ingress || true
```

Expected:

- Legacy object deletion is not stuck in `Terminating`.
- Finalizer cleanup completes even without Ingress.

## Post-Merge Upgrade Validation

Run this on a real upgrade path (cluster already running older release with legacy GSLBs).

### 1. Snapshot before upgrade

```bash
kubectl get gslb.k8gb.absa.oss -A -o yaml > pre-legacy-gslb.yaml
kubectl get gslb.k8gb.io -A -o yaml > pre-canonical-gslb.yaml
kubectl get ingress -A -o yaml > pre-ingress.yaml
```

### 2. Upgrade k8gb

Use your standard Helm upgrade flow for the target version.

### 3. Validate migration outcomes

Request migration for legacy objects you want to migrate in this wave:

```bash
kubectl get gslb.k8gb.absa.oss -A -o json \
  | jq -r '.items[] | [.metadata.namespace, .metadata.name] | @tsv' \
  | while IFS=$'\t' read -r ns name; do
      kubectl label gslb.k8gb.absa.oss -n "$ns" "$name" k8gb.io/migration-requested=true --overwrite
    done
```

Verify every legacy object has matching canonical object:

```bash
kubectl get gslb.k8gb.absa.oss -A -o json \
  | jq -r '.items[] | [.metadata.namespace, .metadata.name] | @tsv' \
  | while IFS=$'\t' read -r ns name; do
      kubectl get gslb.k8gb.io -n "$ns" "$name" >/dev/null || echo "MISSING canonical: $ns/$name"
    done
```

Verify migration labels:

```bash
kubectl get gslb.k8gb.absa.oss -A -o json \
  | jq -r '.items[] | select(.metadata.labels["k8gb.io/migrated-to-k8gb-io"] != "true") | "UNLABELED: \(.metadata.namespace)/\(.metadata.name)"'
```

Check migration events:

```bash
kubectl get events -A --sort-by=.lastTimestamp | grep -E 'LegacyMigrated|LegacyIgnored'
```

### 4. Validate deletion safety on a sampled embedded object

Pick at least one migrated embedded legacy GSLB and delete the legacy object:

```bash
kubectl delete gslb.k8gb.absa.oss -n <ns> <name>
kubectl get ingress -n <ns> <name>
```

Expected:

- Ingress remains present after legacy deletion.

## Cleanup

```bash
kubectl delete ns migration-e2e
```
