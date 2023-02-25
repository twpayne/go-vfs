//go:build !windows
// +build !windows

package vfst_test

import (
	"syscall"
)

func init() {
	syscall.Umask(0o22)
}
