# Proxy External DNS

External DNS needs to communicate with a DNS server outside of the kubernetes cluster to update records. If there is a proxy for egress of the Kubernetes cluster the following should be configured:
```
externaldns:
  extraEnv:
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
