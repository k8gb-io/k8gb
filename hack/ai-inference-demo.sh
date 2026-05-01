#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

NAMESPACE="${NAMESPACE:-ai-inference-demo}"
PROBE_NAMESPACE="${PROBE_NAMESPACE:-default}"
MANIFEST="${MANIFEST:-$ROOT_DIR/deploy/ai-inference-demo/ollama.yaml}"
CONTEXTS="${CONTEXTS:-k3d-test-gslb1:eu:kubernetes k3d-test-gslb2:us:provider-fallback}"
PRIMARY_CONTEXT="${PRIMARY_CONTEXT:-k3d-test-gslb1}"
PROBE_CONTEXT="${PROBE_CONTEXT:-$PRIMARY_CONTEXT}"
PRIMARY_GEO_TAG="${PRIMARY_GEO_TAG:-eu}"
GSLB_DOMAIN="${GSLB_DOMAIN:-cloud.example.com}"
AI_INFERENCE_HOST="${AI_INFERENCE_HOST:-ai-inference.$GSLB_DOMAIN}"
DEPLOYMENT_NAME="${DEPLOYMENT_NAME:-ai-inference-demo}"
PROBE_IMAGE="${PROBE_IMAGE:-curlimages/curl:8.18.0}"
OLLAMA_IMAGE="${OLLAMA_IMAGE:-alpine/ollama:latest}"
OLLAMA_BASE_MODEL="${OLLAMA_BASE_MODEL:-qwen2.5:0.5b}"
OLLAMA_MODEL="${OLLAMA_MODEL:-k8gb-resilient-demo:latest}"
OLLAMA_STORAGE_SIZE="${OLLAMA_STORAGE_SIZE:-2Gi}"
OLLAMA_HOST_PATH="${OLLAMA_HOST_PATH:-/var/lib/k8gb-ai-inference-demo/ollama}"
DEFAULT_ROLLOUT_TIMEOUT="600s"
DEFAULT_PROBE_TIMEOUT="60"
DEFAULT_PROBE_ATTEMPTS="1"
DEFAULT_FAILBACK_REPLICAS="1"
DEFAULT_GSLB_READY_TIMEOUT="240"
PROBE_ATTEMPTS="${PROBE_ATTEMPTS:-$DEFAULT_PROBE_ATTEMPTS}"
PROBE_INTERVAL="${PROBE_INTERVAL:-3}"
PROBE_TIMEOUT="${PROBE_TIMEOUT:-$DEFAULT_PROBE_TIMEOUT}"
CONVERGENCE_ATTEMPTS="${CONVERGENCE_ATTEMPTS:-20}"
CONVERGENCE_INTERVAL="${CONVERGENCE_INTERVAL:-5}"
GSLB_READY_TIMEOUT="${GSLB_READY_TIMEOUT:-$DEFAULT_GSLB_READY_TIMEOUT}"
ROLLOUT_TIMEOUT="${ROLLOUT_TIMEOUT:-$DEFAULT_ROLLOUT_TIMEOUT}"
FAILOVER_REPLICAS="${FAILOVER_REPLICAS:-0}"
FAILBACK_REPLICAS="${FAILBACK_REPLICAS:-$DEFAULT_FAILBACK_REPLICAS}"

if [[ -t 1 ]]; then
  COLOR_RESET=$'\033[0m'
  COLOR_BLUE=$'\033[34m'
  COLOR_GREEN=$'\033[32m'
  COLOR_YELLOW=$'\033[33m'
  COLOR_RED=$'\033[31m'
else
  COLOR_RESET=""
  COLOR_BLUE=""
  COLOR_GREEN=""
  COLOR_YELLOW=""
  COLOR_RED=""
fi

