#!/bin/sh
set -o errexit

DOCKER_REGISTRY_IMAGE=${DOCKER_REGISTRY_IMAGE:-'registry:2'}

# create registry container unless it already exists
reg_name='kind-registry'
reg_port='5000'
running="$(docker inspect -f '{{.State.Running}}' "${reg_name}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    -d --restart=always -p "${reg_port}:5000" --name "${reg_name}" \
    "${DOCKER_REGISTRY_IMAGE}"
fi
