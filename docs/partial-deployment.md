# Partial Deployment

## What happens when a Gslb resource is missing from a cluster

k8gb participation is opt-in per cluster per hostname. Consider a two-cluster setup where Cluster **X** has a `Gslb` resource for `app.cloud.example.com` and Cluster **Y** does not:

1. Cluster **Y**'s k8gb controller has no knowledge of `app.cloud.example.com` and publishes no `localtargets-app.cloud.example.com` A records to its local CoreDNS.
2. NS zone delegation is per-cluster, not per-hostname. Cluster **Y**'s zone delegation registered with the edge DNS is unaffected. However, when Cluster **X** queries Cluster **Y**'s CoreDNS for `localtargets-app.cloud.example.com`, it receives an empty response (NOERROR, no A records) because the records were never published.
3. Cluster **X** collects `localtargets-*` records from all known clusters; for Cluster **Y** it receives an empty result, so only Cluster **X**'s own Ingress node IPs are included in the final answer set.

Clients receive **consistent** DNS responses directed to Cluster **X**, but Cluster **Y** is never included in the GSLB pool for this hostname, regardless of its Pod health.

## Troubleshooting

This section provides guidance for diagnosing situations where a `Gslb` resource is missing from one or more clusters in a multi-cluster k8gb setup.

## Common Causes

| Cause | Description |
|-------|-------------|
| **Incomplete rollout** | `Gslb` resources were deployed to some clusters but not all. |
| **Namespace mismatch** | The `Gslb` resource exists but in a different namespace than k8gb is watching. |
| **Hostname typo** | The `spec.ingress.rules[].host` value in the `Gslb` resource does not exactly match the Ingress host. |

## How to Detect

**1. Verify consistent Gslb coverage across clusters**

Run the following on each cluster and confirm the same hostname appears everywhere you expect multi-cluster load balancing:

```bash
kubectl get gslb -A
```

If a hostname is present in Cluster X but absent in Cluster Y, Cluster Y is excluded from the GSLB pool for that hostname.

**2. Inspect healthyRecords on the Gslb object**

On a cluster that *does* have the `Gslb` resource, check whether remote clusters are contributing endpoints:

```bash
kubectl get gslb <name> -n <namespace> -o jsonpath='{.status.healthyRecords}'
```

If only one cluster's Ingress IPs appear in `status.healthyRecords`, the other clusters are either unhealthy or missing the `Gslb` resource.

**3. Query the remote cluster's CoreDNS directly**

Confirm that a cluster has published `localtargets-*` records by querying its CoreDNS directly. If no records are returned, the `Gslb` resource is absent or not yet reconciled on that cluster:

```bash
# Replace <coredns-ip> with the external IP of the target cluster's CoreDNS service
dig localtargets-app.cloud.example.com @<coredns-ip>
```

An empty answer section (NOERROR, no A records) confirms that the cluster has no `localtargets-*` records for the hostname.
