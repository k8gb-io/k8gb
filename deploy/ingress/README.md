# Ingress controller setup

Here we are installing ingress contoller of expected configuration for e2e testing

### Regerence nginx-controller setup
```
$ helm -n ohmyglb upgrade -i nginx-ingress stable/nginx-ingress -f deploy/coredns/nginx-ingress-values.yaml
```
