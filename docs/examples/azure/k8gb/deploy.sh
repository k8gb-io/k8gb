#!/bin/bash
subscriptionName=""

##Cluster 1
cluster1Name=""
spoke1ResourceGroupName=""

##Cluster 2
cluster2Name=""
spoke2ResourceGroupName=""

#####################
# Deploy to Cluster 1
#####################

# Get credentials
az account set --subscription $subscriptionName
az aks get-credentials --resource-group $spoke1ResourceGroupName --name $cluster1Name

helm repo add k8gb https://www.k8gb.io
helm repo update

kubectl apply -f ../external-dns/external-dns-krb5conf.yaml -n k8gb

helm -n k8gb upgrade -i k8gb k8gb/k8gb --create-namespace -f cluster-we-helm-values.yaml

#####################
# Deploy to Cluster 2
#####################

# Get credentials
az account set --subscription $subscriptionName
az aks get-credentials --resource-group $spoke2ResourceGroupName --name $cluster2Name

helm repo add k8gb https://www.k8gb.io
helm repo update

kubectl apply -f ../external-dns/external-dns-krb5conf.yaml -n k8gb

helm -n k8gb upgrade -i k8gb k8gb/k8gb --create-namespace -f cluster-ne-helm-values.yaml

