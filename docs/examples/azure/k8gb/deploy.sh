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

kubectl apply -f ../external-dns/external-dns-krb5conf.yaml -n k8gb

