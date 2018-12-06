package vfs

// FIXME guard against operations outside path

import (
	"os"
	"path"
	"path/filepath"
	"time"
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

// Chown implements os.Chown.
func (p *PathFS) Chown(name string, uid, gid int) error {
	return p.fs.Chown(p.Join(name), uid, gid)
}

// Chtimes implements os.Chtimes.
func (p *PathFS) Chtimes(name string, atime, mtime time.Time) error {
	return p.fs.Chtimes(p.Join(name), atime, mtime)
}

// Create implements os.Create.
func (p *PathFS) Create(name string) (*os.File, error) {
	return p.fs.Create(p.Join(name))
}

// Lchown implements os.Lchown.
func (p *PathFS) Lchown(name string, uid, gid int) error {
	return p.fs.Lchown(p.Join(name), uid, gid)
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

// OpenFile implements os.OpenFile.
func (p *PathFS) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return p.fs.OpenFile(p.Join(name), flag, perm)
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

// Truncate implements os.Truncate.
func (p *PathFS) Truncate(name string, size int64) error {
	return p.fs.Truncate(p.Join(name), size)
}

// WriteFile implements ioutil.WriteFile.
func (p *PathFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return p.fs.WriteFile(p.Join(filename), data, perm)
}
