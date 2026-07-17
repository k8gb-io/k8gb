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
- Reduced DNS lookup overhead by using `ZoneDelegation` status as the source of DNS server IPs

---

## Dynamic Zones (Feature)

The feature introduces a cluster-scoped resource `ZoneDelegation` that contains:

- `loadBalancedZone`
- `parentZone`
- `dnsZoneNegTTL`
- `doFinalize`

## Why Dynamic Zones?

Pre-configuring zones in anonymous or shared clusters makes CoreDNS authoritative too early, causing `NXDOMAIN` responses.
`ZoneDelegation` lets applications ship their own zone definitions, enabling dynamic activation.

## Deploying Load-Balanced Zones to a Subset of Clusters

In multi-cluster environments, a load-balanced application does not have to be deployed to every cluster in the k8gb network.
The `ZoneDelegation` resource should be deployed together with the application to each cluster that participates in serving
that application's `loadBalancedZone`.

If a cluster does not have a `ZoneDelegation` for a `loadBalancedZone`, k8gb does not configure CoreDNS to serve that zone
from the cluster. The cluster is therefore excluded from that GSLB zone instead of returning inconsistent DNS responses for
an application it does not host.

When different applications are deployed to different cluster subsets, give each application or team its own delegated
sub-zone. For example, deploy `ZoneDelegation` for `app1.gslb.example.com` only in the clusters that host `app1`, and
deploy `ZoneDelegation` for `app2.gslb.example.com` only in the clusters that host `app2`.

## Example

```yaml
apiVersion: k8gb.io/v1beta1
kind: ZoneDelegation
metadata:
  name: test-zone
spec:
  loadBalancedZone: test-zone.cloud.example.com
  parentZone: cloud.example.com
  dnsZoneNegTTL: 30
  doFinalize: true
```

## Setup

`ZoneDelegation` resources can be created directly, and k8gb will reconcile them automatically. For backward compatibility, 
k8gb can also generate `ZoneDelegation` resources from the `dnsZones` configuration in `values.yaml`. Internally, however, 
k8gb operates only with `ZoneDelegation` resources. The recommended approach is therefore to deploy `ZoneDelegation` 
resources explicitly together with the applications that require them, rather than defining zones through Helm values. 

When a `ZoneDelegation` is reconciled, k8gb updates the `k8gb-zone-delegation` ConfigMap with the CoreDNS 
configuration for the delegated `loadBalancedZone`. This ConfigMap is mounted into the k8gb CoreDNS pods and imported 
by the CoreDNS configuration.

> **Important:** When using the ExternalDNS provider, configure `extdns.domainFilters` explicitly and include every 
> `parentZone` that ExternalDNS is expected to manage.
>  ```yaml
> # Include every parent zone managed by ExternalDNS.
> domainFilters:
>    - "example.com"
>  ```

## ZoneDelegation Status
The ZoneDelegation status contains information about all DNS servers participating in zone delegation.
It acts as the single source of truth for DNS server IP addresses used by the k8gb controller,
delegation DNSEndpoints when the ExternalDNS provider is enabled, and application DNSEndpoints.
k8gb reads DNS server IP addresses from this status, which avoids resolving them from edgeDNS
on every reconciliation. The refresh interval is controlled by the k8gb reconcile interval, configured through
`k8gb.reconcileRequeueSeconds` in Helm values or the corresponding `RECONCILE_REQUEUE_SECONDS` environment variable.

```yaml
status:
  dnsServers:
    - name: gslb-ns-eu-cloud2.example.com
      address: 172.18.0.6
    - name: gslb-ns-eu-cloud2.example.com
      address: 172.18.0.7
    - name: gslb-ns-us-cloud2.example.com
      address: 172.18.0.11
    - name: gslb-ns-us-cloud2.example.com
      address: 172.18.0.12
```

## ZoneDelegation Cleanup (WIP)
Use `doFinalize: false` when k8gb should not modify the parent DNS delegation during deletion; the controller only removes this cluster's local CoreDNS configuration.
Use `doFinalize: true` when k8gb is allowed to remove the delegated zone from the parent DNS during cleanup.

### Delayed Finalization
When `doFinalize` is enabled, deleting a `ZoneDelegation` may not remove the object immediately. The object can remain in
a terminating state while the k8gb controller verifies whether the current cluster is the last remaining member of the delegation.
This prevents premature deletion of the delegated zone from the parent DNS in multi-cluster environments.

```yaml
metadata:
  deletionTimestamp: ...
  finalizers:
    - k8gb.io/finalizer
```

### Infoblox provider
When `ZoneDelegation` is removed with the Infoblox provider, the delegated zone is deleted only when the last `ZoneDelegation`
object is removed. The zone is deleted directly by the Infoblox provider, which is part of k8gb.

### ExternalDNS provider
When `ZoneDelegation` is removed with the ExternalDNS provider, delegation `DNSEndpoint` objects must not be deleted immediately.
Instead, deleting a `ZoneDelegation` removes the current cluster target from the delegation `DNSEndpoint` and removes the
zone configuration from the `k8gb-zone-delegation` ConfigMap. The `ZoneDelegation` object then remains blocked by its finalizer.
The last remaining `ZoneDelegation` object no longer has information about any other participating clusters, so it can be finalized and deleted.
The remaining `ZoneDelegation` objects are finalized automatically after they lose information about the last cluster.
Together with the `ZoneDelegation` objects, the delegation `DNSEndpoint` objects are deleted as well. This allows the ExternalDNS provider
to remove the delegated zone from the parent DNS.

## Backward Compatibility
Static zones configured through `k8gb.dnsZones` remain supported. Dynamic Zones are an additional mechanism for environments
where zones should be activated only when `ZoneDelegation` objects are created. For backward compatibility, `ZoneDelegation`
objects generated from `k8gb.dnsZones` or the `DNS_ZONES` environment variable have `doFinalize` set to `false` by default.
This prevents k8gb from deleting delegated zones that were originally configured as static zones.
