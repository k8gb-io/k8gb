# Ingress controller, etcd-operator, CoreDNS helm install

Here we are installing resources which are working in tandem but not directly managed by OhMyGLB operator

### Baremetal(or local Minukube/Kind cluster) nginx-controller setup
```
helm -n ohmyglb upgrade -i nginx-ingress stable/nginx-ingress --set controller.service.type=NodePort --set controller.reportNodeInternalIp=true --set controller.hostNetwork=true
```

### Etcd backend for CoreDNS

```
helm -n ohmyglb upgrade -i etcd-for-coredns stable/etcd-operator --set customResources.createEtcdClusterCRD=true
```

### CoreDNS itself

```
helm -n ohmyglb upgrade -i gslb-coredns stable/coredns -f values.yaml
```
