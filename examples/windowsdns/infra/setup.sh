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
  --command 'helm repo add nginx-stable https://kubernetes.github.io/ingress-nginx && helm repo update && helm upgrade -i -n nginx-ingress --create-namespace nginx-ingress nginx-stable/ingress-nginx --version 4.0.15 -f nginx-ingress-values.yaml' \
  --file infra/nginx-ingress-values.yaml \
  &

# Deploy to Cluster 2
az aks command invoke \
  --resource-group $spoke2ResourceGroupName \
  --name $cluster2Name \
  --command 'helm repo add nginx-stable https://kubernetes.github.io/ingress-nginx && helm repo update && helm upgrade -i -n nginx-ingress --create-namespace nginx-ingress nginx-stable/ingress-nginx --version 4.0.15 -f nginx-ingress-values.yaml' \
  --file infra/nginx-ingress-values.yaml \
  &

wait
