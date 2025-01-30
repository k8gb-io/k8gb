# AWS based deployment with NS1 integration

Here we provide an example of k8gb deployment in AWS context with NS1 as edgeDNS provider

## Reference setup

Two EKS clusters in `eu-west-1` and `us-east-1`.

The EKS setup is identical to [Route53 tutorial](.deploy_route53.md)

Terraform code for cluster reference setup can be found [here](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/route53)

## Deploy k8gb

Use `helm` to deploy stable release from Helm repo

```sh
helm repo add k8gb https://www.k8gb.io
```

Example `values.yaml` configuration files can be found [here](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/ns1)

In `eu-west-1` cluster execute
```sh
helm -n k8gb upgrade -i k8gb k8gb/k8gb --create-namespace -f examples/ns1/k8gb-cluster-ns1-eu-west-1.yaml
```

In `us-east-1` cluster execute
```sh
helm -n k8gb upgrade -i k8gb k8gb/k8gb --create-namespace -f examples/ns1/k8gb-cluster-ns1-us-east-1.yaml
```

Create NS1 secret in each cluster

```sh
export NS1_APIKEY=<ns1-api-key>
make ns1-secret
```
