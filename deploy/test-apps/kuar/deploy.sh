#!/bin/sh

N=${N:-$1}
N=${N:-2}
UPDATE_NGINX=${UPDATE_NGINX:-1}
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

  if [ "$UPDATE_NGINX" == 1 ] ; then
    # update the daemonset with nginx to use the kuar as the default backend (if no Host header is provided)
    helm --kube-context=k3d-test-gslb$c -n k8gb upgrade -i nginx-ingress nginx-stable/ingress-nginx \
        --version 4.0.15 -f ${DIR}/../../ingress/nginx-ingress-values.yaml \
        --set controller.extraArgs.default-backend-service=test-gslb/kuar \
        --set controller.extraArgs.default-server-port=8080 \
        --wait --timeout=2m0s
  fi
  
  # create gslb for the app (ingress will be created)
  kubectl --context k3d-test-gslb$c -n test-gslb apply -f ${DIR}/kuar_failover.yaml
done

echo "\n\nDone. Continue with opening http://localhost\n\n"
