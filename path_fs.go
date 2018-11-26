package vfs

import (
	"os"
	"path"
	"path/filepath"
)

type PathFS struct {
	fs   FS
	path string
}

func NewPathFS(fs FS, path string) *PathFS {
	return &PathFS{
		path: path,
		fs:   fs,
	}
}

func (p *PathFS) Chmod(name string, mode os.FileMode) error {
	return p.fs.Chmod(p.Join(name), mode)
}

func (p *PathFS) Lstat(name string) (os.FileInfo, error) {
	return p.fs.Lstat(p.Join(name))
}

func (p *PathFS) Join(name string) string {
	return filepath.Join(p.path, name)
}

func (p *PathFS) Mkdir(name string, perm os.FileMode) error {
	return p.fs.Mkdir(p.Join(name), perm)
}

func (p *PathFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	return p.fs.ReadDir(p.Join(dirname))
}

func (p *PathFS) ReadFile(filename string) ([]byte, error) {
	return p.fs.ReadFile(p.Join(filename))
}

func (p *PathFS) Readlink(name string) (string, error) {
	return p.fs.Readlink(p.Join(name))
}

func (p *PathFS) Remove(name string) error {
	return p.fs.Remove(p.Join(name))
}

func (p *PathFS) Stat(name string) (os.FileInfo, error) {
	return p.fs.Stat(p.Join(name))
}

func (p *PathFS) Symlink(oldname, newname string) error {
	if path.IsAbs(oldname) {
		oldname = p.Join(oldname)
	}
	return p.fs.Symlink(oldname, p.Join(newname))
}

func (p *PathFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return p.fs.WriteFile(p.Join(filename), data, perm)
}
