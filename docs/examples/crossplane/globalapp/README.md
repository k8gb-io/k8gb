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

## Architecture Diagrams

### Crossplane Function Health Monitoring Flow

![Crossplane Function Health Monitoring](/docs/examples/crossplane/globalapp/assets/crossplane-function-health-monitoring.png)

This diagram illustrates how the Crossplane function dynamically monitors GSLB health status and adjusts management policies. The function observes the GSLB resource status and conditionally sets management policies - when active and healthy, it enables full management (`managementPolicies: ["*"]`), otherwise it switches to observe-only mode (`managementPolicies: ["Observe"]`).

### Active/Passive Cluster Management Policies  

![Active Passive Cluster Policies](/docs/examples/crossplane/globalapp/assets/active-passive-cluster-policies.png)

This diagram shows the management policy differences between Active and Passive clusters. The Active cluster (with green checkmark) has full management policies including Create, Delete, Observe, and Update capabilities. The Passive cluster (with red X) is restricted to Observe-only mode, ensuring it monitors health without making changes until it needs to become active during failover.

## Conference Presentation

This reference PoC architecture was presented at **KubeCon China in Hong Kong 2025**:

- **Presentation Deck**: [Resilient Multiregion Global Control Planes With Crossplane and K8gb](https://docs.google.com/presentation/d/1lCO5k4tTWFTbRdPcx9WzjCS_4aapt116/edit?usp=sharing&ouid=103406682887003571641&rtpof=true&sd=true)
- **Recording**: [Watch on YouTube](https://www.youtube.com/watch?v=L9mRWljLnzw)

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
```

### Local Installation

```
# In the main k8gb repo
make deploy-full-local-setup
```

```
# Check current context. We expect second k8gb cluster to be active after the
# installation
k config current-context
k3d-test-gslb2
```

```
# Install Crossplane
helm install crossplane \
--namespace crossplane-system \
--create-namespace crossplane-stable/crossplane
```

```
# Install Crossplane Providers
k apply -f crossplane-providers
```

```
# Wait for providers readiness
k get providers
NAME                            INSTALLED   HEALTHY   PACKAGE                                                 AGE
provider-azure-cache            True        True      xpkg.upbound.io/upbound/provider-azure-cache:v1         35s
provider-helm                   True        True      xpkg.upbound.io/upbound/provider-helm:v0                35s
provider-kubernetes             True        True      xpkg.upbound.io/upbound/provider-kubernetes:v0          35s
upbound-provider-family-azure   True        True      xpkg.upbound.io/upbound/provider-family-azure:v1.13.0   30s
```

```
# Install ProviderConfigs
k apply -f crossplane-providers/providerconfigs
providerconfig.azure.upbound.io/default created
providerconfig.helm.crossplane.io/default created
providerconfig.kubernetes.crossplane.io/default created
```

```
# Setup Azure Credentials Secret
kubectl create secret generic azure-account-creds -n crossplane-system --from-literal=credentials="$(cat credentials.json)" --dry-run=client -o yaml | kubectl apply -f -
# If you don't have Azure credentials json file, check Quickstart documentation
at https://marketplace.upbound.io/providers/upbound/provider-family-azure/latest
```

```
# Apply demo GlobalApp material
k apply -f definition.yaml,composition.yaml,functions.yaml
compositeresourcedefinition.apiextensions.crossplane.io/globalapps.example.crossplane.io created
Warning: CustomResourceDefinition.apiextensions.k8s.io "GlobalApp.example.crossplane.io" not found
composition.apiextensions.crossplane.io/global-app created
function.pkg.crossplane.io/crossplane-contrib-function-kcl created
function.pkg.crossplane.io/crossplane-contrib-function-auto-ready created
```

```
# Switch context to first k8gb cluster
k config use-context k3d-test-gslb1
Switched to context "k3d-test-gslb1".
# Repeat all installation steps in identical manner on the first cluster!
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
