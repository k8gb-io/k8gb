package version

import "os"

var (
	//Version of K8GB release
	Version = os.Getenv("K8GB_VERSION")
)
