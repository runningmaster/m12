package errors

import (
	"fmt"
	"strings"

	"internal/debug"
	"internal/flag"
	"internal/version"
)

// Locus adds version, file, line, func, etc. from source code
func Locus(err error) error {
	if err == nil {
		return err
	}

	return injectDebug(err)
}

// Locusf replaces fmt.Errorf()
func Locusf(format string, args ...interface{}) error {
	return injectDebug(fmt.Errorf(format, args...))
}

func injectDebug(err error) error {
	if !flag.Debug {
		return err
	}

	ver := version.Stamp.Extended()
	if !strings.Contains(err.Error(), ver) {
		err = fmt.Errorf("%v: version %s", err, ver)
	}

	return fmt.Errorf("%v: at %s", err, debug.FileLineFunc(3))
}
