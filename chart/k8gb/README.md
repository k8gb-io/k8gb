# k8gb

![Version: v0.14.0](https://img.shields.io/badge/Version-v0.14.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v0.14.0](https://img.shields.io/badge/AppVersion-v0.14.0-informational?style=flat-square)

A Helm chart for Kubernetes Global Balancer

```
      _    ___        _
     | | _( _ )  __ _| |__
     | |/ / _ \ / _` | '_ \
     |   < (_) | (_| | |_) |
     |_|\_\___/ \__, |_.__/ .io
                |___/
```

**Homepage:** <https://www.k8gb.io/>

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Andre Baptista Aguas | <andre.aguas@protonmail.com> |  |
| Dinar Valeev | <dinar.valeev@absa.africa> |  |
| Jiri Kremser | <jiri.kremser@gmail.com> |  |
| Michal Kuritka | <kuritka@gmail.com> |  |
| Yury Tsarev | <yury@upbound.io> |  |

## Source Code

* <https://github.com/k8gb-io/k8gb>

## Requirements

Kubernetes: `>= 1.19.0-0`

| Repository | Name | Version |
|------------|------|---------|
| https://coredns.github.io/helm | coredns | 1.39.0 |

For Kubernetes `< 1.19` use this chart and k8gb in version `0.8.8` or lower.

#### Compatibility matrix:

| k8gb                           |      <= 0.8.x      | >= 0.9.0           |
| ------------------------------ | :----------------: | :----------------: |
| Kubernetes <= 1.18             | :white_check_mark: |        :x:         |
| Kubernetes >= 1.19 and <= 1.21 | :white_check_mark: | :white_check_mark: |
| Kubernetes >= 1.22             |        :x:         | :white_check_mark: |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| coredns.deployment.skipConfig | bool | `true` | Skip CoreDNS creation and uses the one shipped by k8gb instead |
| coredns.image.repository | string | `"absaoss/k8s_crd"` | CoreDNS CRD plugin image |
| coredns.image.tag | string | `"v0.1.1"` | image tag |
| coredns.isClusterService | bool | `false` | service: refer to https://www.k8gb.io/docs/service_upgrade.html for upgrading CoreDNS service steps |
| coredns.securityContext | object | `{"capabilities":{"add":[]}}` | Disables all permissions since we don't open privileged ports |
| coredns.servers | list | `[{"plugins":[{"name":"prometheus","parameters":"0.0.0.0:9153"}],"port":5353,"servicePort":53}]` | Only meant to open the correct service and container ports, has no other impact on the coredns configuration |
| coredns.serviceAccount | object | `{"create":true,"name":"coredns"}` | Creates serviceAccount for coredns |
| global.imagePullSecrets | list | `[]` | Reference to one or more secrets to be used when pulling images ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/ |
| infoblox.enabled | bool | `false` | infoblox provider enabled |
| infoblox.gridHost | string | `"10.0.0.1"` | WAPI address |
| infoblox.httpPoolConnections | int | `10` | Size of connections pool |
| infoblox.httpRequestTimeout | int | `20` | Request Timeout in secconds |
| infoblox.sslVerify | bool | `true` | use SSL |
| infoblox.wapiPort | int | `443` | WAPI port |
| infoblox.wapiVersion | string | `"2.3.1"` | WAPI version |
| istio.enabled | bool | `true` | install istio RBAC |
| k8gb.clusterGeoTag | string | `"eu"` | used for places where we need to distinguish between different Gslb instances |
| k8gb.coredns.extraServerBlocks | string | `""` | Extra CoreDNS server blocks |
| k8gb.coredns.extra_plugins | list | `[]` | Extra CoreDNS plugins to be enabled |
| k8gb.deployCrds | bool | `true` | whether it should also deploy the gslb and dnsendpoints CRDs |
| k8gb.deployRbac | bool | `true` | whether it should also deploy the service account, cluster role and cluster role binding |
| k8gb.dnsZones | list | `[{"dnsZoneNegTTL":30,"domain":"cloud.example.com","zone":"example.com"}]` | array of dns zones controlled by gslb§ |
| k8gb.edgeDNSServers[0] | string | `"1.1.1.1"` | use this DNS server as a main resolver to enable cross k8gb DNS based communication |
| k8gb.exposeMetrics | bool | `false` | Exposing metrics |
| k8gb.extGslbClustersGeoTags | string | `"us"` | comma-separated list of external gslb geo tags to pair with |
| k8gb.imageRepo | string | `"docker.io/absaoss/k8gb"` | image repository |
| k8gb.imageTag |  string  | `nil` | image tag defaults to Chart.AppVersion, see Chart.yaml, but can be overrided with imageTag key |
| k8gb.log.format | string | `"simple"` | log format (simple,json) |
| k8gb.log.level | string | `"info"` | log level (panic,fatal,error,warn,info,debug,trace) |
| k8gb.metricsAddress | string | `"0.0.0.0:8080"` | Metrics server address |
| k8gb.nsRecordTTL | int | `30` | TTL of the NS and respective glue record used by external DNS |
| k8gb.podAnnotations | object | `{}` | pod annotations |
| k8gb.podLabels | object | `{}` | pod labels |
| k8gb.reconcileRequeueSeconds | int | `30` | Reconcile time in seconds |
| k8gb.securityContext.allowPrivilegeEscalation | bool | `false` |  |
| k8gb.securityContext.readOnlyRootFilesystem | bool | `true` |  |
| k8gb.securityContext.runAsNonRoot | bool | `true` | For more options consult https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#securitycontext-v1-core |
| k8gb.securityContext.runAsUser | int | `1000` |  |
| k8gb.serviceMonitor | object | `{"enabled":false}` | enable ServiceMonitor |
| k8gb.validatingAdmissionPolicy | object | `{"enabled":false}` | enable validating admission policies |
| openshift.enabled | bool | `false` | Install OpenShift specific RBAC |
| tracing.deployJaeger | bool | `false` | should the Jaeger be deployed together with the k8gb operator? In case of using another OpenTracing solution, make sure that configmap for OTEL agent has the correct exporters set up (`tracing.otelConfig`). |
| tracing.enabled | bool | `false` | if the application should be sending the traces to OTLP collector (env var `TRACING_ENABLED`) |
| tracing.endpoint | string | `"localhost:4318"` | `host:port` where the spans from the applications (traces) should be sent, sets the `OTEL_EXPORTER_OTLP_ENDPOINT` env var This is not the final destination where all the traces are going. Otel collector has its configuration in the associated configmap (`tracing.otelConfig`). |
| tracing.jaegerImage.pullPolicy | string | `"Always"` |  |
| tracing.jaegerImage.repository | string | `"jaegertracing/all-in-one"` | if `tracing.deployJaeger==true` this image will be used in the deployment for Jaeger |
| tracing.jaegerImage.tag | string | `"1.66.0"` |  |
| tracing.otelConfig | string | `nil` | configuration for OTEL collector, this will be represented as configmap called `agent-config` |
| tracing.samplingRatio | string | `nil` | float representing the ratio of how often the span should be kept/dropped (env var `TRACING_SAMPLING_RATIO`) if not specified, the AlwaysSample will be used which is the same as 1.0. `0.1` would mean that 10% of samples will be kept |
| tracing.sidecarImage.pullPolicy | string | `"Always"` |  |
| tracing.sidecarImage.repository | string | `"otel/opentelemetry-collector"` | OpenTelemetry collector into which the k8gb operator sends the spans. It can be further configured to send its data to somewhere else using exporters (Jaeger for instance) |
| tracing.sidecarImage.tag | string | `"0.121.0"` |  |
