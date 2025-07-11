# k8gb strategy options
## RoundRobin
Returns both cluster endpoints in round-robin manner.

## Weight RoundRobin
While roundRobin is fair for all regions, with WeightRoundRobin we can set explicitly how the regions should be loaded with traffic. For example, we can set one region to handle 80% of the traffic, another 20% and a third 0%, so that the last region is practically disabled.

Note:
Please read the [WRR caveats](wrr_caveats.md)for important limitations and recommendations.

## Failover
Pinned to a specified primary cluster until workload on that cluster has no available Pods, upon which the next available cluster's Ingress node IPs will be resolved. When Pods are again available on the primary cluster, the primary cluster will once again be the only eligible cluster for which cluster Ingress node IPs will be resolved

## GeoIP
Similar to `failover` mode, but returns "closest" cluster to the client initiating request. This requires a specially crafted GeoIP database (see [this](https://github.com/k8gb-io/coredns-crd-plugin/tree/main/terratest/geogen) for example) and DNS resolver to support EDNS0 extension (CLIENT-SUBNET in particular). If the client subnet is not in GeoIP Database, all available endpoints are returned
