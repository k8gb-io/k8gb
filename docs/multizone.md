# Multizone Support
Starting with v0.15.0, k8gb supports managing more than one DNS zone in a single deployment.

## Why multizone?
Previously, k8gb could only manage a single DNS zone per deployment. If you needed to handle multiple zones (for example, for several domains or environments), you had to deploy several instances of k8gb, each in a separate namespace, and carefully coordinate their configuration. This approach was difficult to manage and prone to errors, especially with Helm charts and CRDs.

To simplify things, k8gb now supports defining multiple zones directly in your configuration.

### How to configure multiple zones
The new configuration uses the dnsZones field (in your values.yaml or Helm values), where you can define a list of zones to be managed by a single k8gb deployment.

```yaml
k8gb:
  dnsZones:
    - parentZone: "example.com"             # The parent zone which delegates to the GSLB zone (previously edgeDNSZone)
      loadBalancedZone: "cloud.example.com" # The zone actually managed by GSLB (previously dnsZone)
      dnsZoneNegTTL: 30                     # Negative TTL for SOA records
    - parentZone: "example.org"
      loadBalancedZone: "cloud.example.org"
      dnsZoneNegTTL: 30
    - parentZone: "example.org"
      loadBalancedZone: "onprem.example.org"
      dnsZoneNegTTL: 30
```
>**Note:**The field names have also changed for clarity:
>  - `edgeDNSZone` is now called `parentZone`
>  - `dnsZone` is now called `loadBalancedZone`
>
> For backward compatibility, the dnsZone and edgeDNSZone fields are allowed; otherwise, the dnsZones array is used. For valid values, use either dnsZone and edgeDNSZone or dnsZones.We recommend switching to the new syntax for all new deployments.
