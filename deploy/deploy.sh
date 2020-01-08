#!/bin/sh

set -ex

# Deploy CRD and CR
kubectl apply -f deploy/crds

# Deploy Gslb operator
application_order="namespace.yaml
role.yaml
role_binding.yaml
service_account.yaml
operator.yaml"

for file in $application_order
do
    kubectl apply -f "deploy/$file"
done

# Deploy external dependencies 
grep ^helm deploy/coredns/README.md | xargs -n1 -d '\n' |bash
kubectl apply -f deploy/coredns/external-dns.yaml
