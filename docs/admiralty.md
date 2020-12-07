# Integration with Admiralty

Combination of k8gb and admiralty.io provides powerful
global multi-cluster capabilities.

Admiralty will globally schedule. And k8gb will globally balance.

This tutorial covers local end-to-end integration quick start of two projects.

## Deploy Admiralty

Just follow https://admiralty.io/docs/quick_start

## Deploy k8gb to Target clusters

```sh
helm repo add k8gb https://www.k8gb.io

kubectl --context kind-eu create ns k8gb
helm --kube-context kind-eu --namespace k8gb upgrade --install k8gb k8gb/k8gb --set k8gb.clusterGeoTag=eu --set k8gb.extGslbClustersGeoTags=us

kubectl --context kind-us create ns k8gb
helm --kube-context kind-us --namespace k8gb upgrade --install k8gb k8gb/k8gb --set k8gb.clusterGeoTag=us --set k8gb.extGslbClustersGeoTags=eu
```

It will install k8gb instances to both Target clusters.

Notice the GeoTag configuration.

## Create sample workloads on the source cluster

```sh
helm --kube-context kind-cd upgrade --install podinfo podinfo/podinfo --set replicaCount=2 --set ingress.enabled=true
```

## Add pod annotations to enable multi-cluster scheduling

Edit `podinfo` Deployment to add admiralty `multicluster.admiralty.io/elect` scheduling Annotation to Pod spec template

```sh
kubectl edit deploy podinfo
```

The Pod template spec part should look similar to

```yaml
...
  template:
    metadata:
      annotations:
        prometheus.io/port: "9898"
        prometheus.io/scrape: "true"
        multicluster.admiralty.io/elect: ""
      creationTimestamp: null
...
```

Observe that pods were scheduled to Target cluster after then Annotation patch

```sh
kubectl --context kind-us get pod
NAME                             READY   STATUS    RESTARTS   AGE
podinfo-557c458ddb-9qpjf-l29w6   1/1     Running   0          11s

kubectl --context kind-eu get pod
NAME                             READY   STATUS    RESTARTS   AGE
podinfo-557c458ddb-9nkmx-p7z9r   1/1     Running   0          15s
```

Observe that associated Ingress has also followed the Pods to the Target clusters

```sh
kubectl --context kind-eu get ing
NAME      CLASS    HOSTS   ADDRESS   PORTS   AGE
podinfo   <none>   *                 80      10m

kubectl --context kind-eu get ing
NAME      CLASS    HOSTS   ADDRESS   PORTS   AGE
podinfo   <none>   *                 80      10m
```

## Add k8gb annotations to Ingress object to enable global load balancing

Observer that there are no Gslb resources in the Target clusters

```sh
kubectl --context kind-eu get gslb
No resources found in default namespace.

kubectl --context kind-us get gslb
No resources found in default namespace.
```

Add k8gb strategy annotation to the Ingress object in the Source cluster

```sh
kubectl annotate ing podinfo k8gb.io/strategy=roundRobin
```

Observe that Gslb resources were properly created

```sh
kubectl --context kind-eu get gslb
NAME      AGE
podinfo   102s

kubectl --context kind-us get gslb
NAME      AGE
podinfo   69s
```

Note: it's not a real working k8gb setup but just a way to demonstrate multi-cluster
Gslb CR propagation with Admiralty.

Please refer to https://www.k8gb.io/ documentation to create fully working real-life Gslb setup.
