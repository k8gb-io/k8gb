## Changing CoreDNS service type
Since Kubernetes API doesn't allow to change service type. Helm chart provides default label `k8gb-migrated-svc`,
when set during helm upgrade, old service deleted by helm upgrade hook in automatic manner. Here is the label pattern for service deletion `app.kubernetes.io/name=coredns,k8gb-migrated-svc!=true`
