#!/bin/bash
set -o errexit

# Workaround of Make being to smart on up-to-date PHONY targets
# If we execute all of them normal way, then targets from `deploy-second-k8gb`
# will never be executed as they contain the same underlying target as `deploy-first-k8gb`
# but with different variables

NODE_ROLE=${NODE_ROLE:-worker}
ADDITIONAL_TARGETS=${ADDITIONAL_TARGETS:-}
TEST_CURRENT_COMMIT=${TEST_CURRENT_COMMIT:-}

if [ "$TEST_CURRENT_COMMIT" == "yes" ]
then
    ./deploy/registry.sh
    commit_hash=$(git rev-parse --short HEAD)
    operator-sdk build k8gb:${commit_hash}
    docker tag k8gb:${commit_hash} localhost:5000/k8gb:v${commit_hash}
    docker push localhost:5000/k8gb:v${commit_hash}
    export K8GB_IMAGE_REPO=localhost:5000/k8gb
    sed -i "s/${VERSION}/${commit_hash}/g" chart/k8gb/Chart.yaml
fi

make use-second-context
export HOST_ALIAS_IP1=$(kubectl get nodes test-gslb2-${NODE_ROLE} -o custom-columns='IP:status.addresses[0].address' --no-headers)
make use-first-context deploy-first-k8gb deploy-gslb-cr ${ADDITIONAL_TARGETS}
export HOST_ALIAS_IP2=$(kubectl get nodes test-gslb1-${NODE_ROLE} -o custom-columns='IP:status.addresses[0].address' --no-headers)
make use-second-context deploy-second-k8gb deploy-gslb-cr ${ADDITIONAL_TARGETS}

make wait-for-nginx-ingress-ready
make wait-for-gslb-ready

make use-first-context
make wait-for-nginx-ingress-ready
make wait-for-gslb-ready
