//+build !windows

package vfs

import "syscall"

// HostOSFS is the host-specific OSFS.
var HostOSFS = OSFS

var ignoreErrnoInContains = map[syscall.Errno]struct{}{
	syscall.ELOOP:        {},
	syscall.EMLINK:       {},
	syscall.ENAMETOOLONG: {},
	syscall.ENOENT:       {},
	syscall.EOVERFLOW:    {},
}

// relativizePath, on POSIX systems, just returns path.
func relativizePath(path string) string {
	return path
}
