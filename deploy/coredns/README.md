# CoreDNS helm install to act as resolver for GSLB

```
helm -n ohmyglb upgrade -i gslb-coredns stable/coredns -f values.yaml
```
