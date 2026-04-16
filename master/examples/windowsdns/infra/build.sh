#!/bin/bash

#generic configs
subscriptionName="MVP Sponsorship"
dnsIPAddress="10.10.0.6"
AdminUsername="azureuser"
AdminPassword="ChangeMe123456"

#Hub configs
hubResourceGroupName="k8gb-win-hub"
hubvNetName="hub-vnet"
hubvNetRange="10.10.0.0/16"
hubVnetSubnetName="default"
hubSubnetRange="10.10.0.0/24"
hubLocation="westeurope"

#Spoke1 configs
cluster1Name="aks1"
spoke1ResourceGroupName="k8gb-win-spoke1"
spoke1vNetName="spoke1-vnet"
spoke1vNetRange="10.11.0.0/16"
spoke1VnetSubnetName="default"
spoke1SubnetRange="10.11.0.0/20"
spoke1Location="uksouth"

#Spoke2 configs
cluster2Name="aks2"
spoke2ResourceGroupName="k8gb-win-spoke2"
spoke2vNetName="spoke2-vnet"
spoke2vNetRange="10.12.0.0/16"
spoke2VnetSubnetName="default"
spoke2SubnetRange="10.12.0.0/20"
spoke2Location="francecentral"

#Private DNS configurations
privateDnsZone="global.k8gb.local"

az account set --subscription "$subscriptionName"

#Create resource groups to manage azure resources
##### Hub
az group create --resource-group $hubResourceGroupName --location $hubLocation &
##### Spoke1
az group create --resource-group $spoke1ResourceGroupName --location $spoke1Location &
##### Spoke2
az group create --resource-group $spoke2ResourceGroupName --location $spoke2Location &
wait

hubRGId=$(az group show --resource-group $hubResourceGroupName --query id --out tsv)
spoke1RGId=$(az group show --resource-group $spoke1ResourceGroupName --query id --out tsv)
spoke2RGId=$(az group show --resource-group $spoke2ResourceGroupName --query id --out tsv)

#Create Azure Private DNS Zones
az network private-dns zone create -g $spoke1ResourceGroupName -n privatelink.$spoke1Location.azmk8s.io &
az network private-dns zone create -g $spoke2ResourceGroupName -n privatelink.$spoke2Location.azmk8s.io &
wait

spoke1DnsZoneId=$(az network private-dns zone show -g $spoke1ResourceGroupName -n privatelink.$spoke1Location.azmk8s.io --query id --out tsv)
spoke2DnsZoneId=$(az network private-dns zone show -g $spoke2ResourceGroupName -n privatelink.$spoke2Location.azmk8s.io --query id --out tsv)

#Create Virtual Networks to deploy resources
### Hub
az network vnet create \
  --name $hubvNetName \
  --resource-group $hubResourceGroupName \
  --location $hubLocation \
  --address-prefix $hubvNetRange \
  --subnet-name $hubVnetSubnetName \
  --subnet-prefixes $hubSubnetRange \
  &

### Spoke1
az network vnet create \
  --name $spoke1vNetName \
  --resource-group $spoke1ResourceGroupName \
  --location $spoke1Location \
  --address-prefix $spoke1vNetRange \
  --subnet-name $spoke1VnetSubnetName \
  --subnet-prefixes $spoke1SubnetRange \
  --dns-servers $dnsIPAddress \
  &

### Spoke2
az network vnet create \
  --name $spoke2vNetName \
  --resource-group $spoke2ResourceGroupName \
  --location $spoke2Location \
  --address-prefix $spoke2vNetRange \
  --subnet-name $spoke2VnetSubnetName \
  --subnet-prefixes $spoke2SubnetRange \
  --dns-servers $dnsIPAddress \
  &

wait

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
  --name $hubvNetName-$spoke1vNetName \
  --remote-vnet $Spoke1vNetId \
  --resource-group $hubResourceGroupName \
  --vnet-name $hubvNetName \
  --allow-forwarded-traffic \
  --allow-gateway-transit \
  --allow-vnet-access \
  &

