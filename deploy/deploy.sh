#!/bin/sh

set -ex

application_order="namespace.yaml
role.yaml
role_binding.yaml
service_account.yaml
operator.yaml"

for file in $application_order
do
    kubectl apply -f "$file"
done
