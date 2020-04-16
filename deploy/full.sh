#!/bin/bash
# Workaround of Make being to smart on up-to-date PHONY targets
# If we execute all of them normal way, then targets from `deploy-second-ohmyglb`
# will never be executed as they contain the same underlying target as `deploy-first-ohmyglb`
# but with different variables

make use-second-context
export HOST_ALIAS_IP1=$(kubectl get nodes test-gslb2-worker -o custom-columns='IP:status.addresses[0].address' --no-headers)
make use-first-context deploy-first-ohmyglb deploy-gslb-cr $1
export HOST_ALIAS_IP2=$(kubectl get nodes test-gslb1-worker -o custom-columns='IP:status.addresses[0].address' --no-headers)
make use-second-context deploy-second-ohmyglb deploy-gslb-cr $1

make wait-for-nginx-ingress-ready
make wait-for-gslb-ready

make use-first-context
make wait-for-nginx-ingress-ready
make wait-for-gslb-ready
