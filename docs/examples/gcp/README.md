# K8GB with Google Cloud DNS

This example demonstrates how to configure K8GB with Google Cloud DNS as the DNS provider across two GKE clusters in different regions.

## Prerequisites

- Two GKE clusters in different regions (e.g., `europe-west1` and `us-central1`)
- A Cloud DNS managed zone configured in Google Cloud Platform
- Proper IAM permissions configured for external-dns to manage DNS records

## Architecture Overview

This setup creates a multi-region K8GB deployment with:
- **Primary cluster**: `europe-west1` (failover primary)
- **Secondary cluster**: `us-central1` (failover secondary)
- **DNS Zone**: `gcp-test.k8gb.io` (managed by Cloud DNS)
- **Parent Zone**: `k8gb.io` (for delegation)

## Authentication Methods

K8GB with Google Cloud DNS supports two authentication methods:

### 1. Workload Identity (Recommended)

Workload Identity allows your K8GB pods to authenticate as Google Service Accounts without storing credentials.

#### Setup Steps:

1. **Create Google Service Account:**
   ```bash
   gcloud iam service-accounts create k8gb-external-dns \
     --display-name="K8GB External DNS" \
     --project=my-dns-project
   ```

2. **Grant DNS permissions:**
   ```bash
   gcloud projects add-iam-policy-binding my-dns-project \
     --member="serviceAccount:k8gb-external-dns@my-dns-project.iam.gserviceaccount.com" \
     --role="roles/dns.admin"
   ```

3. **Enable Workload Identity binding:**
   ```bash
   gcloud iam service-accounts add-iam-policy-binding \
     k8gb-external-dns@my-dns-project.iam.gserviceaccount.com \
     --role="roles/iam.workloadIdentityUser" \
     --member="serviceAccount:my-gke-project.svc.id.goog[k8gb/k8gb-external-dns]" \
     --project=my-dns-project
   ```

4. **Annotate Kubernetes Service Account** (done automatically by Helm chart):
   ```yaml
   serviceAccount:
     annotations:
       iam.gke.io/gcp-service-account: "k8gb-external-dns@my-dns-project.iam.gserviceaccount.com"
   ```

### 2. Static Service Account Keys (Alternative)

For environments where Workload Identity is not available, you can use static service account keys.

#### Setup Steps:

1. **Create service account and download key:**
   ```bash
   gcloud iam service-accounts create k8gb-external-dns \
     --display-name="K8GB External DNS" \
     --project=my-dns-project

   gcloud projects add-iam-policy-binding my-dns-project \
     --member="serviceAccount:k8gb-external-dns@my-dns-project.iam.gserviceaccount.com" \
     --role="roles/dns.admin"

   gcloud iam service-accounts keys create credentials.json \
     --iam-account=k8gb-external-dns@my-dns-project.iam.gserviceaccount.com \
     --project=my-dns-project
   ```

2. **Create Kubernetes secret:**
   ```bash
   kubectl create secret generic external-dns-gcp-sa \
     --from-file=credentials.json \
     --namespace=k8gb
   ```

3. **Update configuration** to use static credentials (uncomment the env/volumes sections in the YAML files).

## DNS Zone Setup

1. **Create Cloud DNS managed zone:**
   ```bash
   gcloud dns managed-zones create k8gb-io \
     --dns-name="k8gb.io." \
     --description="K8GB parent zone" \
     --project=my-dns-project
   ```

2. **Note the nameservers:**
   ```bash
   gcloud dns record-sets list --zone="k8gb-io" --name="k8gb.io." --type="NS"
   ```

3. **Configure delegation** in your domain registrar to point to Google Cloud DNS nameservers.

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
   kubectl apply -f ../../../deploy/crds/test-namespace-ingress.yaml

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

## Configuration Parameters

### Key External-DNS Parameters for Google Cloud DNS:

- `provider.name: google` - Specifies Google Cloud DNS provider
- `extraArgs.google-project` - GCP project containing DNS zones (required if different from cluster project)
- `extraArgs.google-zone-visibility` - Filter zones by visibility (`public` or `private`)
- `domainFilters` - Limit external-dns to specific parent domains
- `txtPrefix` - Prefix for TXT records (should include cluster geo tag)
- `txtOwnerId` - Unique identifier for this external-dns instance

### Cross-Project Setup

When your GKE clusters are in one project but Cloud DNS zones are in another:

```yaml
extraArgs:
  google-project: "my-dns-project"  # Project containing DNS zones
```

Ensure the service account has `dns.admin` role in the DNS project and the appropriate Workload Identity bindings.

## Verification

1. **Check K8GB status:**
   ```bash
   kubectl get gslb -n test-gslb
   kubectl describe gslb test-gslb-failover -n test-gslb
   ```

2. **Verify DNS records:**
   ```bash
   # Check that DNS records were created
   gcloud dns record-sets list --zone="k8gb-io" --filter="name~'gcp-test.k8gb.io'"

   # Test DNS resolution
   nslookup failover.gcp-test.k8gb.io 8.8.8.8
   nslookup roundrobin.gcp-test.k8gb.io 8.8.8.8
   ```

3. **Test failover:**
   ```bash
   # Scale down the primary cluster application
   kubectl scale deployment frontend-podinfo --replicas=0 -n test-gslb

   # Wait a few minutes and verify DNS now points to secondary cluster
   nslookup failover.gcp-test.k8gb.io 8.8.8.8
   ```

## Troubleshooting

### Common Issues:

1. **Permission denied errors:**
   - Verify service account has `roles/dns.admin` permission
   - Check Workload Identity binding is correct
   - Ensure `google-project` parameter points to the correct DNS project

2. **DNS records not created:**
   - Check external-dns logs: `kubectl logs -l app.kubernetes.io/name=external-dns -n k8gb`
   - Verify domain filters match your DNS zone
   - Ensure TXT ownership records don't conflict between clusters

3. **Cross-cluster communication issues:**
   - Verify both clusters can resolve each other's CoreDNS services
   - Check firewall rules allow traffic between cluster regions
   - Confirm `extGslbClustersGeoTags` match the other cluster's `clusterGeoTag`

### Useful Commands:

```bash
# Check external-dns logs
kubectl logs -l app.kubernetes.io/name=external-dns -n k8gb -f

# Check K8GB operator logs
kubectl logs -l name=k8gb -n k8gb -f

# View DNSEndpoint resources created by K8GB
kubectl get dnsendpoint -n k8gb

# Check Cloud DNS records via gcloud
gcloud dns record-sets list --zone="k8gb-io"
```

## Cost Considerations

- Cloud DNS pricing is based on managed zones and queries
- Network Load Balancer usage for CoreDNS services
- Cross-region network traffic between clusters
- Consider using Private Google Access to reduce egress costs

For production deployments, review Google Cloud DNS [pricing](https://cloud.google.com/dns/pricing) and optimize based on your query patterns and zone requirements.