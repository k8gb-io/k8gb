# k8gb strategy options
## roundRobin
Returns both cluster endpoints in round-robin manner.

## failover
Pinned to a specified primary cluster until workload on that cluster has no available Pods, upon which the next available cluster's Ingress node IPs will be resolved. When Pods are again available on the primary cluster, the primary cluster will once again be the only eligible cluster for which cluster Ingress node IPs will be resolved

## geoip
Similar to `failover` mode, but returns "closest" cluster to the client initiating request. This requires a specially crafted GeoIP database (see [this](https://github.com/k8gb-io/coredns-crd-plugin/tree/main/terratest/geogen) for example) and DNS resolver to support EDNS0 extension (CLIENT-SUBNET in particular). If the client subnet is not in GeoIP Database, all available endpoints are returned
