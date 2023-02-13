#!/bin/bash
subscriptionName="MBCPANLS02"

##Cluster 1
cluster1Name="cscm2lakslabK8GB001"
spoke1ResourceGroupName="RGCM2LILABK8GB002"

##Cluster 2
cluster2Name="cscm2lakslabK8GB001"
spoke2ResourceGroupName="RGCM2LILABK8GB002"

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