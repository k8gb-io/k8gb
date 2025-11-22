# K8GB with Google Cloud DNS - Example

This example demonstrates how to configure K8GB with Google Cloud DNS across two GKE clusters in different regions.

> **📖 For detailed setup instructions, architecture explanation, and troubleshooting, see [Google Cloud DNS Provider Documentation](../../provider_gcp.md)**

## Quick Overview

This example deploys K8GB across two GKE clusters:
- **Primary cluster**: `europe-west1` (failover primary)
- **Secondary cluster**: `us-central1` (failover secondary)
- **DNS Zone**: `gcp-test.k8gb.io` (load-balanced zone hosted by CoreDNS)
- **Parent Zone**: `k8gb.io` (Cloud DNS managed zone for delegation)

## Prerequisites

1. Two GKE clusters in different regions with Workload Identity enabled
2. A Cloud DNS managed zone (see [provider_gcp.md - Prerequisites](../../provider_gcp.md#prerequisites))
3. Service account with `roles/dns.admin` (see [provider_gcp.md - Authentication](../../provider_gcp.md#authentication-methods))

## Deployment

1. **Deploy K8GB on the first cluster:**
   ```bash
   # Configure kubectl for europe-west1 cluster
   gcloud container clusters get-credentials my-cluster-eu \
     --region=europe-west1 --project=my-gke-project

   # Install K8GB
   helm upgrade --install k8gb k8gb/k8gb \
     --namespace=k8gb --create-namespace \
     --values=k8gb-cluster-gcp-europe-west1.yaml
   ```

2. **Deploy K8GB on the second cluster:**
   ```bash
   # Configure kubectl for us-central1 cluster
   gcloud container clusters get-credentials my-cluster-us \
     --region=us-central1 --project=my-gke-project

   # Install K8GB
   helm upgrade --install k8gb k8gb/k8gb \
     --namespace=k8gb --create-namespace \
     --values=k8gb-cluster-gcp-us-central1.yaml
   ```

3. **Deploy test application and GSLB resources:**
   ```bash
   # Create test namespace (do this on both clusters)
   kubectl apply -f ../../../deploy/gslb/test-namespace-ingress.yaml

   # Deploy podinfo via Helm (do this on both clusters)
   helm repo add podinfo https://stefanprodan.github.io/podinfo
   helm upgrade --install frontend --namespace test-gslb \
     -f ../../../deploy/test-apps/podinfo/podinfo-values.yaml \
     --set ui.message="$(kubectl config current-context)" \
     podinfo/podinfo --version 6.9.2

   # Apply GSLB configurations (do this on both clusters)
   kubectl apply -f test-gslb-failover.yaml
   kubectl apply -f test-gslb-roundrobin.yaml
   ```

## Configuration Files

This example includes:
- `k8gb-cluster-gcp-europe-west1.yaml` - Configuration for EU cluster
- `k8gb-cluster-gcp-us-central1.yaml` - Configuration for US cluster
- `test-gslb-failover.yaml` - Failover strategy example
- `test-gslb-roundrobin.yaml` - Round-robin strategy example

### Key Configuration Notes

- **Workload Identity**: Configured by default in the example YAML files
- **Static Credentials**: Commented out alternative (uncomment env/volumes sections if needed)
- **Cross-Project Setup**: Use `extraArgs.google-project` if DNS zone is in a different project
- **Geo Tags**: Must be unique per cluster (`europe-west1`, `us-central1`)
- **TXT Ownership**: Each cluster needs unique `txtPrefix` and `txtOwnerId`

See [provider_gcp.md](../../provider_gcp.md) for detailed configuration parameters.

## Verification

```bash
# Check K8GB status
kubectl get gslb -n test-gslb
kubectl describe gslb test-gslb-failover -n test-gslb

# Verify NS delegation records in Cloud DNS
gcloud dns record-sets list --zone="k8gb-io" --filter="name~'gcp-test.k8gb.io'"

# Test DNS resolution
nslookup failover.gcp-test.k8gb.io 8.8.8.8
nslookup roundrobin.gcp-test.k8gb.io 8.8.8.8
```

**For detailed verification steps and troubleshooting, see:**
- [Verification](../../provider_gcp.md#verification)
- [Troubleshooting](../../provider_gcp.md#monitoring-and-troubleshooting)

## Testing Failover

```bash
# Scale down the primary cluster application
kubectl scale deployment frontend-podinfo --replicas=0 -n test-gslb

# Wait for health check to detect failure (~30-60 seconds)
# Verify DNS now points to secondary cluster
nslookup failover.gcp-test.k8gb.io 8.8.8.8
```

## Additional Resources

- [Google Cloud DNS Provider Documentation](../../provider_gcp.md) - Complete reference
- [K8GB Introduction](../../intro.md) - How K8GB works
