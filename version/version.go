package version

import "os"

var (
	//Version of KGB release
	Version = os.Getenv("KGB_VERSION")
)
