# Partial Deployment Troubleshooting

This page provides operational guidance for diagnosing and resolving situations where a `Gslb` resource is missing from one or more clusters in a multi-cluster k8gb setup. For a description of how k8gb behaves in this scenario, see [use case 4 in the Getting Started guide](intro.md#4-partial-deployment---gslb-resource-missing-from-a-cluster).

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
