package route53_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type ResourceRecord struct {
	Value string `json:"Value"`
}

type ResourceRecordSet struct {
	Name            string           `json:"Name"`
	Type            string           `json:"Type"`
	TTL             int64            `json:"TTL,omitempty"`
	ResourceRecords []ResourceRecord `json:"ResourceRecords"`
}

type RecordSetsOutput struct {
	ResourceRecordSets []ResourceRecordSet `json:"ResourceRecordSets"`
}

type ChangeBatchItem struct {
	Action            string            `json:"Action"`
	ResourceRecordSet ResourceRecordSet `json:"ResourceRecordSet"`
}

type ChangeBatch struct {
	Changes []ChangeBatchItem `json:"Changes"`
}

func runCmd(t *testing.T, dir string, env []string, name string, args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		t.Logf("Command failed: %s %s\nStdout: %s\nStderr: %s", name, strings.Join(args, " "), outBuf.String(), errBuf.String())
		t.Fatalf("Failed to run %s: %v", name, err)
	}
	return strings.TrimSpace(outBuf.String())
}

func runCmdOutput(dir string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, errBuf.String())
	}
	return strings.TrimSpace(outBuf.String()), nil
}

func TestRoute53Integration(t *testing.T) {
	awsKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecret := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsKey == "" || awsSecret == "" {
		t.Skip("Skipping Route53 integration test: AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY not set")
	}

	workDir, err := filepath.Abs("../../")
	require.NoError(t, err)

	gopath, err := runCmdOutput("", "go", "env", "GOPATH")
	if err == nil && gopath != "" {
		gopathBin := filepath.Join(gopath, "bin")
		_ = os.Setenv("PATH", fmt.Sprintf("%s:%s", gopathBin, os.Getenv("PATH")))
	}

	opentoFuDir := filepath.Join(workDir, "dns-provider-test", "route53", "opentofu")

	var zoneID string

	// Defer cleanup execution
	defer func() {
		t.Log("---- Starting Teardown & Cleanup ----")
		if zoneID != "" {
			t.Logf("Cleaning up remaining record sets in Hosted Zone %s...", zoneID)
			recordsRaw, err := runCmdOutput("", "aws", "route53", "list-resource-record-sets", "--hosted-zone-id", zoneID, "--output", "json")
			if err == nil && recordsRaw != "" {
				var recordOut RecordSetsOutput
				if err := json.Unmarshal([]byte(recordsRaw), &recordOut); err == nil {
					var toDelete []ResourceRecordSet
					for _, rr := range recordOut.ResourceRecordSets {
						nameClean := strings.TrimSuffix(rr.Name, ".")
						if (rr.Type == "NS" && nameClean == "k8gb.io") || (rr.Type == "SOA" && nameClean == "k8gb.io") {
							continue
						}
						toDelete = append(toDelete, rr)
					}
					if len(toDelete) > 0 {
						batch := ChangeBatch{}
						for _, item := range toDelete {
							batch.Changes = append(batch.Changes, ChangeBatchItem{
								Action:            "DELETE",
								ResourceRecordSet: item,
							})
						}
						batchJSON, _ := json.Marshal(batch)
						changeRaw, err := runCmdOutput("", "aws", "route53", "change-resource-record-sets", "--hosted-zone-id", zoneID, "--change-batch", string(batchJSON), "--query", "ChangeInfo.Id", "--output", "text")
						if err == nil && changeRaw != "" && changeRaw != "None" {
							runCmdOutput("", "aws", "route53", "wait", "record-sets-changed", "--id", changeRaw)
						}
					}
				}
			}
		}

		t.Log("Destroying OpenTofu infrastructure...")
		runCmdOutput(opentoFuDir, "tofu", "destroy", "-auto-approve", "-var=dns_zone_name=k8gb.io")

		os.Remove(filepath.Join(workDir, "dns-provider-test", "route53", "values.yaml"))
		os.Remove(filepath.Join(workDir, "credentials"))

		t.Log("Destroying local k3d clusters...")
		runCmdOutput(workDir, "make", "destroy-full-local-setup")
		t.Log("---- Teardown Finished ----")
	}()

	t.Log("---- Step 1: Provisioning Route53 Zone via OpenTofu ----")
	runCmd(t, opentoFuDir, nil, "tofu", "init", "-upgrade")
	runCmd(t, opentoFuDir, nil, "tofu", "apply", "-auto-approve", "-var=dns_zone_name=k8gb.io")

	t.Log("---- Step 2: Extracting Canonical Route53 Hosted Zone ID ----")
	zoneID = runCmd(t, opentoFuDir, nil, "tofu", "output", "-raw", "zone_id")
	require.NotEmpty(t, zoneID, "Hosted Zone ID from OpenTofu output must not be empty")
	t.Logf("Canonical Hosted Zone ID: %s", zoneID)

	t.Log("---- Step 3: Resolving Route53 Zone Nameserver ----")
	edgeDNSRaw := runCmd(t, workDir, nil, "aws", "route53", "list-resource-record-sets", "--hosted-zone-id", zoneID, "--query", "ResourceRecordSets[?Type == 'NS'].ResourceRecords[0]", "--output", "text")
	edgeDNSServer := strings.TrimSuffix(strings.TrimSpace(edgeDNSRaw), ".")
	require.NotEmpty(t, edgeDNSServer, "Resolved Edge DNS server must not be empty")
	t.Logf("Edge DNS Server: %s", edgeDNSServer)

	t.Log("---- Step 4: Generating Shared Helm Values File ----")
	templatePath := filepath.Join(workDir, "dns-provider-test", "route53", "values-template.yaml")
	templateContent, err := os.ReadFile(templatePath)
	require.NoError(t, err)

	valuesContent := strings.ReplaceAll(string(templateContent), "DNS_SERVER_TODO", edgeDNSServer)
	valuesContent = strings.ReplaceAll(valuesContent, "ZONE_ID_TODO", zoneID)

	valuesPath := filepath.Join(workDir, "dns-provider-test", "route53", "values.yaml")
	err = os.WriteFile(valuesPath, []byte(valuesContent), 0644)
	require.NoError(t, err)

	t.Log("---- Step 5: Creating Local k3d Clusters & Building Snapshot Images ----")
	runCmd(t, workDir, nil, "make", "create-local-clusters")
	runCmd(t, workDir, nil, "make", "release-images")

	t.Log("---- Step 6: Deploying AWS Credentials Secret to Clusters ----")
	credContent := fmt.Sprintf("[default]\naws_access_key_id = %s\naws_secret_access_key = %s\n", awsKey, awsSecret)
	credPath := filepath.Join(workDir, "credentials")
	err = os.WriteFile(credPath, []byte(credContent), 0600)
	require.NoError(t, err)

	runCmd(t, workDir, nil, "kubectl", "create", "ns", "k8gb", "--context", "k3d-test-gslb1")
	runCmd(t, workDir, nil, "kubectl", "create", "ns", "k8gb", "--context", "k3d-test-gslb2")
	runCmd(t, workDir, nil, "kubectl", "create", "secret", "generic", "external-dns-secret-aws", "-n", "k8gb", "--from-file", credPath, "--context", "k3d-test-gslb1")
	runCmd(t, workDir, nil, "kubectl", "create", "secret", "generic", "external-dns-secret-aws", "-n", "k8gb", "--from-file", credPath, "--context", "k3d-test-gslb2")

	t.Log("---- Step 7: Deploying k8gb Operator ----")
	deployEnv := []string{
		"VALUES_YAML=dns-provider-test/route53/values.yaml",
		"K8GB_LOCAL_VERSION=test",
		"DEPLOY_APPS=false",
	}
	runCmd(t, workDir, deployEnv, "make", "deploy-test-version")

	t.Log("---- Step 8: Asserting A Records on Canonical Zone ID ----")
	require.Eventually(t, func() bool {
		recordsRaw, err := runCmdOutput(workDir, "aws", "route53", "list-resource-record-sets", "--hosted-zone-id", zoneID, "--query", "ResourceRecordSets[?Type == 'A']", "--output", "json")
		if err != nil || recordsRaw == "" {
			return false
		}
		var recordOut []ResourceRecordSet
		if err := json.Unmarshal([]byte(recordsRaw), &recordOut); err != nil {
			return false
		}

		euIPs, usIPs := 0, 0
		for _, rr := range recordOut {
			if strings.Contains(rr.Name, "gslb-ns-eu-cloud.k8gb.io") {
				euIPs = len(rr.ResourceRecords)
			}
			if strings.Contains(rr.Name, "gslb-ns-us-cloud.k8gb.io") {
				usIPs = len(rr.ResourceRecords)
			}
		}

		t.Logf("Polling A Records... EU IP count: %d/2, US IP count: %d/2", euIPs, usIPs)
		return euIPs == 2 && usIPs == 2
	}, 3*time.Minute, 10*time.Second, "Timed out waiting for EU and US A records to be published with 2 IPs each")

	t.Log("---- Step 9: Asserting NS Records on Canonical Zone ID ----")
	require.Eventually(t, func() bool {
		recordsRaw, err := runCmdOutput(workDir, "aws", "route53", "list-resource-record-sets", "--hosted-zone-id", zoneID, "--query", "ResourceRecordSets[?Type == 'NS']", "--output", "json")
		if err != nil || recordsRaw == "" {
			return false
		}
		var recordOut []ResourceRecordSet
		if err := json.Unmarshal([]byte(recordsRaw), &recordOut); err != nil {
			return false
		}

		hasRootNS := false
		hasCloudNS := false
		nsTargetFound := false

		for _, rr := range recordOut {
			if strings.TrimSuffix(rr.Name, ".") == "k8gb.io" {
				hasRootNS = true
			}
			if strings.Contains(rr.Name, "cloud.k8gb.io") {
				hasCloudNS = true
				for _, val := range rr.ResourceRecords {
					if strings.Contains(val.Value, "gslb-ns-eu-cloud.k8gb.io") || strings.Contains(val.Value, "gslb-ns-us-cloud.k8gb.io") {
						nsTargetFound = true
					}
				}
			}
		}

		t.Logf("Polling NS Records... RootNS: %v, CloudNS: %v, DelegationTargetFound: %v", hasRootNS, hasCloudNS, nsTargetFound)
		return hasRootNS && hasCloudNS && nsTargetFound
	}, 1*time.Minute, 5*time.Second, "Timed out waiting for delegation NS records to be published")

	t.Log("---- Route53 Integration Test Passed Successfully ----")
}
