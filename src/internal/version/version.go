//go:generate go run version_generator.go

package version

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	major      = 2
	minor      = 0
	patch      = 5
	prerelease = ""
)

// String returns the version according to http://semver.org/
func String() string {
	v := fmt.Sprintf("%d.%d.%d-%s", major, minor, patch, prerelease)
	if strings.HasSuffix(v, "-") {
		return v[:len(v)-1]
	}
	return v
}

// WithBuildInfo returns the version with build metadata
func WithBuildInfo() string {
	return fmt.Sprintf("%s+%s.%s", String(), BuildTime, GitCommit)
}

// AppName return application name
func AppName() string {
	return filepath.Base(os.Args[0])
}
