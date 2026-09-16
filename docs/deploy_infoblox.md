# General deployment with Infoblox integration

## Prologue

For simplicity let's assume that you operate two geographically distributed clusters you want to enable global load-balancing for. In this example, two local clusters will represent those two distributed clusters.

* Let's switch the context to the first cluster
```sh
export KUBECONFIG=eu-cluster
```

* Copy the default `values.yaml` from k8gb chart to any convenient location, e.g.
```sh
cp chart/k8gb/values.yaml ~/k8gb/eu-cluster.yaml
```

* Modify the example configuration. Important parameters described below:
  * `dnsZone` - this zone will be delegated to the `edgeDNS` in your environment. E.g. `yourzone.edgedns.com`
  * `edgeDNSZone` - this zone will be automatically configured by k8gb to delegate to `dnsZone` and will make k8gb controlled nodes act as authoritative server for this zone. E.g. `edgedns.com`
  * `parentZoneDNSServers` stable DNS servers in your environment that is controlled by edgeDNS provider e.g. Infoblox so k8gb instances will be able to talk to each other through automatically created DNS names
  * `clusterGeoTag` to geographically tag your cluster. We are operating `eu` cluster in this example
  * `extGslbClustersGeoTags` contains Geo tag of the cluster(s) to talk with when k8gb is deployed to multiple clusters. Imagine your second cluster is `us` so we tag it accordingly
  * `infoblox.enabled: true` to enable automated zone delegation configuration at edgeDNS provider. You don't need it for local testing and can optionally be skipped. Meanwhile, in this section we will cover a fully operational end-to-end scenario.
The other parameters do not need to be modified unless you want to do something special. E.g. to use images from private registry

  * Export Infoblox related information in the shell.
```sh
export WAPI_USERNAME=<WAPI_USERNAME>
export WAPI_PASSWORD=<WAPI_PASSWORD>
```

* Create the Infoblox secret which is used by k8gb to configure edgeDNS by running:
```sh
kubectl create ns k8gb
make infoblox-secret
```

* Expose associated k8gb CoreDNS service for DNS traffic on worker nodes.
  > Check [this document](exposing_dns.md) for detailed information.

* Let's deploy k8gb to the first cluster. Most of the helper commands are abstracted by GNU `make`. If you want to look under the hood please check the `Makefile`. In general, standard Kubernetes/Helm commands are used. Point deployment mechanism to your custom `values.yaml`
```sh
make deploy-gslb-operator VALUES_YAML=~/k8gb/eu-cluster.yaml
```

* It should deploy k8gb pretty quickly. Let's check the pod status
```sh
 kubectl -n k8gb get pod
NAME                                                       READY   STATUS     RESTARTS   AGE
k8gb-76cc56b55-t779s                                       1/1     Running    0          39s
k8gb-coredns-799984c646-qz88m                              1/1     Running    0          41s
```

* Deploy k8gb to the second cluster by repeating the same steps with the exception of:
  * Switch context to 2nd cluster with `export KUBECONFIG=us-cluster`
  * Create another custom `values.yaml` with `cp ~/k8gb/eu-cluster.yaml ~/k8gb/us-cluster.yaml`
  * Create another geo tag to enable cross cluster communication:
    * `clusterGeoTag` becomes `us`
    * `extGslbClustersGeoTags` becomes `eu`
  * Run the installation pointing to new values file `make deploy-gslb-operator VALUES_YAML=~/k8gb/us-cluster.yaml`

* When your 2nd cluster is ready by checking with `kubectl -n k8gb get pod`, we can proceed with the sample application installation

