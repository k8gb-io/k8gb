#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)

# Keep the checks executable on a developer machine without sudo or an agent.
systemctl() {
  [[ "$*" == 'is-active --quiet agent.service' && "$scenario" != stopped ]]
}
sudo() {
  case "$*" in
    'cat /home/agent/agent.status')
      [[ "$scenario" != missing-status ]] || return 1
      if [[ "$scenario" == uninitialized ]]; then
        echo 'Error in Initialization'
      else
        echo Initialized
      fi
      ;;
    'cat /home/agent/agent.log')
      case "$scenario" in
        missing-log) return 1 ;;
        empty-log) return 0 ;;
        reverted) echo 'Reverted changes' ;;
        *) echo 'done' ;;
      esac
      ;;
    'iptables -C OUTPUT -o eth0 -j REJECT')
      [[ "$scenario" != runner-firewall-missing ]]
      ;;
    'iptables -C DOCKER-USER -i docker0 -j REJECT')
      [[ "$scenario" != container-firewall-missing ]]
      ;;
    *) return 2 ;;
  esac
}
dig() {
  [[ "$*" == '+time=3 +tries=1 +short A example.com' ]] || return 2
  case "$scenario" in
    dns-error) return 1 ;;
    dns-empty) return 0 ;;
    dns-allowed) echo 192.0.2.1 ;;
    *) echo 54.185.253.63 ;;
  esac
}
curl() {
  case "${*: -1}" in
    https://github.com) [[ "$scenario" != allowed-unreachable ]] ;;
    http://github.com)
      case "$scenario" in
        http-allowed) return 0 ;;
        http-dns-error) return 6 ;;
        http-timeout) return 28 ;;
        *) return 7 ;;
      esac
      ;;
    *) return 2 ;;
  esac
}
export -f systemctl sudo dig curl

export scenario=healthy
bash "$script_dir/verify-harden-runner.sh"
for scenario in stopped missing-status uninitialized missing-log empty-log reverted \
  runner-firewall-missing container-firewall-missing dns-error dns-empty dns-allowed \
  allowed-unreachable http-allowed http-dns-error http-timeout; do
  if bash "$script_dir/verify-harden-runner.sh" >/dev/null 2>&1; then
    echo "FAIL: accepted $scenario"
    exit 1
  fi
  echo "PASS: rejected $scenario"
done
