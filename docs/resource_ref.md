# GSLB ResourceRef Support
Starting from v0.15.0, k8gb introduces a much simpler way to link a GSLB resource to an Ingress object in Kubernetes. You no longer need to duplicate the Ingress configuration in your GSLB definition—instead, you can simply reference an existing Ingress. 
This makes your Ingress the single source of truth for application routing.

K8GB supports the following ingress resources:

* [Kubernetes Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/)
* [Kubernetes LoadBalancer Service](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer)
* [Istio Virtual Service](https://istio.io/latest/docs/reference/config/networking/virtual-service/)
* Gateway API Resources(to be released in **v0.17.0**):
  - [HTTPRoute](https://gateway-api.sigs.k8s.io/api-types/httproute/)
  - [GRPCRoute](https://gateway-api.sigs.k8s.io/api-types/grpcroute/)
  - [TCPRoute](https://gateway-api.sigs.k8s.io/guides/tcp/)
  - [UDPRoute](https://gateway-api.sigs.k8s.io/reference/1.4/spec/?h=udproute#udproute)
  - [TLSRoute](https://gateway-api.sigs.k8s.io/geps/gep-2643/?h=tls#tlsroute-tls-passthrough)

## 1. Declaration by Name
The simplest way is to directly specify the name of the resource you want to reference in your GSLB. The namespace will be automatically taken from the GSLB’s namespace.

Ingress:
```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: playground-failover
  namespace: playground
spec:
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    name: playground-failover-ingress
```

LoadBalancer Service:
```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: playground-failover
  namespace: playground
  annotations:
    k8gb.io/hostname: "myapp.example.com"
spec:
  resourceRef:
    apiVersion: v1
    kind: Service
    name: playground-failover-lbservice
```

Istio Virtual Service:
```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: playground-failover
  namespace: playground
spec:
  resourceRef:
    apiVersion: networking.istio.io/v1
    kind: VirtualService
    name: playground-failover-virtualservice
```

GatewayAPI HTTPRoute:
```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: playground-failover
  namespace: playground
spec:
  resourceRef:
    apiVersion: gateway.networking.k8s.io/v1
    kind: HTTPRoute
    name: playground-failover-httproute
```

GatewayAPI GRPCRoute:
```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: playground-failover
  namespace: playground
spec:
  resourceRef:
    apiVersion: gateway.networking.k8s.io/v1
    kind: GRPCRoute
    name: playground-failover-grpcroute
```

GatewayAPI TCPRoute:
```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: failover-tcproute
  namespace: playground
  annotations:
    k8gb.io/hostname: gatewayapi-tcproute.cloud.example.com
spec:
  resourceRef:
    apiVersion: gateway.networking.k8s.io/v1alpha2
    kind: TCPRoute
    name: failover-tcproute
```

GatewayAPI UDPRoute:
```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: failover-udproute
  namespace: playground
  annotations:
    k8gb.io/hostname: gatewayapi-udproute.cloud.example.com
spec:
  resourceRef:
    apiVersion: gateway.networking.k8s.io/v1alpha2
    kind: UDPRoute
    name: failover-udproute
```

GatewayAPI TLSRoute:
```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: failover-tlsroute
  namespace: playground
spec:
  resourceRef:
    apiVersion: gateway.networking.k8s.io/v1alpha3
    kind: TLSRoute
    name: failover-tlsroute
```


## 2. Declaration by Label
Alternatively, you can reference the ingress resource by label. This approach is useful when you need more flexibility—for example, 
in CI/CD pipelines. It is required that only one resource in the namespace matches the label; otherwise, k8gb will return 
an error.

Here we show only an example for Ingress resources, but the same applies for Istio and GatewayAPI integrations.

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
