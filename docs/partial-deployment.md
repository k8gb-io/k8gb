# Partial Deployment Troubleshooting

## What happens when a Gslb resource is missing from a cluster

When `ZoneDelegation` is present in cluster **Y** but a `Gslb` resource for `app.cloud.example.com` is absent:

1. Cluster **Y** remains authoritative for the configured load-balanced zone because `ZoneDelegation` continues to publish its NS and glue records.
2. Cluster **Y** does not publish application or `localtargets-*` records for `app.cloud.example.com`.
3. The edge DNS continues to delegate queries for the load-balanced zone to cluster **Y** because the NS and glue records published by `ZoneDelegation` are unaffected. Queries for `app.cloud.example.com` that reach cluster **Y**'s CoreDNS receive NXDOMAIN or an empty answer because no records for that hostname were ever published.
4. Cluster **X** cannot discover healthy application targets from cluster **Y**, so those targets are not included in responses produced by cluster **X**.

Clients can therefore receive inconsistent DNS responses depending on which authoritative cluster answers the query. A partial `Gslb` deployment should not be treated as a safe way to exclude a cluster.

Application participation is configured per hostname through `Gslb`, while DNS delegation is configured separately per load-balanced zone through `ZoneDelegation`. Operators should deploy the `Gslb` resource consistently across every delegated cluster expected to serve the hostname.

## Common Causes

| Cause | Description |
|-------|-------------|
| **Incomplete rollout** | `Gslb` resources were deployed to some clusters but not all. |
| **Namespace mismatch** | The `Gslb` resource exists but in a different namespace than k8gb is watching. |
| **Hostname typo** | The `spec.ingress.rules[].host` value in the `Gslb` resource does not exactly match the Ingress host. |

## How to Detect

- Verify that the expected `ZoneDelegation` exists in every cluster intended to be authoritative for the load-balanced zone.
- Run `kubectl get gslb -A` in each delegated cluster and confirm that the application hostname is present.
- Query the edge DNS for the configured load-balanced zone, for example: `dig NS cloud.example.com @<edge-dns>`.
- Query `app.cloud.example.com` and its generated `localtargets-*` record directly against each authoritative cluster DNS server returned by the NS lookup. Different answers identify an incomplete cluster configuration.
