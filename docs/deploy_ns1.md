# AWS based deployment with NS1 integration

Here we provide an example of k8gb deployment in AWS context with NS1 as edgeDNS provider

## Reference setup

Two EKS clusters in `eu-west-1` and `us-east-1`.

Terraform code for cluster reference setup can be found [here](https://github.com/AbsaOSS/k8gb/tree/master/docs/examples/route53)

The EKS setup is identical to [Route53 tutorial](/docs/deploy_route53.md)

## Deploy k8gb

Example values.yaml override configs can be found [here](https://github.com/AbsaOSS/k8gb/tree/master/docs/examples/ns1/)

You can use `helm` to deploy stable release from helm repo

```sh
helm repo add k8gb https://www.k8gb.io
```

Alternatively, use make target to deploy right from the git repository

```sh
make deploy-gslb-operator VALUES_YAML=./docs/examples/ns1/k8gb-cluster-ns1-eu-west-1.yaml

#switch kubectl context to us-east-1

make deploy-gslb-operator VALUES_YAML=./docs/examples/ns1/k8gb-cluster-ns1-us-east-1.yaml
```
