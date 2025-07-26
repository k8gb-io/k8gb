# Dynamic GeoTags

_**Note:**
Dynamic GeoTags currently work only with parent DNS servers running Infoblox (WAPI-based integration).
Other parent DNS solutions, such as those managed by ExternalDNS, are not supported for dynamic GeoTag discovery.
If your parent zone is managed by a different DNS provider (e.g., Route53, CoreDNS, or ExternalDNS), you must use the 
static configuration method and set the `extGslbClustersGeoTags` value explicitly in your `values.yaml`._


From v0.15.0, k8gb makes it easier to configure and manage cluster GeoTags.

### What is a GeoTag?
A GeoTag is a short identifier (for example: `eu`, `us`, `za`) that uniquely marks each k8gb cluster’s location or role. GeoTags are essential for k8gb’s global DNS-based routing and failover logic.

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

## What’s new with Dynamic GeoTags?

Previously, any change in the list of external clusters (extGslbClustersGeoTags) required you to update and restart all k8gb pods across all clusters, which was inconvenient and error prone especially as the number of clusters grew.

**Dynamic GeoTags** allow k8gb to discover external GeoTags directly from DNS (from NS records on the parent zone), without the need to keep all values manually in sync.

If the `extGslbClustersGeoTags` value is empty, k8gb will attempt to extract external GeoTags dynamically at runtime.

- This reduces manual configuration and operational overhead.
- You can add or remove clusters without having to update and restart all existing k8gb instances.
- It’s especially useful for dynamic, cloud-native, or large-scale multi-cluster environments.

### Example (`values.yaml`):
```yaml
k8gb:
  clusterGeoTag: "eu"
  extGslbClustersGeoTags: ""  # leave empty to enable dynamic discovery
```

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
