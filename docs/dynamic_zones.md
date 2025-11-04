# Dynamic Zones (Feature)

This feature introduces new type of `Cluster Scoped` resource `ZoneDelegation` which carries DNS Zone information. Such as `parentZone`, `loadBalancedZone` and `dnsZoneNegTTL`

THere are cases where clusters are anonymous or shared across multiple tenants. When controller starts, we don't know what zones it will potentially participate into.

If we try to guess all potential loadBalancedZones in advance, prior application deployed into such cluster, k8gb's CoreDNS starts to be authoritative to such zone and be part of such zones query.
This leads to `NXDOMAIN` responses, since application is not yet deployed there.

`ZoneDelegation` resource helps us to align application deployment (or say `Team` onboard) and delegated zone creation, where zone is shipped together with `Appliation` and prevents `NXDOMAIN` to happen. And
will be reconciled once added.

For example:
```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: ZoneDelegation
metadata:
  generation: 1
  name: test-zone
spec:
  dnsZoneNegTTL: 30
  loadBalancedZone: test-zone.cloud.example.com
  parentZone: cloud.example.com 
```

Setup:
Set `k8gb.feature.dynamicZones` to true via helm chart values. This will add a new empty ConfigMap named `dynamic-zones`. Every `ZoneDelegation` reconcile loop, `k8gb` will
get all ZoneDelegations and create configmap key per `ZoneDelegation` like:
```yaml
apiVersion: v1
data:
  test-zone-cloud-example-com.conf: |2-
    test-zone.cloud.example.com:5353 {
      import k8gbplugins
    }
kind: ConfigMap
metadata:
  annotations:
    meta.helm.sh/release-name: k8gb
    meta.helm.sh/release-namespace: k8gb
  labels:
    app.kubernetes.io/managed-by: Helm
  name: coredns-dynamic
```
This configmap is mounted into CoreDNS pod and imported with import plugin:
```yaml
apiVersion: v1
data:
  Corefile: |-
    (k8gbplugins) {
        errors
        health
        reload 30s 15s
        ready
        prometheus 0.0.0.0:9153
        forward . /etc/resolv.conf
        k8s_crd {
            filter k8gb.absa.oss/dnstype=local
            negttl 30
            loadbalance weight
        }
    }
    static-zone.cloud.example.com:5353 {
        import k8gbplugins
    }
    import ../dynamic/*.conf
```

# ZoneDelegation Status
Holds an info about all DNSServers participating into zone delegation.

# ZoneDelegation Cleanup (WIP)
`ZoneDelegation` is protected by Finalizer, where on object removal, controller is responsible to clean up own reference in zone delegation and delete delegation completely if
current member is the last one standing.
