## Changing CoreDNS service type
Since Kubernetes API doesn't allow to change service type. Helm chart provides default label `k8gb-migrated-svc`,
when set during helm upgrade, old service deleted by helm upgrade hook in automatic manner. Here is the label pattern for service deletion `app.kubernetes.io/name=coredns,k8gb-migrated-svc!=true`

## Legacy API migration validation

For manual acceptance checks of legacy GSLB migration (`k8gb.absa.oss` -> `k8gb.io`), including pre-merge and post-merge upgrade validation cases, see [Legacy GSLB Migration Manual Acceptance Tests](migration_acceptance.md).

## Controlled legacy migration model

To reduce upgrade blast radius, legacy migration is controlled explicitly per legacy object:

- Legacy API (`k8gb.absa.oss/v1beta1`) remains supported during transition.
- Legacy reconcile behavior remains active before migration.
- Legacy objects emit deprecation warning events.
- Migration starts only when label `k8gb.io/migration-requested=true` is set.
- Migration completion is marked by `k8gb.io/migrated-to-k8gb-io=true`.
- New features are implemented only in canonical `k8gb.io/v1beta1`.

### Reconcile decision table

| Legacy labels on object | Legacy reconcile mode | Migration action | Expected operator behavior |
|---|---|---|---|
| no migration labels | full legacy reconcile | none | continue normal operation; plan migration |
| `k8gb.io/migration-requested=true`, not migrated | legacy reconcile + migration | execute one-way migration | switch edits to canonical `k8gb.io` object |
| `k8gb.io/migrated-to-k8gb-io=true` | compatibility/read-only | none | treat legacy object as compatibility artifact |
| both labels set | compatibility/read-only | none | optional cleanup of request label |
