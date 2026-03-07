## Changing CoreDNS service type
Since Kubernetes API doesn't allow to change service type. Helm chart provides default label `k8gb-migrated-svc`,
when set during helm upgrade, old service deleted by helm upgrade hook in automatic manner. Here is the label pattern for service deletion `app.kubernetes.io/name=coredns,k8gb-migrated-svc!=true`

## Legacy API migration validation

For manual acceptance checks of legacy GSLB migration (`k8gb.absa.oss` -> `k8gb.io`), including pre-merge and post-merge upgrade validation cases, see [Legacy GSLB Migration Manual Acceptance Tests](migration_acceptance.md).
