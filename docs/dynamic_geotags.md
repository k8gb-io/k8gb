# Dynamic GeoTags

_**Note:**
Dynamic GeoTags work with Infoblox out of the box. ExternalDNS providers need
`k8gb.extDNSNsMerge` plus a second ExternalDNS instance that shares `txtOwnerId`
across clusters; see below._

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

If the `extGslbClustersGeoTags` value is empty, k8gb extracts external GeoTags at runtime when Infoblox is enabled, or when ExternalDNS NS-merge is enabled (`k8gb.extDNSNsMerge`).

- This reduces manual configuration and operational overhead.
- You can add or remove clusters without having to update and restart all existing k8gb instances.
- It’s especially useful for dynamic, cloud-native, or large-scale multi-cluster environments.

### Example (`values.yaml`):
```yaml
k8gb:
  clusterGeoTag: "eu"
  extGslbClustersGeoTags: ""  # leave empty to enable dynamic discovery
  extDNSNsMerge: true         # required for ExternalDNS-based providers
```

## ExternalDNS Dynamic GeoTags

Investigation tracked in [#2464](https://github.com/k8gb-io/k8gb/issues/2464).

### How static GeoTags work with ExternalDNS

Each cluster embeds ExternalDNS with a unique `txtOwnerId` (for example `k8gb-eu`, `k8gb-us`). Zone delegation is published as a `DNSEndpoint` that contains:

1. One shared NS RRset for the load-balanced zone (for example `cloud.example.com`), listing every known `gslb-ns-<geotag>-...` target.
2. One glue A record for *this* cluster’s NS name only.

Because every cluster is configured with the full static peer list, the cluster that first owns the shared NS RRset already publishes the complete target set. Peer clusters only need to own their own glue A records. Content converges even though ownership of the NS RRset is single-writer.

### How Infoblox Dynamic GeoTags work

The in-tree Infoblox provider does a WAPI read-modify-write on the delegated zone: it keeps remote NS targets, replaces its own glue, and writes the union back. There is no per-cluster TXT ownership model, so every cluster can safely extend the NS set. Dynamic discovery then reads that merged NS set from DNS.

### Why a unique-owner ExternalDNS deadlocks

If `extGslbClustersGeoTags` is empty and every ExternalDNS instance has its own `txtOwnerId`:

1. Cluster `eu` starts alone and publishes NS targets `[gslb-ns-eu-...]`. Its ExternalDNS instance owns that NS RRset.
2. Cluster `us` joins and can publish its glue A (`gslb-ns-us-...`) under its own owner id.
3. Cluster `us` cannot append `gslb-ns-us-...` to the shared NS RRset owned by `k8gb-eu`.
4. Cluster `eu`’s next dynamic dig still only sees `[eu]`, so it never learns about `us`.

Related upstream attempt: [kubernetes-sigs/external-dns#5619](https://github.com/kubernetes-sigs/external-dns/pull/5619) (`merge-strategy: merge-targets`) — closed without merge (design hold on deletes/updates / multi-ownership).

### Workaround: shared-owner NS ExternalDNS

k8gb can publish the NS RRset on a **separate** `DNSEndpoint` labeled `k8gb.io/dnstype=extdns-ns`. A second ExternalDNS instance, with the **same** `txtOwnerId` on every cluster, watches only that label and applies `union(live parent NS, this cluster)`. Glue A records stay on the unique-owner instance (`k8gb.io/dnstype=extdns`), so one cluster cannot delete another cluster’s glue.

Enable it on every cluster:

```yaml
k8gb:
  clusterGeoTag: "eu"
  extGslbClustersGeoTags: ""
  extDNSNsMerge: true
extdns:
  enabled: true
  txtOwnerId: "k8gb-eu"          # unique per cluster (glue A)
  txtPrefix: "k8gb-eu-"
  labelFilter: "k8gb.io/dnstype=extdns"
```

Deploy a second ExternalDNS (same provider credentials, same `domainFilters`) on **every** cluster:

```yaml
# extra ExternalDNS — values must be identical across clusters except image/resources
txtOwnerId: "k8gb-ns"            # shared; must NOT include the cluster geotag
txtPrefix: "k8gb-ns-"
labelFilter: "k8gb.io/dnstype=extdns-ns"
managedRecordTypes: ["NS"]
policy: sync
sources: ["crd"]
```

k8gb then:

- Discovers peers from parent NS (ANSWER and AUTHORITY sections).
- Always unions the local geotag so a joiner can extend the set.
- Writes the union to the `extdns-ns` endpoint.
- Leaves glue A on the unique-owner endpoint.

On non-last finalize, k8gb removes this cluster from the shared NS targets but does **not** delete the NS endpoint (so the departing ExternalDNS does not wipe the whole RRset). The last cluster deletes both endpoints.

#### Migration from static ExternalDNS geotags

The existing unique-owner instance already owns the NS RRset. Enabling `extDNSNsMerge` moves NS off that endpoint, so unique-owner ExternalDNS deletes the old NS (and its ownership TXT). The shared-owner instance then creates NS. Expect a short gap (one ExternalDNS interval, often ~20s). Schedule this during a maintenance window, or keep static `extGslbClustersGeoTags` until you can accept that gap.

If you do not run the second ExternalDNS, leave `extDNSNsMerge` false and keep `extGslbClustersGeoTags` explicit.

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
 - For ExternalDNS-based edge DNS, enable `k8gb.extDNSNsMerge` and a shared-owner NS ExternalDNS, or set `extGslbClustersGeoTags` explicitly.
