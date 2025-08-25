# Automatic GSLB Creation for LoadBalancer Services

## Overview

k8gb now supports automatic GSLB creation for LoadBalancer services, similar to how it works for Ingress resources. This feature allows you to simply add k8gb annotations to your LoadBalancer services, and k8gb will automatically create and manage the corresponding GSLB resources.

## How It Works

When you create a LoadBalancer service with the appropriate k8gb annotations, the k8gb controller will:

1. **Detect the service** with k8gb annotations
2. **Automatically create a GSLB resource** that references the service
3. **Manage the GSLB lifecycle** including updates and deletion
4. **Create DNS records** in your configured DNS provider (Infoblox, etc.)

## Required Annotations

### `k8gb.io/strategy` (Required)
Specifies the GSLB strategy to use:
- `roundRobin` - Distributes traffic across all healthy endpoints
- `failover` - Routes traffic to primary location, fails over to secondary
- `geoip` - Routes traffic based on client geographic location

### `k8gb.io/hostname` (Required)
The hostname that will be used for the LoadBalancer service. This is required because LoadBalancer services don't inherently have hostname information.

## Optional Annotations

### `k8gb.io/primary-geotag` (Required for failover strategy)
Specifies the primary geotag for failover strategy. Must be set when using `k8gb.io/strategy: failover`.

### `k8gb.io/dns-ttl-seconds` (Optional)
DNS TTL in seconds. Defaults to 30 seconds if not specified.

### `k8gb.io/exposed-ip-addresses` (Optional)
Explicitly specify IP addresses for the service. If not provided, k8gb will use the LoadBalancer IP addresses from the service status.

## Example Usage

### Round Robin Strategy

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
  namespace: default
  annotations:
    k8gb.io/strategy: "roundRobin"
    k8gb.io/hostname: "my-app.example.com"
spec:
  type: LoadBalancer
  ports:
    - port: 80
      targetPort: 8080
  selector:
    app: my-app
```

### Failover Strategy

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
  namespace: default
  annotations:
    k8gb.io/strategy: "failover"
    k8gb.io/hostname: "my-app.example.com"
    k8gb.io/primary-geotag: "eu"
    k8gb.io/dns-ttl-seconds: "60"
spec:
  type: LoadBalancer
  ports:
    - port: 80
      targetPort: 8080
  selector:
    app: my-app
```

### GeoIP Strategy

```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
  namespace: default
  annotations:
    k8gb.io/strategy: "geoip"
    k8gb.io/hostname: "my-app.example.com"
    k8gb.io/dns-ttl-seconds: "120"
spec:
  type: LoadBalancer
  ports:
    - port: 80
      targetPort: 8080
  selector:
    app: my-app
```

## Generated GSLB Resource

When you create a LoadBalancer service with k8gb annotations, k8gb will automatically create a GSLB resource that looks like this:

```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: my-app-service
  namespace: default
  annotations:
    k8gb.io/hostname: "my-app.example.com"
  ownerReferences:
    - apiVersion: v1
      kind: Service
      name: my-app-service
      uid: <service-uid>
spec:
  resourceRef:
    apiVersion: v1
    kind: Service
    name: my-app-service
  strategy:
    type: roundRobin
    dnsTtlSeconds: 30
```

## Key Features

### Automatic Lifecycle Management
- **Creation**: GSLB is automatically created when you add k8gb annotations to a LoadBalancer service
- **Updates**: Changes to service annotations are automatically reflected in the GSLB
- **Deletion**: When the service is deleted, the GSLB is automatically cleaned up

### Validation
- Only LoadBalancer services are processed
- Services must have both `k8gb.io/strategy` and `k8gb.io/hostname` annotations
- Services owned by existing GSLB resources are ignored to prevent loops

### Compatibility
- Works alongside existing manual GSLB creation
- Compatible with all existing k8gb features (DNS providers, strategies, etc.)
- No breaking changes to existing functionality

## Migration from Manual GSLB

If you currently have manually created GSLB resources for LoadBalancer services, you can:

1. **Keep existing GSLB resources** - they will continue to work as before
2. **Migrate to automatic creation** by:
   - Adding the required annotations to your LoadBalancer services
   - Deleting the manual GSLB resources
   - k8gb will automatically recreate them

## Troubleshooting

### GSLB Not Created
- Ensure the service is of type `LoadBalancer`
- Verify both `k8gb.io/strategy` and `k8gb.io/hostname` annotations are present
- Check k8gb controller logs for any errors

### DNS Records Not Created
- Verify your DNS provider (Infoblox, etc.) is properly configured
- Check that the LoadBalancer service has external IPs assigned
- Review k8gb controller logs for DNS provider errors

### Strategy Validation Errors
- For `failover` strategy, ensure `k8gb.io/primary-geotag` is specified
- Verify strategy values are valid (`roundRobin`, `failover`, `geoip`)
- Check DNS TTL values are positive integers

## Comparison with Ingress Integration

| Feature | Ingress | LoadBalancer Service |
|---------|---------|---------------------|
| Automatic GSLB Creation | ✅ | ✅ |
| Required Hostname | No (from Ingress spec) | Yes (`k8gb.io/hostname`) |
| Service Type | Any | LoadBalancer only |
| IP Source | Ingress status or annotation | Service status or annotation |
| Strategy Support | All | All |
| DNS Provider Support | All | All |

This feature provides the same convenience as Ingress integration while being specifically tailored for LoadBalancer services that need explicit hostname configuration.
