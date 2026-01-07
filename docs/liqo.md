# Integration with Liqo

You can provide powerful global multi-cluster capabilities by combining k8gb and [liqo.io](https://docs.liqo.io).

In this tutorial, you will learn how to leverage Liqo and K8GB to deploy and expose a multi-cluster application through a *global ingress*.
More in detail, this enables improved load balancing and distribution of the external traffic towards the application replicated across multiple clusters.

Liqo will globally schedule workloads and provide east-west connectivity, while K8GB will globally balance user traffic providing north-south connectivity over the multi-cluster and/or multi-provider environment.

The figure below outlines the high-level scenario, with a client consuming an application from either cluster 1 (e.g., located in EU) or cluster 2 (e.g., located in the US), based on the endpoint returned by the DNS server.

![Global Ingress Overview](images/gslb-liqo-integration.drawio.svg)

## Setup Environment

Checkout the [liqo docs](https://docs.liqo.io/en/v0.5.4/examples/global-ingress.html) to get the environment setup script and to get more details.
It creates the k3d clusters required for the K8GB playground as described in [Local playground for testing and development](local.md) and installs Liqo over them.

## Peer the clusters

To proceed, first generate a new *peer command* from the *gslb-us* cluster:

```bash
PEER_US=$(liqoctl generate peer-command --only-command --kubeconfig $KUBECONFIG_US)
```

And then, run the generated command from the *gslb-eu* cluster:

```bash
echo "$PEER_US" | bash
```

## Deploy an application

First, create a hosting namespace in the *gslb-eu* cluster, and offload it to the remote cluster through Liqo.

```bash
kubectl create namespace podinfo
liqoctl offload namespace podinfo --namespace-mapping-strategy EnforceSameName
```

At this point, it is possible to deploy the *podinfo* helm chart in the `podinfo` namespace:

```bash
helm upgrade --install podinfo --namespace podinfo podinfo/podinfo \
    -f https://raw.githubusercontent.com/liqotech/liqo/master/examples/global-ingress/manifests/values/podinfo.yaml
```

This chart creates a *Deployment* with a *custom affinity* to ensure that the two frontend replicas are scheduled on different nodes and clusters:

```yaml
affinity:
  nodeAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
      nodeSelectorTerms:
      - matchExpressions:
        - key: node-role.kubernetes.io/control-plane
          operator: DoesNotExist
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
    - labelSelector:
        matchExpressions:
        - key: app.kubernetes.io/name
          operator: In
          values:
          - podinfo
      topologyKey: "kubernetes.io/hostname"
```

Additionally, it creates an *Ingress* resource configured with the ingress_annotations.md. This annotation will [instruct the K8GB Global Ingress Controller](ingress_annotations.md) to distribute the traffic across the different clusters.
ingress_annotations.md is an HTTP service, you can contact it using the *curl* command.
Use the `-v` option to understand which of the nodes is being targeted.

You need to use the DNS server in order to resolve the hostname to the IP address of the service.
To this end, create a pod in one of the clusters (it does not matter which one) overriding its DNS configuration.

```bash
HOSTNAME="liqo.cloud.example.com"
K8GB_COREDNS_IP=$(kubectl get svc k8gb-coredns -n k8gb -o custom-columns='IP:spec.clusterIP' --no-headers)

kubectl run -it --rm curl --restart=Never --image=curlimages/curl:7.82.0 --command \
    --overrides "{\"spec\":{\"dnsConfig\":{\"nameservers\":[\"${K8GB_COREDNS_IP}\"]},\"dnsPolicy\":\"None\"}}" \
    -- curl $HOSTNAME -v
```
