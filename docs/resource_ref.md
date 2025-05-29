# GSLB ResourceRef Support
Starting from v0.15.0, k8gb introduces a much simpler way to link a GSLB resource to an Ingress object in Kubernetes. You no longer need to duplicate the Ingress configuration in your GSLB definition—instead, you can simply reference an existing Ingress. 
This makes your Ingress the single source of truth for application routing.

## 1. Declaration by Name
The simplest way is to directly specify the name of the Ingress you want to reference in your GSLB. The namespace will be automatically taken from the GSLB’s namespace.

```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: playground-failover
spec:
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    name: playground-failover-ingress
```

## 2. Declaration by Label
Alternatively, you can reference the Ingress by label. This approach is useful when you need more flexibility—for example, 
in CI/CD pipelines. It is required that only one Ingress in the namespace matches the label; otherwise, k8gb will return 
an error.

```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: playground-failover
spec:
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    matchLabels:
      app: playground-failover
```

## 3. Embedded Declaration (Legacy)
For backward compatibility, you can still use the original way where the Ingress configuration is embedded directly inside the GSLB resource. This method will continue to work, but we recommend switching to reference-based configuration for simpler management and to avoid configuration drift.

```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: failover-playground-embedded
spec:
  ingress:
    ingressClassName: nginx
    rules:
      - host: failover-playground-embedded.cloud.example.com
        http:
          paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: frontend-podinfo
                port:
                  name: http
  strategy:
    type: failover
    dnsTtlSeconds: 5
    primaryGeoTag: "eu"
```

> **Note:**
If the Ingress is created automatically by a GSLB resource, in addition to an ownerReference, it will also be marked with the label:
`app.k8gb.io/managed-by: gslb`. This makes it easy to identify Ingresses managed by k8gb.
