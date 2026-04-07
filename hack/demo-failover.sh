#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

NAMESPACE="${NAMESPACE:-test-gslb}"
PRIMARY_CONTEXT="${PRIMARY_CONTEXT:-k3d-test-gslb1}"
SECONDARY_CONTEXT="${SECONDARY_CONTEXT:-k3d-test-gslb2}"
PRIMARY_GEO_TAG="${PRIMARY_GEO_TAG:-eu}"
SECONDARY_GEO_TAG="${SECONDARY_GEO_TAG:-us}"
PROBE_CONTEXT="${PROBE_CONTEXT:-$PRIMARY_CONTEXT}"
GSLB_NAME="${GSLB_NAME:-failover-ingress}"
DEPLOYMENT_NAME="${DEPLOYMENT_NAME:-frontend-podinfo}"
FAILOVER_HOST="${FAILOVER_HOST:-}"
FAILOVER_MANIFEST="${FAILOVER_MANIFEST:-$ROOT_DIR/deploy/gslb/k8gb.absa.oss_v1beta1_gslb_cr_failover_ingress_ref.yaml}"
ROLLOUT_TIMEOUT="${ROLLOUT_TIMEOUT:-180s}"
PROBE_COUNT="${PROBE_COUNT:-3}"
PROBE_INTERVAL="${PROBE_INTERVAL:-1}"
TRANSITION_ATTEMPTS="${TRANSITION_ATTEMPTS:-20}"
TRANSITION_INTERVAL="${TRANSITION_INTERVAL:-3}"
GSLB_RESOURCE="gslbs.k8gb.io"
INGRESS_RESOURCE="ingress.networking.k8s.io"
DEFAULT_FAILOVER_HOST="failover.example.com"

if [[ -t 1 ]]; then
  COLOR_BOLD=$'\033[1m'
  COLOR_RESET=$'\033[0m'
  COLOR_DIM=$'\033[2m'
  COLOR_BLUE=$'\033[34m'
  COLOR_CYAN=$'\033[36m'
  COLOR_GREEN=$'\033[32m'
  COLOR_YELLOW=$'\033[33m'
  COLOR_RED=$'\033[31m'
else
  COLOR_BOLD=""
  COLOR_RESET=""
  COLOR_DIM=""
  COLOR_BLUE=""
  COLOR_CYAN=""
  COLOR_GREEN=""
  COLOR_YELLOW=""
  COLOR_RED=""
fi

usage() {
  cat <<EOF
Usage:
  $(basename "$0") status [--yaml]
  $(basename "$0") probe
  $(basename "$0") failover
  $(basename "$0") failback [replicas]
  $(basename "$0") 0|1|2|3|4|5|6
  $(basename "$0") all [--init] [--yaml]
  $(basename "$0") run [--init] [--yaml]

Environment overrides:
  NAMESPACE=$NAMESPACE
  PRIMARY_CONTEXT=$PRIMARY_CONTEXT
  SECONDARY_CONTEXT=$SECONDARY_CONTEXT
  PRIMARY_GEO_TAG=$PRIMARY_GEO_TAG
  SECONDARY_GEO_TAG=$SECONDARY_GEO_TAG
  PROBE_CONTEXT=$PROBE_CONTEXT
  GSLB_NAME=$GSLB_NAME
  DEPLOYMENT_NAME=$DEPLOYMENT_NAME
  FAILOVER_HOST=${FAILOVER_HOST:-<auto>}
  PROBE_COUNT=$PROBE_COUNT
  PROBE_INTERVAL=$PROBE_INTERVAL
  TRANSITION_ATTEMPTS=$TRANSITION_ATTEMPTS
  TRANSITION_INTERVAL=$TRANSITION_INTERVAL
EOF
}

timestamp() {
  date +"%H:%M:%S"
}

header() {
  printf "\n%s=== %s ===%s\n" "$COLOR_BLUE" "$1" "$COLOR_RESET"
}

section() {
  printf "\n%s[%s]%s %s\n" "$COLOR_BLUE" "$(timestamp)" "$COLOR_RESET" "$1"
}

