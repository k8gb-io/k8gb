#!/bin/bash

#generic configs
subscriptionName=""


#Hub configs
hubResourceGroupName=""
hubvNetName=""
hubVnetSubnetName=""
hubLocation=""

#Spoke1 configs
cluster1Name=""
spoke1ResourceGroupName=""
spoke1vNetName=""
spoke1VnetSubnetName=""
spoke1Location=""

#Spoke2 configs
cluster2Name=""
spoke2ResourceGroupName=""
spoke2vNetName=""
spoke2VnetSubnetName=""
spoke2Location=""

# Private DNS configurations
privateDnsZone=""

#################
# Set subscription
##################
az account set --subscription $subscriptionName

#################
# Delete clusters
##################
az aks delete -n $cluster2Name -g $spoke2ResourceGroupName --yes
az aks delete -n $cluster1Name -g $spoke1ResourceGroupName --yes

#################
# Delete clusters identities
##################
cluster1IdentityId=$(az identity show -n $cluster1Name-identity -g $spoke1ResourceGroupName --query id --out tsv) 
az identity delete --ids $cluster1IdentityId -n $cluster1Name-identity -g $spoke1ResourceGroupName

cluster2IdentityId=$(az identity show -n $cluster2Name-identity -g $spoke2ResourceGroupName --query id --out tsv) 
az identity delete --ids $cluster2IdentityId -n $cluster2Name-identity -g $spoke2ResourceGroupName

#################
# Delete Link Private DNS zone with the Hub Vnet
##################
az network private-dns link vnet delete -n $privateDnsZone-$hubVnetName -g $hubResourceGroupName -z $privateDnsZone

#################
# Delete peers between Spoke Vnets and Hub Vnet
##################
# Hub - Spoke 2
##################
az network vnet peering delete -n $spoke2vNetName-$hubVnetName -g $spoke2ResourceGroupName --vnet-name $spoke2vNetName
az network vnet peering delete -n $hubVnetName-$spoke2vNetName -g $hubResourceGroupName --vnet-name $hubVnetName
##################
# Hub - Spoke 1
##################
az network vnet peering delete -n $spoke1vNetName-$hubVnetName -g $spoke1ResourceGroupName --vnet-name $spoke1vNetName
az network vnet peering delete -n $hubVnetName-$spoke1vNetName -g $hubResourceGroupName --vnet-name $hubVnetName

#################
# Delete vnets
##################
az network vnet delete -g $spoke2ResourceGroupName -n $spoke2vNetName
az network vnet delete -g $spoke1ResourceGroupName -n $spoke1vNetName
az network vnet delete -g $hubResourceGroupName -n $hubVnetName

#################
# Delete private zone
##################
az network private-dns zone delete -g $hubResourceGroupName -n $privateDnsZone

#################
# Delete resource groups
##################
az group delete -n $spoke2ResourceGroupName
az group delete -n $spoke1ResourceGroupName
az group delete -n $hubResourceGroupName