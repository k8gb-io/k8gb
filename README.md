# Oh My GLB

## Project Health

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Actions Status](https://github.com/AbsaOSS/ohmyglb/workflows/build/badge.svg)](https://github.com/AbsaOSS/ohmyglb/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/AbsaOSS/ohmyglb)](https://goreportcard.com/report/github.com/AbsaOSS/ohmyglb)

A Global Service Load Balancing solution with a focus on having cloud native qualities and work natively in a Kubernetes context.

## Motivation and Architecture

Please see the extended documentation [here](/docs/index.md)

## Installation and Configuration

### Local Playground

####  Deploy local cluster

```sh
$ make deploy-local-cluster
```
Creates local [kind](https://github.com/kubernetes-sigs/kind) cluster
with several workers for more realistic setup.

#### Deploy local ingress

```sh
$ make deploy-local-ingress
```
Creates local nginx ingress with deployment config similar to Rancher.
IP addresses of workers will be exposed as Ingress addresses.
It will create proper environment for Ohmyglb testing

#### Deploy gslb operator

```sh
$ make deploy-gslb-operator
```
Operator is packaged as a helm chart at [chart/ohmyglb](/chart/ohmyglb) and its
configuration is controlled by [chart/ohmyglb/values.yaml](/chart/ohmyglb/values.yaml)

This step will deploy the operator and its dependencies:

* `EtcdCluster` controlled by [etcd-operator](https://github.com/helm/charts/tree/master/stable/etcd-operator)
* Dedicated [CoreDNS](https://coredns.io/) which uses this etcd cluster as a backend
* [external-dns](https://github.com/kubernetes-sigs/external-dns) with CRD as the source
* ohmyglb controller

Follow the chart notes instructions to check the deployment status.

#### Deploy gslb Custom Resource

```sh
$ make deploy-gslb-cr
```
Creates example `gslb` custom resource of sample configuration

Check testing `gslb` status with
```sh
$ kubectl -n test-gslb describe gslb test-gslb
```

#### Deploy sample workload

```sh
$ make deploy-test-apps
```
It will deploy sample [podinfo](https://github.com/stefanprodan/podinfo) application
matching with `app3` in `gslb` configuration.

After successful deployment you should observe healthy status of `app3` in `gslb` status

```sh
$ kubectl -n test-gslb describe gslb test-gslb
...
  Service Health:
    app1.cloud.example.com:  NotFound
    app2.cloud.example.com:  Unhealthy
    app3.cloud.example.com:  Healthy
```

#### Deploy full local setup

To deploy two cross communicating `ohmyglb` enabled clusters with testing application on top, execute
```
$ make deploy-full-local-setup
```
