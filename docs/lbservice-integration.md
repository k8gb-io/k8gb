# LoadBalancer Service Integration

This document describes how to use k8gb GSLB with Kubernetes Services of type `LoadBalancer`.

## Overview

Unlike Ingress resources which contain hostname information in their `spec.rules[].host` fields, LoadBalancer services don't have built-in hostname information. Therefore, when using LoadBalancer services with k8gb GSLB, you must explicitly specify the desired hostname via annotation.

## Usage

### 1. Create a LoadBalancer Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
  namespace: default
  labels:
    app: my-app
spec:
  type: LoadBalancer
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
      name: http
  selector:
    app: my-app
```

### 2. Create a GSLB Resource

```yaml
apiVersion: k8gb.io/v1beta1
kind: Gslb
metadata:
  name: my-app-gslb
  namespace: default
  annotations:
    k8gb.io/hostname: "myapp.example.com"  # Required: Specify your desired hostname
spec:
  resourceRef:
    apiVersion: v1
    kind: Service
    matchLabels:
      app: my-app  # Must match the service labels
  strategy:
    type: failover
    primaryGeoTag: "eu"
    dnsTtlSeconds: 30
  splitBrainThresholdSeconds: 300
```

## Key Points

### Required Annotation

The `k8gb.io/hostname` annotation is **required** for LoadBalancer service GSLBs. This annotation specifies the hostname that k8gb will use for DNS record generation.

### DNS Zone Validation

The specified hostname must match one of the configured DNS zones in your k8gb deployment. For example:

- If your k8gb has `DNS_ZONES=example.com:cloud.example.com:300`
- Your hostname must contain `cloud.example.com` (e.g., `myapp.cloud.example.com`)

### Service Labels

The `resourceRef.matchLabels` in the GSLB must match the labels on your LoadBalancer service. This is how k8gb finds the service to manage.

## Examples

### Example 1: Standard DNS Zone

```yaml
# k8gb deployment configuration
env:
- name: DNS_ZONES
  value: "example.com:cloud.example.com:300"

# GSLB
metadata:
  annotations:
    k8gb.io/hostname: "myapp.cloud.example.com"
```

### Example 2: Custom Domain

```yaml
# k8gb deployment configuration
env:
- name: DNS_ZONES
  value: "mycompany.com:lb.mycompany.com:300"

# GSLB
metadata:
  annotations:
    k8gb.io/hostname: "hello.lb.mycompany.com"
```

### Example 3: Multiple Hostnames

If you need multiple hostnames for the same service, you can create multiple GSLB resources:

```yaml
# GSLB 1
apiVersion: k8gb.io/v1beta1
kind: Gslb
metadata:
  name: my-app-gslb-primary
  annotations:
    k8gb.io/hostname: "myapp.example.com"
spec:
  resourceRef:
    apiVersion: v1
    kind: Service
    matchLabels:
      app: my-app
  strategy:
    type: failover
    primaryGeoTag: "eu"

---
# GSLB 2
apiVersion: k8gb.io/v1beta1
kind: Gslb
metadata:
  name: my-app-gslb-alias
  annotations:
    k8gb.io/hostname: "www.myapp.example.com"
spec:
  resourceRef:
    apiVersion: v1
    kind: Service
    matchLabels:
      app: my-app
  strategy:
    type: failover
    primaryGeoTag: "eu"
```

## Error Handling

### Missing Hostname Annotation

If you forget to specify the `k8gb.io/hostname` annotation:

```
LoadBalancer service GSLB my-app-gslb requires k8gb.io/hostname annotation
```

### Empty Hostname Annotation

If the annotation is present but empty:

```
LoadBalancer service GSLB my-app-gslb has empty k8gb.io/hostname annotation
```

### DNS Zone Mismatch

If the hostname doesn't match any configured DNS zones:

```
ingress host myapp.example.com does not match delegated zone [cloud.example.com]
```

## Comparison with Other Integrations

| Integration Type | Hostname Source | Configuration |
|------------------|-----------------|---------------|
| **Ingress** | `spec.rules[].host` | Automatic |
| **Istio VirtualService** | `spec.hosts[]` | Automatic |
| **LoadBalancer Service** | `k8gb.io/hostname` annotation | **Manual** |

## Benefits

1. **Explicit Control**: You have full control over the hostname
2. **Flexibility**: Works with any DNS zone configuration
3. **Consistency**: Follows k8gb annotation patterns
4. **Clarity**: Makes it obvious what hostname will be used

## Limitations

1. **Manual Configuration**: Requires explicit hostname specification
2. **DNS Zone Dependency**: Hostname must match configured DNS zones
3. **Single Hostname**: Each GSLB can only specify one hostname (use multiple GSLBs for multiple hostnames)
