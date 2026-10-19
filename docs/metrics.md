# Metrics

K8GB generates [Prometheus][prometheus]-compatible metrics, in addition to the standard [controller-runtime][controller-runtime-metrics] metrics.

When either `k8gb.exposeMetrics` or `k8gb.serviceMonitor.enabled` is set in the Helm chart, a `k8gb-metrics` Service is created. It exposes a single `metrics` port, taken from `k8gb.metricsAddress` (default `0.0.0.0:8080`), which serves the `/metrics` endpoint and can be scraped by 3rd party tools.

If [Prometheus Operator][prometheus-operator] is deployed into the cluster, set `k8gb.serviceMonitor.enabled` to create a [ServiceMonitor][service-monitor] so that the metrics are discovered and monitored automatically.

[prometheus]: https://prometheus.io/
[prometheus-operator]: https://github.com/coreos/prometheus-operator
[service-monitor]: https://github.com/coreos/prometheus-operator#customresourcedefinitions
[controller-runtime-metrics]: https://book.kubebuilder.io/reference/metrics.html

## Available metrics

The k8gb exposes several metrics to help you monitor the health and behavior.

| Metric | Type | Description | Labels |
|---|:---:|---|---|
| `k8gb_gslb_errors_total` | Counter | Number of errors | `namespace`, `name` |
| `k8gb_gslb_healthy_records` | Gauge | Number of healthy records observed by k8gb. | `namespace`, `name` |
| `k8gb_gslb_reconciliation_loops_total` | Counter | Number of successful reconciliation loops. | `namespace`, `name` |
| `k8gb_gslb_service_status_num` | Gauge | Number of managed hosts observed by k8gb. | `namespace`, `name`, `status` |
| `k8gb_gslb_status_count_for_failover` | Gauge | Gslb status count for Failover strategy. | `namespace`, `name`, `status` |
| `k8gb_gslb_status_count_for_geoip` | Gauge | Gslb status count for GeoIP strategy. | `namespace`, `name`, `status` |
| `k8gb_gslb_status_count_for_roundrobin` | Gauge | Gslb status count for RoundRobin strategy. | `namespace`, `name`, `status` |
| `k8gb_infoblox_heartbeat_errors_total` | Counter | Number of k8gb Infoblox TXT record errors. | `namespace`, `name` |
| `k8gb_infoblox_heartbeats_total` | Counter | Number of k8gb Infoblox heartbeat TXT record updates. | `namespace`, `name` |
| `k8gb_infoblox_request_duration` | Histogram | Duration of the HTTP request to Infoblox API in seconds. | `request`, `success` |
| `k8gb_infoblox_zone_update_errors_total` | Counter | Number of k8gb Infoblox zone update errors. | `namespace`, `name` |
| `k8gb_infoblox_zone_updates_total` | Counter | Number of k8gb Infoblox zone updates. | `namespace`, `name` |
| `k8gb_endpoint_status_num` | Gauge | Number of targets in DNS endpoint. | `namespace`, `name`, `dns_name` |
| `k8gb_runtime_info` | Gauge | K8gb runtime info. | `namespace`, `k8gb_version`, <br>`go_version`, `arch`, `os`, `git_sha` |
| `k8gb_gslb_healthy_local_records` | Gauge | Number of local cluster IPs in healthy records.| `namespace`, `name`, `geotag`|

## OpenTracing

Optionally k8gb operator can expose traces in OpenTelemetry format to any available OTEL compliant tracing solution. Consult the [following page](traces.md) for more details.
