//+build windows

package vfs

import (
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows"
)

// HostOSFS is the host-specific OSFS.
var HostOSFS = WindowsOSFS{}

// relativizePath, on Windows, strips any leading volume name from path and
// replaces backslashes with slashes.
func relativizePath(path string) string {
	if volumeName := filepath.VolumeName(path); volumeName != "" {
		path = path[len(volumeName):]
	}
	return filepath.ToSlash(path)
}

func shouldSkipSystemError(err syscall.Errno) bool {
	return err == windows.ERROR_CANT_RESOLVE_FILENAME
}