usage() {
  cat <<EOF
Usage:
  $(basename "$0") deploy
  $(basename "$0") status
  $(basename "$0") probe
  $(basename "$0") failover
  $(basename "$0") failback
  $(basename "$0") run
  $(basename "$0") logs
  $(basename "$0") delete

Environment overrides:
  CONTEXTS="$CONTEXTS"
    Format: "context:region:backend context:region:backend"
  MANIFEST=$MANIFEST
  PRIMARY_CONTEXT=$PRIMARY_CONTEXT
  PROBE_CONTEXT=$PROBE_CONTEXT
  PRIMARY_GEO_TAG=$PRIMARY_GEO_TAG
  GSLB_DOMAIN=$GSLB_DOMAIN
  AI_INFERENCE_HOST=$AI_INFERENCE_HOST
  NAMESPACE=$NAMESPACE
  PROBE_NAMESPACE=$PROBE_NAMESPACE
  OLLAMA_IMAGE=$OLLAMA_IMAGE
  OLLAMA_BASE_MODEL=$OLLAMA_BASE_MODEL
  OLLAMA_MODEL=$OLLAMA_MODEL
  OLLAMA_STORAGE_SIZE=$OLLAMA_STORAGE_SIZE
  OLLAMA_HOST_PATH=$OLLAMA_HOST_PATH
  PROBE_ATTEMPTS=$PROBE_ATTEMPTS
  PROBE_INTERVAL=$PROBE_INTERVAL
  PROBE_TIMEOUT=$PROBE_TIMEOUT
  CONVERGENCE_ATTEMPTS=$CONVERGENCE_ATTEMPTS
  CONVERGENCE_INTERVAL=$CONVERGENCE_INTERVAL
  GSLB_READY_TIMEOUT=$GSLB_READY_TIMEOUT
EOF
}

timestamp() {
  date +"%H:%M:%S"
}

info() {
  printf "%s[%s]%s %s\n" "$COLOR_BLUE" "$(timestamp)" "$COLOR_RESET" "$1"
}

success() {
  printf "%s[%s]%s %s\n" "$COLOR_GREEN" "$(timestamp)" "$COLOR_RESET" "$1"
}

warn() {
  printf "%s[%s]%s %s\n" "$COLOR_YELLOW" "$(timestamp)" "$COLOR_RESET" "$1"
}

die() {
  printf "%s[%s]%s %s\n" "$COLOR_RED" "$(timestamp)" "$COLOR_RESET" "$1" >&2
  exit 1
}

need_bin() {
  command -v "$1" >/dev/null 2>&1 || die "Missing required binary: $1"
}

validate_config() {
  if [[ -n "${AI_DEMO_MODE:-}" && "${AI_DEMO_MODE}" != "ollama" ]]; then
    die "AI_DEMO_MODE is no longer supported. This demo always runs lightweight Ollama inference."
  fi
}

context_name() {
  printf '%s' "${1%%:*}"
}

context_region() {
  local pair=$1
  pair="${pair#*:}"
  printf '%s' "${pair%%:*}"
}

context_backend() {
  local pair=$1
  pair="${pair#*:}"
  if [[ "$pair" == *:* ]]; then
    printf '%s' "${pair#*:}"
  else
    printf 'inference-endpoint'
  fi
}

render_manifest() {
  local region=$1
  local backend=$2

  sed \
    -e "s/__GSLB_DOMAIN__/$GSLB_DOMAIN/g" \
    -e "s/__AI_DEMO_REGION__/$region/g" \
    -e "s/__AI_DEMO_BACKEND__/$backend/g" \
    -e "s/__AI_DEMO_PRIMARY_GEO_TAG__/$PRIMARY_GEO_TAG/g" \
    -e "s#__OLLAMA_IMAGE__#$OLLAMA_IMAGE#g" \
    -e "s#__OLLAMA_BASE_MODEL__#$OLLAMA_BASE_MODEL#g" \
    -e "s#__OLLAMA_MODEL__#$OLLAMA_MODEL#g" \
    -e "s#__OLLAMA_STORAGE_SIZE__#$OLLAMA_STORAGE_SIZE#g" \
    -e "s#__OLLAMA_HOST_PATH__#$OLLAMA_HOST_PATH#g" \
    "$MANIFEST"
}

