//go:generate go run gen.go

package version

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Stamp is implementetion of Stamper interface
var Stamp Stamper

// Stamper is implemented the information that specifies the company name,
// application name, copyright, version number, language edition of a program, etc.
type Stamper interface {
	AppName() string
	Extended() string
	Reset(int, int, int, string)
	fmt.Stringer
}

func init() {
	Stamp = stamp{}
}

type stamp struct{}

// FIXME (for testing golint)
func (s stamp) AppName() string {
	return filepath.Base(os.Args[0])
}

// String returns the version according to http://semver.org/
func (s stamp) String() string {
	v := strings.Join(
		[]string{
			strconv.Itoa(major),
			strconv.Itoa(minor),
			strconv.Itoa(patch),
		},
		".",
	)
	if prerelease != "" {
		v = fmt.Sprintf("%s-%s", v, prerelease)
	}
	return v
}

// Extended returns the version with build metadata
func (s stamp) Extended() string {
	return fmt.Sprintf("%s+%s.%s", s.String(), buildtime, gitcommit)
}

// Reset changes default autogenerated main version values
func (s stamp) Reset(mjr, mnr, ptch int, pre string) {
	r := strings.NewReplacer(
		"-", "",
		"+", "",
	)
	major, minor, patch, prerelease = mjr, mnr, ptch, r.Replace(pre)
}
