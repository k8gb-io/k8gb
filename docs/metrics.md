# Metrics

K8GB generates [Prometheus][prometheus]-compatible metrics.
Metrics endpoints are exposed via `-metrics` service in operator namespace and can be scraped by 3rd party tools:

``` yaml
spec:
...
  ports:
  - name: http-metrics
    port: 8383
    protocol: TCP
    targetPort: 8383
  - name: cr-metrics
    port: 8686
    protocol: TCP
    targetPort: 8686
```

Metrics can be also automatically discovered and monitored by [Prometheus Operator][prometheus-operator] via automatically generated [ServiceMonitor][service-monitor] CRDs , in case if [Prometheus Operator][prometheus-operator]  is deployed into the cluster.

### General metrics

[controller-runtime][controller-runtime-metrics] standard metrics, extended with K8GB operator-specific metrics listed below:

#### `healthy_records`

Number of healthy records observed by K8GB.

Example:

```yaml
# HELP k8gb_gslb_healthy_records Number of healthy records observed by K8GB.
# TYPE k8gb_gslb_healthy_records gauge
k8gb_gslb_healthy_records{name="test-gslb",namespace="test-gslb"} 6
```

#### `ingress_hosts_per_status`

Number of ingress hosts per status (NotFound, Healthy, Unhealthy), observed by K8GB.

Example:

```yaml
# HELP k8gb_gslb_ingress_hosts_per_status Number of managed hosts observed by K8GB.
# TYPE k8gb_gslb_ingress_hosts_per_status gauge
k8gb_gslb_ingress_hosts_per_status{name="test-gslb",namespace="test-gslb",status="Healthy"} 1
k8gb_gslb_ingress_hosts_per_status{name="test-gslb",namespace="test-gslb",status="NotFound"} 1
k8gb_gslb_ingress_hosts_per_status{name="test-gslb",namespace="test-gslb",status="Unhealthy"} 2
```

Served on `0.0.0.0:8383/metrics` endpoint

### Custom resource specific metrics

Info metrics, automatically exposed by operator based on the number of the current instances of an operator's custom resources in the cluster.

Example:

```yaml
# HELP gslb_info Information about the Gslb custom resource.
# TYPE gslb_info gauge
gslb_info{namespace="test-gslb",gslb="test-gslb"} 1
```

Served on `0.0.0.0:8686/metrics` endoint

[prometheus]: https://prometheus.io/
[prometheus-operator]: https://github.com/coreos/prometheus-operator
[service-monitor]: https://github.com/coreos/prometheus-operator#customresourcedefinitions
[controller-runtime-metrics]: https://book.kubebuilder.io/reference/metrics.html
