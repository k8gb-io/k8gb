#!/usr/bin/env bash
set -euo pipefail

# A disposable k3d reproducer for #2529; no changes to the user's kubeconfig.
repo_root=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")/.." && pwd)
cluster="k8gb-localhost-$$"
work_dir=$(mktemp -d)
export KUBECONFIG="$work_dir/kubeconfig"
trap 'k3d cluster delete "$cluster"; rm -rf "$work_dir"' EXIT
k3d cluster create "$cluster" --image "rancher/k3s:${1:-v1.34.3-k3s1}" \
  --no-lb --k3s-arg '--disable=traefik,metrics-server,servicelb@server:*' \
  --kubeconfig-update-default=false --kubeconfig-switch-context=false --wait
k3d kubeconfig get "$cluster" > "$KUBECONFIG"

# Model the agent's positive sinkhole answer with a documentation-only address.
kubectl apply -f - <<'YAML'
apiVersion: v1
kind: ConfigMap
metadata:
  name: coredns-custom
  namespace: kube-system
data:
  localhost.server: |
    localhost:53 {
      log
      template IN A {
        answer "{{ .Name }} 60 IN A 192.0.2.1"
      }
      template IN AAAA {
        rcode NOERROR
      }
    }
YAML
kubectl -n kube-system rollout restart deployment/coredns
kubectl -n kube-system rollout status deployment/coredns --timeout=120s

helm repo add podinfo https://stefanprodan.github.io/podinfo
helm install frontend podinfo/podinfo --version 5.1.1
kubectl wait --for=jsonpath='{.status.phase}'=Running pod \
  -l app.kubernetes.io/name=frontend-podinfo --timeout=120s
kubectl exec deployment/frontend-podinfo -- nslookup -type=A localhost. | tee "$work_dir/old-dns.log"
grep -F 192.0.2.1 "$work_dir/old-dns.log"
kubectl exec deployment/frontend-podinfo -- podcli check http 127.0.0.1:9898/readyz
if kubectl exec deployment/frontend-podinfo -- env GODEBUG=netdns=go+2 \
  podcli check http localhost:9898/readyz > "$work_dir/old-probe.log" 2>&1; then
  echo 'FAIL: Podinfo 5.1.1 unexpectedly resolved localhost locally'
  exit 1
fi
cat "$work_dir/old-probe.log"
grep -F 'hostLookupOrder(localhost) = dns,files' "$work_dir/old-probe.log"
grep -F 'check failed' "$work_dir/old-probe.log"

helm upgrade frontend podinfo/podinfo --version 5.2.0 \
  -f "$repo_root/deploy/test-apps/podinfo/podinfo-values.yaml" --wait --timeout=120s
kubectl exec deployment/frontend-podinfo -- nslookup -type=A localhost. | tee "$work_dir/dns.log"
grep -F 192.0.2.1 "$work_dir/dns.log"
for endpoint in healthz readyz; do
  kubectl exec deployment/frontend-podinfo -- env GODEBUG=netdns=go+2 \
    podcli check http "localhost:9898/$endpoint"
done
echo 'PASS: unchanged localhost probes succeed while DNS still returns the sinkhole'
