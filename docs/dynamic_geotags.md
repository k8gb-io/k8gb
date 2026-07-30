# Dynamic GeoTags

_**Note:**
Dynamic GeoTags discovery currently works only with Infoblox. ExternalDNS based providers are not supported
because their TXT ownership model prevents reliable GeoTag discovery across clusters._

### What is a GeoTag?
A GeoTag is a short identifier (for example: `eu`, `us`, `za`) that uniquely marks each k8gb cluster’s location or role. 
GeoTags are essential for k8gb’s global DNS-based routing and failover logic.

### How to configure GeoTags
You configure GeoTags directly in your values.yaml (usually when installing or upgrading k8gb via Helm):

```yaml
k8gb:
  # Unique GeoTag for this cluster. Common values: "eu", "us", "za"
  clusterGeoTag: "eu"

  # Comma-separated list of GeoTags for all other k8gb clusters in your GSLB network.
  # Example: "eu,us,za"
  extGslbClustersGeoTags: "eu,us"
```
- `clusterGeoTag`: This value must be unique for each cluster. It identifies this k8gb instance.
- `extGslbClustersGeoTags`: This is a comma-separated list of all external k8gb clusters’ GeoTags.

Both values are set in values.yaml and are automatically passed to the k8gb Pod as environment variables by the Helm chart.

## Dynamic GeoTags

Previously, any change in the list of external clusters (extGslbClustersGeoTags) required you to update and restart all k8gb pods across all clusters, which was inconvenient and error prone especially as the number of clusters grew.

**Dynamic GeoTags** allow k8gb to discover external GeoTags directly from DNS (from NS records on the parent zone), without the need to keep all values manually in sync.

If the `extGslbClustersGeoTags` value is empty, k8gb will attempt to extract external GeoTags dynamically at runtime **when the Infoblox provider is enabled**.

- This reduces manual configuration and operational overhead.
- You can add or remove clusters without having to update and restart all existing k8gb instances.
- It’s especially useful for dynamic, cloud-native, or large-scale multi-cluster environments.

### Example (`values.yaml`):
```yaml
k8gb:
  clusterGeoTag: "eu"
  extGslbClustersGeoTags: ""  # leave empty to enable dynamic discovery (Infoblox only)
```

## Why ExternalDNS providers cannot use Dynamic GeoTags yet

Investigation tracked in [#2464](https://github.com/k8gb-io/k8gb/issues/2464). Summary:

### How static GeoTags work with ExternalDNS today

Each cluster embeds ExternalDNS with a unique `txtOwnerId` (for example `k8gb-eu`, `k8gb-us`). Zone delegation is published as a `DNSEndpoint` that contains:

1. One shared NS RRset for the load-balanced zone (for example `cloud.example.com`), listing every known `gslb-ns-<geotag>-...` target.
2. One glue A record for *this* cluster’s NS name only.

Because every cluster is configured with the full static peer list, the cluster that first owns the shared NS RRset already publishes the complete target set. Peer clusters only need to own their own glue A records. Content converges even though ownership of the NS RRset is single-writer.

### How Infoblox Dynamic GeoTags work

The in-tree Infoblox provider does a WAPI read-modify-write on the delegated zone: it keeps remote NS targets, replaces its own glue, and writes the union back. There is no per-cluster TXT ownership model, so every cluster can safely extend the NS set. Dynamic discovery then reads that merged NS set from DNS.

### Why enabling Dynamic GeoTags for ExternalDNS deadlocks

If `extGslbClustersGeoTags` is empty:

1. Cluster `eu` starts alone and publishes NS targets `[gslb-ns-eu-...]`. Its ExternalDNS instance owns that NS RRset.
2. Cluster `us` joins and can publish its glue A (`gslb-ns-us-...`) under its own owner id.
3. Cluster `us` cannot append `gslb-ns-us-...` to the shared NS RRset owned by `k8gb-eu`.
4. Cluster `eu`’s next dynamic dig still only sees `[eu]`, so it never learns about `us`.

Result: peer discovery never converges. This is an ExternalDNS multi-owner limitation, not a missing dig in k8gb.

Related upstream attempt: [kubernetes-sigs/external-dns#5619](https://github.com/kubernetes-sigs/external-dns/pull/5619) (`merge-strategy: merge-targets`) — closed without merge (design hold on deletes/updates / multi-ownership). Also noted in [#1967](https://github.com/k8gb-io/k8gb/issues/1967).

### What would unblock ExternalDNS support

Dynamic GeoTags for ExternalDNS-backed providers become feasible only when ExternalDNS (or an equivalent coordination layer) can merge NS targets across distinct `txtOwnerId` values, including correct delete semantics when the last owner removes its target. Until then:

- Use Infoblox for Dynamic GeoTags, or
- Keep `extGslbClustersGeoTags` explicitly set for Route53, Cloudflare, Azure DNS, GCP Cloud DNS, RFC2136, and other ExternalDNS providers.

Secondary readiness items once multi-owner merge exists:

- Each cluster should publish only its own NS target (plus local glue A) and rely on ExternalDNS merge for the union.
- NS discovery must accept records from both ANSWER and AUTHORITY sections (referral responses from cloud DNS parents).
- k8gb should refresh local `DNSEndpoint` objects when parent DNS changes (reverse sync).

## Important considerations
Dynamic GeoTags provide convenience and flexibility, but it’s important to understand their impact on your DNS infrastructure:

> ⚠️ **WARNING:**
Enabling dynamic GeoTags adds two extra DNS queries per reconciliation for each GSLB resource. In most cases, this overhead is negligible, but with a large number of GSLBs and short reconciliation intervals, your DNS server could become overwhelmed.

If you experience high DNS query load or see signs of DNS server saturation, you can mitigate the issue as follows:

**1. Increase the number of DNS servers:**
Add more DNS servers to the list `.Values.k8gb.edgeDNSServers` in your `values.yaml`.
k8gb will choose one DNS server at random (round-robin) for each reconciliation, distributing the load.

**2. Increase the reconciliation interval**
Raise the value of `.Values.k8gb.reconcileRequeueSeconds` in your `values.yaml`. By increasing this interval, you reduce how often k8gb triggers reconciliations, which directly decreases the DNS query rate.
You can also set it to `0`—in this case, reconciliation will only occur in response to changes in GSLB, Ingress, DNSEndpoint, or during initial bootstrap.

**3. Revert to static external GeoTags**
If dynamic GeoTags are not suitable for your environment, you can switch back to using the `.Values.k8gb.extGslbClustersGeoTags` to explicitly define the list of remote cluster tags, disabling dynamic discovery.

**In summary:**
 - Dynamic GeoTags simplify configuration but come with extra DNS queries per GSLB. 
 - For most environments, this is not an issue. 
 - For very large-scale or highly sensitive environments, use the mitigations above to prevent DNS overload.
 - For ExternalDNS-based edge DNS, Dynamic GeoTags are not supported; set `extGslbClustersGeoTags` explicitly.
