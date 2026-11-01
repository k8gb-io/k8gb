#!/usr/bin/env bash
# Collect state before Chainsaw cleans up a failed test. Collection errors are
# retained in the artifact files so an unavailable cluster does not hide evidence
# from the other cluster or replace the original test failure.
set -eu

namespace=${1:?Usage: collect-chainsaw-debug.sh NAMESPACE}
repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
artifact_root=${CHAINSAW_ARTIFACTS_DIR:-"$repo_root/tmp/chainsaw"}
mkdir -p "$artifact_root/$namespace"
snapshot=$(mktemp -d "$artifact_root/$namespace/failure-XXXXXX")

collect_cluster() {
  local cluster=$1 output kubeconfig
  local -a kubectl_cmd
  output="$snapshot/$cluster"
  mkdir -p "$output"
  kubeconfig="$repo_root/chainsaw/kubeconfig/$cluster.config"
  if [[ ! -f "$kubeconfig" ]]; then
    echo "Missing kubeconfig: $kubeconfig" > "$output/collection-error.log"
    return
  fi

  kubectl_cmd=(kubectl --kubeconfig="$kubeconfig" --request-timeout=5s)
  "${kubectl_cmd[@]}" get gslbs.k8gb.absa.oss,gslbs.k8gb.io,dnsendpoints.externaldns.k8s.io,ingresses,services,pods \
    -n "$namespace" -o yaml > "$output/resources.yaml" 2> "$output/resources.stderr" || true
  "${kubectl_cmd[@]}" get zonedelegations.k8gb.io -o yaml \
    > "$output/zonedelegations.yaml" 2> "$output/zonedelegations.stderr" || true
  "${kubectl_cmd[@]}" get events -n "$namespace" --sort-by=.metadata.creationTimestamp \
    > "$output/events.log" 2>&1 || true
  "${kubectl_cmd[@]}" logs -n k8gb -l name=k8gb --all-containers=true --tail=1000 --timestamps=true \
    > "$output/controller.log" 2>&1 || true
}

for cluster in eu us; do
  collect_cluster "$cluster" &
done
wait

echo "Chainsaw failure diagnostics: $snapshot"