wait_for_gslb() {
  local context=$1
  local deadline=$((SECONDS + GSLB_READY_TIMEOUT))
  local hosts healthy_records

  info "Waiting for Gslb status in $context"
  while (( SECONDS < deadline )); do
    hosts="$(kubectl --context "$context" -n "$NAMESPACE" get gslbs.k8gb.io "$DEPLOYMENT_NAME" -o jsonpath='{.status.hosts}' 2>/dev/null || true)"
    healthy_records="$(kubectl --context "$context" -n "$NAMESPACE" get gslbs.k8gb.io "$DEPLOYMENT_NAME" -o jsonpath='{.status.healthyRecords}' 2>/dev/null || true)"
    if [[ -n "$hosts" && -n "$healthy_records" ]]; then
      return 0
    fi
    sleep 3
  done

  kubectl --context "$context" -n "$NAMESPACE" get gslbs.k8gb.io "$DEPLOYMENT_NAME" -o yaml 2>/dev/null || true
  die "Timed out waiting for Gslb status in $context"
}

deploy_demo() {
  need_bin kubectl

  for pair in $CONTEXTS; do
    local context region backend
    context="$(context_name "$pair")"
    region="$(context_region "$pair")"
    backend="$(context_backend "$pair")"

    info "Deploying AI inference demo to $context region=$region backend=$backend"
    render_manifest "$region" "$backend" | kubectl --context "$context" apply -f -
    kubectl --context "$context" -n "$NAMESPACE" rollout status "deployment/$DEPLOYMENT_NAME" --timeout="$ROLLOUT_TIMEOUT"
    wait_for_gslb "$context"
  done

  success "AI inference demo deployed. Probe host: http://$AI_INFERENCE_HOST/v1/chat/completions"
}

delete_demo() {
  need_bin kubectl

  for pair in $CONTEXTS; do
    local context region backend
    context="$(context_name "$pair")"
    region="$(context_region "$pair")"
    backend="$(context_backend "$pair")"

    info "Deleting AI inference demo from $context"
    if kubectl --context "$context" get namespace "$NAMESPACE" >/dev/null 2>&1; then
      kubectl --context "$context" -n "$NAMESPACE" delete \
        gslbs.k8gb.io "$DEPLOYMENT_NAME" \
        ingress "$DEPLOYMENT_NAME" \
        service "$DEPLOYMENT_NAME" \
        deployment "$DEPLOYMENT_NAME" \
        pvc ai-inference-demo-ollama \
        --ignore-not-found
    fi
    kubectl --context "$context" delete pv "ai-inference-demo-ollama-$region" --ignore-not-found
    kubectl --context "$context" delete namespace "$NAMESPACE" --ignore-not-found --wait=false
  done
}

coredns_ip() {
  local context=$1
  kubectl --context "$context" -n k8gb get svc k8gb-coredns -o jsonpath='{.spec.clusterIP}'
}

ingress_ips() {
  local context=$1
  kubectl --context "$context" -n "$NAMESPACE" get ingress "$DEPLOYMENT_NAME" -o jsonpath='{.status.loadBalancer.ingress[*].ip}'
}

secondary_context() {
  local pair context

  for pair in $CONTEXTS; do
    context="$(context_name "$pair")"
    if [[ "$context" != "$PRIMARY_CONTEXT" ]]; then
      printf '%s' "$context"
      return 0
    fi
  done

  die "Could not find a secondary context in CONTEXTS=$CONTEXTS"
}

