#!/bin/bash

#generic configs
subscriptionName=""
windowsDnsServers=""

#Hub configs
hubResourceGroupName=""
hubvNetName=""
hubvNetRange=""
hubVnetSubnetName=""
hubSubnetRange=""
hubLocation=""

#Spoke1 configs
cluster1Name=""
spoke1ResourceGroupName=""
spoke1vNetName=""
spoke1vNetRange=""
spoke1VnetSubnetName=""
spoke1SubnetRange=""
spoke1Location=""

#Spoke2 configs
cluster2Name=""
spoke2ResourceGroupName=""
spoke2vNetName=""
spoke2vNetRange=""
spoke2VnetSubnetName=""
spoke2SubnetRange=""
spoke2Location=""

#Private DNS configurations
privateDnsZone=""

az account set --subscription $subscriptionName
#Create resource groups to manage azure resources
##### Hub
az group create --resource-group $hubResourceGroupName --location $hubLocation --tags solution=shared
##### Spoke1
az group create --resource-group $spoke1ResourceGroupName --location $spoke1Location --tags solution=shared
##### Spoke2
az group create --resource-group $spoke2ResourceGroupName --location $spoke2Location --tags solution=shared

#### The Private DNS Zone will be deployed on the same resource group for the Hub. While this is acceptable for Lab environment, for other scenarios, a specific resource group should be created
#Create an Azure Private DNS Zone*
az network private-dns zone create -g $hubResourceGroupName -n $privateDnsZone

#Create Virtual Networks to deploy resources
### Hub
az network vnet create \
  --name $hubVnetName \
  --resource-group $hubResourceGroupName \
  --location $hubLocation \
  --address-prefix $hubvNetRange \
  --subnet-name $hubVnetSubnetName \
  --subnet-prefixes $hubSubnetRange \
  --dns-servers $windowsDnsServers

### Spoke1
az network vnet create \
  --name $spoke1vNetName \
  --resource-group $spoke1ResourceGroupName \
  --location $spoke1Location \
  --address-prefix $spoke1vNetRange \
  --subnet-name $spoke1VnetSubnetName \
  --subnet-prefixes $spoke1SubnetRange \
  --dns-servers $windowsDnsServers

### Spoke2
az network vnet create \
  --name $spoke2vNetName \
  --resource-group $spoke2ResourceGroupName \
  --location $spoke2Location \
  --address-prefix $spoke2vNetRange \
  --subnet-name $spoke2VnetSubnetName \
  --subnet-prefixes $spoke2SubnetRange \
  --dns-servers $windowsDnsServers

#  Fetch Vnet resources id's - Needed for cross resource group peering
HubvNetId=$(az network vnet show \
  --resource-group $hubResourceGroupName \
  --name $hubvNetName \
  --query id --out tsv)

Spoke1vNetId=$(az network vnet show \
  --resource-group $spoke1ResourceGroupName \
  --name $spoke1vNetName \
  --query id --out tsv)

Spoke2vNetId=$(az network vnet show \
  --resource-group $spoke2ResourceGroupName \
  --name $spoke2vNetName \
  --query id --out tsv)

#  Peer Spoke Vnets with Hub Vnet - A peering requires configuration from both Vnets to each other, in order to be fully in sync
### Hub <-> Spoke1
az network vnet peering create \
   --name $hubVnetName-$spoke1vNetName \
   --remote-vnet $Spoke1vNetId \
   --resource-group $hubResourceGroupName \
   --vnet-name $hubVnetName \
   --allow-forwarded-traffic \
   --allow-gateway-transit \
   --allow-vnet-access

az network vnet peering create \
   --name $spoke1vNetName-$hubVnetName \
   --remote-vnet $HubvNetId \
   --resource-group $spoke1ResourceGroupName \
   --vnet-name $spoke1vNetName \
   --allow-forwarded-traffic \
   --allow-gateway-transit \
   --allow-vnet-access

### Hub <-> Spoke2
az network vnet peering create \
   --name $hubVnetName-$spoke2vNetName \
   --remote-vnet $Spoke2vNetId \
   --resource-group $hubResourceGroupName \
   --vnet-name $hubVnetName \
   --allow-forwarded-traffic \
   --allow-gateway-transit \
   --allow-vnet-access

az network vnet peering create \
   --name $spoke2vNetName-$hubVnetName \
   --remote-vnet $HubvNetId \
   --resource-group $spoke2ResourceGroupName \
   --vnet-name $spoke2vNetName \
   --allow-forwarded-traffic \
   --allow-gateway-transit \
   --allow-vnet-access

# Link Private DNS zone with the Hub Vnet
az network private-dns link vnet create \
    --resource-group $hubResourceGroupName \
    --name $privateDnsZone-$hubVnetName \
    --registration-enabled false \
    --virtual-network $HubvNetId \
    --zone-name $privateDnsZone

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
  --resource-group $spoke1ResourceGroupName

cluster1IdentityId=$(az identity show -n $cluster1Name-identity -g $spoke1ResourceGroupName --query id --out tsv) 

cluster1IdentityClientId=$(az identity show -n $cluster1Name-identity -g $spoke1ResourceGroupName --query clientId --out tsv)

#required permissions for the identity
az role assignment create --role "Network Contributor" --assignee $cluster1Name-identity -g $spoke1ResourceGroupName

# cluster identity for cluster B
az identity create \
  --name $cluster2Name-identity \
  --resource-group $spoke2ResourceGroupName

cluster2IdentityId=$(az identity show -n $cluster2Name-identity -g $spoke2ResourceGroupName --query id --out tsv) 

cluster2IdentityClientId=$(az identity show -n $cluster2Name-identity -g $spoke2ResourceGroupName --query clientId --out tsv)

#required permissions for the identity
az role assignment create --role "Network Contributor" --assignee $cluster2Name-identity -g $spoke2ResourceGroupName

# Create AKS clusters on Spoke regions. In order for this command to be execute, the user should have specific permissions. Check the documentation 
# Cluster A
az aks create \
  --resource-group $spoke1ResourceGroupName \
  --name $cluster1Name \
  --generate-ssh-keys \
  --vm-set-type VirtualMachineScaleSets \
  --os-sku CBLMariner \
  --node-vm-size Standard_DS2_v2 \
  --vnet-subnet-id $Spoke1vNetSubnetId \
  --load-balancer-sku standard \
  --enable-managed-identity \
  --assign-identity $cluster1IdentityId \
  --network-plugin azure \
  --node-count 1 \
  --zones 1

# Cluster B
az aks create \
  --resource-group $spoke2ResourceGroupName \
  --name $cluster2Name \
  --generate-ssh-keys \
  --vm-set-type VirtualMachineScaleSets \
  --os-sku CBLMariner \
  --node-vm-size Standard_DS2_v2 \
  --vnet-subnet-id $Spoke2vNetSubnetId \
  --load-balancer-sku standard \
  --enable-managed-identity \
  --assign-identity $cluster2IdentityId \
  --network-plugin azure \
  --node-count 1 \
  --zones 1