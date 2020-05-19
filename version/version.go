package version

import "os"

var (
	Version = os.Getenv("OHMYGLB_VERSION")
)
