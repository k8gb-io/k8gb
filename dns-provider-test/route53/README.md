# Test Amazon Route 53 integration from a local cluster

This is a guide how to test the Amazon Route 53 integration of K8GB

## Pre-requisites

### Opentofu

Install opentofu to provision the infrastructure
```
brew install opentofu
```

### AWS account

First you will need an AWS account, if you don't have one already you can get started with a [free tier account](https://aws.amazon.com/free).
Afterwards, create access keys in the portal `IAM Dashboard -> My security credentials -> Access keys` and use then to login using `aws configure`. This will make the credentials available to terraform.

## Run the script

```
./test.sh
```

The script will:
* create a DNS zone and a IAM user in AWS.
* create a local cluster
* deploy k8gb in the local cluster
* verify that external DNS operator creates the correct records in AWS
