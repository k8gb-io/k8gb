package e2e

import (
	"os"
	"testing"

	f "github.com/operator-framework/operator-sdk/pkg/test"
)

func TestMain(m *testing.M) {
	if os.Getenv("GITHUB_WORKFLOW") != "" {
		return
	}
	f.MainEntry(m)
}
