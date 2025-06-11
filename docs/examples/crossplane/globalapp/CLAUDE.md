# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

This is a Crossplane demonstration for "Resilient Multiregion Global Control Planes With Crossplane and K8gb" showcasing how to build resilient, scalable multicluster environments using Crossplane's declarative infrastructure provisioning integrated with k8gb for DNS-based failover.

### Core Concept

The example demonstrates a reference architecture for creating a Crossplane-based Global Control Plane with Active/Passive setup:
- **Active/Passive Control Planes**: Multiple regions with one active control plane and passive standby
- **DNS-based Failover**: k8gb provides automatic DNS failover when active region fails
- **GSLB Health Monitoring**: Crossplane observes GSLB resources to track control plane health
- **Automated Failover**: Passive control plane transitions to Active during failures

### Key Components

- **Composition Pipeline**: Uses KCL function to observe GSLB resources and report health status
- **GSLB Observation**: Monitors k8gb.absa.oss/v1beta1 Gslb resources in observe-only mode
- **Health Status Integration**: Extracts serviceHealth from GSLB and updates XR status for failover decisions
- **Multiregion Coordination**: Enables Crossplane to make intelligent decisions based on regional health

### File Structure

- `composition.yaml`: Crossplane Composition with KCL function for GSLB health monitoring
- `definition.yaml`: CompositeResourceDefinition with GSLB status schema for failover logic
- `xr.yaml`: Example XR instance representing a control plane endpoint
- `test.yaml`: Expected Kubernetes Object for GSLB observation setup
- `functions.yaml`: KCL function package for health monitoring logic
- `provider-*.yaml`: Kubernetes provider configurations for multicluster access

## Common Commands

```bash
# Render the composition locally for testing
make run
# or
crossplane render xr.yaml composition.yaml functions.yaml -r

# Test with Docker environment
make run-in-docker
```

## Development Notes

The KCL function implements the core failover monitoring logic:
1. **Observe GSLB**: Creates Kubernetes Object to monitor GSLB resource health
2. **Health Assessment**: Extracts serviceHealth status from k8gb
3. **Status Propagation**: Updates XR status to reflect regional control plane health (Healthy/UNHEALTHY)
4. **Failover Enablement**: Provides health data for automated Active/Passive transitions

**Important**: The domain is fully parameterized via `spec.hostname` for production flexibility.

The composition uses managementPolicies: ["Observe"] for read-only health monitoring without modifying the underlying GSLB resources, enabling safe observation across multiple regions.