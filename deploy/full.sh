#!/bin/bash
set -o errexit

# Workaround of Make being to smart on up-to-date PHONY targets
# If we execute all of them normal way, then targets from `deploy-second-ohmyglb`
# will never be executed as they contain the same underlying target as `deploy-first-ohmyglb`
# but with different variables

NODE_ROLE=${NODE_ROLE:-worker}
ADDITIONAL_TARGETS=${ADDITIONAL_TARGETS:-}
TEST_CURRENT_COMMIT=${TEST_CURRENT_COMMIT:-}

if [ "$TEST_CURRENT_COMMIT" == "yes" ]
then
    ./deploy/registry.sh
    commit_hash=$(git rev-parse --short HEAD)
    operator-sdk build ohmyglb:${commit_hash}
    docker tag ohmyglb:${commit_hash} localhost:5000/ohmyglb:${commit_hash}
    docker push localhost:5000/ohmyglb:${commit_hash}
    export OHMYGLB_IMAGE=localhost:5000/ohmyglb:${commit_hash}
fi

make use-second-context
export HOST_ALIAS_IP1=$(kubectl get nodes test-gslb2-${NODE_ROLE} -o custom-columns='IP:status.addresses[0].address' --no-headers)
export REGION_ARG="eu"
make use-first-context deploy-first-ohmyglb deploy-gslb-cr ${ADDITIONAL_TARGETS}
export HOST_ALIAS_IP2=$(kubectl get nodes test-gslb1-${NODE_ROLE} -o custom-columns='IP:status.addresses[0].address' --no-headers)
export REGION_ARG="us"
make use-second-context deploy-second-ohmyglb deploy-gslb-cr ${ADDITIONAL_TARGETS}

make wait-for-nginx-ingress-ready
make wait-for-gslb-ready

make use-first-context
make wait-for-nginx-ingress-ready
make wait-for-gslb-ready