info() {
  printf "%s[%s]%s %s\n" "$COLOR_DIM" "$(timestamp)" "$COLOR_RESET" "$1"
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

pause_demo() {
  if [[ "${DEMO_NO_PAUSE:-0}" == "1" ]]; then
    return 0
  fi

  printf "\n%s[Press ENTER to continue]%s\n" "$COLOR_YELLOW" "$COLOR_RESET"
  read -r
}

need_bin() {
  command -v "$1" >/dev/null 2>&1 || die "Missing required binary: $1"
}

have_bin() {
  command -v "$1" >/dev/null 2>&1
}

run_cmd() {
  info "\$ $*"
  "$@"
}

kubectl_named_jsonpath() {
  local context=$1
  local resource=$2
  local name=$3
  local jsonpath=$4

  kubectl --context "$context" -n "$NAMESPACE" get "$resource" "$name" -o "jsonpath=${jsonpath}" 2>/dev/null || true
}

resource_exists() {
  local context=$1
  local resource=$2
  local name=$3

  kubectl --context "$context" -n "$NAMESPACE" get "$resource" "$name" >/dev/null 2>&1
}

deployment_replicas() {
  local context=$1
  kubectl_named_jsonpath "$context" "deployment" "$DEPLOYMENT_NAME" '{.spec.replicas}'
}

deployment_ready() {
  local context=$1
  kubectl_named_jsonpath "$context" "deployment" "$DEPLOYMENT_NAME" '{.status.readyReplicas}'
}

gslb_summary_value() {
  local context=$1
  local jsonpath=$2
  kubectl_named_jsonpath "$context" "$GSLB_RESOURCE" "$GSLB_NAME" "$jsonpath"
}

ingress_summary_value() {
  local context=$1
  local jsonpath=$2
  kubectl_named_jsonpath "$context" "$INGRESS_RESOURCE" "$GSLB_NAME" "$jsonpath"
}

first_csv_item() {
  local value=${1:-}
  value="${value%%,*}"
  printf '%s' "$value"
}

resolve_failover_host() {
  local context=${1:-$PROBE_CONTEXT}
  local host

  if [[ -n "$FAILOVER_HOST" ]]; then
    printf '%s' "$FAILOVER_HOST"
    return 0
  fi

  host="$(first_csv_item "$(gslb_summary_value "$context" '{.status.hosts}')")"
  if [[ -n "$host" ]]; then
    printf '%s' "$host"
    return 0
  fi

  host="$(ingress_summary_value "$context" '{.spec.rules[0].host}')"
  if [[ -n "$host" ]]; then
    printf '%s' "$host"
    return 0
  fi

  printf '%s' "$DEFAULT_FAILOVER_HOST"
}

show_cluster_status() {
  local context=$1
  local label=$2
  local deploy_state strategy primary status_geo hosts role
  local replicas ready

  if resource_exists "$context" "deployment" "$DEPLOYMENT_NAME"; then
    replicas="$(deployment_replicas "$context")"
    ready="$(deployment_ready "$context")"
    deploy_state="${ready:-0}/${replicas:-0}"
  else
    deploy_state="missing"
  fi

  if resource_exists "$context" "$GSLB_RESOURCE" "$GSLB_NAME"; then
    strategy="$(gslb_summary_value "$context" '{.spec.strategy.type}')"
    primary="$(gslb_summary_value "$context" '{.spec.strategy.primaryGeoTag}')"
    status_geo="$(gslb_summary_value "$context" '{.status.geoTag}')"
    hosts="$(gslb_summary_value "$context" '{.status.hosts}')"
  elif resource_exists "$context" "$INGRESS_RESOURCE" "$GSLB_NAME"; then
    strategy="missing"
    primary="missing"
    status_geo="missing"
    hosts="$(ingress_summary_value "$context" '{.spec.rules[0].host}')"
  else
    strategy="missing"
    primary="missing"
    status_geo="missing"
    hosts="missing"
  fi

  if [[ "$status_geo" != "missing" && "$primary" != "missing" ]]; then
    if [[ "$status_geo" == "$primary" ]]; then
      role="primary"
    else
      role="secondary"
    fi
  else
    role="missing"
  fi

  printf "%-7s  context=%-17s deploy=%-8s  strategy=%-8s role=%-9s geoTag=%-7s host=%s\n" \
    "$label" \
    "$context" \
    "$deploy_state" \
    "${strategy:--}" \
    "${role:--}" \
    "${status_geo:--}" \
    "${hosts:--}"
}

show_gslb_yaml() {
  local context=$1
  local manifest

  if ! manifest="$(kubectl --context "$context" -n "$NAMESPACE" get "$GSLB_RESOURCE" "$GSLB_NAME" -o yaml 2>/dev/null)"; then
    warn "Unable to fetch ${GSLB_RESOURCE}/${GSLB_NAME} from ${context}"
    return 1
  fi

  if have_bin yq; then
    printf '%s\n' "$manifest" \
      | yq eval '{"metadata": {"name": .metadata.name, "namespace": .metadata.namespace}, "spec": {"strategy": .spec.strategy, "resourceRef": .spec.resourceRef}, "status": {"geoTag": .status.geoTag, "hosts": .status.hosts, "serviceHealth": .status.serviceHealth, "healthyRecords": .status.healthyRecords}}' -
  else
    printf '%s\n' "$manifest"
  fi
}

show_ingress_yaml() {
  local context=$1
  local manifest

  if ! manifest="$(kubectl --context "$context" -n "$NAMESPACE" get "$INGRESS_RESOURCE" "$GSLB_NAME" -o yaml 2>/dev/null)"; then
    warn "Unable to fetch ${INGRESS_RESOURCE}/${GSLB_NAME} from ${context}"
    return 1
  fi

  if have_bin yq; then
    printf '%s\n' "$manifest" \
      | yq eval '{"metadata": {"name": .metadata.name, "namespace": .metadata.namespace}, "spec": {"ingressClassName": .spec.ingressClassName, "rules": .spec.rules}}' -
  else
    printf '%s\n' "$manifest"
  fi
}

show_manual_steps() {
  local host

  host="$(resolve_failover_host "$PROBE_CONTEXT")"

  section "Manual walkthrough"
  cat <<EOF
1. Inspect both clusters:
   kubectl --context ${PRIMARY_CONTEXT} -n ${NAMESPACE} get ${GSLB_RESOURCE} ${GSLB_NAME} -o yaml | yq
   kubectl --context ${SECONDARY_CONTEXT} -n ${NAMESPACE} get ${GSLB_RESOURCE} ${GSLB_NAME} -o yaml | yq
2. Probe the failover host:
   kubectl --context ${PROBE_CONTEXT} -n ${NAMESPACE} run -it --rm failover-probe --restart=Never --image=busybox --overrides '{"spec":{"dnsConfig":{"nameservers":["<k8gb-coredns-ip>"]},"dnsPolicy":"None"}}' -- wget -qO- http://${host}
3. Trigger failover:
   kubectl --context ${PRIMARY_CONTEXT} -n ${NAMESPACE} scale deployment/${DEPLOYMENT_NAME} --replicas=0
4. Restore primary:
   kubectl --context ${PRIMARY_CONTEXT} -n ${NAMESPACE} scale deployment/${DEPLOYMENT_NAME} --replicas=$(determine_restore_replicas)
EOF
}

print_recap_diagram() {
  local host primary_replicas secondary_replicas

  host="$(resolve_failover_host "$PROBE_CONTEXT")"
  primary_replicas="$(deployment_ready "$PRIMARY_CONTEXT")"
  secondary_replicas="$(deployment_ready "$SECONDARY_CONTEXT")"

  primary_replicas="${primary_replicas:-0}"
  secondary_replicas="${secondary_replicas:-0}"

  cat <<EOF
${COLOR_BOLD}${COLOR_CYAN}                    FAILOVER RECAP${COLOR_RESET}

  host: ${COLOR_BOLD}${host}${COLOR_RESET}

  ${COLOR_BOLD}Act 1${COLOR_RESET}  normal traffic
      ${COLOR_CYAN}client${COLOR_RESET}
        ${COLOR_DIM}|${COLOR_RESET}
        ${COLOR_DIM}v${COLOR_RESET}
   ${COLOR_BOLD}${host}${COLOR_RESET}
        ${COLOR_DIM}|${COLOR_RESET}
        ${COLOR_GREEN}+----------------------------> [${PRIMARY_GEO_TAG}] primary${COLOR_RESET}   (${primary_replicas} ready)
        ${COLOR_DIM}|${COLOR_RESET}
        ${COLOR_DIM}'---------------------------->${COLOR_RESET} [${SECONDARY_GEO_TAG}] secondary (${secondary_replicas} ready, idle)

  ${COLOR_BOLD}Act 2-3${COLOR_RESET}  primary scaled to zero
      ${COLOR_CYAN}client${COLOR_RESET}
        ${COLOR_DIM}|${COLOR_RESET}
        ${COLOR_DIM}v${COLOR_RESET}
   ${COLOR_BOLD}${host}${COLOR_RESET}
        ${COLOR_DIM}|${COLOR_RESET}
        ${COLOR_RED}x----------------------------> [${PRIMARY_GEO_TAG}] primary${COLOR_RESET}   (0 ready)
        ${COLOR_DIM}|${COLOR_RESET}
        ${COLOR_YELLOW}'===========================> [${SECONDARY_GEO_TAG}] serving failover traffic${COLOR_RESET}

  ${COLOR_BOLD}Act 4-5${COLOR_RESET}  primary restored
      ${COLOR_CYAN}client${COLOR_RESET}
        ${COLOR_DIM}|${COLOR_RESET}
        ${COLOR_DIM}v${COLOR_RESET}
   ${COLOR_BOLD}${host}${COLOR_RESET}
        ${COLOR_DIM}|${COLOR_RESET}
        ${COLOR_GREEN}+----------------------------> [${PRIMARY_GEO_TAG}] primary${COLOR_RESET}   (traffic restored)
        ${COLOR_DIM}|${COLOR_RESET}
        ${COLOR_DIM}'---------------------------->${COLOR_RESET} [${SECONDARY_GEO_TAG}] secondary (standby)
EOF
}

show_status() {
  local show_yaml=${1:-0}

  section "Current Global State"
  show_cluster_status "$PRIMARY_CONTEXT" "$PRIMARY_GEO_TAG"
  show_cluster_status "$SECONDARY_CONTEXT" "$SECONDARY_GEO_TAG"
  info "Probe host: $(resolve_failover_host "$PROBE_CONTEXT")"

  if [[ "$show_yaml" == "1" ]]; then
    section "Primary GSLB excerpt"
    show_gslb_yaml "$PRIMARY_CONTEXT"
    section "Secondary GSLB excerpt"
    show_gslb_yaml "$SECONDARY_CONTEXT"
  fi
}

probe_overrides() {
  local context=$1
  local nameserver

  nameserver="$(kubectl --context "$context" -n k8gb get svc k8gb-coredns -o jsonpath='{.spec.clusterIP}')"
  printf '{"spec":{"dnsConfig":{"nameservers":["%s"]},"dnsPolicy":"None"}}' "$nameserver"
}

extract_json_field() {
  local field=$1
  local payload=$2

  printf '%s\n' "$payload" | sed -n "s/.*\"${field}\":[[:space:]]*\"\([^\"]*\)\".*/\1/p" | head -n1
}

normalize_region() {
  local message=${1:-}
  local normalized

  normalized="$(printf '%s' "$message" | sed -E 's/^[[:space:]]+//; s/[[:space:]]+$//; s/^.*:[[:space:]]*//')"
  case "$normalized" in
    eu|us)
      printf '%s' "$normalized"
      ;;
    *)
      printf '%s' "$message"
      ;;
  esac
}

