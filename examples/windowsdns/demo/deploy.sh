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
  --command 'helm repo add podinfo https://stefanprodan.github.io/podinfo && helm repo update && helm upgrade --install podinfo podinfo/podinfo -n podinfo --create-namespace -f aks1-podinfo-values.yaml && kubectl apply -f podinfo-ingress.yaml' \
  --file demo/aks1-podinfo-values.yaml \
  --file demo/podinfo-ingress.yaml \
  &

# Deploy to Cluster 2
az aks command invoke \
  --resource-group $spoke2ResourceGroupName \
  --name $cluster2Name \
  --command 'helm repo add podinfo https://stefanprodan.github.io/podinfo && helm repo update && helm upgrade --install podinfo podinfo/podinfo -n podinfo --create-namespace -f aks2-podinfo-values.yaml && kubectl apply -f podinfo-ingress.yaml' \
  --file demo/aks2-podinfo-values.yaml \
  --file demo/podinfo-ingress.yaml \
  &

wait