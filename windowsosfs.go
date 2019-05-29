//+build windows

package vfs

import (
	"io/ioutil"
	"os"

	"github.com/hectane/go-acl"
)

type WindowsOSFS struct {
	osfs
}

// Chmod implements os.Chmod.
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
