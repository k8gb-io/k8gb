package test

import (
	"github.com/stretchr/testify/require"
	"k8gbterratest/utils"
	"testing"
)

func TestMultipleDelegatedZonesOnStartup(t *testing.T) {
	const basePath = "../examples/failover-playground"
	const host = "playground-failover.cloud.example.com"

	workflowEU := utils.NewWorkflow(t, "k3d-test-gslb1", 5053)
	workflowUS := utils.NewWorkflow(t, "k3d-test-gslb2", 5054)
	abstractTestMultipleDelegatedZones(t,
		workflowEU.Enrich(basePath, host, utils.IngressEmbedded),
		workflowUS.Enrich(basePath, host, utils.IngressEmbedded),
	)
}

func abstractTestMultipleDelegatedZones(t *testing.T, workflowEU, workflowUS *utils.Workflow) {
	const euGeoTag = "eu"
	const usGeoTag = "us"

	instanceEU, err := workflowEU.Start()
	require.NoError(t, err)
	defer instanceEU.Kill()

	instanceUS, err := workflowUS.Start()
	require.NoError(t, err)
	defer instanceUS.Kill()

	t.Run("", func(t *testing.T) {
		err = instanceEU.WaitForAppIsRunning()
		require.NoError(t, err)
		err = instanceUS.WaitForAppIsRunning()
		require.NoError(t, err)
	})

}
