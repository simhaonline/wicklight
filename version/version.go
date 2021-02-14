package version

import (
	"os"
	"runtime"
)

// Version variables
var (
	Version   string
	BuildTime string
)

// Print print version
func Print() string {
	return os.Args[0] + ": " + Version + " build: " + BuildTime + " (platform: " + runtime.GOOS + "-" + runtime.GOARCH + ")."
}

// Help print help
func Help() string {
	return "\tUsage: " + os.Args[0] + " config.yaml"
}
