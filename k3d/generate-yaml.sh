#!/bin/bash
# use this simple script for generating additional config files for k3d
# note: make sure the generated ports doesn't collide with your local environment
# for instance on mac port 88 might be taken by kdc

[[ $# != 1 ]] && echo "Usage: $0 <how_many>" && exit 1
[[ -z "${1##*[!0-9]*}" ]] && echo "'$1' is not a positive integer" && exit 1
UPTO=$1
DIR="${DIR:-$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )}"

echo Generating configs for $UPTO clusters..
for c in $(seq $UPTO); do
    export CLUSTER_INDEX=$(( 1 + $c ))
    export PORT_HTTP=$(( 80 + $c ))
    export PORT_HTTPS=$(( 443 + $c ))
    export PORT_PROM=$(( 9090 + $c ))
    export PORT_DNS=$(( 5053 + $c ))
    cat ${DIR}/gslb.yaml.tmpl | envsubst > ${DIR}/test-gslb${CLUSTER_INDEX}.yaml
done
