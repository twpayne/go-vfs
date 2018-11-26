package vfs

import (
	"io/ioutil"
	"os"
)

type osfs struct{}

var OSFS = &osfs{}

func (osfs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

func (osfs) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

func (osfs) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

func (osfs) ReadDir(dirname string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirname)
}

func (osfs) ReadFile(dirname string) ([]byte, error) {
	return ioutil.ReadFile(dirname)
}

func (osfs) Readlink(name string) (string, error) {
	return os.Readlink(name)
}

func (osfs) Remove(name string) error {
	return os.Remove(name)
}

func (osfs) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (osfs) Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

func (osfs) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return ioutil.WriteFile(filename, data, perm)
}
