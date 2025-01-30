#!/bin/bash

#generic configs
subscriptionName="MVP Sponsorship"
dnsZone="k8gb-kubeconeu2023.com"

#Spoke1 configs
cluster1Name="aks1"
spoke1ResourceGroupName="k8gb-az-spoke1"
spoke1vNetName="spoke1-vnet"
spoke1vNetRange="10.11.0.0/16"
spoke1VnetSubnetName="default"
spoke1SubnetRange="10.11.0.0/20"
spoke1Location="uksouth"

#Spoke2 configs
cluster2Name="aks2"
spoke2ResourceGroupName="k8gb-az-spoke2"
spoke2vNetName="spoke2-vnet"
spoke2vNetRange="10.12.0.0/16"
spoke2VnetSubnetName="default"
spoke2SubnetRange="10.12.0.0/20"
spoke2Location="francecentral"

az account set --subscription "$subscriptionName"

#Create resource groups to manage azure resources
az group create --resource-group $spoke1ResourceGroupName --location $spoke1Location &
az group create --resource-group $spoke2ResourceGroupName --location $spoke2Location &
wait

spoke1RGId=$(az group show --resource-group $spoke1ResourceGroupName --query id --out tsv)
spoke2RGId=$(az group show --resource-group $spoke2ResourceGroupName --query id --out tsv)

#Create Virtual Networks to deploy resources
### Spoke1
az network vnet create \
  --name $spoke1vNetName \
  --resource-group $spoke1ResourceGroupName \
  --location $spoke1Location \
  --address-prefix $spoke1vNetRange \
  --subnet-name $spoke1VnetSubnetName \
  --subnet-prefixes $spoke1SubnetRange \
  &

### Spoke2
az network vnet create \
  --name $spoke2vNetName \
  --resource-group $spoke2ResourceGroupName \
  --location $spoke2Location \
  --address-prefix $spoke2vNetRange \
  --subnet-name $spoke2VnetSubnetName \
  --subnet-prefixes $spoke2SubnetRange \
  &

wait

#  Fetch Vnet resources id's - Needed for cross resource group peering
Spoke1vNetId=$(az network vnet show \
  --resource-group $spoke1ResourceGroupName \
  --name $spoke1vNetName \
  --query id --out tsv)

Spoke2vNetId=$(az network vnet show \
  --resource-group $spoke2ResourceGroupName \
  --name $spoke2vNetName \
  --query id --out tsv)

#Fetch Subnet Id from the Spoke Vnets to deploy AKS Clusters
Spoke1vNetSubnetId=$(az network vnet subnet show \
  --resource-group $spoke1ResourceGroupName \
  --name $spoke1VnetSubnetName \
  --vnet-name $spoke1vNetName \
  --query id --out tsv)

Spoke2vNetSubnetId=$(az network vnet subnet show \
  --resource-group $spoke2ResourceGroupName \
  --name $spoke2VnetSubnetName \
  --vnet-name $spoke2vNetName \
  --query id --out tsv)

# Create managed identity for the clusters
# In order to create the cluster on the specified Vnet and subnets that were created, this identity requires to have Network Contributor role on the resource group

# cluster identity for cluster A
az identity create \
  --name $cluster1Name-identity \
  --resource-group $spoke1ResourceGroupName \
  &

# cluster identity for cluster B
az identity create \
  --name $cluster2Name-identity \
  --resource-group $spoke2ResourceGroupName \
  &

wait

cluster1IdentityId=$(az identity show -n $cluster1Name-identity -g $spoke1ResourceGroupName --query id --out tsv)
cluster1IdentityPrincipalId=$(az identity show -n $cluster1Name-identity -g $spoke1ResourceGroupName --query principalId --out tsv)
cluster2IdentityId=$(az identity show -n $cluster2Name-identity -g $spoke2ResourceGroupName --query id --out tsv)
cluster2IdentityPrincipalId=$(az identity show -n $cluster2Name-identity -g $spoke2ResourceGroupName --query principalId --out tsv)

#required permissions for identities
az role assignment create --role "Contributor" --assignee $cluster1IdentityPrincipalId --scope $spoke1RGId &
az role assignment create --role "Contributor" --assignee $cluster2IdentityPrincipalId --scope $spoke2RGId &

wait

# Create AKS clusters on Spoke regions. In order for this command to be execute, the user should have specific permissions. Check the documentation
# Cluster A
az aks create \
  --resource-group $spoke1ResourceGroupName \
  --name $cluster1Name \
  --os-sku AzureLinux \
  --node-vm-size Standard_DS2_v2 \
  --enable-managed-identity \
  --assign-identity $cluster1IdentityId \
  --vnet-subnet-id $Spoke1vNetSubnetId \
  --network-plugin azure \
  --enable-cluster-autoscaler \
  --min-count 1 \
  --max-count 3 \
  &

# Cluster B
az aks create \
  --resource-group $spoke2ResourceGroupName \
  --name $cluster2Name \
  --os-sku AzureLinux \
  --node-vm-size Standard_DS2_v2 \
  --enable-managed-identity \
  --assign-identity $cluster2IdentityId \
  --vnet-subnet-id $Spoke2vNetSubnetId \
  --network-plugin azure \
  --enable-cluster-autoscaler \
  --min-count 1 \
  --max-count 3 \
  &

wait

