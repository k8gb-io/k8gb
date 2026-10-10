#!/usr/bin/env bash
set -euo pipefail

# These checks match the community agent shipped by Harden Runner v2.21.1.
trap 'echo "::error::Harden Runner is not enforcing its block policy"' ERR
systemctl is-active --quiet agent.service
test "$(sudo cat /home/agent/agent.status)" = Initialized
agent_log=$(sudo cat /home/agent/agent.log)
test -n "$agent_log"
if [[ "$agent_log" == *"Reverted changes"* ]]; then
  echo "::error::Harden Runner reverted its changes"
  exit 1
fi

# Check both runner and container firewall rules, even after Docker/k3d setup.
sudo iptables -C OUTPUT -o eth0 -j REJECT
sudo iptables -C DOCKER-USER -i docker0 -j REJECT

# An unlisted example domain must get the agent's sinkhole, not a DNS failure.
test "$(dig +time=3 +tries=1 +short A example.com)" = 54.185.253.63

# GitHub HTTPS is allowed, but a fresh HTTP connection must be refused (curl 7).
curl --noproxy '*' --connect-timeout 5 --max-time 15 --silent --show-error \
  --fail --output /dev/null https://github.com
if curl --noproxy '*' --connect-timeout 5 --max-time 15 --silent --show-error \
  --output /dev/null http://github.com; then
  echo "::error::Harden Runner allowed an unlisted destination port"
  exit 1
else
  test "$?" -eq 7
fi
echo "Harden Runner is initialized and block enforcement is active"
