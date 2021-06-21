# K8GB terratests
Terratests are another layer in k8gb testing. They test the state of k8gb on tuple of clusters and examine the behavior
for individual use cases. They replace manual tests; install GSLB in a separate namespace and check if everything
behaves as expected.

Terratest is basically a separate project inside k8gb. The tests are running as GitHub Actions when interacting
with remote repo, or you can run them locally with `make terratest`. To run it locally, you must have k3d clusters
with the k8gb operator installed. Read [Quick Start](../README.md#quick-start) if you have not already done so.

Terratests consist of three directories:

- `/test/` contains the individual tests. Each test usually creates its own namespace where it installs its own
  version of GSLB. The test then runs against such instance. For optimization reasons, all tests are running
  in parallel. Keep this on mind!
- `/examples/` contains yaml configurations (GSLB, ingresses) for individual tests.
- `/utils/` contains a common framework that makes writing tests easier.

## Terratest Common Framework
First of all, nothing forces you to use the framework. There are certainly many cases where it is better not to use
the framework and instead interact with the cluster directly using terratest (which is itself a powerful framework).
On the other hand, the ability to quickly spin up clusters with test applications and then easily read ip addresses
from DNS, without a deeper knowledge of k8gb or terratest is not a bad thing either.

The terratest framework is located in `/utils/` and contains a fluent-style configuration of the GSLB cluster instance.

### Workflows
First of all, we create a workflow instance that includes a namespace, and may include a test application or a GSLB object.
We create the instance by calling the `NewWorkflow` function with the name of the cluster where the instance will run and t
he port from which the cluster is accessed. If you followed [Quick Start](../README.md#quick-start), the clusters will
be `k3d-test-gslb1:5053` and `k3d-test-gslb2:5054`. This is optionally followed by the GSLB configuration from the yaml
file and the test application [podinfo](https://github.com/stefanprodan/podinfo). The `Start` function creates the
cluster resources and returns the workflow instance, the `Kill` function deletes the instance along with all resources
in the test cluster.
```go
instance, err := utils.NewWorkflow(t, "k3d-test-gslb1", 5053).
    WithGslb(gslbPath, host).
    WithTestApp("us").
    Start()
require.NoError(t, err)
defer instance1.Kill()
```

#### Workflow functions
Once the instance is created, we can change the state of the test application `instance.StopTestApp()`, 
`instance.StartTestApp()`. We can also read the ingress IPs `GetIngressIPs` local targets `GetLocalTargets()`, or call `Dig()` 
against the CoreDNS cluster where the instance is running. Finally, I would like to mention `HitTestApp` which is a function 
that makes an http request to a test application and returns its status.

### Workflow status
Terratests are difficult in principle due to race conditions. For example, it may happen that the test application 
has not started yet and you read invalid localtargets (see `WaitForAppIsRunning` in the following section).
For this reason, the instance has a `GetStatus` function that returns various information about the running instance.
`GetStatus(annotation).String()` returns formatted json.

#### WaitForGSLB, WaitForExpected, WaitForAppIsRunning
**WaitForGSLB()** and **WaitForExpected** do basically the same thing. They wait for the state on the called instance 
to match the expected state. If this doesn't happen within some time, it returns an error. If the function succeeds, 
it returns the expected (unordered) slice of IP addresses.

 - `WaitForExpected(expectedIPs []string) ([]string, error)` - The function will wait until Dig on the given workflow instance 
   returns the expected IP addresses (remember, the workflow instance is running against a specific cluster).

  - `WaitForGSLB(instances ...*Instance) ([]string, error)` - The function works similarly to `WaitForExpected`, but instead 
    of receiving the IP address slice, it receives the slice of other workflow instances and then calculates the expected IP 
    address slice from their LocalTargets. In practice it looks like this: `instance1.WaitForGSLB(instance2, instance3)` produces: 
    `desiredIPList := instance1.GetLocalTargets() + instance2.GetLocalTargets() + instance3.GetLocalTargets()`.

  - `WaitForAppIsRunning()` - The function is usually used during initialization, when you need to make sure that the test application is actually running. 
    After that you can read the cluster values to be used for further assertion. 

### Demo
As a demo we use a short test that prepares two clusters with a test application. Then it kills the application on the
first cluster and examines what happens.
```go
// common_framework_test.go
func TestCommonFramework(t *testing.T) {
	t.Parallel()
	const host = "terratest-failover.cloud.example.com"
	const gslbPath = "../examples/failover.yaml"

	// create namespace with failover Gslb on k3d-test-gslb1
	instanceFailoverEU, err := utils.NewWorkflow(t, "k3d-test-gslb1", 5053).
		WithGslb(gslbPath, host).
		WithTestApp("eu").
		Start()
	require.NoError(t, err)
	defer instanceFailoverEU.Kill()

	// create namespace with failover Gslb on k3d-test-gslb2
	instanceFailoverUS, err := utils.NewWorkflow(t, "k3d-test-gslb2", 5054).
		WithGslb(gslbPath, host).
		WithTestApp("us").
		Start()
	require.NoError(t, err)
	defer instanceFailoverUS.Kill()

	t.Run("ensure clusters are ready", func(t *testing.T) {
		err = instanceFailoverEU.WaitForAppIsRunning()
		require.NoError(t, err)
		err = instanceFailoverUS.WaitForAppIsRunning()
		require.NoError(t, err)
	})

	localTargetsEU := instanceFailoverEU.GetLocalTargets() // e.g: [10.43.78.154, 10.43.78.155]
	localTargetsUS := instanceFailoverUS.GetLocalTargets() // e.g: [10.43.150.206, 10.43.150.207]

	t.Run("stop podinfo on the first cluster", func(t *testing.T) {
		// digging EU, US cluster returns EU localtargets
		require.True(t, utils.EqualStringSlices(instanceFailoverEU.Dig(), localTargetsEU))
		require.True(t, utils.EqualStringSlices(instanceFailoverUS.Dig(), localTargetsEU))

		instanceFailoverEU.StopTestApp()
		err = instanceFailoverEU.WaitForExpected(localTargetsUS)
		require.NoError(t, err)
		err = instanceFailoverUS.WaitForExpected(localTargetsUS)
		require.NoError(t, err)

		// digging EU, US cluster returns US localtargets
		require.True(t, utils.EqualStringSlices(instanceFailoverEU.Dig(), localTargetsUS))
		require.True(t, utils.EqualStringSlices(instanceFailoverUS.Dig(), localTargetsUS))
		require.Empty(t, instanceFailoverEU.GetLocalTargets())
		require.True(t, utils.EqualStringSlices(instanceFailoverUS.GetLocalTargets(), localTargetsUS))
	})
}
```

### Troubleshoot
In my experience, most of the bugs come from not upgrading the local clusters. Consider running `make reset upgrade-candidate` 
before you start writing a test. The framework is still under development and there may be some 
functionality that needs to be extended. Feel free to contribute and create pull requests / issues. 
