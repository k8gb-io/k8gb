# General deployment with Cloudflare integration

In this guide, we will demonstrate how to configure k8gb to integrate with
Cloudflare for automated zone delegation configuration.

## Initial setup

As a prerequisite, we will need two Kubernetes clusters where you want to deploy
k8gb and enable global load balancing between them.

You can reuse local clusters from the [Infoblox tutorial](../docs/deploy_infoblox.html),
the EKS-based setup from [Route53 tutorial](../docs/deploy_route53.md)
or any Kubernetes deployment method that is convenient to you.

The specific Kubernetes deployment method is not essential for the focus of this documentation guide.

For simplicity, we will assume that clusters have simple 'eu' and 'us' geotags.

## Deploy k8gb with Cloudflare integration enabled

Use `helm` to deploy a stable release from the Helm repo.

```sh
helm repo add k8gb https://www.k8gb.io
```

Example `values.yaml` configuration files can be found [here](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/cloudflare)

Remember to change the zone-related values to point configuration to your own DNS zone.

```yaml
k8gb:
  dnsZone: "cloudflare-test.k8gb.io"
  # -- main zone which would contain gslb zone to delegate
  edgeDNSZone: "k8gb.io" # main zone which would contain gslb zone to delegate
```

### Cloudflare-specific configuration

Let's look closer at the Cloudflare section of the configuration examples.

```yaml
cloudflare:
  # -- Enable Cloudflare provider
  enabled: true
  # -- Cloudflare Zone ID
  zoneID: cdebf92e613133e4bb176a14a9c5b730
  # -- Configure how many DNS records to fetch per request
  # see https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/cloudflare.md#throttling
  dnsRecordsPerPage: 5000
```

Follow
https://developers.cloudflare.com/fundamentals/setup/find-account-and-zone-ids/
to find your `zoneID`

### Install the k8gb helm chart in each cluster

In `eu` cluster, execute
```sh
helm -n k8gb upgrade -i k8gb k8gb/k8gb --create-namespace -f ./docs/examples/cloudflare/k8gb-cluster-cloudflare-eu.yaml
```

In `us` cluster, execute
```sh
helm -n k8gb upgrade -i k8gb k8gb/k8gb --create-namespace -f ./docs/examples/cloudflare/k8gb-cluster-cloudflare-us.yaml
```

### Create a Cloudflare secret in each cluster

```sh
kubectl -n k8gb create secret generic cloudflare --from-literal=token=<api-secret>
```

Note: you can create Cloudflare API tokens at https://dash.cloudflare.com/profile/api-tokens

### Create test Ingress and Gslb resources

Now we can test the setup with a pretty standard Gslb resource configuration.

```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: test-gslb-failover
  namespace: test-gslb
spec:
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    matchLabels:
      app: test-gslb-failover
  strategy:
    dnsTtlSeconds: 60 # Minimum for non-Enterprise Cloudflare https://developers.cloudflare.com/dns/manage-dns-records/reference/ttl/
    primaryGeoTag: eu
    splitBrainThresholdSeconds: 300
    type: failover
```

The only unusual thing here is `spec.strategy.dnsTtlSeconds` that should be of a
minimum 60-second value in case you are operating a non-Enterprise Cloudflare
subscription. The lower values will be rejected by Cloudflare API.

Apply the Gslb and Ingress resources to each cluster.

```sh
kubectl apply -f ./docs/examples/cloudflare/test-gslb-failover.yaml
```

## Check Zone Delegation configuration

As a result of the setup, you should observe DNSEndpoint automatically created,
similar to the one below:

```yaml
$ kubectl -n k8gb get dnsendpoints.externaldns.k8s.io k8gb-ns-extdns -o yaml
apiVersion: externaldns.k8s.io/v1alpha1
kind: DNSEndpoint
metadata:
  annotations:
    k8gb.absa.oss/dnstype: extdns
  creationTimestamp: "2023-11-12T19:55:20Z"
  generation: 3
  name: k8gb-ns-extdns
  namespace: k8gb
  resourceVersion: "5851"
  uid: 5d240eb8-1c19-48c3-bf69-508545f52ea4
spec:
  endpoints:
  - dnsName: cloudflare-test.k8gb.io
    recordTTL: 60
    recordType: NS
    targets:
    - gslb-ns-eu-cloudflare-test.k8gb.io
    - gslb-ns-us-cloudflare-test.k8gb.io
  - dnsName: gslb-ns-us-cloudflare-test.k8gb.io
    recordTTL: 60
    recordType: A
    targets:
    - 172.26.0.8
    - 172.26.0.9
```

On the Cloudflare dashboard side, you should observe that NS and glue A records are
automatically created:

![Cloudflare dashboard with Zone Delegation records](/docs/images/k8gb-cloudflare.png)

## Troubleshooting

If something is not working as expected with the integration, check the logs of
the externalDNS pod that is responsible for the creation of the DNS records
with Cloudflare API.

```sh
kubectl -n k8gb logs -f deploy/external-dns
```
