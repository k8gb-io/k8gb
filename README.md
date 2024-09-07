<p align="center" class="disable-logo">
<a href="#"><img src="https://raw.githubusercontent.com/cncf/artwork/master/projects/k8gb/icon/color/k8gb-icon-color.svg" width="200"/></a>
</p>
<h1 align="center" class="disable-logo" style="margin-top: 0;">K8GB - Kubernetes Global Balancer<a href="https://www.k8gb.io"></h1>
<p align="center"><a href="https://landscape.cncf.io/?item=orchestration-management--coordination-service-discovery--k8gb">CNCF Sandbox Project</a> | <a href="https://github.com/orgs/k8gb-io/projects/2/views/2">Roadmap</a> | <a href="https://cloud-native.slack.com/archives/C021P656HGB">Join #k8gb on CNCF Slack</a></p>

[![License: MIT](https://img.shields.io/badge/License-Apache_2.0-yellow.svg)](https://opensource.org/licenses/Apache-2.0)
[![Build Status](https://github.com/k8gb-io/k8gb/workflows/Golang%20lint,%20golic,%20gokart%20and%20test/badge.svg?branch=master)](https://github.com/k8gb-io/k8gb/actions?query=workflow%3A%22Golang%20lint,%20golic,%20gokart%20and%20test%22+branch%3Amaster)
[![Terratest Status](https://github.com/k8gb-io/k8gb/workflows/Terratest/badge.svg?branch=master)](https://github.com/k8gb-io/k8gb/actions?query=workflow%3ATerratest+branch%3Amaster)
[![CodeQL](https://github.com/k8gb-io/k8gb/workflows/CodeQL/badge.svg?branch=master)](https://github.com/k8gb-io/k8gb/actions?query=workflow%3ACodeQL+branch%3Amaster)
[![Go Report Card](https://goreportcard.com/badge/github.com/k8gb-io/k8gb)](https://goreportcard.com/report/github.com/k8gb-io/k8gb)
[![Helm Publish](https://github.com/k8gb-io/k8gb/actions/workflows/helm_publish.yaml/badge.svg)](https://github.com/k8gb-io/k8gb/actions/workflows/helm_publish.yaml)
[![KubeLinter](https://github.com/k8gb-io/k8gb/workflows/KubeLinter/badge.svg?branch=master)](https://github.com/k8gb-io/k8gb/actions?query=workflow%3AKubeLinter+branch%3Amaster)
[![Docker Pulls](https://img.shields.io/docker/pulls/absaoss/k8gb)](https://hub.docker.com/r/absaoss/k8gb)
[![Artifact HUB](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/k8gb)](https://artifacthub.io/packages/search?repo=k8gb)
[![doc.crds.dev](https://img.shields.io/badge/doc-crds-purple)](https://doc.crds.dev/github.com/k8gb-io/k8gb)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb.svg?type=shield)](https://app.fossa.com/projects/custom%2B162%2Fgithub.com%2Fk8gb-io%2Fk8gb?ref=badge_shield)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/4866/badge)](https://bestpractices.coreinfrastructure.org/projects/4866)
[![CLOMonitor](https://img.shields.io/endpoint?url=https://clomonitor.io/api/projects/cncf/k8gb/badge)](https://clomonitor.io/projects/cncf/k8gb)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/k8gb-io/k8gb/badge)](https://securityscorecards.dev/viewer/?uri=github.com/k8gb-io/k8gb)

A Global Service Load Balancing solution with a focus on having cloud native qualities and work natively in a Kubernetes context.

Just a single Gslb CRD to enable the Global Load Balancing:

```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: test-gslb-failover
  namespace: test-gslb
spec:
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    matchLabels: # ingresses.networking.k8s.io resource selector
      app: test-gslb-failover
  strategy:
    type: failover # Global load balancing strategy
    primaryGeoTag: eu-west-1 # Primary cluster geo tag
```

[Global load balancing](https://cloud.redhat.com/blog/global-load-balancer-approaches), commonly referred to as GSLB (Global Server Load Balancing) solutions, has been typically the domain of proprietary network software and hardware vendors and installed and managed by siloed network teams.

k8gb is a completely open source, cloud native, global load balancing solution for Kubernetes.

k8gb focuses on load balancing traffic across geographically dispersed Kubernetes clusters using multiple load balancing [strategies](./docs/strategy.md) to meet requirements such as region failover for high availability.

Global load balancing for any Kubernetes Service can now be enabled and managed by any operations or development teams in the same Kubernetes native way as any other custom resource.

## Key Differentiators

* Load balancing is based on timeproof DNS protocol which is perfect for global scope and extremely reliable
* No dedicated management cluster and no single point of failure
* Kubernetes native application health checks utilizing status of Liveness and Readiness probes for load balancing decisions
* Configuration with a single Kubernetes CRD of Gslb kind

## Quick Start

Simply run

```sh
make deploy-full-local-setup
```

It will deploy two local [k3s](https://k3s.io/) clusters via [k3d](https://k3d.io/), [expose associated CoreDNS service for UDP DNS traffic](./docs/exposing_dns.md)), and install k8gb with test applications and two sample Gslb resources on top.

This setup is adapted for local scenarios and works without external DNS provider dependency.

Consult with [local playground](/docs/local.md) documentation to learn all the details of experimenting with local setup.

Optionally, you can run `make deploy-prometheus` and check the metrics on the test clusters (http://localhost:9080, http://localhost:9081).

## Motivation and Architecture

k8gb was born out of the need for an open source, cloud native GSLB solution at Absa Group in South Africa.

As part of the bank's wider container adoption running multiple, geographically dispersed Kubernetes clusters, the need for a global load balancer that was driven from the health of Kubernetes Services was required and for which there did not seem to be an existing solution.

Yes, there are proprietary network software and hardware vendors with GSLB solutions and products, however, these were costly, heavyweight in terms of complexity and adoption, and were not Kubernetes native in most cases, requiring dedicated hardware or software to be run outside of Kubernetes.

This was the problem we set out to solve with k8gb.

Born as a completely open source project and following the popular Kubernetes operator pattern, k8gb can be installed in a Kubernetes cluster and via a Gslb custom resource, can provide independent GSLB capability to any Ingress or Service in the cluster, without the need for handoffs and coordination between dedicated network teams.

k8gb commoditizes GSLB for Kubernetes, putting teams in complete control of exposing Services across geographically dispersed Kubernetes clusters across public and private clouds.

k8gb requires no specialized software or hardware, relying completely on other OSS/CNCF projects, has no single point of failure, and fits in with any existing Kubernetes deployment workflow (e.g. GitOps, Kustomize, Helm, etc.) or tools.

Please see the extended architecture documentation [here](/docs/index.md)

Internal k8gb architecture and its components are described [here](/docs/components.md)

## Installation and Configuration Tutorials

* [General deployment with Infoblox integration](/docs/deploy_infoblox.md)
* [AWS based deployment with Route53 integration](/docs/deploy_route53.md)
* [AWS based deployment with NS1 integration](/docs/deploy_ns1.md)
* [Using Azure Public DNS provider](/docs/deploy_azuredns.md)
* [Azure based deployment with Windows DNS integration](/docs/deploy_windowsdns.md)
* [General deployment with Cloudflare integration](/docs/deploy_cloudflare.md)
* [Seamless DDNS Integration with Bind9 and other RFC2136-Compatible DNS Environments](/docs/provider_rfc2136.md)
* [Local playground for testing and development](/docs/local.md)
* [Local playground with Kuar web app](/docs/local-kuar.md)
* [Metrics](/docs/metrics.md)
* [Traces](/docs/traces.md)
* [Ingress annotations](/docs/ingress_annotations.md)
* [Integration with Admiralty](/docs/admiralty.md)
* [Integration with Liqo](/docs/liqo.md)
* [Integration with Rancher Fleet](/docs/rancher.md)

## Adopters

A list of publicly known users of the K8GB project can be found in [ADOPTERS.md](/ADOPTERS.md).
We encourage all users of K8GB to add themselves to this list!

## Production Readiness

k8gb is very well tested with the following environment options

| Type                             | Implementation                                                               |
|----------------------------------|------------------------------------------------------------------------------|
| Kubernetes Version               | for k8s `< 1.19` use k8gb `<= 0.8.8`; since k8s `1.19` use `0.9.0` or newer  |
| Environment                      | Self-managed, AWS(EKS) [*](#clarify)                                         |
| Ingress Controller               | NGINX, AWS Load Balancer Controller [*](#clarify)                            |
| EdgeDNS                          | Infoblox, Route53, NS1                                                       |

<a name="clarify"></a>* We only mention solutions where we have tested and verified a k8gb installation.
If your Kubernetes version or Ingress controller is not included in the table above, it does not mean that k8gb will not work for you. k8gb is architected to run on top of any compliant Kubernetes cluster and Ingress controller.

## Presentations Featuring k8gb

[//]: # (Table is generated with the help of https://www.tablesgenerator.com/markdown_tables#)

| **KubeCon EU 2024** [![](https://img.youtube.com/vi/MsQ0E7SYNPo/0.jpg)](https://www.youtube.com/watch?v=MsQ0E7SYNPo "K8gb: Reliable Global Service Load Balancing without vendor lock-in \| Project Lightning Talk") | **KubeCon NA 2023** [![](https://img.youtube.com/vi/4qJDkw5YGqM/0.jpg)](https://www.youtube.com/watch?v=4qJDkw5YGqM "KubeCon NA 2023: Take It to the Edge: Creating a Globally Distributed Ingress with Istio & K8gb - Jimmi Dyson, D2iQ") |
|---|---|
| **KCDBengaluru 2023** [![](https://img.youtube.com/vi/vrDCUIVyc4g/0.jpg)](https://www.youtube.com/watch?v=vrDCUIVyc4g "Kubernetes Community Days Bengaluru 2023: Cloud Native Multi Cluster/Multicloud Global Load Balancer for Kubernetes") | **KubeCon EU 2023** [![](https://img.youtube.com/vi/U46hlF0Z3xs/0.jpg)](https://www.youtube.com/watch?v=U46hlF0Z3xs "KubeCon EU 2023: Recovering from Regional Failures at Cloud Native Speeds") |
| **KubeCon NA 2021** [![](https://img.youtube.com/vi/-lkKZRdv81A/0.jpg)](https://www.youtube.com/watch?v=-lkKZRdv81A "KubeCon NA 2021: Cloud Native Global Load Balancer for Kubernetes") | **FOSDEM 2022** [![](https://img.youtube.com/vi/1UTWxf7PQis/0.jpg)](https://www.youtube.com/watch?v=1UTWxf7PQis "FOSDEM 2022: Cloud Native Global Load Balancer for Kubernetes") |
| **NS1 INS1GHTS** [![](https://img.youtube.com/vi/T_4EiAqwevI/0.jpg)](https://www.youtube.com/watch?v=T_4EiAqwevI "INS1GHTS: Cloud Native Global Load Balancer for Kubernetes") | **Crossplane Community Day** [![](https://img.youtube.com/vi/5l4Xf_Q8ybY/0.jpg)](https://www.youtube.com/watch?v=5l4Xf_Q8ybY "Crossplane Community Day Europe: Scaling Kubernetes Global Balancer with Crossplane") |
| **#29 DoK Community** [![](https://img.youtube.com/vi/MluFlwPFZws/hqdefault.jpg)](https://www.youtube.com/watch?v=MluFlwPFZws "#29 DoK Community: How Absa Developed Cloud Native Global Load Balancer for Kubernetes") | **AWS Containers from the Couch show** [![](https://img.youtube.com/vi/5pe3ezSnVI8/hqdefault.jpg)](https://www.youtube.com/watch?v=5pe3ezSnVI8 "AWS Containers from the Couch") |
| **OpenShift Commons Briefings** [![](https://img.youtube.com/vi/5DhO9C2NCrk/0.jpg)](https://www.youtube.com/watch?v=5DhO9C2NCrk "OpenShift Commons Briefings") | **Demo at Kubernetes SIG Multicluster** [![](https://img.youtube.com/vi/jeUeRQM-ZyM/0.jpg)](https://www.youtube.com/watch?v=jeUeRQM-ZyM "Kubernetes SIG Multicluster") |

You can also find recordings from our community meetings at [k8gb youtube channel](https://www.youtube.com/channel/UCwvtktvdZu_pg-t-INvuW5g).

## Online Publications Featuring k8gb

* https://oilbeater.com/en/2024/04/18/k8gb-best-cloudnative-gslb/
* https://www.redhat.com/en/blog/global-load-balancing-red-hat-openshift-k8gb
* https://andrewbaker.ninja/2021/01/22/external-k8gb-presentation-to-kubernetes-sig-multicluster/

## And Even Books Featuring k8gb :)

| **Kubernetes - An Enterprise Guide - Second Edition** [![](https://m.media-amazon.com/images/I/81zq0mNn-WL._AC_UY436_FMwebp_QL65_.jpg)](https://www.amazon.com/Kubernetes-Enterprise-Effectively-containerize-applications/dp/1803230037 "Kubernetes - An Enterprise Guide - Second Edition: Effectively containerize applications, integrate enterprise systems, and scale applications in your enterprise") | **Kubernetes – An Enterprise Guide - Third Edition** [![](https://m.media-amazon.com/images/I/71mWBgaJMRL._AC_UY436_FMwebp_QL65_.jpg)]( https://www.amazon.com/Kubernetes-Enterprise-Effectively-containerize-applications-ebook/dp/B0CT8M958T/ "Kubernetes – An Enterprise Guide: Effectively containerize applications, integrate enterprise systems, and scale applications in your enterprise 3rd Edition") |
|---|---|

## Contributing

See [CONTRIBUTING](/CONTRIBUTING.md)
