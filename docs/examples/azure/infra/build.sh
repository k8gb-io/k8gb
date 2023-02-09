#!/bin/bash

subscriptionName="MBCPANLS02"
#Hub configs
hubResourceGroupName="RGCM2LILABK8GB001"
hubvNetName="VNCM2LINK8GB001"
hubVnetSubnetName="VNCM2LINK8GB001-APP01"

#Spoke1 configs
clusterName="cscm2lakslabK8GB002"

#Spoke2 configs

az account set --subscription $subscriptionName
#Create resource groups to manage azure resources
##### Hub
az group create --resource-group $hubResourceGroupName --location westeurope --tags solution=shared
##### Spoke1
az group create --resource-group RGCM2LILABK8GB002 --location westeurope --tags solution=shared
##### Spoke2
az group create --resource-group RGCM2LILABK8GB003 --location northeurope --tags solution=shared

#### The Private DNS Zone will be deployed on the same resource group for the Hub. While this is acceptable for Lab environment, for other scenarios, a specific resource group should be created
#Create an Azure Private DNS Zone*
az network private-dns zone create -g $hubResourceGroupName -n lab.mbcp.cloud

#Create Virtual Networks to deploy resources
### Hub
az network vnet create \
  --name VNCM2LINK8GB001 \
  --resource-group RGCM2LILABK8GB001 \
  --location westeurope \
  --address-prefix 10.121.240.0/20 \
  --subnet-name VNCM2LINK8GB001-HUB \
  --subnet-prefixes 10.121.240.0/22 \
  --dns-servers 10.121.64.196 10.121.64.197

### Spoke1
az network vnet create \
  --name VNCM2LINK8GB002 \
  --resource-group RGCM2LILABK8GB002 \
  --location westeurope \
  --address-prefix 10.111.240.0/20 \
  --subnet-name VNCM2LINK8GB002-APP01 \
  --subnet-prefixes 10.111.240.0/22 \
  --dns-servers 10.121.64.196 10.121.64.197

### Spoke2
az network vnet create \
  --name VNCM2LINK8GB003 \
  --resource-group RGCM2LILABK8GB003 \
  --location northeurope \
  --address-prefix 10.101.240.0/20 \
  --subnet-name VNCM2LINK8GB003-APP01 \
  --subnet-prefixes 10.101.240.0/22 \
  --dns-servers 10.121.64.196 10.121.64.197

#  Fetch Vnet resources id's - Needed for cross resource group peering
HubvNetId=$(az network vnet show \
  --resource-group RGCM2LILABK8GB001 \
  --name VNCM2LINK8GB001 \
  --query id --out tsv)

Spoke1vNetId=$(az network vnet show \
  --resource-group RGCM2LILABK8GB002 \
  --name VNCM2LINK8GB002 \
  --query id --out tsv)

Spoke2vNetId=$(az network vnet show \
  --resource-group RGCM2LILABK8GB003 \
  --name VNCM2LINK8GB003 \
  --query id --out tsv)

#  Peer Spoke Vnets with Hub Vnet - A peering requires configuration from both Vnets to each other, in order to be fully in sync
### Hub <-> Spoke1
az network vnet peering create \
   --name VNCM2LINK8GB001-VNCM2LINK8GB002 \
   --remote-vnet $Spoke1vNetId \
   --resource-group RGCM2LILABK8GB001 \
   --vnet-name VNCM2LINK8GB001 \
   --allow-forwarded-traffic \
   --allow-gateway-transit \
   --allow-vnet-access

az network vnet peering create \
   --name VNCM2LINK8GB001-VNCM2LINK8GB002 \
   --remote-vnet $HubvNetId \
   --resource-group RGCM2LILABK8GB002 \
   --vnet-name VNCM2LINK8GB002 \
   --allow-forwarded-traffic \
   --allow-gateway-transit \
   --allow-vnet-access

### Hub <-> Spoke2
az network vnet peering create \
   --name VNCM2LINK8GB001-VNCM2LINK8GB003 \
   --remote-vnet $Spoke2vNetId \
   --resource-group RGCM2LILABK8GB001 \
   --vnet-name VNCM2LINK8GB001 \
   --allow-forwarded-traffic \
   --allow-gateway-transit \
   --allow-vnet-access

