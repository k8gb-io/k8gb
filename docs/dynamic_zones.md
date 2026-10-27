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
When `doFinalize` is enabled, deleting a `ZoneDelegation` may not remove the object immediately.
The finalizer remains until the provider's cleanup requirements are satisfied. With ExternalDNS,
local withdrawal can finish while other clusters continue serving the zone.

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
Deleting a `ZoneDelegation` removes its local CoreDNS zone configuration and removes only
that cluster's NS target and glue A record from its delegation `DNSEndpoint`.
Because ExternalDNS publishes changes asynchronously, k8gb retains the finalizer until every
parent authority discovered through the configured DNS servers confirms that the local NS target
and glue A record are absent. Configured servers may be recursive resolvers; k8gb discovers the
parent zone NS records and resolves their IPv4/IPv6 addresses before querying those authorities
without recursion, using reachable addresses for each nameserver. A configured server answering authoritatively for the parent is queried directly.
Only a parent delegation response or an authoritative negative response with the parent SOA
can confirm withdrawal. Empty nonauthoritative replies and DNS errors postpone cleanup. Other clusters can continue serving the zone, and their presence
does not prevent the local `ZoneDelegation` from being deleted and subsequently recreated.

The delegation `DNSEndpoint` retains peer NS targets while those peers remain active. It is
marked with `k8gb.io/withdrawn-zone-delegation` and `k8gb.io/withdrawn-parent-zone` so cleanup
continues independently of the deleted `ZoneDelegation`, including after controller restarts.
Deferred cleanup discovers the saved parent's current authorities and queries peer glue without
recursion. It deletes the retained endpoint only when every authority confirms peer glue absence
with an authoritative negative answer carrying the parent SOA. Cached recursive negatives, empty
nonauthoritative replies, and referrals do not prove absence; cleanup retains the endpoint and retries.
ExternalDNS can then remove its remaining delegation records. Recreating the delegation reuses
the endpoint and clears both cleanup markers.

Endpoints marked by an older version without `k8gb.io/withdrawn-parent-zone` are retained safely.
To resume their deferred cleanup, add that annotation with the original `ZoneDelegation.spec.parentZone`
value, or recreate the delegation so its next withdrawal persists the parent automatically.

The last participating cluster still requests deletion of its whole delegation `DNSEndpoint`.
Delegations with `doFinalize: false` retain their existing behavior.

## Backward Compatibility
Static zones configured through `k8gb.dnsZones` remain supported. Dynamic Zones are an additional mechanism for environments
where zones should be activated only when `ZoneDelegation` objects are created. For backward compatibility, `ZoneDelegation`
objects generated from `k8gb.dnsZones` or the `DNS_ZONES` environment variable have `doFinalize` set to `false` by default.
This prevents k8gb from deleting delegated zones that were originally configured as static zones.

## See Also

- [Partial Deployment Troubleshooting](partial-deployment.md): what happens when a `Gslb` resource is missing from one or more delegated clusters.
