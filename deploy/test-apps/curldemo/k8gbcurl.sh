#!/bin/sh

# Small script to continously poll demo podinfo fqdn to show k8gb in action
# $1 - nameserver to use usually k8gb-coredns service ip
# $2 - test fqdn to resolve in demo loops

# checks
[[ $# != 2 ]] && echo "Usage: $0 <nameserver> <url>" && exit 1

DELAY="${DELAY:-5}"

if [ "$DEBUG" == 1 ]
then
   set -x
fi

if [ "$1" != '--local' ]
then
    echo "nameserver $1" > /etc/resolv.conf
fi

url="$2"

while true
do
  curl -k -s -w "%{stderr}\n%{http_code}\n" --location --request GET "${url}" | grep message
  sleep ${DELAY}
  echo "[`date`] ..."
done
