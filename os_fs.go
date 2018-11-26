package vfs

import (
	"io/ioutil"
	"os"
)

type osfs struct{}

// OSFS is the FS that calls os and ioutil functions directly.
var OSFS = &osfs{}

// Chmod implements os.Chmod.
func (osfs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

// Lstat implements os.Lstat.
func (osfs) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

// Mkdir implements os.Mkdir.
func (osfs) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// ReadDir implenents ioutil.ReadDir.
func (osfs) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

// ReadFile implements ioutil.ReadFile.
func (osfs) ReadFile(dirname string) ([]byte, error) {
	return ioutil.ReadFile(dirname)
}

// Readlink implments os.Readlink.
func (osfs) Readlink(name string) (string, error) {
	return os.Readlink(name)
}

// Remove implements os.Remove.
func (osfs) Remove(name string) error {
	return os.Remove(name)
}

// Stat implements os.Stat.
func (osfs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// Symlink implements os.Symlink.
func (osfs) Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

// WriteFile implements ioutil.WriteFile.
func (osfs) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}
