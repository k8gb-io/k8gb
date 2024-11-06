# Contribfest exercise Kubecon NA 2024

In this exercise we will extend K8GB's ingress integrations.
The project already supports integrations for `Ingress` and `Istio VirtualService` resources. In this session we will add support for `Kubernetes Services` of the type `LoadBalancer`.

To support you with the implementation we provide unit tests and e2e tests, and the places where you should add your code are marked with `FIXME`.

In case of any doubts we are here to help, don't hesitate to ask questions.

## Intro

Let's start with an introduction to become familiar with the problem we are trying to solve.

### Getting familiar with the GSLB resource

The GSLB resource is used to connect a load balancing strategy with an application running on the cluster. In this exercise we will reference the application using a `Service` of type `LoadBalancer`. The GSLB resource will look as follows:
```yaml
apiVersion: k8gb.absa.oss/v1beta1
kind: Gslb
metadata:
  name: lb-service
  namespace: test-gslb
spec:
  resourceRef:
    apiVersion: v1
    kind: Service
    matchLabels:
      app: lb-service
  strategy:
    type: roundRobin
    splitBrainThresholdSeconds: 300
    dnsTtlSeconds: 30
```

Strategies are not part of the scope of this exercise. But if you are curious you can read more about them [here](https://www.k8gb.io/docs/strategy.html).
We are interested in the `resourceRef` block, where a `Service` is referenced using label selectors.

### Reconciliation loop

The K8GB controller has a reconcile function that is responsible for reconciling the desired state of the GSLB resource with the real world.
The code can be found in [controllers/gslb_controller_reconciliation.go](../../controllers/gslb_controller_reconciliation.go).
We are particularly interested in this block:
```golang
	refResolver, err := refresolver.New(gslb, r.Client)

	servers, err := refResolver.GetServers(r.Config.DNSZone)
	gslb.Status.Servers = servers

	loadBalancerExposedIPs, err := refResolver.GetGslbExposedIPs(r.Config.EdgeDNSServers)
	gslb.Status.LoadBalancer.ExposedIPs = loadBalancerExposedIPs
```

First, a `refResolver` is instantiated from a `GSLB` resource. Then, the `refResolver` is used to fetch the `servers` (domain name and `Kubernetes Service` powering the application) and `exposedIPs`. This is the only information the controller needs to know about the application. The rest of the code is completely decoupled from the type of ingress used.

This means we only need to implement the functions above to contribute with a new integration, everything else will automagically work.
Here is the interface that can be found at [controller/refresolver/refresolver.go](../../controllers/refresolver/refresolver.go):
```golang
type GslbReferenceResolver interface {
	// GetServers retrieves GSLB the server configuration
	GetServers(dnsZone string) ([]*k8gbv1beta1.Server, error)
	// GetGslbExposedIPs retrieves the load balancer IP address of the GSLB
	GetGslbExposedIPs(edgeDNSServers utils.DNSList) ([]string, error)
}
```

Let's proceed with the implementation!

## Implementation

### Instantiating a `lbservice refResolver` (10 mins)

First, let's handle the constructor. We want to instantiate a `GslbReferenceResolver` of the correct type.

Navigate to [controllers/refresolver/](../../controllers/refresolver/) and run:
```bash
go test
```
You will see the following error:
```bash
actual  : *errors.errorString(&errors.errorString{s:"APIVersion:v1, Kind:Service not supported"})
```

Have a look at [controllers/refresolver/refresolver.go](../../controllers/refresolver/refresolver.go) and replace the `FIXME` comment with a solution. When the test passes move on to the next section.

<details>
  <summary>Solution</summary>

  ```golang
	if gslb.Spec.ResourceRef.Kind == "Service" && gslb.Spec.ResourceRef.APIVersion == "v1" {
		return lbservice.NewReferenceResolver(gslb, k8sClient)
	}
  ```

</details>

<br />

### Creating a `lbservice.Reference` resolver object (20 mins)

Now let's implement the logic of the new `lbservice.ReferenceResolver`. This can be done in the function `NewReferenceResolver` in [controllers/refresolver/lbservice/lbservice.go](../../controllers/refresolver/lbservice/lbservice.go).

In the constructor we want to fetch the `Service` passed in the resourceRef block and store it in the struct, so that we can use it later.

For the sake of this exercise you can assume that there is exactly one `Service` matching the label selector.

<details>
  <summary>Solution</summary>

  ```golang
  func NewReferenceResolver(gslb *k8gbv1beta1.Gslb, k8sClient client.Client) (*ReferenceResolver, error) {
    serviceList := &corev1.ServiceList{}

    // retrieve services
    selector, err := metav1.LabelSelectorAsSelector(&gslb.Spec.ResourceRef.LabelSelector)
    if err != nil {
      return nil, err
    }
    opts := &client.ListOptions{
      LabelSelector: selector,
      Namespace:     gslb.Namespace,
    }
    err = k8sClient.List(context.TODO(), serviceList, opts)
    if err != nil {
      return nil, err
    }

    // filter for type LoadBalancer
    lbServices := []corev1.Service{}
    for _, svc := range serviceList.Items {
      if svc.Spec.Type == v1.ServiceTypeLoadBalancer {
        lbServices = append(lbServices, svc)
      }
    }

    return &ReferenceResolver{
      lbService: &lbServices[0],
    }, nil
  }
  ```
</details>

### Implementing the ReferenceResolver interface

We can now implement the `ReferenceResolver` interface that we saw before.

Navigate to [controllers/refresolver/lbservice](../../controllers/refresolver/lbservice/) and run:
```bash
go test
```
You will see 3 tests failing. Let's fix them!

#### Fetching the Servers (10 mins)

Let's first gather information about the servers running your application. You can see the datastructure in [api/v1beta1/gslb_types.go](../../api/v1beta1/gslb_types.go), it looks as follows:
```golang
type Server struct {
	// Hostname exposed by the GSLB
	Host string `json:"host,omitempty"`
	// Kubernetes Services backing the load balanced application
	Services []*NamespacedName `json:"services,omitempty"`
}
```

Try to implement the function `GetServers` by using the `service` that you fetched in the previous step and the `dnsZone` passed as an argument to the method.

Once you are done the test `TestGetServers/single_server` should no longer be failing.

<details>
  <summary>Solution</summary>

  ```golang
  func (rr *ReferenceResolver) GetServers(dnsZone string) ([]*k8gbv1beta1.Server, error) {
    host := fmt.Sprintf("%s.%s", rr.lbService.Name, dnsZone)
    server := &k8gbv1beta1.Server{
      Host: host,
      Services: []*k8gbv1beta1.NamespacedName{
        {
          Name:      rr.lbService.Name,
          Namespace: rr.lbService.Namespace,
        },
      },
    }
    return []*k8gbv1beta1.Server{server}, nil
  }
  ```
</details>

#### Fetching the exposed IP addresses (10 mins)

The last step is to fetch the IP addresses under which your application is available.
In this case the application is exposed using a `Service` of type `LoadBalancer`, therefore we can fetch the IP addresses from its Status. Go ahead and try to implement the function `GetGslbExposedIPs`.

Note: You can assume that the status of the `LoadBalancer Service` is populated with IP addresses and not with hostnames.

All tests should be passing once you are done.

<details>
  <summary>Solution</summary>

  ```golang
  func (rr *ReferenceResolver) GetGslbExposedIPs(edgeDNSServers utils.DNSList) ([]string, error) {
    gslbIngressIPs := []string{}
    for _, ip := range rr.lbService.Status.LoadBalancer.Ingress {
      gslbIngressIPs = append(gslbIngressIPs, ip.IP)
    }
    return gslbIngressIPs, nil
  }
  ```
</details>

### Testing e2e (15 mins)

All unit tests pass! The implementation is now complete! It is time to make sure everything is working e2e with a real failover scenario.

Let's setup a couple of local clusters to test our code.

First, install the prerequisites:
* [kubectl](https://kubernetes.io/docs/tasks/tools/)
* [helm3](https://helm.sh/docs/intro/install/)
* [k3d](https://k3d.io/v5.7.4/#other-installers)
* [chainsaw](https://kyverno.io/blog/2023/12/12/kyverno-chainsaw-the-ultimate-end-to-end-testing-tool/#install-chainsaw)

You should now be able to setup a local cluster running your code by running the following command:
```bash
K8GB_LOCAL_VERSION=test FULL_LOCAL_SETUP_WITH_APPS=false make deploy-full-local-setup
```

Once the cluster is setup you can run the e2e tests for the `failover` and `roundrobin` scenarios:
```bash
make chainsaw
```

While the tests run you can take a quick peak at the logic in [chainsaw/tests/failover-playground-lbservice/chainsaw-test.yaml](chainsaw/tests/failover-playground-lbservice/chainsaw-test.yaml) and [chainsaw/tests/roundrobin-lbservice/chainsaw-test.yaml](chainsaw/tests/roundrobin-lbservice/chainsaw-test.yaml).

Hopefully all tests are passing :)

# Closing

Awesome, you completed the exercise! You now understand deeply how to integrate ingress configuration in K8GB and are ready to contribute with the new integrations! We are thrilled to have you in the community!
