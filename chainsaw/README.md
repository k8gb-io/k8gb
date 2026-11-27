# Useful commands

Run a single test:

`chainsaw test --config ./config.yaml --values ./values.yaml --test-dir ./tests/<testname>`

CI explicitly enables `OwnerReferencesPermissionEnforcement` on both application
clusters. For the same checks locally, add this argument when creating each
cluster (existing clusters must be recreated):

```sh
--k3s-arg '--kube-apiserver-arg=enable-admission-plugins=OwnerReferencesPermissionEnforcement@server:*'
```

The ownership regression tests exercise the controller with the chart's normal
service account and RBAC:

- `rbac-owner-reference-admission` checks fresh canonical Gslb-owned Ingresses and
  DNSEndpoints. Restricted namespace credentials must fail blocking-owner creation
  and owner removal, proving admission enforcement is active. It also checks
  Ingress-owned Gslb creation and status updates through the annotation flow.
- `legacy-migration-embedded` starts with controller-created legacy dependents,
  then checks Ingress owner removal and DNSEndpoint ownership transfer without
  recreating either object.
- `legacy-deletion-embedded` models a protected legacy object deleted before
  migration completes and checks both dependents survive with their owners
  detached and the migration finalizer removed.

The existing `ingress-playground-resourceref` test also checks the deprecated
annotation flow, where the Ingress owns the generated Gslb. It requires permission
to update Ingress finalizers even outside OpenShift.