probe_once() {
  local context=$1
  local name response host

  host="$(resolve_failover_host "$context")"

  name="failover-probe-$(date +%s)-$RANDOM"
  response="$(
    kubectl --context "$context" -n "$NAMESPACE" run "$name" \
      --attach \
      --quiet \
      --rm \
      --restart=Never \
      --image=busybox \
      --overrides "$(probe_overrides "$context")" \
      --command -- wget -qO- "http://${host}" 2>/dev/null || true
  )"

  printf '%s' "$response"
}

probe_loop() {
  local count=${1:-$PROBE_COUNT}
  local interval=${2:-$PROBE_INTERVAL}
  local context=${3:-$PROBE_CONTEXT}
  local i response message region host

  host="$(resolve_failover_host "$context")"

  section "Probing ${host} from ${context}"
  for ((i = 1; i <= count; i++)); do
    response="$(probe_once "$context")"
    message="$(extract_json_field "message" "$response")"
    region="$(normalize_region "$message")"

    if [[ -n "$message" ]]; then
      if [[ -n "$message" && "$region" != "$message" ]]; then
        printf "%2d/%d  region=%-3s raw=%s\n" "$i" "$count" "${region:--}" "$message"
      else
        printf "%2d/%d  region=%-3s\n" "$i" "$count" "${region:--}"
      fi
    else
      printf "%2d/%d  raw=%s\n" "$i" "$count" "${response:-<empty>}"
    fi

    if [[ "$i" -lt "$count" ]]; then
      sleep "$interval"
    fi
  done
}

