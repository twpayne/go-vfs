//+build windows

package vfs

import (
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows"
)

// HostOSFS is the host-specific OSFS.
var HostOSFS = WindowsOSFS{}

var ignoreErrnoInContains = map[syscall.Errno]struct{}{
	syscall.ELOOP:                       {},
	syscall.EMLINK:                      {},
	syscall.ENAMETOOLONG:                {},
	syscall.ENOENT:                      {},
	syscall.EOVERFLOW:                   {},
	windows.ERROR_CANT_RESOLVE_FILENAME: {},
}

// relativizePath, on Windows, strips any leading volume name from path and
// replaces backslashes with slashes.
func relativizePath(path string) string {
	if volumeName := filepath.VolumeName(path); volumeName != "" {
		path = path[len(volumeName):]
	}
	return filepath.ToSlash(path)
}
