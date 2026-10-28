# Installation and Configuration Tutorials

This section provides comprehensive tutorials for deploying and configuring K8GB with various DNS providers and environments.

> For new deployments, prefer [Gateway API](https://www.k8gb.io/latest/resource_ref/index.md), supported since k8gb v0.17. Examples that specifically use ingress-nginx are legacy: that controller was [retired in March 2026](https://www.kubernetes.dev/blog/2025/11/12/ingress-nginx-retirement/). Kubernetes Ingress and k8gb's Ingress integration remain supported with maintained controllers.

## DNS Provider Integrations

- [General deployment with Infoblox integration](https://www.k8gb.io/latest/deploy_infoblox/index.md)
- [AWS based deployment with Route53 integration](https://www.k8gb.io/latest/deploy_route53/index.md)
- [AWS based deployment with NS1 integration](https://www.k8gb.io/latest/deploy_ns1/index.md)
- [Google Cloud DNS provider](https://www.k8gb.io/latest/provider_gcp/index.md)
- [Using Azure Public DNS provider](https://www.k8gb.io/latest/deploy_azuredns/index.md)
- [Azure based deployment with Windows DNS integration](https://www.k8gb.io/latest/deploy_windowsdns/index.md)
- [General deployment with Cloudflare integration](https://www.k8gb.io/latest/deploy_cloudflare/index.md)
- [Seamless DDNS Integration with Bind9 and other RFC2136-Compatible DNS Environments](https://www.k8gb.io/latest/provider_rfc2136/index.md)

## Development and Testing

- [Local playground for testing and development](https://www.k8gb.io/latest/local/index.md)
- [Local playground with Kuar web app](https://www.k8gb.io/latest/local-kuar/index.md)
- [AI Inference Resilience Demo](https://www.k8gb.io/latest/ai-inference-demo/index.md)
- [Rollback Process](https://www.k8gb.io/latest/rollback_procedures/index.md)

## Monitoring and Observability

- [Metrics](https://www.k8gb.io/latest/metrics/index.md)
- [Traces](https://www.k8gb.io/latest/traces/index.md)

## Configuration

- [Resource References](https://www.k8gb.io/latest/resource_ref/index.md)
- [Address Discovery](https://www.k8gb.io/latest/address_discovery/index.md)
- [Exposed Hostnames](https://www.k8gb.io/latest/exposed-hostnames/index.md)
- [Dynamic Geotags](https://www.k8gb.io/latest/dynamic_geotags/index.md)
- [Dynamic Zones](https://www.k8gb.io/latest/dynamic_zones/index.md)
- [Multi-zone Setup](https://www.k8gb.io/latest/multizone/index.md)
- [Exposing DNS](https://www.k8gb.io/latest/exposing_dns/index.md)

## Platform Integrations

- [Integration with Admiralty](https://www.k8gb.io/latest/admiralty/index.md)
- [Integration with Liqo](https://www.k8gb.io/latest/liqo/index.md)
- [Integration with Rancher Fleet](https://www.k8gb.io/latest/rancher/index.md)
- [Crossplane Global Control Plane](https://www.k8gb.io/latest/crossplane_globalapp/index.md)
- [LoadBalancer Service Integration](https://www.k8gb.io/latest/lbservice-integration/index.md)

## Advanced Topics

- [Service Upgrade](https://www.k8gb.io/latest/service_upgrade/index.md)
- [Partial Deployment Troubleshooting](https://www.k8gb.io/latest/partial-deployment/index.md)
- [Migration Acceptance Tests](https://www.k8gb.io/latest/migration_acceptance/index.md)
- [Split-Brain Scenarios](https://www.k8gb.io/latest/splitbrain/index.md)
- [WRR Caveats](https://www.k8gb.io/latest/wrr_caveats/index.md)
- [External DNS Proxy](https://www.k8gb.io/latest/proxy_externaldns/index.md)
