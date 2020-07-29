#!/bin/sh

# Small script to continously poll demo podinfo fqdn to show k8gb in action
# $1 - nameserver to use usually k8gb-coredns service ip
# $2 - test fqdn to resolve in demo loops

echo "nameserver $1" > /etc/resolv.conf
while true
do
  curl -s -w "%{stderr}\n%{http_code}\n" --location --request GET "$2" |grep message
  sleep 5
  echo "[`date`] ..."
done
