# Test Amazon Route 53 integration from a local cluster

This is a guide how to test the Amazon Route 53 integration of K8GB

## AWS infrastructure

### AWS account

First you will need an AWS account, if you don't have one already you can get started with a [free tier account](https://aws.amazon.com/free).
Afterwards, create access keys in the portal `IAM Dashboard -> My security credentials -> Access keys` and use then to login using `aws configure`. This will make the credentials available to terraform.

### DNS Zone and service principal

The next step is to create a DNS zone and a service principal that allows K8GB to modify records in the zone.
You can use the terraform code provided in the `terraform` folder to get started. You will be prompted with the name of the DNS zone. The name needs to be unique in AWS, but you don't need to own the zone for the purpose of this guide:
```
$ cd terraform
$ terraform init
$ terraform apply
var.dns_zone_name
  Name of the DNS zone

  Enter a value: k8gb.io
```

### Create local clusters

We have everything we need from AWS, we can now create a local cluster.
Navigate to the home of the k8gb repo run the following command. It will create the clusters `k3d-test-gslb1` and `k3d-test-gslb2`, and install k8gb from the branch you are on:
```
K8GB_LOCAL_VERSION=test FULL_LOCAL_SETUP_WITH_APPS=false make deploy-full-local-setup
```

The k8gb image must be stored in an environment variable to be used later in the helm chart. Double check if your architecture is correct:
```
VERSION=$(git fetch --force --tags &> /dev/null ; git describe --tags --abbrev=0)
COMMIT_HASH=$(git rev-parse --short HEAD)
IMAGE="$VERSION-$COMMIT_HASH-arm64"
```

### Connect K8GB to AWS

With this configuration K8GB is using the upstream DNS server running on the local cluster `k3d-edgedns`. We want to point it to the DNS infrastructure we created in AWS.

To do that we will need to create a secret on both clusters, on the `k8gb` namespace with the name `external-dns-secret-aws`. The format of the secret is documented in the [external dns docs](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md#static-credentials). You can create it using:
```
SECRET_ACCESS_KEY=$(aws iam create-access-key --user-name "externaldns")
cat <<-EOF > credentials

[default]
aws_access_key_id = $(echo $SECRET_ACCESS_KEY | jq -r '.AccessKey.AccessKeyId')
aws_secret_access_key = $(echo $SECRET_ACCESS_KEY | jq -r '.AccessKey.SecretAccessKey')
EOF
```
Now apply the secret to both of the clusters:
```
kubectl create secret generic external-dns-secret-aws -n k8gb --from-file credentials --context k3d-test-gslb1
kubectl create secret generic external-dns-secret-aws -n k8gb --from-file credentials --context k3d-test-gslb2
```

### Create application

Finally, we can create a GSLB resouce that will trigger a reconciliation loop of the controller and configure DNS name delegation on AWS.
To do that we will need to configure the DNS zone we create on K8GB:
```
# replace with your zone
EDGE_DNS_ZONE="k8gb.io"
```
```
DNS_ZONE="cloud.${EDGE_DNS_ZONE}"
ZONE_ID=$(aws route53 list-hosted-zones --query "HostedZones[?Name == 'k8gb.io.'].Id" --output text)
EDGE_DNS_SERVER=$(aws route53 list-resource-record-sets --hosted-zone-id "$ZONE_ID" --query "ResourceRecordSets[?Type == 'NS'].ResourceRecords[0]" --output text | sed 's/\.$//')
```

```
cd ../helm
helm package -u . > /dev/null && helm template k8gb k8gb-v0.1.0.tgz -n k8gb -f values.yaml -f values-eu.yaml --set "k8gb.k8gb.imageTag=$IMAGE" --set "k8gb.k8gb.dnsZone=$DNS_ZONE" --set "k8gb.k8gb.edgeDNSZone=$EDGE_DNS_ZONE" --set "k8gb.k8gb.edgeDNSServers[0]=$EDGE_DNS_SERVER" --set "k8gb.extdns.domainFilters[0]=$EDGE_DNS_ZONE" > manifests-eu.yaml
helm package -u . > /dev/null && helm template k8gb k8gb-v0.1.0.tgz -n k8gb -f values.yaml -f values-us.yaml --set "k8gb.k8gb.imageTag=$IMAGE" --set "k8gb.k8gb.dnsZone=$DNS_ZONE" --set "k8gb.k8gb.edgeDNSZone=$EDGE_DNS_ZONE" --set "k8gb.k8gb.edgeDNSServers[0]=$EDGE_DNS_SERVER" --set "k8gb.extdns.domainFilters[0]=$EDGE_DNS_ZONE" > manifests-us.yaml

kubectl apply -f manifests-eu.yaml --context k3d-test-gslb1
kubectl apply -f manifests-us.yaml --context k3d-test-gslb2
```

### Verify zone delegation in AWS

And voila, our local clusters are now integrated with AWS Route 53. We can quickly verify everything is working.

In AWS we should find the following records (the IP addresses may be different depending on the allocation by docker):
| Name    | Type | Value |
| -------- | ------- |  ------- |
| cloud  | NS    | gslb-ns-eu-cloud.k8gb.io gslb-ns-us-cloud.k8gb.io
| gslb-ns-eu-cloud | A     | 172.18.0.6 172.18.0.7
| gslb-ns-us-cloud    | A    | 172.18.0.10 172.18.0.11
```
aws route53 list-resource-record-sets --hosted-zone-id "$ZONE_ID" --query "ResourceRecordSets[?Type == 'NS']"
aws route53 list-resource-record-sets --hosted-zone-id "$ZONE_ID" --query "ResourceRecordSets[?Type == 'A']"
```

You can also fetch the records using the following DNS query:
```
dig @${EDGE_DNS_SERVER} cloud.k8gb.io
...
;; AUTHORITY SECTION:
cloud.k8gb.io.		5	IN	NS	gslb-ns-eu-cloud.k8gb.io.
cloud.k8gb.io.		5	IN	NS	gslb-ns-us-cloud.k8gb.io.
...
```

Unfortunately the A records cannot be queried because they are private IP addresses and AWS does not return them in a public DNS zone, but this is enough for testing.
