# Integration with Rancher Fleet

The K8gb has been modified to be easily deployed using [Rancher Fleet](https://fleet.rancher.io/). All you need to supply is a 
[fleet.yaml](https://fleet.rancher.io/ref-fleet-yaml) file  and possibly expose the labels on your cluster.

## Deploy k8gb to Target clusters
The following shows the rancher application that will be installed on the target cluster.  The values `k8gb-dnsZone`, 
`k8gb-clusterGeoTag`, `k8gb-extGslbClustersGeoTags` will be taken from the labels that are set on the cluster.

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
  version: v0.11.4
  releaseName: k8gb
  values:
    k8gb:
      dnsZones:
        - parentZone: cloud.example.com
          loadBalancedZone: global.fleet.clusterLabels.k8gb-dnsZone
      parentZoneDNSServers:
        - "1.2.3.4"
        - "5.6.7.8"
      clusterGeoTag: global.fleet.clusterLabels.k8gb-clusterGeoTag
      extGslbClustersGeoTags: global.fleet.clusterLabels.k8gb-extGslbClustersGeoTags
      log:
        format: simple
```
