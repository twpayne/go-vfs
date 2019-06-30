// +build windows

package vfs

import (
	"path/filepath"
	"strings"
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

// trimPrefix, on Windows, trims prefix from path and returns an absolute path.
// prefix must be a /-separated path.
func trimPrefix(path, prefix string) (string, error) {
	trimmedPath, err := filepath.Abs(strings.TrimPrefix(filepath.ToSlash(path), prefix))
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(trimmedPath), nil
}
