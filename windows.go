//+build windows

package vfs

import (
	"syscall"

	"golang.org/x/sys/windows"
)

// HostOSFS is the host-specific OSFS.
var HostOSFS = WindowsOSFS{}

func shouldSkipSystemError(err syscall.Errno) bool {
	return err == windows.ERROR_CANT_RESOLVE_FILENAME
}
