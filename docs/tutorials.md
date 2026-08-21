# Installation and Configuration Tutorials

This section provides comprehensive tutorials for deploying and configuring K8GB with various DNS providers and environments.

> **North-south data plane:** for new installs prefer [Gateway API](resource_ref.md) (k8gb `resourceRef` since v0.17).
> Many older tutorials still show Kubernetes Ingress / NGINX examples; treat those as legacy — Ingress NGINX was
> [retired in March 2026](https://www.kubernetes.dev/blog/2025/11/12/ingress-nginx-retirement/).

## DNS Provider Integrations

* [General deployment with Infoblox integration](deploy_infoblox.md)
* [AWS based deployment with Route53 integration](deploy_route53.md)
* [AWS based deployment with NS1 integration](deploy_ns1.md)
* [Using Azure Public DNS provider](deploy_azuredns.md)
* [Azure based deployment with Windows DNS integration](deploy_windowsdns.md)
* [General deployment with Cloudflare integration](deploy_cloudflare.md)
* [Seamless DDNS Integration with Bind9 and other RFC2136-Compatible DNS Environments](provider_rfc2136.md)

## Development and Testing

* [Local playground for testing and development](local.md)
* [Local playground with Kuar web app](local-kuar.md)

## Monitoring and Observability

* [Metrics](metrics.md)
* [Traces](traces.md)

## Configuration

* [Resource References](resource_ref.md)
* [Address Discovery](address_discovery.md)
* [Dynamic Geotags](dynamic_geotags.md)
* [Dynamic Zones](dynamic_zones.md)
* [Multi-zone Setup](multizone.md)
* [Exposing DNS](exposing_dns.md)

## Platform Integrations

* [Integration with Admiralty](admiralty.md)
* [Integration with Liqo](liqo.md)
* [Integration with Rancher Fleet](rancher.md)

## Advanced Topics

* [Service Upgrade](service_upgrade.md)
* [WRR Caveats](wrr_caveats.md)
* [External DNS Proxy](proxy_externaldns.md)