az network vnet peering create \
   --name VNCM2LINK8GB001-VNCM2LINK8GB003 \
   --remote-vnet $HubvNetId \
   --resource-group RGCM2LILABK8GB003 \
   --vnet-name VNCM2LINK8GB003 \
   --allow-forwarded-traffic \
   --allow-gateway-transit \
   --allow-vnet-access

# Link Private DNS zone with the Hub Vnet
az network private-dns link vnet create \
    --resource-group rgcm2liglb001 \
    --name lab.mbcp.cloud-k8gb-hub \
    --registration-enabled false \
    --virtual-network $HubvNetId \
    --zone-name lab.mbcp.cloud

#Fetch Subnet Id from the Spoke Vnets to deploy AKS Clusters
Spoke1vNetSubnetId=$(az network vnet subnet show \
  --resource-group RGCM2LILABK8GB002 \
  --name VNCM2LINK8GB002-APP01 \
  --vnet-name VNCM2LINK8GB002 \
  --query id --out tsv)

Spoke2vNetSubnetId=$(az network vnet subnet show \
  --resource-group RGCM2LILABK8GB003 \
  --name VNCM2LINK8GB003-APP01 \
  --vnet-name VNCM2LINK8GB003 \
  --query id --out tsv)

# Create managed identity for the clusters
# In order to create the cluster on the specified Vnet and subnets that were created, this identity requires to have Network Contributor role on the resource group
# cluster identity for cluster A
az identity create \
  --name cscm2lakslabK8GB001-identity \
  --resource-group RGCM2LILABK8GB002

clusterAIdentityId=$(az identity show -n cscm2lakslabK8GB001-identity -g RGCM2LILABK8GB002 --query id --out tsv) 

clusterAIdentityClientId=$(az identity show -n cscm2lakslabK8GB001-identity -g RGCM2LILABK8GB002 --query clientId --out tsv)

#required permissions
az role assignment create --role "Network Contributor" --assignee cscm2lakslabK8GB001-identity -g RGCM2LILABK8GB002

# cluster identity for cluster B
az identity create \
  --name cscm2lakslabK8GB002-identity \
  --resource-group RGCM2LILABK8GB003

clusterBIdentityId=$(az identity show -n cscm2lakslabK8GB002-identity -g RGCM2LILABK8GB003 --query id --out tsv) 

clusterBIdentityClientId=$(az identity show -n cscm2lakslabK8GB002-identity -g RGCM2LILABK8GB003 --query clientId --out tsv)

# required permissions
az role assignment create --role "Network Contributor" --assignee cscm2lakslabK8GB002-identity -g RGCM2LILABK8GB003

# Create AKS clusters on Spoke regions. In order for this command to be execute, the user should have specific permissions. Check the documentation 
# Cluster A
az aks create \
  --resource-group RGCM2LILABK8GB002 \
  --name cscm2lakslabK8GB001 \
  --generate-ssh-keys \
  --vm-set-type VirtualMachineScaleSets \
  --os-sku CBLMariner \
  --node-vm-size Standard_DS2_v2 \
  --vnet-subnet-id $Spoke1vNetSubnetId \
  --load-balancer-sku standard \
  --enable-managed-identity \
  --assign-identity $clusterAIdentityId \
  --network-plugin azure \
  --node-count 1 \
  --zones 1

# Cluster B
az aks create \
  --resource-group RGCM2LILABK8GB003 \
  --name cscm2lakslabK8GB002 \
  --generate-ssh-keys \
  --vm-set-type VirtualMachineScaleSets \
  --os-sku CBLMariner \
  --node-vm-size Standard_DS2_v2 \
  --vnet-subnet-id $Spoke2vNetSubnetId \
  --load-balancer-sku standard \
  --enable-managed-identity \
  --assign-identity $clusterBIdentityId \
  --network-plugin azure \
  --node-count 1 \
  --zones 1