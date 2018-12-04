package vfs

import (
	"os"
	"syscall"
)

// A ReadOnlyFS operates on an existing FS, but any methods that
// modify the FS return an error.
type ReadOnlyFS struct {
	fs FS
}

// NewReadOnlyFS returns a new *ReadOnlyFS operating on fs.
func NewReadOnlyFS(fs FS) *ReadOnlyFS {
	return &ReadOnlyFS{
		fs: fs,
	}
}

// Chmod implements os.Chmod.
func (r *ReadOnlyFS) Chmod(name string, mode os.FileMode) error {
	return &os.PathError{
		Op:   "Chmod",
		Path: name,
		Err:  syscall.EPERM,
	}
}

// Lstat implements os.Lstat.
func (r *ReadOnlyFS) Lstat(name string) (os.FileInfo, error) {
	return r.fs.Lstat(name)
}

// Mkdir implements os.Mkdir.
func (r *ReadOnlyFS) Mkdir(name string, perm os.FileMode) error {
	return &os.PathError{
		Op:   "Mkdir",
		Path: name,
		Err:  syscall.EPERM,
	}
}

// Open implements os.Open.
func (r *ReadOnlyFS) Open(name string) (*os.File, error) {
	return r.fs.Open(name)
}

// ReadDir implenents ioutil.ReadDir.
func (r *ReadOnlyFS) ReadDir(dirname string) ([]os.FileInfo, error) {
	return r.fs.ReadDir(dirname)
}

// ReadFile implements ioutil.ReadFile.
func (r *ReadOnlyFS) ReadFile(filename string) ([]byte, error) {
	return r.fs.ReadFile(filename)
}

// Readlink implments os.Readlink.
func (r *ReadOnlyFS) Readlink(name string) (string, error) {
	return r.fs.Readlink(name)
}

// Remove implements os.Remove.
func (r *ReadOnlyFS) Remove(name string) error {
	return &os.PathError{
		Op:   "Remove",
		Path: name,
		Err:  syscall.EPERM,
	}
}

// RemoveAll implements os.RemoveAll.
func (r *ReadOnlyFS) RemoveAll(name string) error {
	return &os.PathError{
		Op:   "RemoveAll",
		Path: name,
		Err:  syscall.EPERM,
	}
}

// Rename implements os.Rename.
func (r *ReadOnlyFS) Rename(oldpath, newpath string) error {
	return &os.PathError{
		Op:   "Rename",
		Path: oldpath,
		Err:  syscall.EPERM,
	}
}

// Stat implements os.Stat.
func (r *ReadOnlyFS) Stat(name string) (os.FileInfo, error) {
	return r.fs.Stat(name)
}

// Symlink implements os.Symlink.
func (r *ReadOnlyFS) Symlink(oldname, newname string) error {
	return &os.PathError{
		Op:   "Symlink",
		Path: newname,
		Err:  syscall.EPERM,
	}
}

// WriteFile implements ioutil.WriteFile.
func (r *ReadOnlyFS) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return &os.PathError{
		Op:   "WriteFile",
		Path: filename,
		Err:  syscall.EPERM,
	}
}
