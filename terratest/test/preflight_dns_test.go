package test

/*
Copyright 2025 The k8gb Contributors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// TestLocalDNSTCPConnectivity verifies the local k3d load balancer forwards TCP/53
// from localhost:5053 (and :5054 for two clusters) to the CoreDNS Service.
// This is a preflight to provide clear diagnostics when "dig +tcp @localhost -p 5053" fails.
func TestLocalDNSTCPConnectivity(t *testing.T) {
	dial := func(port int) error {
		address := fmt.Sprintf("127.0.0.1:%d", port)
		conn, err := net.DialTimeout("tcp", address, 3*time.Second)
		if err != nil {
			return err
		}
		_ = conn.Close()
		return nil
	}

	if err := dial(settings.Port1); err != nil {
		t.Fatalf("cannot connect to localhost:%d over TCP: %v\n\nHints:\n- Ensure k3d >= v5.3.0\n- Recreate local setup: 'make destroy-full-local-setup && make deploy-full-local-setup'\n- Verify CoreDNS Service is LoadBalancer: 'make verify-dns-lb'\nSee docs/local.md (Troubleshooting: localhost:5053/5054 +tcp fails).", settings.Port1, err)
	}

	if settings.ClustersNumber >= 2 {
		if err := dial(settings.Port2); err != nil {
			t.Fatalf("cannot connect to localhost:%d over TCP: %v\n\nHints:\n- Ensure k3d >= v5.3.0\n- Recreate local setup: 'make destroy-full-local-setup && make deploy-full-local-setup'\n- Verify CoreDNS Service is LoadBalancer: 'make verify-dns-lb'\nSee docs/local.md (Troubleshooting: localhost:5053/5054 +tcp fails).", settings.Port2, err)
		}
	}
}


