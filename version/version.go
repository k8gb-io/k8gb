package version

import "os"

var (
	//Version of OhMyGlb release
	Version = os.Getenv("OHMYGLB_VERSION")
)
