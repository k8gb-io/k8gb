# External DNS behind a proxy

External DNS needs to communicate with a DNS server outside of the Kubernetes cluster to update records. If a proxy is used for egress from the Kubernetes cluster, configure the embedded ExternalDNS subchart under the `extdns` key (the chart alias for `external-dns`):

```yaml
extdns:
  env:
  - name: HTTPS_PROXY
    value: http://proxy.example.com:8080
  extraVolumes:
  - name: ca-bundle
    secret:
      secretName: ca-proxy
  extraVolumeMounts:
  - name: ca-bundle
    mountPath: /etc/ssl/certs
    readOnly: true
```

The `HTTPS_PROXY` environment variable should contain the address of the proxy.
The volume mount should contain the proxy CA certificate so that the container can trust the proxy.

> Use `extdns.env` (not `extraEnv`) — that matches the upstream [external-dns Helm chart](https://github.com/kubernetes-sigs/external-dns/tree/master/charts/external-dns) values.
