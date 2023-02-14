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

### Deploy Ingress Controller
kubectl apply -f ingress/namespace.yaml
helm repo add --force-update nginx-stable https://kubernetes.github.io/ingress-nginx
helm -n nginx-ingress upgrade -i nginx-ingress nginx-stable/ingress-nginx \
                --version 4.0.15 -f nginx-ingress-values.yaml


#####################
# Deploy to Cluster 2
#####################

# Get credentials
az account set --subscription $subscriptionName
az aks get-credentials --resource-group $spoke2ResourceGroupName --name $cluster2Name

### Deploy Ingress Controller
kubectl apply -f ingress/namespace.yaml
helm repo add --force-update nginx-stable https://kubernetes.github.io/ingress-nginx
helm -n nginx-ingress upgrade -i nginx-ingress nginx-stable/ingress-nginx \
                --version 4.0.15 -f nginx-ingress-values.yaml