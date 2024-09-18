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
| https://k8gb-io.github.io/coredns-helm | coredns | 1.15.3 |

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
| azuredns.authSecretName | string | `"external-dns-secret-azure"` | External-dns secret name which contains Azure credentials. See https://github.com/k8gb-io/external-dns/blob/master/docs/tutorials/azure.md#configuration-file for expected format |
| azuredns.createAuthSecret.aadClientId | string | `"myAadClientId"` | Azure client ID that is associated with the Service Principal. |
| azuredns.createAuthSecret.aadClientSecret | string | `"myAadClientSecret"` | Azure client secret that is associated with the Service Principal. |
| azuredns.createAuthSecret.enabled | bool | `true` | Create an authentication secret for Azure DNS based on the values below alternatively, you can create the secret manually and pass its name in the `azuredns.authSecretName` value |
| azuredns.createAuthSecret.resourceGroup | string | `"myDnsResourceGroup"` | Azure Resource Group which holds the Azure DNS Zone (which is defined as 'k8gb.edgeDNSZone') |
| azuredns.createAuthSecret.subscriptionId | string | `"mySubscriptionId"` | subscription id which holds the Azure DNS zone |
| azuredns.createAuthSecret.tenantId | string | `"myTenantId"` | Azure tenant ID which holds the managed identity |
| azuredns.createAuthSecret.useManagedIdentityExtension | bool | `false` | Use either AKS Kubelet Identity or AAD Pod Identities |
| azuredns.createAuthSecret.useWorkloadIdentityExtension | bool | `false` | Use AKS workload identity extension |
| azuredns.createAuthSecret.userAssignedIdentityID | string | `"myUserAssignedIdentityID"` | Client id from the Managed identitty when using the AAD Pod Identities |
| azuredns.enabled | bool | `false` |  |
| cloudflare.dnsRecordsPerPage | int | `5000` | Configure how many DNS records to fetch per request see https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/cloudflare.md#throttling |
| cloudflare.enabled | bool | `false` | Enable Cloudflare provider |
| cloudflare.zoneID | string | `"replaceme"` | Cloudflare Zone ID follow https://developers.cloudflare.com/fundamentals/setup/find-account-and-zone-ids/ to find your zoneID value |
| coredns.deployment.skipConfig | bool | `true` | Skip CoreDNS creation and uses the one shipped by k8gb instead |
| coredns.image.repository | string | `"absaoss/k8s_crd"` | CoreDNS CRD plugin image |
| coredns.image.tag | string | `"v0.1.0"` | image tag |
| coredns.isClusterService | bool | `false` | service: refer to https://www.k8gb.io/docs/service_upgrade.html for upgrading CoreDNS service steps |
| coredns.serviceAccount | object | `{"create":true,"name":"coredns"}` | Creates serviceAccount for coredns |
| externaldns.dnsPolicy | string | `"ClusterFirst"` | `.spec.template.spec.dnsPolicy` for ExternalDNS deployment |
| externaldns.extraEnv | list | `[]` | extra environment variables |
| externaldns.extraVolumeMounts | list | `[]` | extra volume mounts |
| externaldns.extraVolumes | list | `[]` | extra volumes |
| externaldns.image | string | `"ghcr.io/k8gb-io/external-dns:v0.13.4-azure-ns-multiarch"` | external-dns image repo:tag It is important to use the image from k8gb external-dns fork to get the full functionality. See links below https://github.com/k8gb-io/external-dns https://github.com/k8gb-io/external-dns/pkgs/container/external-dns |
| externaldns.interval | string | `"20s"` | external-dns sync interval |
| externaldns.securityContext.fsGroup | int | `65534` | For ExternalDNS to be able to read Kubernetes and AWS token files |
| externaldns.securityContext.runAsNonRoot | bool | `true` |  |
| externaldns.securityContext.runAsUser | int | `1000` | For more options consult https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#securitycontext-v1-core |
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
| k8gb.dnsZone | string | `"cloud.example.com"` | dnsZone controlled by gslb |
| k8gb.dnsZoneNegTTL | int | `300` | Negative TTL for SOA record |
| k8gb.edgeDNSServers | list | `["1.1.1.1"]` | host/ip[:port] format is supported here where port defaults to 53 |
| k8gb.edgeDNSServers[0] | string | `"1.1.1.1"` | use this DNS server as a main resolver to enable cross k8gb DNS based communication |
| k8gb.edgeDNSZone | string | `"example.com"` | main zone which would contain gslb zone to delegate |
| k8gb.exposeMetrics | bool | `false` | Exposing metrics |
| k8gb.extGslbClustersGeoTags | string | `"us"` | comma-separated list of external gslb geo tags to pair with |
| k8gb.imageRepo | string | `"docker.io/absaoss/k8gb"` | image repository |
| k8gb.imageTag |  string  | `nil` | image tag defaults to Chart.AppVersion, see Chart.yaml, but can be overrided with imageTag key |
| k8gb.log.format | string | `"simple"` | log format (simple,json) |
| k8gb.log.level | string | `"info"` | log level (panic,fatal,error,warn,info,debug,trace) |
| k8gb.metricsAddress | string | `"0.0.0.0:8080"` | Metrics server address |
| k8gb.podAnnotations | object | `{}` | pod annotations |
| k8gb.podLabels | object | `{}` | pod labels |
| k8gb.reconcileRequeueSeconds | int | `30` | Reconcile time in seconds |
| k8gb.securityContext.allowPrivilegeEscalation | bool | `false` |  |
| k8gb.securityContext.readOnlyRootFilesystem | bool | `true` |  |
| k8gb.securityContext.runAsNonRoot | bool | `true` | For more options consult https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#securitycontext-v1-core |
| k8gb.securityContext.runAsUser | int | `1000` |  |
| k8gb.serviceMonitor | object | `{"enabled":false}` | enable ServiceMonitor |
| k8gb.splitBrainCheck | bool | `false` | Enable SplitBrain check (Infoblox only) |
| ns1.enabled | bool | `false` | Enable NS1 provider |
| ns1.ignoreSSL | bool | `false` | optional custom NS1 API endpoint for on-prem setups endpoint: https://api.nsone.net/v1/ |
| openshift.enabled | bool | `false` | Install OpenShift specific RBAC |
| rfc2136.enabled | bool | `false` |  |
| rfc2136.rfc2136Opts[0].host | string | `"host.k3d.internal"` |  |
| rfc2136.rfc2136Opts[1].port | int | `1053` |  |
| rfc2136.rfc2136auth.gssTsig.enabled | bool | `false` |  |
| rfc2136.rfc2136auth.gssTsig.gssTsigCreds[0].kerberos-username | string | `"ad-user"` |  |
| rfc2136.rfc2136auth.gssTsig.gssTsigCreds[1].kerberos-password | string | `"ad-user-pass"` |  |
| rfc2136.rfc2136auth.gssTsig.gssTsigCreds[2].kerberos-realm | string | `"ad-domain-realm"` |  |
| rfc2136.rfc2136auth.gssTsig.kerberosConfigMap | string | `"kerberos-config"` |  |
| rfc2136.rfc2136auth.insecure.enabled | bool | `false` |  |
| rfc2136.rfc2136auth.tsig.enabled | bool | `true` |  |
| rfc2136.rfc2136auth.tsig.tsigCreds[0].tsig-secret-alg | string | `"hmac-sha256"` |  |
| rfc2136.rfc2136auth.tsig.tsigCreds[1].tsig-keyname | string | `"externaldns-key"` |  |
| route53.assumeRoleArn | string | `nil` | specify IRSA Role in AWS ARN format for assume role permissions or disable it by setting to `null` |
| route53.enabled | bool | `false` | Enable Route53 provider |
| route53.hostedZoneID | string | `"ZXXXSSS"` | Route53 ZoneID |
| route53.irsaRole | string | `"arn:aws:iam::111111:role/external-dns"` | specify IRSA Role in AWS ARN format or disable it by setting to `null` |
| route53.secret | string | `nil` | alternatively specify the secret name with static credentials for IAM user (make sure this user has limited privileges) this can be useful when IRSA is not present or when using say Azure cluster and Route53 docs: https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md#create-iam-user-and-attach-the-policy |
| tracing.deployJaeger | bool | `false` | should the Jaeger be deployed together with the k8gb operator? In case of using another OpenTracing solution, make sure that configmap for OTEL agent has the correct exporters set up (`tracing.otelConfig`). |
| tracing.enabled | bool | `false` | if the application should be sending the traces to OTLP collector (env var `TRACING_ENABLED`) |
| tracing.endpoint | string | `"localhost:4318"` | `host:port` where the spans from the applications (traces) should be sent, sets the `OTEL_EXPORTER_OTLP_ENDPOINT` env var This is not the final destination where all the traces are going. Otel collector has its configuration in the associated configmap (`tracing.otelConfig`). |
| tracing.jaegerImage.pullPolicy | string | `"Always"` |  |
| tracing.jaegerImage.repository | string | `"jaegertracing/all-in-one"` | if `tracing.deployJaeger==true` this image will be used in the deployment for Jaeger |
| tracing.jaegerImage.tag | string | `"1.61.0"` |  |
| tracing.otelConfig | string | `nil` | configuration for OTEL collector, this will be represented as configmap called `agent-config` |
| tracing.samplingRatio | string | `nil` | float representing the ratio of how often the span should be kept/dropped (env var `TRACING_SAMPLING_RATIO`) if not specified, the AlwaysSample will be used which is the same as 1.0. `0.1` would mean that 10% of samples will be kept |
| tracing.sidecarImage.pullPolicy | string | `"Always"` |  |
| tracing.sidecarImage.repository | string | `"otel/opentelemetry-collector"` | OpenTelemetry collector into which the k8gb operator sends the spans. It can be further configured to send its data to somewhere else using exporters (Jaeger for instance) |
| tracing.sidecarImage.tag | string | `"0.109.0"` |  |
