//+build !windows

package vfs

import "syscall"

// HostOSFS is the host-specific OSFS.
var HostOSFS = OSFS

func shouldSkipSystemError(err syscall.Errno) bool {
	return false
}
