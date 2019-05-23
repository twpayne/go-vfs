//+build windows

package vfs

import (
	"os"

	acl "github.com/hectane/go-acl"
)

// HostOSFS is the host-specific OSFS.
var HostOSFS = WindowsOSFS{}

type WindowsOSFS struct {
	osfs
}

func (WindowsOSFS) Chmod(name string, mode os.FileMode) error {
	return acl.Chmod(name, mode)
}
