//+build windows

package vfs

import (
	"io/ioutil"
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

// WriteFile implements ioutil.WriteFile.
func (fs WindowsOSFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	err := ioutil.WriteFile(filename, data, perm)
	if err != nil {
		return err
	}
	return fs.Chmod(filename, perm)
}
