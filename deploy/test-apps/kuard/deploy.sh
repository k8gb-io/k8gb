#!/bin/sh

N=${N:-$1}
N=${N:-2}
[[ ! "$N" =~ ^[0-9]{1,2}$ ]] && echo "error: Not a number: ${N}" >&2 && exit 1
[[ "$DEBUG" == 1 ]] && set -x

for c in $(seq 1 $N); do
  echo "\nCluster ${c}:"
  # deploy and expose kuar app
  kubectl --context k3d-test-gslb$c -n test-gslb create deployment kuar \
    --image=gcr.io/kuar-demo/kuard-amd64:blue \
    --port 8080 \
    --dry-run=client -o yaml | kubectl apply -f -
  
  # add custom dns server
  DNS_IP=$(kubectl --context k3d-test-gslb$c get svc k8gb-coredns -n k8gb -o custom-columns='IP:spec.clusterIP' --no-headers)
  kubectl --context k3d-test-gslb$c -n test-gslb patch deployment kuar \
    -p "{\"spec\":{\"template\":{\"spec\":{\"dnsConfig\":{\"nameservers\":[\"${DNS_IP}\"]},\"dnsPolicy\":\"None\"}}}}"

  # expose the service
  kubectl --context k3d-test-gslb$c -n test-gslb expose deployment kuar \
    --port 8080 \
    --dry-run=client -o yaml | kubectl apply -f -

  # update the daemonset with nginx to use the kuard as the default backend (if no Host header is provided)
  helm --kube-context=k3d-test-gslb$c -n k8gb upgrade -i nginx-ingress nginx-stable/ingress-nginx \
      --version 4.0.15 -f ../../ingress/nginx-ingress-values.yaml \
      --set controller.extraArgs.default-backend-service=test-gslb/kuar \
      --set controller.extraArgs.default-server-port=8080

  # create gslb for the app (ingress will be created)
  kubectl --context k3d-test-gslb$c -n test-gslb apply -f ./kuar_failover.yaml
done

echo "\n\nDone. Continue with opening http://localhost"
