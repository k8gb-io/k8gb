#!/bin/bash
subscriptionName="MVP Sponsorship"

##Cluster 1
cluster1Name="aks1"
spoke1ResourceGroupName="k8gb-win-spoke1"

##Cluster 2
cluster2Name="aks2"
spoke2ResourceGroupName="k8gb-win-spoke2"

az account set --subscription "$subscriptionName"

# Deploy to Cluster 1
az aks command invoke \
  --resource-group $spoke1ResourceGroupName \
  --name $cluster1Name \
  --command 'kubectl apply -n k8gb -f external-dns-krb5conf.yaml -f rfc2136-tsig-secret.yaml && helm upgrade -i -n k8gb --create-namespace k8gb ../../../chart/k8gb -f aks1-helm-values.yaml' \
  --file k8gb/aks1-helm-values.yaml \
  --file external-dns/external-dns-krb5conf.yaml \
  --file external-dns/rfc2136-tsig-secret.yaml \
  &

# Deploy to Cluster 2
az aks command invoke \
  --resource-group $spoke2ResourceGroupName \
  --name $cluster2Name \
  --command 'kubectl apply -n k8gb -f external-dns-krb5conf.yaml -f rfc2136-tsig-secret.yaml && helm upgrade -i -n k8gb --create-namespace k8gb ../../../chart/k8gb -f aks2-helm-values.yaml' \
  --file k8gb/aks2-helm-values.yaml \
  --file external-dns/external-dns-krb5conf.yaml \
  --file external-dns/rfc2136-tsig-secret.yaml \
  &

wait
