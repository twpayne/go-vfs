package vfs

import (
	"io/ioutil"
	"os"
	"time"
)

type osfs struct{}

// OSFS is the FS that calls os and ioutil functions directly.
var OSFS = &osfs{}

// Chmod implements os.Chmod.
func (osfs) Chmod(name string, mode os.FileMode) error {
	return os.Chmod(name, mode)
}

// Chown implements os.Chown.
func (osfs) Chown(name string, uid, gid int) error {
	return os.Chown(name, uid, gid)
}

// Chtimes implements os.Chtimes.
func (osfs) Chtimes(name string, atime, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

// Lstat implements os.Lstat.
func (osfs) Lstat(name string) (os.FileInfo, error) {
	return os.Lstat(name)
}

// Mkdir implements os.Mkdir.
func (osfs) Mkdir(name string, perm os.FileMode) error {
	return os.Mkdir(name, perm)
}

// Open implements os.Open.
func (osfs) Open(name string) (*os.File, error) {
	return os.Open(name)
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

// RemoveAll implements os.RemoveAll.
func (osfs) RemoveAll(name string) error {
	return os.RemoveAll(name)
}

// Rename implements os.Rename.
func (osfs) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
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