wait_for_replicas() {
  local context=$1
  local expected=$2
  local attempts=30
  local sleep_seconds=2
  local current

  for ((i = 1; i <= attempts; i++)); do
    current="$(deployment_ready "$context")"
    current="${current:-0}"
    if [[ "$current" == "$expected" ]]; then
      success "Deployment ${DEPLOYMENT_NAME} on ${context} is now ${current}/${expected} ready"
      return 0
    fi
    info "Waiting for ${DEPLOYMENT_NAME} on ${context}: ready=${current}, target=${expected}"
    sleep "$sleep_seconds"
  done

  warn "Timed out waiting for ${DEPLOYMENT_NAME} on ${context} to reach ${expected} ready replicas"
  return 1
}

wait_for_message() {
  local expected=$1
  local context=${2:-$PROBE_CONTEXT}
  local attempts=${3:-$TRANSITION_ATTEMPTS}
  local interval=${4:-$TRANSITION_INTERVAL}
  local response message region

  section "Waiting for traffic to resolve to ${expected}"
  for ((i = 1; i <= attempts; i++)); do
    response="$(probe_once "$context")"
    message="$(extract_json_field "message" "$response")"
    region="$(normalize_region "$message")"
    if [[ -n "$message" && "$region" != "$message" ]]; then
      info "Probe ${i}/${attempts}: region=${region:--} raw=${message}"
    else
      info "Probe ${i}/${attempts}: region=${region:--}"
    fi
    if [[ "$region" == "$expected" ]]; then
      success "Observed failover response from ${expected}"
      return 0
    fi
    sleep "$interval"
  done

  warn "Did not observe ${expected} within ${attempts} probes"
  return 1
}

