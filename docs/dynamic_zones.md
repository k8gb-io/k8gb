# Dynamic Zones

## High-Level Summary
Dynamic Zones allow k8gb to serve DNS zones only when applications are deployed, solving the problem of premature `NXDOMAIN`
responses in anonymous, multi-tenant, or shared clusters. By introducing a new cluster-scoped resource, `ZoneDelegation`, 
k8gb dynamically detects which zones a cluster should serve and configures CoreDNS automatically.

This provides:
- Accurate DNS behavior (no premature `NXDOMAIN`)
- Declarative onboarding
- Automatic CoreDNS reconfiguration without restarts
- Safe cleanup when zones are removed

---

# Dynamic Zones (Feature)

The feature introduces a cluster-scoped resource `ZoneDelegation` that contains:

- `loadBalancedZone`
- `parentZone`
- `dnsZoneNegTTL`

## Why Dynamic Zones?

Pre-configuring zones in anonymous or shared clusters makes CoreDNS authoritative too early, causing `NXDOMAIN` responses.
`ZoneDelegation` lets applications ship their own zone definitions, enabling dynamic activation.

## Example

```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: ZoneDelegation
metadata:
  name: test-zone
spec:
  loadBalancedZone: test-zone.cloud.example.com
  parentZone: cloud.example.com
  dnsZoneNegTTL: 30
```

## Setup

Set `k8gb.dynamicZones` to true via helm chart values. This will add a new empty ConfigMap named `dynamic-zones`. 
Every `ZoneDelegation` reconcile loop, `k8gb` will get all ZoneDelegations and create configmap key per `ZoneDelegation` like:

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
  name: k8gb-zone-delegation
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
`ZoneDelegation` is protected by Finalizer, where on object removal, controller is responsible to clean up own reference 
in zone delegation and delete delegation completely if current member is the last one standing.
