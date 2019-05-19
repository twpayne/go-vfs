// +build windows

package vfs

import (
	"syscall"

	"golang.org/x/sys/windows"
)

// return true if the error should be skipped, false otherwise
func shouldSkipSystemError(err syscall.Errno) bool {
	return err == windows.ERROR_CANT_RESOLVE_FILENAME
}