scale_primary() {
  local replicas=$1

  section "Scaling ${DEPLOYMENT_NAME} on ${PRIMARY_CONTEXT} to ${replicas}"
  run_cmd kubectl --context "$PRIMARY_CONTEXT" -n "$NAMESPACE" scale "deployment/${DEPLOYMENT_NAME}" --replicas "$replicas"
  if [[ "$replicas" -gt 0 ]]; then
    run_cmd kubectl --context "$PRIMARY_CONTEXT" -n "$NAMESPACE" rollout status "deployment/${DEPLOYMENT_NAME}" --timeout "$ROLLOUT_TIMEOUT"
  fi
  wait_for_replicas "$PRIMARY_CONTEXT" "$replicas" || true
}

apply_failover_manifest() {
  section "Applying failover manifest to both clusters"
  run_cmd kubectl --context "$SECONDARY_CONTEXT" apply -f "$FAILOVER_MANIFEST"
  run_cmd kubectl --context "$PRIMARY_CONTEXT" apply -f "$FAILOVER_MANIFEST"
}

determine_restore_replicas() {
  local replicas

  replicas="$(deployment_replicas "$PRIMARY_CONTEXT")"
  replicas="${replicas:-1}"
  if [[ "$replicas" -lt 1 ]]; then
    replicas=1
  fi

  printf '%s' "$replicas"
}

