# Resilient Multiregion Global Control Planes With Crossplane and K8gb

This example demonstrates how to build resilient, scalable multicluster environments using Crossplane's declarative infrastructure provisioning integrated with k8gb for DNS-based failover.

## Overview

This reference architecture showcases a Crossplane-based Global Control Plane with Active/Passive setup:

- **Active/Passive Control Planes**: Multiple regions with one active control plane and passive standby
- **DNS-based Failover**: k8gb provides automatic DNS failover when active region fails
- **GSLB Health Monitoring**: Crossplane observes GSLB resources to track control plane health
- **Automated Failover**: Passive control plane transitions to Active during failures

## Key Components

- **Composition Pipeline**: Uses KCL function to observe GSLB resources and report health status
- **GSLB Observation**: Monitors k8gb.absa.oss/v1beta1 Gslb resources in observe-only mode
- **Health Status Integration**: Extracts serviceHealth from GSLB and updates XR status for failover decisions
- **Multiregion Coordination**: Enables Crossplane to make intelligent decisions based on regional health

## Files

- `composition.yaml`: Crossplane Composition with KCL function for GSLB health monitoring
- `definition.yaml`: CompositeResourceDefinition with GSLB status schema for failover logic
- `xr.yaml`: Example XR instance representing a control plane endpoint
- `xr-auto.yaml`: Auto-mode XR example with intelligent GSLB-based failover
- `xr-passive.yaml`: Passive-mode XR example for standby regions
- `functions.yaml`: KCL function package for health monitoring logic
- `provider-*.yaml`: Provider configurations for multicluster access

## Usage

### Local Testing

```shell
# Render the composition locally
make run

# Or using crossplane render directly
crossplane render xr.yaml composition.yaml functions.yaml -r

# Test with Docker environment
make run-in-docker
```

### Deploy Examples

```shell
# Deploy the basic XR
kubectl apply -f xr.yaml

# Deploy auto-mode XR with intelligent failover
kubectl apply -f xr-auto.yaml

# Deploy passive-mode XR for standby regions
kubectl apply -f xr-passive.yaml
```

## Architecture Details

The KCL function implements the core failover monitoring logic:

1. **Observe GSLB**: Creates Kubernetes Object to monitor GSLB resource health
2. **Health Assessment**: Extracts serviceHealth status from k8gb  
3. **Status Propagation**: Updates XR status to reflect regional control plane health (Healthy/UNHEALTHY)
4. **Failover Enablement**: Provides health data for automated Active/Passive transitions

The composition uses `managementPolicies: ["Observe"]` for read-only health monitoring without modifying the underlying GSLB resources, enabling safe observation across multiple regions.

## Important Notes

- The domain is fully parameterized via `spec.hostname` for production flexibility
- Health monitoring is performed in observe-only mode to avoid conflicts across regions
- GSLB status updates drive automated failover decisions between Active/Passive control planes
- Auto-apply policy mode enables intelligent management policy switching based on GSLB health
