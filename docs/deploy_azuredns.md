<h1 align="center" style="margin-top: 0;">Using Azure Public DNS provider</h1>

This document outlines how to configure k8gb to use the Azure Public DNS provider. Azure Private DNS is not supported as it does not support NS records at this time. For private DNS scenarios in Azure, please refer to the [Windows DNS](deploy_windowsdns.md) documentation and consider implementing it using VM-based DNS services such as Windows DNS or BIND.

## Sample solution

In this sample solution we will deploy two private AKS clusters in different regions. A workload will be deployed to both clusters and exposed to the internet with the help of k8gb and Azure Public DNS.

## Reference Setup

The reference setup includes two private AKS clusters that can be deployed on two different regions for load balancing or to provide a failover solution.

Configurable resources:

* Resource groups
* VNet and subnets
* Managed Identity
* Clusters

## Run the sample

* To run the provided sample, please use the provided Makefile [here](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/azure/).
    * Deploys all the required infrastructure and configurations
    * Before executing, please fill all the local variables in the scripts with the correct naming for the resources in order to avoid having problems with your Azure policies
    * Scripts will use Az CLI, please ensure that it is installed and logged when trying to execute the command
        * [Microsoft Learn](https://learn.microsoft.com/en-us/cli/azure/install-azure-cli "Install Az CLI")

### Deploy infrastructure

This action will create resource groups, vnets and private AKS clusters to run all required workloads

```sh
make deploy-infra
```

### Setup clusters

Install required Ingress controller in both clusters in order to deploy K8GB and demo application

```sh
make setup-clusters
```

### Install K8gb

This action will install K8gb in both clusters using the provided [sample](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/azure/) values.yaml for each cluster. Please ensure that the are correctly updated before execution

```sh
make deploy-k8gb
```

### Deploy the credentials for Azure DNS

In this example, we can use a registered app in Microsoft Entra ID and it's corresponding Client ID / Client Secret to authenticate with the Azure DNS zone. Deploy on both clusters a secret called `external-dns-secret-azure` following [External DNS's documentation](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/azure.md#configuration-file).

### Install demo app

Deploys the sample Podinfo workload with failover GLSB configured using annotations in the Ingress resource [samples](https://github.com/k8gb-io/k8gb/tree/master/docs/examples/azure/demo/).
Ensure that the hosts on the samples are correctly updated before execution

```sh
make deploy-demo
```

### Destroy lab

* Destroys the lab environment created for this sample

```sh
make destroy-infra
```
