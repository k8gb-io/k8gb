# Test Cloud DNS integration from a local cluster

This is a guide how to test the Cloud DNS integration of K8GB

## GCP infrastructure

### GCP account

First you will need a GCP account and a project, if you don't have one already you can get started with a [free account](https://cloud.google.com/free?hl=en).
Then assign it a billing account and enable the `Cloud DNS API`.

Afterwards, make your credentials availabe to terraform by running `gcloud init` and `gcloud auth application-default login`.

### DNS Zone and service principal

The next step is to create a DNS zone and a service account that allows K8GB to modify records in the zone.
You can use the terraform code provided in the `terraform` folder to get started. You will be prompted with a project id and the name of the DNS zone. The name needs to be unique in GCP, but you don't need to own the zone for the purpose of this guide:
```
$ cd terraform
$ terraform init
$ terraform apply
var.dns_zone_name
  Name of the DNS zone (must end with a dot)

  Enter a value: k8gb.io.
```

### Create local clusters

We have everything we need from GCP, we can now create a local cluster.
Navigate to the home of the k8gb repo run the following command. It will create the clusters `k3d-test-gslb1` and `k3d-test-gslb2`, and install k8gb from the branch you are on:
```
K8GB_LOCAL_VERSION=test FULL_LOCAL_SETUP_WITH_APPS=false make deploy-full-local-setup
```

### Connect K8GB to GCP

At this moment K8GB is using the upstream DNS server running on the local cluster `k3d-edgedns`. We want to point it to the DNS infrastructure we created in GCP.

To do that we will need to create a secret on both clusters, on the `k8gb` namespace with the name `external-dns-secret-gcp`. The format of the secret is documented in the [external dns docs](https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/gke.md#static-credentials). If you are in your `terraform` folder you can create it using:
```
gcloud iam service-accounts keys create credentials.json --iam-account "$(terraform output --raw service_account_email)"
```
Now apply the secret to both of the clusters:
```
kubectl create secret generic external-dns-secret-gcp -n k8gb --from-file credentials.json --context k3d-test-gslb1
kubectl create secret generic external-dns-secret-gcp -n k8gb --from-file credentials.json --context k3d-test-gslb2
```

### Create application

Finally, we can create a GSLB resouce that will trigger a reconciliation loop of the controller and configure DNS name delegation on GCP.
To do that we will need to configure the DNS zone we create on K8GB:
```
# replace with your zone
EDGE_DNS_ZONE="k8gb.io"
```
```
DNS_ZONE="cloud.${EDGE_DNS_ZONE}"
EDGE_DNS_SERVER=$(gcloud dns record-sets list --zone="k8gb" --name="${EDGE_DNS_ZONE}." --type="NS" --format="value(rrdatas[0])" | sed 's/\.$//')
```

```
cd ../helm
helm package -u . > /dev/null && helm template k8gb k8gb-v0.1.0.tgz -n k8gb -f values.yaml -f values-eu.yaml --set "k8gb.k8gb.dnsZone=$DNS_ZONE" --set "k8gb.k8gb.edgeDNSZone=$EDGE_DNS_ZONE" --set "k8gb.k8gb.edgeDNSServers[0]=$EDGE_DNS_SERVER" > manifests-eu.yaml
helm package -u . > /dev/null && helm template k8gb k8gb-v0.1.0.tgz -n k8gb -f values.yaml -f values-us.yaml --set "k8gb.k8gb.dnsZone=$DNS_ZONE" --set "k8gb.k8gb.edgeDNSZone=$EDGE_DNS_ZONE" --set "k8gb.k8gb.edgeDNSServers[0]=$EDGE_DNS_SERVER" > manifests-us.yaml

kubectl apply -f manifests-eu.yaml --context k3d-test-gslb1
kubectl apply -f manifests-us.yaml --context k3d-test-gslb2
```

### Verify zone delegation in GCP

And voila, our local clusters are now integrated with GCP. We can quickly verify everything is working.

In GCP we should find the following records (the IP addresses may be different depending on the allocation by docker):
| Name    | Type | Value |
| -------- | ------- |  ------- |
| cloud  | NS    | gslb-ns-eu-cloud.k8gb.io gslb-ns-us-cloud.k8gb.io
| gslb-ns-eu-cloud | A     | 172.18.0.6 172.18.0.7
| gslb-ns-us-cloud    | A    | 172.18.0.10 172.18.0.11
```
gcloud dns record-sets list --zone="k8gb"
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

Unfortunately the A records cannot be queried because they are private IP addresses and Google does not return them in a public DNS zone, but this is enough for testing.
