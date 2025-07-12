<!-- omit in toc -->
# Local playground for testing and development

- [Environment prerequisites](#environment-prerequisites)
- [Running project locally](#running-project-locally)
- [Verify installation](#verify-installation)
- [Run integration tests](#run-integration-tests)
- [Cleaning](#cleaning)
- [Sample demo](#sample-demo)
  - [Round Robin](#round-robin)
  - [Failover](#failover)

 ---
**NOTE**

This tutorial relies on some makefile targets, to be able to fully understand what's happening under the covers, check the Makefile
[source code](https://github.com/k8gb-io/k8gb/blob/master/Makefile).
Or you may want to run a target first with `-n` switch that will print what it is going to do (example: `make -n test-round-robin`).
For more user-centric targets in that makefile consult `make help`.

---

## Environment prerequisites

- [Install **Go 1.22.3**](https://golang.org/dl/)

- [Install **Git**](https://git-scm.com/downloads)

- [Install **Docker**](https://docs.docker.com/get-docker/)
  > Ensure you are able to push/pull from your docker registry

  > To run multiple clusters, reserve 8GB of memory*
    ![](/docs/images/docker_settings.png)
      <div>
        <sup><i>* above screenshot is provided for <strong>Docker for Mac</strong>, options for other Docker distributions may vary
        </i></sup>
      </div>

 - [install **kubectl**](https://kubernetes.io/docs/tasks/tools/install-kubectl/) to operate k8s clusters

 - [install **helm3**](https://helm.sh/docs/intro/install/) to deploy k8gb and related test workloads

 - [install **k3d**](https://k3d.io/#installation) to run local [k3s](https://k3s.io/) clusters (minimum v5.3.0 version is required)

 - [install **golangci-lint**](https://golangci-lint.run/usage/install/#local-installation) for code quality checks

## Running project locally

To spin-up a local environment using two k3s clusters and deploy a test application to both clusters, execute the command below:
```sh
make deploy-full-local-setup
```

## Verify installation

If local setup runs well, check if clusters are correctly installed

```sh
kubectl cluster-info --context k3d-edgedns && kubectl cluster-info --context k3d-test-gslb1 && kubectl cluster-info --context k3d-test-gslb2
```

Cluster [test-gslb1](https://github.com/k8gb-io/k8gb/tree/master/k3d/test-gslb1.yaml) is exposing external DNS on default port `:5053`
while [test-gslb2](https://github.com/k8gb-io/k8gb/tree/master/k3d/test-gslb2.yaml) on port `:5054`.

Cluster [edgedns](https://github.com/k8gb-io/k8gb/tree/master/k3d/edge-dns.yaml) runs BIND and acts as EdgeDNS holding Delegated Zone for out test setup and answers
on port `:1053`.

```sh
dig @localhost -p 1053 roundrobin.cloud.example.com +short +tcp
```
Should return ***two A records*** from both clusters (IP addresses and order may differ):
```
172.20.0.2
172.20.0.5
172.20.0.4
172.20.0.6
```

You can verify that correct IP addresses of all the nodes in both clusters were populated:
```sh
for c in k3d-test-gslb{1,2}; do kubectl get no -ocustom-columns="NAME:.metadata.name,IP:status.addresses[0].address" --context $c; done
```

Returns a result similar to:
```
NAME                      IP
k3d-test-gslb1-agent-0    172.20.0.2
k3d-test-gslb1-server-0   172.20.0.4
NAME                      IP
k3d-test-gslb2-server-0   172.20.0.6
k3d-test-gslb2-agent-0    172.20.0.5
```

Or you can ask specific CoreDNS instance for its local targets:
```sh
dig -p 5053 +tcp @localhost localtargets-roundrobin.cloud.example.com && \
dig -p 5054 +tcp @localhost localtargets-roundrobin.cloud.example.com
```
As expected result you should see **two A records** divided between both clusters.
```sh
...
...
;; ANSWER SECTION:
localtargets-roundrobin.cloud.example.com. 30 IN A 172.20.0.4
localtargets-roundrobin.cloud.example.com. 30 IN A 172.20.0.2
...
...
localtargets-roundrobin.cloud.example.com. 30 IN A 172.20.0.5
localtargets-roundrobin.cloud.example.com. 30 IN A 172.20.0.6
```
Both clusters have [podinfo](https://github.com/stefanprodan/podinfo) installed on the top.
Run following command and check if you get two json responses.
```sh
curl localhost:80 -H "Host:roundrobin.cloud.example.com" && curl localhost:81 -H "Host:roundrobin.cloud.example.com"
```

## Run integration tests

There is wide range of scenarios which **GSLB** provides and all of them are covered within [tests](https://github.com/k8gb-io/k8gb/tree/master/terratest).
To check whether everything is running properly execute [terratest](https://terratest.gruntwork.io/) :

```sh
make terratest
```

## Cleaning

Clean up your local development clusters with
```sh
make destroy-full-local-setup
```

## Sample demo

### Round Robin

Both clusters have [podinfo](https://github.com/stefanprodan/podinfo) installed on the top, where each
cluster has been tagged to serve a different region. In this demo we will hit podinfo by `wget -qO - roundrobin.cloud.example.com` and depending
on the region, podinfo will return **us** or **eu**. In the current round robin implementation IP addresses are randomly picked.
See [Gslb manifest with round robin strategy](https://github.com/k8gb-io/k8gb/tree/master/deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr_roundrobin_ingress.yaml)

Try to run the following command several times and watch the `message` field.
```sh
make test-round-robin
```
As expected result you should see podinfo message changing

```text
{
  "hostname": "frontend-podinfo-856bb46677-8p45m",
  ...
  "message": "us",
  ...
}
```
```text
{
  "hostname": "frontend-podinfo-856bb46677-8p45m",
  ...
  "message": "eu",
  ...
}
```

### Failover

Both clusters have [podinfo](https://github.com/stefanprodan/podinfo) installed on the top where each
cluster has been tagged to serve a different region. In this demo we will hit podinfo by `wget -qO - failover.cloud.example.com` and depending
on whether podinfo is running inside the cluster it returns only **eu** or **us**.
See [Gslb manifest with failover strategy](https://github.com/k8gb-io/k8gb/tree/master/deploy/crds/k8gb.absa.oss_v1beta1_gslb_cr_failover_ingress.yaml)

Switch GLSB to failover mode:
```sh
make init-failover
```
Now both clusters are running in failover mode and podinfo is running on both of them.
Run several times command below and watch `message` field.
```sh
make test-failover
```
You will see only **eu** podinfo is responsive:
```text
{
  "hostname": "frontend-podinfo-856bb46677-8p45m",
  ...
  "message": "eu",
  ...
}
```
Stop podinfo on **current (eu)** cluster:
```
make stop-test-app
```
Several times hit application again
```sh
make test-failover
```
As expected result you should see only podinfo from **second cluster (us)** is responding:
```text
{
  "hostname": "frontend-podinfo-856bb46677-v5nll",
  ...
  "message": "us",
  ...
}
```
It might happen that podinfo will be unavailable for a while due to
[DNS sync interval](https://github.com/k8gb-io/k8gb/pull/81) and default k8gb DNS TTL of 30 seconds
```text
wget: server returned error: HTTP/1.1 503 Service Temporarily Unavailable
```
Start podinfo again on **current (eu)** cluster:
```sh
make start-test-app
```
and hit several times hit podinfo:
```sh
make test-failover
```
After DNS sync interval is over **eu** will be back
```text
{
  "hostname": "frontend-podinfo-6945c9ddd7-xksrc",
  ...
  "message": "eu",
  ...
}
```
Optionally you can switch GLSB back to round-robin mode
```sh
make init-round-robin
```
