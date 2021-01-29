// +build !windows

package vfst

import "syscall"

func init() {
	syscall.Umask(0o22)
}
