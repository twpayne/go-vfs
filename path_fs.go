package vfs

// FIXME guard against operations outside path

import (
	"os"
	"path"
	"path/filepath"
)

// A PathFS operates on an existing FS, but prefixes all names with a path.
type PathFS struct {
	fs   FS
	path string
}

// NewPathFS returns a new *PathFS operating on fs and prefixing all names with
// path.
func NewPathFS(fs FS, path string) *PathFS {
	return &PathFS{
		path: path,
		fs:   fs,
	}
}

// Chmod implements os.Chmod.
func (p *PathFS) Chmod(name string, mode os.FileMode) error {
	return p.fs.Chmod(p.Join(name), mode)
}

// Lstat implements os.Lstat.
func (p *PathFS) Lstat(name string) (os.FileInfo, error) {
	return p.fs.Lstat(p.Join(name))
}

// Join returns p's path joined with name.
func (p *PathFS) Join(name string) string {
	return filepath.Join(p.path, name)
}

// Mkdir implements os.Mkdir.
func (p *PathFS) Mkdir(name string, perm os.FileMode) error {
	return p.fs.Mkdir(p.Join(name), perm)
}

// Open implements os.Open.
func (p *PathFS) Open(name string) (*os.File, error) {
	return p.fs.Open(p.Join(name))
}

// ReadDir implenents ioutil.ReadDir.
func (p *PathFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	return p.fs.ReadDir(p.Join(dirname))
}

// ReadFile implements ioutil.ReadFile.
func (p *PathFS) ReadFile(filename string) ([]byte, error) {
	return p.fs.ReadFile(p.Join(filename))
}

// Readlink implments os.Readlink.
func (p *PathFS) Readlink(name string) (string, error) {
	return p.fs.Readlink(p.Join(name))
}

// Remove implements os.Remove.
func (p *PathFS) Remove(name string) error {
	return p.fs.Remove(p.Join(name))
}

// RemoveAll implements os.RemoveAll.
func (p *PathFS) RemoveAll(name string) error {
	return p.fs.RemoveAll(p.Join(name))
}

// Rename implements os.Rename.
func (p *PathFS) Rename(oldpath, newpath string) error {
	return p.fs.Rename(p.Join(oldpath), p.Join(newpath))
}

// Stat implements os.Stat.
func (p *PathFS) Stat(name string) (os.FileInfo, error) {
	return p.fs.Stat(p.Join(name))
}

// Symlink implements os.Symlink.
func (p *PathFS) Symlink(oldname, newname string) error {
	if path.IsAbs(oldname) {
		oldname = p.Join(oldname)
	}
	return p.fs.Symlink(oldname, p.Join(newname))
}

// WriteFile implements ioutil.WriteFile.
func (p *PathFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return p.fs.WriteFile(p.Join(filename), data, perm)
}
