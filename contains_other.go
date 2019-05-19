// +build !windows

package vfs

import "syscall"

func shouldSkipSystemError(err syscall.Errno) bool {
	return false
}
