#!/bin/bash

#generic configs
subscriptionName="MBCPANLS02"
windowsDnsServers="10.121.64.196"

#Hub configs
hubResourceGroupName="RGCM2LILABK8GB001"
hubvNetName="VNCM2LINK8GB001"
hubvNetRange="10.121.240.0/20"
hubVnetSubnetName="VNCM2LINK8GB001-APP01"
hubSubnetRange="10.121.240.0/22"
hubLocation="westeurope"

#Spoke1 configs
cluster1Name="cscm2lakslabK8GB001"
spoke1ResourceGroupName="RGCM2LILABK8GB002"
spoke1vNetName="VNCM2LINK8GB002"
spoke1vNetRange="10.111.240.0/20"
spoke1VnetSubnetName="VNCM2LINK8GB002-APP01"
spoke1SubnetRange="10.111.240.0/22"
spoke1Location="westeurope"

#Spoke2 configs
cluster2Name="cscm2lakslabK8GB002"
spoke2ResourceGroupName="RGCM2LILABK8GB003"
spoke2vNetName="VNCM2LINK8GB003"
spoke2vNetRange="10.101.240.0/20"
spoke2VnetSubnetName="VNCM2LINK8GB003-APP01"
spoke2SubnetRange="10.101.240.0/22"
spoke2Location="northeurope"

# Private DNS configurations
privateDnsZone="lab.mbcp.cloud"

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