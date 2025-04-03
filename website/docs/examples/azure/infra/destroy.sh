#!/bin/bash

#generic configs
subscriptionName="MVP Sponsorship"

#Spoke1 configs
spoke1ResourceGroupName="k8gb-az-spoke1"

#Spoke2 configs
spoke2ResourceGroupName="k8gb-az-spoke2"

#################
# Set subscription
##################
az account set --subscription "$subscriptionName"

#################
# Delete resource groups
##################
az group delete -n $spoke1ResourceGroupName -y &
az group delete -n $spoke2ResourceGroupName -y &

wait