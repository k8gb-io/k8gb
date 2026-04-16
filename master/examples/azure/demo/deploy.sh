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

kubectl create ns podinfo --context aks1-admin &
kubectl create ns podinfo --context aks2-admin &
wait 
helm repo add podinfo https://stefanprodan.github.io/podinfo
helm repo update 
helm upgrade --install podinfo podinfo/podinfo -n podinfo --create-namespace -f demo/aks1-podinfo-values.yaml --kube-context aks1-admin &
helm upgrade --install podinfo podinfo/podinfo -n podinfo --create-namespace -f demo/aks2-podinfo-values.yaml --kube-context aks2-admin &
wait
kubectl apply -f demo/podinfo-ingress.yaml --context aks1-admin &
kubectl apply -f demo/podinfo-ingress.yaml --context aks2-admin &
wait
