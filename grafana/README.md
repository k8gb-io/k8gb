## Sample Grafana Dashboards

This directory contains semi-generated dashboards for Grafana. It's using kubebuilder's 
plugin introduced in [pr #2858](https://github.com/kubernetes-sigs/kubebuilder/pull/2858).

In order to modify the custom metrics, change the config in [`custom-metrics/config.yaml`](./custom-metrics/config.yaml) and in project root run:

```bash
kubebuilder edit --plugins grafana.kubebuilder.io/v1-alpha
```
