#!/bin/bash
subscriptionName="MVP Sponsorship"

##Cluster 1
cluster1Name="aks1"
spoke1ResourceGroupName="k8gb-az-spoke1"

##Cluster 2
cluster2Name="aks2"
spoke2ResourceGroupName="k8gb-az-spoke2"

az account set --subscription "$subscriptionName"

# get credentials
az aks get-credentials -g $spoke1ResourceGroupName -n $cluster1Name --admin --overwrite-existing
az aks get-credentials -g $spoke2ResourceGroupName -n $cluster2Name --admin --overwrite-existing

# Deploy nginx
helm repo add nginx-stable https://kubernetes.github.io/ingress-nginx 
helm repo update 
helm upgrade -i -n nginx-ingress --create-namespace nginx-ingress nginx-stable/ingress-nginx --version 4.0.15 -f infra/nginx-ingress-values.yaml --kube-context aks1-admin &
helm upgrade -i -n nginx-ingress --create-namespace nginx-ingress nginx-stable/ingress-nginx --version 4.0.15 -f infra/nginx-ingress-values.yaml --kube-context aks2-admin &
wait
