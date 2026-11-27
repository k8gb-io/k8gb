#!/bin/sh

N=${N:-$1}
N=${N:-2}
UPDATE_TRAEFIK=${UPDATE_TRAEFIK:-0}
DIR="${DIR:-$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )}"
[[ ! "$N" =~ ^[0-9]{1,2}$ ]] && echo "error: Not a number: ${N}" >&2 && exit 1
[[ "$DEBUG" == 1 ]] && set -x

for c in $(seq 1 $N); do
  echo "\nCluster ${c}:"
  # deploy kuar app and svc
  kubectl --context k3d-test-gslb$c -n test-gslb apply -f ${DIR}/kuar.yaml
  
  # add custom dns server
  DNS_IP=$(kubectl --context k3d-test-gslb$c get svc k8gb-coredns -n k8gb -o custom-columns='IP:spec.clusterIP' --no-headers)
  kubectl --context k3d-test-gslb$c -n test-gslb patch deployment kuar \
    -p "{\"spec\":{\"template\":{\"spec\":{\"dnsConfig\":{\"nameservers\":[\"${DNS_IP}\"]},\"dnsPolicy\":\"None\"}}}}"

  if [ "$UPDATE_TRAEFIK" == 1 ] ; then
    # Reinstall Traefik in case additional configuration is needed
    helm --kube-context=k3d-test-gslb$c -n k8gb upgrade -i traefik traefik/traefik \
        -f ${DIR}/../../ingress/traefik-ingress-values.yaml \
        --wait --timeout=2m0s
  fi
  
  # create gslb for the app (ingress will be created)
  kubectl --context k3d-test-gslb$c -n test-gslb apply -f ${DIR}/kuar_failover.yaml
done

echo "\n\nDone. Continue with opening http://localhost\n\n"
