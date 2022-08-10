# Traces

K8GB can create and send spans/traces to either Jaeger which is supported also on the Helm chart level or to
any other OTEL compliant solution. We don't recommend sending the tracing data directly to the tracer, because
of the possible vendor lock-in. OTEL collector can be deployed as a side-car container for k8gb pod and
forward all the traces to a configurable sink or event multiple sinks for redundancy.

### Architecture

We are not opinionated about the OpenTracing vendors. In the following diagram `X` can be Jaeger, Zipkin, LigthStep, Grafana's Tempo, Instana, etc. It can also be another OTEL collector to form more [sophisticated](https://opentelemetry.io/docs/collector/configuration/) pipeline setup.

Sidecar use-case:

```text
+--------------+               +------------+             +----------+
|     k8gb     |               |     X      |    http     |   User   |
| ------------ |     otlp      |            +------------>|          |
| OTEL sidecar +-------------->|            |             |          |
+--------------+               +------------+             +----------+
```

### Deployment

By default the tracing is disabled and no sidecar container is being created during the k8gb deployment. To
enable the tracing, one has to set the `tracing.enabled=true` in Helm Chart. This will create the sidecar container for k8gb deployment, tweaks couple of env vars there. It will create the [configmap](https://github.com/k8gb-io/k8gb/blob/master/chart/k8gb/templates/otel/otel-config.yaml) for OTEL sidecar. This configuration of OTEL collector can be overriden by `tracing.otelConfig`.

If you need something quickly up and running, make sure that `tracing.deployJaeger` is also set to `true`.
In this scenario you will end up also with Jaeger deployed and service for it. To be able to access it one can continue with:

```bash
kubectl port-forward svc/jaeger-collector 16686
open http://localhost:16686
```

Also both sidecar container image and jaeger deployment's container image can be tweaked by
`tracing.{sidecarImage,jaegerImage}.{repository,tag,pullPolicy}`.

### Custom Architecture

In case you have already a OTEL collector present in the Kubernetes cluster and you do not want to introduce
a new one, you can deploy also the following topology:

```text
+--------------+               +----------------+             +------------+             +----------+
|     k8gb     |               | OTEL collector |    otlp     |     X      |    http     |   User   |
|              |     otlp      |                +------------>|            +------------>|          |
|              +-------------->|                |             |            |             |          |
+--------------+               +----------------+             +------------+             +----------+
```

However, we don't support this use-case on the Helm chart level so you are on your own with the setup. Nonetheless, it should be relatively straightforward. All you have to do is set the following env vars for k8gb deployment:
 - `TRACING_ENABLED` (set it to `true`)
 - `OTEL_EXPORTER_OTLP_ENDPOINT` (`host:port` of OTEL collector from the ASCII diagram)
 - `TRACING_SAMPLING_RATIO` (optional, float representing the % ratio of how many traces are being collected)


#### Distributed Tracing & Context Propagation

K8gb is a k8s controller/operator so in nature event-based system, where the invocation of the request is done
using creating custom resource (`gslb`) or `Ingress` with certain annotation. In more traditional
micro-service world, the context propagation between traced systems is done using HTTP headers. Provided that
the comm client is also instrumented with OpenTracing bits one doesn't have to call the `Extract` and 
`Inject` on his own. However, for the CRD space nothing has been introduced so far and having the tracing 
key-value metadata stored in each custom resource would be overkill here. Not speaking about increasing the 
overall complexity of such a system.

If k8gb had a REST api, it would be a very low-hanging fruit on the other hand.

As for the propagation of context down for the calls that k8gb does, this may make sense for direct HTTP 
calls to external systems such as Infoblox or Route53. Then provided that those systems also support 
OpenTracing and they have the context propagation done right on their part, one could see the full insight 
into the requests and examine what takes most of the time or where the issue happened. As for communication 
with `ExternalDNS`, it's again on the "CRD level" -> hard to achieve + the ExternalDNS operator is not traced.