* We will use well known testing community app of [podinfo](https://github.com/stefanprodan/podinfo)
```sh
helm repo add podinfo https://stefanprodan.github.io/podinfo
kubectl create ns test-gslb
helm upgrade --install podinfo --namespace test-gslb --set ui.message="us" podinfo/podinfo
```
As you can see above we did set special geo tag message in podinfo configuration matching cluster geo tag. It is just for demonstration purposes.

* Check that podinfo is running
```sh
kubectl -n test-gslb get pod
NAME                       READY   STATUS    RESTARTS   AGE
podinfo-5cfcdc9c45-jbg96   1/1     Running   0          2m18s
```

* Let's create Gslb CRD to enable global load balancing for this application. Notice the podinfo Service name
```sh
kubectl -n test-gslb get svc
NAME            TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)             AGE
podinfo         ClusterIP   10.96.250.84    <none>        9898/TCP,9999/TCP   9m39s
```

* Create a custom resource `~/k8gb/podinfogslb.yaml` describing an `Ingress` and a `Gslb` as per the sample below:
```yaml
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: podinfo
  namespace: test-gslb
  labels:
    app: podinfo
spec:
  ingressClassName: nginx
  rules:
    - host: podinfo.cloud.example.com
      http:
        paths:
        - path: /
          pathType: Prefix
          backend:
            service:
              name: podinfo # This should point to Service name of testing application
              port:
                name: http
---
apiVersion: k8gb.io/v1beta1
kind: Gslb
metadata:
  name: podinfo
  namespace: test-gslb
spec:
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    matchLabels:
      app: podinfo
  strategy:
    type: roundRobin
```

* And apply the resource in the target app namespace
```sh
kubectl -n test-gslb apply -f podinfogslb.yaml
gslb.k8gb.io/podinfo created
```

* Check Gslb resource
```sh
kubectl -n test-gslb get gslb
NAME      AGE
podinfo   39s
```

* Check Gslb resource status (spec reflects `resourceRef` + `strategy`; status is filled by the controller):
```sh
kubectl -n test-gslb get gslb podinfo -o yaml
apiVersion: k8gb.io/v1beta1
kind: Gslb
metadata:
  name: podinfo
  namespace: test-gslb
spec:
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    matchLabels:
      app: podinfo
  strategy:
    type: roundRobin
status:
  geoTag: us
  healthyRecords:
    podinfo.cloud.example.com:
    - 172.18.0.10
    - 172.18.0.7
    - 172.18.0.8
  serviceHealth:
    podinfo.cloud.example.com: Healthy
```

* In the output above you should see that Gslb detected the `Healthy` status of the referenced `podinfo` Ingress / Service

* Check that internal k8gb DNS servers are responding accordingly on this cluster
  * Pick one of the worker nodes to test with
    ```sh
    kubectl get nodes -o wide
    NAME                       STATUS   ROLES           AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION     CONTAINER-RUNTIME
    test-gslb2-control-plane   Ready    control-plane   53m   v1.29.4   172.18.0.9    <none>        Ubuntu 22.04.4 LTS   6.5.0-linuxkit     containerd://1.7.13
    test-gslb2-worker          Ready    <none>          52m   v1.29.4   172.18.0.8    <none>        Ubuntu 22.04.4 LTS   6.5.0-linuxkit     containerd://1.7.13
    test-gslb2-worker2         Ready    <none>          52m   v1.29.4   172.18.0.7    <none>        Ubuntu 22.04.4 LTS   6.5.0-linuxkit     containerd://1.7.13
    test-gslb2-worker3         Ready    <none>          52m   v1.29.4   172.18.0.10   <none>        Ubuntu 22.04.4 LTS   6.5.0-linuxkit     containerd://1.7.13
    ```
  * Use `dig` to make a DNS query to it
    ```sh
    dig +short @172.18.0.10 podinfo.cloud.example.com
    172.18.0.8
    172.18.0.10
    172.18.0.7
    ```
  * One of your workers should already return DNS responses constructed by Gslb based on service health information
  * If edgeDNS was configured you can query your standard infra DNS directly and it should return the same
    ```sh
    dig +short podinfo.cloud.example.com
    172.18.0.8
    172.18.0.10
    172.18.0.7
    ```
* Now it's time to deploy this application to the first `eu` cluster. The steps and configuration are exactly the same. Just changing `ui.message` to `eu`
```sh
kubectl create ns test-gslb
helm upgrade --install podinfo --namespace test-gslb --set ui.message="eu" podinfo/podinfo
```

* Apply exactly the same Gslb definition
```sh
kubectl -n test-gslb apply -f podinfogslb.yaml
```

* Check the Gslb resource status.
```sh
kubectl -n test-gslb get gslb podinfo -o yaml
apiVersion: k8gb.io/v1beta1
kind: Gslb
metadata:
  name: podinfo
  namespace: test-gslb
spec:
  resourceRef:
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    matchLabels:
      app: podinfo
  strategy:
    type: roundRobin
status:
  geoTag: eu
  healthyRecords:
    podinfo.cloud.example.com:
    - 172.18.0.3
    - 172.18.0.5
    - 172.18.0.6
    - 172.18.0.8
    - 172.18.0.10
    - 172.18.0.7
  serviceHealth:
    podinfo.cloud.example.com: Healthy
```

* Ideally you should already see that `Healthy Records` of `podinfo.cloud.example.com` return the records from __both__ of the clusters. Otherwise, give it a couple of minutes to sync up.

* Now you can check the DNS responses the same way as before.
```sh
dig +short podinfo.cloud.example.com
172.18.0.8
172.18.0.5
172.18.0.10
172.18.0.7
172.18.0.6
172.18.0.3
```

* And for the final end-to-end test, we can use `curl` to query the application
```sh
curl -s podinfo.cloud.example.com | grep message
  "message": "eu",

curl -s podinfo.cloud.example.com | grep message
  "message": "us",

curl -s podinfo.cloud.example.com | grep message
  "message": "us",

curl -s podinfo.cloud.example.com | grep message
  "message": "eu",
```

* As you can see specially marked `podinfo` returns different geo tags showing us the Global Round Robin strategy is working as expected

Hope you enjoyed the ride!

If anything unclear or is going wrong, feel free to contact us at https://github.com/k8gb-io/k8gb/issues. We will appreciate any feedback/bug report and Pull Requests are welcome.

For more advanced technical documentation and fully automated local installation steps, see [Local playground for testing and development](local.md).
