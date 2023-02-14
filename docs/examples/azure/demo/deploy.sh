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

# Create namespace
kubectl apply -f namespace.yaml

# Deploy Pod info demo with local values file
helm repo add podinfo https://stefanprodan.github.io/podinfo

helm upgrade --install frontend --namespace podinfo -f podinfo-values-we.yaml

kubectl apply -f podinfo-ingress-we.yaml

#####################
# Deploy to Cluster 2
#####################

# Get credentials
az account set --subscription $subscriptionName
az aks get-credentials --resource-group $spoke2ResourceGroupName --name $cluster2Name

# Create namespace
kubectl apply -f namespace.yaml

# Deploy Pod info demo with local values file
helm repo add podinfo https://stefanprodan.github.io/podinfo

helm upgrade --install frontend --namespace podinfo -f podinfo-values-ne.yaml

kubectl apply -f podinfo-ingress-ne.yaml