phase0() {
  local show_yaml=${1:-0}

  header "Act 0: Current Global State"
  show_status "$show_yaml"

  section "Primary GSLB excerpt"
  show_gslb_yaml "$PRIMARY_CONTEXT"

  section "Primary Ingress excerpt"
  show_ingress_yaml "$PRIMARY_CONTEXT"

  show_manual_steps
}

phase1() {
  header "Act 1: Primary Serving Traffic"
  probe_loop
}

phase2() {
  header "Act 2: Trigger Failover"
  scale_primary 0
  show_status "${1:-0}"
}

phase3() {
  header "Act 3: Observe Secondary Takeover"
  wait_for_message "$SECONDARY_GEO_TAG" "$PROBE_CONTEXT" || true
  probe_loop
}

phase4() {
  local replicas=${1:-1}
  local show_yaml=${2:-0}

  header "Act 4: Restore Primary"
  scale_primary "$replicas"
  show_status "$show_yaml"
}

phase5() {
  header "Act 5: Observe Failback"
  wait_for_message "$PRIMARY_GEO_TAG" "$PROBE_CONTEXT" || true
  probe_loop
}

phase6() {
  header "Act 6: Recap"
  show_status
  section "ASCII recap"
  print_recap_diagram
}

run_demo() {
  local init=${1:-0}
  local show_yaml=${2:-0}
  local original_replicas

  original_replicas="$(determine_restore_replicas)"

  if [[ "$init" == "1" ]]; then
    apply_failover_manifest
    scale_primary "$original_replicas"
  fi

  phase0 "$show_yaml"
  pause_demo
  phase1
  pause_demo
  phase2 "$show_yaml"
  pause_demo
  phase3
  pause_demo
  phase4 "$original_replicas" "$show_yaml"
  pause_demo
  phase5
  pause_demo
  phase6
}

main() {
  local command=all
  local show_yaml=0
  local init=0

  need_bin kubectl

  if [[ $# -gt 0 ]]; then
    case "$1" in
      -h|--help)
        usage
        exit 0
        ;;
      *)
        command=$1
        shift
        ;;
    esac
  fi

  while [[ $# -gt 0 ]]; do
    case "$1" in
      --yaml)
        show_yaml=1
        ;;
      --init)
        init=1
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        break
        ;;
    esac
    shift
  done

  case "$command" in
    0)
      phase0 "$show_yaml"
      ;;
    1)
      phase1
      ;;
    2)
      phase2 "$show_yaml"
      ;;
    3)
      phase3
      ;;
    4)
      phase4 "${1:-1}" "$show_yaml"
      ;;
    5)
      phase5
      ;;
    6)
      phase6
      ;;
    all)
      run_demo "$init" "$show_yaml"
      ;;
    status)
      show_status "$show_yaml"
      ;;
    probe)
      probe_loop
      ;;
    failover)
      scale_primary 0
      show_status "$show_yaml"
      wait_for_message "$SECONDARY_GEO_TAG" "$PROBE_CONTEXT" || true
      ;;
    failback)
      scale_primary "${1:-1}"
      show_status "$show_yaml"
      wait_for_message "$PRIMARY_GEO_TAG" "$PROBE_CONTEXT" || true
      ;;
    run)
      DEMO_NO_PAUSE=1 run_demo "$init" "$show_yaml"
      ;;
    *)
      usage
      die "Unknown command: ${command}"
      ;;
  esac
}

main "$@"
