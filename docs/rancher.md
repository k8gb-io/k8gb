# Integration with Rancher Fleet

The K8gb has been modified to be easily deployed using [Rancher Fleet](https://fleet.rancher.io/). All you need to supply is a
[fleet.yaml](https://fleet.rancher.io/ref-fleet-yaml) file and possibly expose the labels on your cluster.

## Deploy k8gb to Target clusters

The following shows the Rancher/Fleet application that will be installed on the target cluster. Cluster-specific values are
taken from labels set on each Fleet cluster:

| Cluster label | Helm value | Meaning |
| --- | --- | --- |
| `k8gb-dnsZone` | `k8gb.dnsZones[].loadBalancedZone` | GSLB / load-balanced zone for that cluster |
| `k8gb-clusterGeoTag` | `k8gb.clusterGeoTag` | Geo tag of this cluster |
| `k8gb-extGslbClustersGeoTags` | `k8gb.extGslbClustersGeoTags` | Comma-separated geo tags of peer clusters |

`parentZone` is usually shared across clusters (hardcoded below). Use `edgeDNSServers` for resolvers that can answer
the parent zone — that is the only supported chart key for this purpose.

```yaml
# fleet.yaml
defaultNamespace: k8gb
kustomize:
  dir: overlays/kustomization
labels:
  bundle: k8gb
helm:
  repo: https://www.k8gb.io
  chart: k8gb
  version: v0.20.0
  releaseName: k8gb
  values:
    k8gb:
      dnsZones:
        - parentZone: example.com
          loadBalancedZone: global.fleet.clusterLabels.k8gb-dnsZone
      edgeDNSServers:
        - "1.2.3.4"
        - "5.6.7.8"
      clusterGeoTag: global.fleet.clusterLabels.k8gb-clusterGeoTag
      extGslbClustersGeoTags: global.fleet.clusterLabels.k8gb-extGslbClustersGeoTags
      log:
        format: simple
```
