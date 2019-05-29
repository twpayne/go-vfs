//+build !windows

package vfs

import "syscall"

// HostOSFS is the host-specific OSFS.
var HostOSFS = OSFS

// relativizePath, on POSIX systems, just returns path.
func relativizePath(path string) string {
	return path
}

func shouldSkipSystemError(err syscall.Errno) bool {
	return false
}
