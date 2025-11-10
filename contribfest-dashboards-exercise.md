# ContribFest Exercise: Grafana dashboards & k8gb metrics

## Overview

This exercise encourages the community to improve k8gb's dashboards the metrics exposed by the k8gb controller.

## The Problem

The k8gb maintainers prepared some dashboards a couple of years ago. But they are now outdated and don't track all the Golden Signals.


### Setup

Create a local cluster with Prometheus and Grafana:
```
$ K8GB_LOCAL_VERSION=test FULL_LOCAL_SETUP_WITH_APPS=true make deploy-full-local-setup
$ make deploy-prometheus
$ make deploy-grafana
```

Login to Grafana at localhost:3000, the credentials are `admin` `admin`.

### Your Task

- Analyse the already existing dashboards (Dashboards view)
- Explore the metrics that are being collected (Explore view)
- Create a new dashboard that coverts the Four Golden Signals of Observability from an platform engineer's perspective (some inspiration: https://povilasv.me/how-to-monitor-kubernetes-controllers/) & contribute it to the project
- Think about k8gb both as a platform engineer and as a user, is there additional telemetry that you would like to see? Suggest the collection of additional metrics by creating an issue, or implement it yourself if you are up to the challenge.