az network vnet peering create \
  --name $spoke1vNetName-$hubvNetName \
  --remote-vnet $HubvNetId \
  --resource-group $spoke1ResourceGroupName \
  --vnet-name $spoke1vNetName \
  --allow-forwarded-traffic \
  --allow-gateway-transit \
  --allow-vnet-access \
  &

### Hub <-> Spoke2
az network vnet peering create \
  --name $hubvNetName-$spoke2vNetName \
  --remote-vnet $Spoke2vNetId \
  --resource-group $hubResourceGroupName \
  --vnet-name $hubvNetName \
  --allow-forwarded-traffic \
  --allow-gateway-transit \
  --allow-vnet-access \
  &

az network vnet peering create \
  --name $spoke2vNetName-$hubvNetName \
  --remote-vnet $HubvNetId \
  --resource-group $spoke2ResourceGroupName \
  --vnet-name $spoke2vNetName \
  --allow-forwarded-traffic \
  --allow-gateway-transit \
  --allow-vnet-access \
  &

# Link Private DNS zones
az network private-dns link vnet create \
  --resource-group $spoke1ResourceGroupName \
  --name privatelink.$spoke1Location.azmk8s.io-$hubvNetName \
  --registration-enabled false \
  --virtual-network $HubvNetId \
  --zone-name privatelink.$spoke1Location.azmk8s.io \
  &

az network private-dns link vnet create \
  --resource-group $spoke2ResourceGroupName \
  --name privatelink.$spoke2Location.azmk8s.io-$hubvNetName \
  --registration-enabled false \
  --virtual-network $HubvNetId \
  --zone-name privatelink.$spoke2Location.azmk8s.io \
  &

wait

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

# Create NSG for DC
az network nsg create --name dc-nsg --resource-group $hubResourceGroupName
az network nsg rule create --nsg-name dc-nsg --resource-group $hubResourceGroupName --name AllowRDP --priority 100 --access Allow --direction Inbound --source-address-prefixes "10.0.0.0/8" --destination-port-ranges 3389 --protocol Tcp &
az network nsg rule create --nsg-name dc-nsg --resource-group $hubResourceGroupName --name AllowDNS --priority 110 --access Allow --direction Inbound --source-address-prefixes "10.0.0.0/8" --destination-port-ranges 53 --protocol "*" &
wait

nsgId=$(az network nsg show --resource-group $hubResourceGroupName --name dc-nsg --query id --out tsv)

# Create DC VM on the Hub Vnet
az vm create \
    --resource-group $hubResourceGroupName \
    --name dc \
    --size Standard_DS2_v2 \
    --image Win2022Datacenter \
    --admin-username $AdminUsername \
    --admin-password $AdminPassword \
    --private-ip-address $dnsIPAddress \
    --nsg $nsgId \
    --public-ip-address ""

# Create Bastion
az network public-ip create --resource-group $hubResourceGroupName --name $hubvNetName-ip --sku Standard --location $hubLocation
az network vnet subnet create --name AzureBastionSubnet --resource-group $hubResourceGroupName --vnet-name $hubvNetName --address-prefix "10.10.1.0/26"
az network bastion create --name $hubvNetName-bastion --public-ip-address $hubvNetName-ip --resource-group $hubResourceGroupName --vnet-name $hubvNetName --location $hubLocation --sku Basic &

# Promote VM to DC
az vm run-command invoke \
    --resource-group $hubResourceGroupName \
    --name dc \
    --command-id RunPowerShellScript \
    --scripts @./infra/dc.ps1 \
    &

wait

# Update hub VNet to point do DNS in DC VM
az network vnet update --resource-group $hubResourceGroupName --name $hubvNetName --dns-servers $dnsIPAddress

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
  --enable-private-cluster \
  --private-dns-zone $spoke1DnsZoneId \
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
  --enable-private-cluster \
  --private-dns-zone $spoke2DnsZoneId \
  --network-plugin azure \
  --enable-cluster-autoscaler \
  --min-count 1 \
  --max-count 3 \
  &

wait

