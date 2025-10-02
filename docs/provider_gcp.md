# Google Cloud DNS Provider

This document describes how to configure K8GB with Google Cloud DNS as the DNS provider.

## Overview

K8GB integrates with Google Cloud DNS through the [external-dns](https://github.com/kubernetes-sigs/external-dns) project using the upstream Helm chart. This integration uses a hybrid DNS architecture where:

- **Google Cloud DNS** manages the parent zone and NS delegation records
- **CoreDNS** (in K8GB clusters) hosts and serves the actual load-balanced zones
- **External-DNS** creates NS delegation records in Cloud DNS pointing to CoreDNS endpoints

This approach provides fast failover and reduces cloud DNS costs while maintaining reliable zone delegation.

## Architecture

### DNS Flow

```
example.com (Google Cloud DNS managed zone)
├── NS app.example.com → gslb-ns-us-central1.example.com (created by external-dns)
├── NS app.example.com → gslb-ns-europe-west1.example.com (created by external-dns)
├── A gslb-ns-us-central1.example.com → <coredns-cluster1-ip> (glue record)
└── A gslb-ns-europe-west1.example.com → <coredns-cluster2-ip> (glue record)

app.example.com (CoreDNS hosted zone - not in Cloud DNS)
├── A failover.app.example.com → <active-cluster-ip> (managed by K8GB/CoreDNS)
├── A roundrobin.app.example.com → <all-cluster-ips> (managed by K8GB/CoreDNS)
└── TXT ownership records for coordination
```

### Key Components

1. **External-DNS**: Creates NS delegation records in Google Cloud DNS
2. **CoreDNS**: Hosts the actual load-balanced DNS zones
3. **K8GB Controller**: Manages DNS records in CoreDNS and coordinates between clusters

## Prerequisites

### Google Cloud Platform Setup

1. **GCP Project**: A GCP project with the Cloud DNS API enabled
2. **Parent DNS Zone**: A Cloud DNS managed zone for your parent domain (e.g., `example.com`)
3. **Service Account**: A Google Service Account with DNS admin permissions
4. **GKE Clusters** (recommended): GKE clusters with Workload Identity enabled

#### Initial GCP Setup Commands

```bash
# Set your project and domain variables
export PROJECT_ID="my-dns-project"
export DOMAIN="example.com"
export ZONE_NAME=$(echo $DOMAIN | sed 's/\./-/g')  # example-com

# Set the default project
gcloud config set project $PROJECT_ID

# Enable required APIs
gcloud services enable dns.googleapis.com
gcloud services enable container.googleapis.com  # If using GKE

# Create the parent DNS managed zone
gcloud dns managed-zones create $ZONE_NAME \
  --dns-name="$DOMAIN." \
  --description="Parent zone for K8GB delegation" \
  --visibility=public

# View the assigned nameservers
gcloud dns managed-zones describe $ZONE_NAME \
  --format="value(nameServers)" | tr ';' '\n'
```

**Note**: For testing purposes, you don't need to own the domain or configure nameservers at your registrar. The Cloud DNS zone just needs to exist for K8GB to create delegation records. For production use, you would need to configure your domain registrar to use the Google Cloud DNS nameservers listed above.

### Required Permissions

The service account used by external-dns needs:
- `roles/dns.admin` - To create NS delegation records in the parent Cloud DNS zone

## Authentication Methods

### Workload Identity (Recommended for GKE)

#### Setup:

1. **Create Google Service Account:**
   ```bash
   gcloud iam service-accounts create k8gb-external-dns \
     --display-name="K8GB External DNS" \
     --project=DNS_PROJECT_ID
   ```

2. **Grant DNS permissions:**
   ```bash
   gcloud projects add-iam-policy-binding DNS_PROJECT_ID \
     --member="serviceAccount:k8gb-external-dns@DNS_PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/dns.admin"
   ```

3. **Create Workload Identity binding:**
   ```bash
   gcloud iam service-accounts add-iam-policy-binding \
     k8gb-external-dns@DNS_PROJECT_ID.iam.gserviceaccount.com \
     --role="roles/iam.workloadIdentityUser" \
     --member="serviceAccount:GKE_PROJECT_ID.svc.id.goog[k8gb/k8gb-external-dns]" \
     --project=DNS_PROJECT_ID
   ```

### Static Service Account Keys (Alternative)

1. **Create service account and generate key:**
   ```bash
   gcloud iam service-accounts create k8gb-external-dns \
     --project=DNS_PROJECT_ID

   gcloud projects add-iam-policy-binding DNS_PROJECT_ID \
     --member="serviceAccount:k8gb-external-dns@DNS_PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/dns.admin"

   gcloud iam service-accounts keys create credentials.json \
     --iam-account=k8gb-external-dns@DNS_PROJECT_ID.iam.gserviceaccount.com
   ```

2. **Create Kubernetes secret:**
   ```bash
   kubectl create secret generic external-dns-gcp-sa \
     --from-file=credentials.json \
     --namespace=k8gb
   ```

## Configuration

### Basic Configuration (Workload Identity)

```yaml
k8gb:
  dnsZones:
    - loadBalancedZone: "app.example.com"  # Zone hosted by CoreDNS
      parentZone: "example.com"            # Zone managed by Cloud DNS
  clusterGeoTag: "us-central1"
  extGslbClustersGeoTags: "europe-west1"

extdns:
  enabled: true
  fullnameOverride: "k8gb-external-dns"
  provider:
    name: google
  serviceAccount:
    annotations:
      iam.gke.io/gcp-service-account: "k8gb-external-dns@DNS_PROJECT_ID.iam.gserviceaccount.com"
  txtPrefix: "k8gb-us-central1-"
  txtOwnerId: "k8gb-app.example.com-us-central1"
  domainFilters:
    - "example.com"  # Must match the parent zone in Cloud DNS
  extraArgs:
    google-project: "DNS_PROJECT_ID"  # Project containing the parent zone

coredns:
  serviceType: LoadBalancer  # Exposes CoreDNS for NS delegation
```

### Static Credentials Configuration

```yaml
extdns:
  enabled: true
  provider:
    name: google
  env:
    - name: GOOGLE_APPLICATION_CREDENTIALS
      value: /etc/secrets/service-account/credentials.json
  extraVolumes:
    - name: google-service-account
      secret:
        secretName: external-dns-gcp-sa
  extraVolumeMounts:
    - name: google-service-account
      mountPath: /etc/secrets/service-account/
      readOnly: true
  extraArgs:
    google-project: "DNS_PROJECT_ID"
```

### Advanced Configuration Options

```yaml
extdns:
  extraArgs:
    # Required: GCP project containing the parent DNS zone
    google-project: "my-dns-project"

    # Optional: Filter zones by visibility
    google-zone-visibility: "public"  # or "private"

    # Optional: Logging configuration
    log-level: "info"
    log-format: "json"  # Recommended for GCP environments

    # Optional: DNS record management
    policy: "sync"  # or "upsert-only"
```

## Multi-Region Deployment

### Two-Region Setup Example

**US Cluster Configuration:**
```yaml
k8gb:
  dnsZones:
    - loadBalancedZone: "app.example.com"
      parentZone: "example.com"
  clusterGeoTag: "us-central1"
  extGslbClustersGeoTags: "europe-west1"

extdns:
  enabled: true
  provider:
    name: google
  txtPrefix: "k8gb-us-central1-"
  txtOwnerId: "k8gb-app.example.com-us-central1"
  domainFilters:
    - "example.com"
  # ... authentication config
```

**Europe Cluster Configuration:**
```yaml
k8gb:
  dnsZones:
    - loadBalancedZone: "app.example.com"  # Same zone name
      parentZone: "example.com"
  clusterGeoTag: "europe-west1"
  extGslbClustersGeoTags: "us-central1"

extdns:
  enabled: true
  provider:
    name: google
  txtPrefix: "k8gb-europe-west1-"  # Different prefix
  txtOwnerId: "k8gb-app.example.com-europe-west1"  # Different owner
  domainFilters:
    - "example.com"
  # ... authentication config
```

## How It Works

1. **Deployment**: When K8GB starts, external-dns detects the CoreDNS service LoadBalancer
2. **NS Record Creation**: External-dns automatically creates NS records in the Cloud DNS parent zone pointing to nameserver names (e.g., `gslb-ns-us-central1.example.com`)
3. **Glue Record Creation**: External-dns creates A records (glue records) for the nameserver names pointing to CoreDNS LoadBalancer IPs
4. **Zone Hosting**: CoreDNS hosts the `app.example.com` zone and responds to DNS queries
5. **Load Balancing**: K8GB manages A records within CoreDNS based on cluster health and strategy
6. **Coordination**: Multiple K8GB clusters coordinate through DNS queries to each other's CoreDNS

## Verification

### Check NS Delegation

```bash
# Verify NS records were created in Cloud DNS parent zone
gcloud dns record-sets list --zone="example-com" --name="app.example.com." --type="NS"

# Should show nameserver names like:
# app.example.com. NS 300 gslb-ns-us-central1.example.com. gslb-ns-europe-west1.example.com.

# Verify glue A records were created
gcloud dns record-sets list --zone="example-com" --name="gslb-ns-us-central1.example.com." --type="A"
gcloud dns record-sets list --zone="example-com" --name="gslb-ns-europe-west1.example.com." --type="A"

# Should show CoreDNS LoadBalancer IPs
```

### Test DNS Resolution

```bash
# Test that delegation works
nslookup failover.app.example.com 8.8.8.8

# Test direct CoreDNS query
nslookup failover.app.example.com <coredns-cluster-ip>

# Verify both clusters respond
dig roundrobin.app.example.com @<coredns-cluster1-ip>
dig roundrobin.app.example.com @<coredns-cluster2-ip>
```

## Monitoring and Troubleshooting

### Key Log Sources

```bash
# External-DNS logs (NS record management)
kubectl logs -l app.kubernetes.io/name=external-dns -n k8gb -f

# K8GB operator logs (CoreDNS record management)
kubectl logs -l name=k8gb -n k8gb -f

# CoreDNS logs (DNS query handling)
kubectl logs -l app.kubernetes.io/name=coredns -n k8gb -f
```

### Common Issues

1. **NS records not created in Cloud DNS:**
   ```
   Error: Failed to create NS record
   ```
   - Verify service account has `roles/dns.admin` in the DNS project
   - Check that `domainFilters` matches the parent zone name
   - Ensure CoreDNS LoadBalancer has external IPs assigned

2. **CoreDNS not reachable:**
   ```
   Error: DNS resolution timeout
   ```
   - Verify CoreDNS service type is LoadBalancer
   - Check that firewall rules allow UDP/53 traffic to CoreDNS
   - Confirm LoadBalancer provisioning completed successfully

3. **Cross-cluster coordination fails:**
   ```
   Warning: Could not resolve peer cluster
   ```
   - Verify NS delegation is working from both clusters
   - Check that both clusters can reach each other's CoreDNS services
   - Confirm `extGslbClustersGeoTags` configuration matches

### Debug Commands

```bash
# Check CoreDNS service and endpoints
kubectl get svc k8gb-coredns -n k8gb
kubectl get endpoints k8gb-coredns -n k8gb

# Verify Cloud DNS records
gcloud dns record-sets list --zone="example-com"

# Test DNS delegation chain
dig +trace failover.app.example.com

# Check K8GB DNSEndpoint resources
kubectl get dnsendpoint -n k8gb
```

## Best Practices

### DNS Configuration

1. **Use meaningful geo tags** that reflect actual geographic regions
2. **Configure appropriate TTLs** for your failover requirements
3. **Monitor DNS propagation** especially during initial setup

### Security

1. **Prefer Workload Identity** over static service account keys
2. **Use zone-level permissions** when possible instead of project-level
3. **Monitor DNS changes** through Cloud Audit Logs
4. **Secure CoreDNS endpoints** with appropriate network policies

### Performance

1. **Deploy CoreDNS with anti-affinity** to ensure high availability
2. **Use regional load balancers** for CoreDNS services when possible
3. **Monitor DNS query latency** and optimize based on usage patterns
4. **Consider private zones** for internal-only applications

### Cost Optimization

1. **Minimize Cloud DNS usage** - only NS records are stored in Cloud DNS
2. **Use appropriate LoadBalancer tiers** for CoreDNS services
3. **Monitor cross-region traffic** costs between clusters

## Examples

Complete working examples are available in the [docs/examples/gcp](../examples/gcp/) directory:

- [k8gb-cluster-gcp-europe-west1.yaml](../examples/gcp/k8gb-cluster-gcp-europe-west1.yaml)
- [k8gb-cluster-gcp-us-central1.yaml](../examples/gcp/k8gb-cluster-gcp-us-central1.yaml)
- [test-gslb-failover.yaml](../examples/gcp/test-gslb-failover.yaml)
- [test-gslb-roundrobin.yaml](../examples/gcp/test-gslb-roundrobin.yaml)

See the [README](../examples/gcp/README.md) in the examples directory for a complete setup walkthrough.

## References

- [Google Cloud DNS Documentation](https://cloud.google.com/dns/docs)
- [External-DNS Google Provider](https://kubernetes-sigs.github.io/external-dns/latest/docs/tutorials/gke/)
- [GKE Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity)
- [K8GB Architecture Documentation](../intro.md)