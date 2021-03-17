# AWS based deployment with Route53 integration

Here we provide an example of k8gb deployment in AWS context with Route53 as edgeDNS provider.

## Reference Setup

Two EKS clusters in `eu-west-1` and `us-east-1`.

Terraform code for cluster reference setup can be found [here](https://github.com/AbsaOSS/k8gb/tree/master/docs/examples/route53)

Feel free to reuse this code fully or partially and adapt for your existing scenario
things like IRSA(IAM Roles for Service Accounts).

## Install Ingress Controller

```sh
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.40.2/deploy/static/provider/aws/deploy.yaml
```

## Deploy k8gb

Example helm configuration files can be found [here](https://github.com/AbsaOSS/k8gb/tree/master/docs/examples/route53/k8gb)

Modify them to reflect your `dnsZone`, `edgeDNSZone`, valid `hostedZoneID` and `irsaRole` ARN.

Clone k8gb repository and use `helm` with custom values

```sh
git clone https://github.com/AbsaOSS/k8gb.git
cd k8gb

helm repo add k8gb https://www.k8gb.io
helm repo update

kubectl create ns k8gb

#switch kubectl context to eu-west-1
helm -n k8gb upgrade -i k8gb k8gb/k8gb -f ./docs/examples/route53/k8gb/k8gb-cluster-eu-west-1.yaml

#switch kubectl context to us-east-1
helm -n k8gb upgrade -i k8gb k8gb/k8gb -f ./docs/examples/route53/k8gb/k8gb-cluster-us-east-1.yaml
```

## Test

*Note*: here and for all occurrences below whenever we speak about application to *each*
cluster, we assume that you switch kubctl context and apply the same command to all clusters.

* Deploy test application to *each* cluster.

```sh
make deploy-test-apps
```

* Modify sample [Gslb CR](https://github.com/AbsaOSS/k8gb/tree/master/docs/examples/route53/k8gb/gslb-failover.yaml) to reflect desired `.spec.ingress.rules[0].host` FQDN

* Apply Gslb CR to *each* cluster

```sh
kubectl apply -f ./docs/examples/route53/k8gb-failover.yaml
```

* Check Gslb status.

```sh
kubectl -n test-gslb get gslb test-gslb-failover -o yaml
```

* Check route53 entries.

```sh
aws route53 list-resource-record-sets --hosted-zone-id $YOUR_HOSTED_ZONE_ID
```

You should see that `gslb-ns-$dnsZone-$geotag` NS and glue A records were created to
automatically configure DNS zone delegation.

* Check test application availability.

```sh
curl -s failover.test.k8gb.io| grep message
  "message": "eu-west-1",
```

 Replace `failover.test.k8gb.io` with the domain you specified in Gslb spec.

Notice that traffic was routed to `eu-west-1`.

* Emulate the failure in `eu-west-1`

```sh
kubectl -n test-gslb scale deploy frontend-podinfo --replicas=0
```

* Observe Gslb status change.

```sh
k -n test-gslb get gslb test-gslb-failover -o yaml | grep status -A6
status:
  geoTag: us-east-1
  healthyRecords:
    failover.test.k8gb.io:
    - 35.168.91.100
  serviceHealth:
    failover.test.k8gb.io: Healthy
```

IP in healthyRecords should change to the IP address of NLB in `us-east-1`

* Check failover to `us-east-1`

```sh
curl -s failover.test.k8gb.io| grep message
  "message": "us-east-1",
```

Notice that traffic is properly failed over to `us-east-1`

* Experiment

Now you can scale `eu-west-1` back and observe that traffic is routed back to the primary cluster.

In addition, you can test `roundRobin` load balancing strategy, which is spreading the traffic
over the clusters in active-active mode.
