package helpers

import (
	"io"
)

type ilogger interface {
	Error(format string, args ...interface{})
}

// DeferClose handle error for closer interface
func DeferClose(f io.Closer, l ilogger) func() {
	return func() {
		err := f.Close()
		if err != nil {
			l.Error("error", err)
		}
	}
}
