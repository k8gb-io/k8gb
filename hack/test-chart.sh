#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
CHART_DIR="${ROOT_DIR}/chart/k8gb"
VALUES_FILE="$(mktemp)"
HELM_HOME="$(mktemp -d)"
trap 'rm -f "${VALUES_FILE}"; rm -rf "${HELM_HOME}"' EXIT

export HELM_CACHE_HOME="${HELM_HOME}/cache"
export HELM_CONFIG_HOME="${HELM_HOME}/config"
export HELM_DATA_HOME="${HELM_HOME}/data"

cat >"${VALUES_FILE}" <<'EOF'
k8gb:
  coredns:
    extraServerBlocks: |
      delegated.{{ .Release.Namespace }}:5353 {
        errors
        forward . 192.0.2.53
      }
EOF

helm repo add coredns https://coredns.github.io/helm >/dev/null
helm repo add external-dns https://kubernetes-sigs.github.io/external-dns >/dev/null
helm dependency build "${CHART_DIR}" >/dev/null

rendered="$(helm template extra-server-blocks "${CHART_DIR}" \
  --namespace chart-test \
  --values "${VALUES_FILE}" \
  --show-only templates/coredns/cm.yaml)"

server_block_count="$(grep -Fc 'delegated.chart-test:5353 {' <<<"${rendered}" || true)"
if [[ "${server_block_count}" -ne 1 ]]; then
  echo "expected the configured global CoreDNS server block exactly once, got ${server_block_count}" >&2
  exit 1
fi

grep -Fq 'forward . 192.0.2.53' <<<"${rendered}" || {
  echo "expected the configured global CoreDNS server block contents" >&2
  exit 1
}

echo "PASS: k8gb.coredns.extraServerBlocks renders in the Corefile"
