#!/bin/bash

[ "${DEBUG}" == 1 ] && set -x

TOOL_VERSION=${TOOL_VERSION:-"0.5.4"}
TOOL_REPO=${TOOL_REPO:-"AbsaOSS"}
DIR="$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

main() {
    # checks
    [[ $# != 1 ]] && echo "Usage: $0 <version> # provide version in x.y.z format" && exit 1
    _VERSION=$1
    _VERSION=${_VERSION#"v"}
    _OS=$(go env GOOS)
    _ARCH=$(go env GOARCH)

    # download olm-bundle if not present locally
    if ! which olm-bundle > /dev/null; then
        [ -f ${DIR}/olm-bundle ] || downloadOlmBundle
        OLM_BINARY="${DIR}/olm-bundle"
    else
        OLM_BINARY="olm-bundle"
    fi

    # if master is set as version don't do the checkout and use the latest annotated tag (~release) as the
    # desired version for the olm bundle
    if [[ ${_VERSION} == "master" ]]; then
        _VERSION=$(git describe --abbrev=0 --tags)
        _VERSION=${_VERSION#"v"}
    else
        git checkout v${_VERSION}
    fi
    PREVIOUS_VERSION=${PREVIOUS_VERSION:-$(git describe --abbrev=0 --tags v${_VERSION}^)}

    generate
}

generate() {
    echo "    containerImage: docker.io/absaoss/k8gb:v${_VERSION}" >> ${DIR}/annotations.yaml.tmpl
    cd ${DIR}/../chart/k8gb && helm dependency update && cd -
    helm -n placeholder template ${DIR}/../chart/k8gb \
        --name-template=k8gb \
        --set k8gb.securityContext.runAsNonRoot=null \
        --set k8gb.securityContext.runAsUser=null | ${OLM_BINARY} \
            --chart-file-path=${DIR}/../chart/k8gb/Chart.yaml \
            --version=v${_VERSION} \
            --replaces-version=${PREVIOUS_VERSION} \
            --helm-chart-overrides \
            --output-dir ${DIR}
    git checkout ${DIR}/annotations.yaml.tmpl
}

downloadOlmBundle() {
    curl -Lo ${DIR}/olm-bundle https://github.com/${TOOL_REPO}/olm-bundle/releases/download/v${TOOL_VERSION}/olm-bundle_${_OS}-${_ARCH}
    chmod +x ${DIR}/olm-bundle
}

main $@