probe_loop() {
  local context=$1
  local dns_ip=$2
  local expected_remote_ips=${3:-}
  local stop_after_success=${4:-false}
  local pod="k8gb-ai-probe-$RANDOM"
  local rc

  kubectl --context "$context" -n "$PROBE_NAMESPACE" delete pod "$pod" --ignore-not-found >/dev/null 2>&1 || true

  kubectl --context "$context" -n "$PROBE_NAMESPACE" run "$pod" \
    --restart=Never \
    --image="$PROBE_IMAGE" \
    --env="AI_HOST=$AI_INFERENCE_HOST" \
    --env="MODEL=$OLLAMA_MODEL" \
    --env="ATTEMPTS=$PROBE_ATTEMPTS" \
    --env="INTERVAL=$PROBE_INTERVAL" \
    --env="TIMEOUT=$PROBE_TIMEOUT" \
    --env="EXPECT_REMOTE_IPS=$expected_remote_ips" \
    --env="STOP_AFTER_SUCCESS=$stop_after_success" \
    --overrides "{\"spec\":{\"dnsConfig\":{\"nameservers\":[\"$dns_ip\"]},\"dnsPolicy\":\"None\"}}" \
    --command -- sleep 3600 >/dev/null

  if ! kubectl --context "$context" -n "$PROBE_NAMESPACE" wait "pod/$pod" --for=condition=Ready --timeout=60s >/dev/null; then
    kubectl --context "$context" -n "$PROBE_NAMESPACE" logs "$pod" || true
    kubectl --context "$context" -n "$PROBE_NAMESPACE" delete pod "$pod" --ignore-not-found --wait=false >/dev/null
    return 1
  fi

  # Variables inside the single-quoted script expand in the temporary probe pod.
  set +e
  # shellcheck disable=SC2016
  kubectl --context "$context" -n "$PROBE_NAMESPACE" exec "$pod" -- /bin/sh -c '
      attempt=1
      successes=0
      while [ "$attempt" -le "$ATTEMPTS" ]; do
        printf "\nAttempt %s/%s\n" "$attempt" "$ATTEMPTS"
        body="{\"model\":\"$MODEL\",\"stream\":false,\"max_tokens\":64,\"messages\":[{\"role\":\"user\",\"content\":\"Where did this inference request run? Answer using the exact tokens from your system instructions.\"}]}"
        if metrics="$(curl -sS --max-time "$TIMEOUT" -o /tmp/k8gb-ai-response.json -w "http_code=%{http_code} remote_ip=%{remote_ip}" \
          -H "Content-Type: application/json" \
          -d "$body" \
          "http://$AI_HOST/v1/chat/completions")"; then
          printf "%s\n" "$metrics"
          response="$(cat /tmp/k8gb-ai-response.json)"
          content="$(printf "%s" "$response" | sed -n "s/.*\"content\":\"\([^\"]*\)\".*/\1/p")"
          if [ -n "$content" ]; then
            printf "content=%s\n" "$content"
          else
            printf "%s\n" "$response"
          fi
          remote_ip="$(printf "%s" "$metrics" | sed -n "s/.*remote_ip=\([^ ]*\).*/\1/p")"
          route_matches=1
          if [ -n "$EXPECT_REMOTE_IPS" ]; then
            case " $EXPECT_REMOTE_IPS " in
              *" $remote_ip "*) ;;
              *)
                route_matches=0
                printf "route_status=waiting_for_expected_remote_ip expected=\"%s\"\n" "$EXPECT_REMOTE_IPS"
                ;;
            esac
          fi
          case "$metrics" in
            http_code=2*) [ "$route_matches" -eq 1 ] && successes=$((successes + 1)) ;;
          esac
          if [ "$successes" -gt 0 ] && [ "$STOP_AFTER_SUCCESS" = "true" ]; then
            exit 0
          fi
        else
          printf "request failed; k8gb may still be converging\n"
        fi
        if [ "$attempt" -lt "$ATTEMPTS" ]; then
          sleep "$INTERVAL"
        fi
        attempt=$((attempt + 1))
      done
      if [ "$successes" -eq 0 ]; then
        exit 1
      fi
    '
  rc=$?
  set -e

  kubectl --context "$context" -n "$PROBE_NAMESPACE" delete pod "$pod" --ignore-not-found --wait=false >/dev/null
  return "$rc"
}

probe_demo() {
  need_bin kubectl

  local expected_context=${1:-}
  local stop_after_success=${2:-false}
  local dns_ip expected_remote_ips
  dns_ip="$(coredns_ip "$PROBE_CONTEXT")"
  [[ -n "$dns_ip" ]] || die "Could not find k8gb CoreDNS service IP in $PROBE_CONTEXT"
  if [[ -n "$expected_context" ]]; then
    expected_remote_ips="$(ingress_ips "$expected_context")"
    [[ -n "$expected_remote_ips" ]] || die "Could not find ingress IPs for $DEPLOYMENT_NAME in $expected_context"
    info "Expecting $AI_INFERENCE_HOST to route to ingress IPs: $expected_remote_ips"
  fi

  info "Probing $AI_INFERENCE_HOST through k8gb CoreDNS $dns_ip in $PROBE_CONTEXT"
  probe_loop "$PROBE_CONTEXT" "$dns_ip" "$expected_remote_ips" "$stop_after_success"
}

wait_for_convergence() {
  local reason=$1
  local expected_context=$2
  local old_attempts=$PROBE_ATTEMPTS
  local old_interval=$PROBE_INTERVAL

  PROBE_ATTEMPTS=$CONVERGENCE_ATTEMPTS
  PROBE_INTERVAL=$CONVERGENCE_INTERVAL
  info "Waiting for $reason convergence with up to $CONVERGENCE_ATTEMPTS probes"
  probe_demo "$expected_context" true
  PROBE_ATTEMPTS=$old_attempts
  PROBE_INTERVAL=$old_interval
}

scale_primary() {
  local replicas=$1

  need_bin kubectl
  info "Scaling $DEPLOYMENT_NAME in primary context $PRIMARY_CONTEXT to $replicas replicas"
  kubectl --context "$PRIMARY_CONTEXT" -n "$NAMESPACE" scale "deployment/$DEPLOYMENT_NAME" --replicas="$replicas"

  if [[ "$replicas" != "0" ]]; then
    kubectl --context "$PRIMARY_CONTEXT" -n "$NAMESPACE" rollout status "deployment/$DEPLOYMENT_NAME" --timeout="$ROLLOUT_TIMEOUT"
  fi
}

status_demo() {
  need_bin kubectl

  for pair in $CONTEXTS; do
    local context
    context="$(context_name "$pair")"

    printf "\n%s[%s]%s %s\n" "$COLOR_BLUE" "$(timestamp)" "$COLOR_RESET" "$context"
    kubectl --context "$context" -n "$NAMESPACE" get deploy,pod,svc,ingress,gslbs.k8gb.io,pvc 2>/dev/null || true
  done
}

logs_demo() {
  need_bin kubectl

  for pair in $CONTEXTS; do
    local context
    context="$(context_name "$pair")"

    printf "\n%s[%s]%s %s logs\n" "$COLOR_BLUE" "$(timestamp)" "$COLOR_RESET" "$context"
    kubectl --context "$context" -n "$NAMESPACE" logs "deployment/$DEPLOYMENT_NAME" --tail=100 --prefix=true 2>/dev/null || true
  done
}

run_demo() {
  deploy_demo
  status_demo
  wait_for_convergence "primary routing" "$PRIMARY_CONTEXT"
  scale_primary "$FAILOVER_REPLICAS"
  warn "Primary endpoint is down."
  wait_for_convergence "failover" "$(secondary_context)"
  scale_primary "$FAILBACK_REPLICAS"
  warn "Primary endpoint restored."
  wait_for_convergence "failback" "$PRIMARY_CONTEXT"
}

validate_config

case "${1:-}" in
  deploy)
    deploy_demo
    ;;
  status)
    status_demo
    ;;
  probe)
    probe_demo
    ;;
  failover)
    scale_primary "$FAILOVER_REPLICAS"
    ;;
  failback)
    scale_primary "$FAILBACK_REPLICAS"
    ;;
  run)
    run_demo
    ;;
  logs)
    logs_demo
    ;;
  delete)
    delete_demo
    ;;
  -h|--help|help|"")
    usage
    ;;
  *)
    usage
    exit 1
    ;;
esac
