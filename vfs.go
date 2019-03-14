// Package vfs provides an abstraction of the os and ioutil packages that is
// easy to test.
package vfs

import (
	"os"
	"path/filepath"
	"sort"
	"time"
)

// A MkdirStater implements all the functionality needed by MkdirAll.
type MkdirStater interface {
	Mkdir(name string, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
}

// A LstatReadDirer implements all the functionality needed by Walk.
type LstatReadDirer interface {
	Lstat(name string) (os.FileInfo, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
}

// An FS is an abstraction over commonly-used functions in the os and ioutil
// packages.
type FS interface {
	Chmod(name string, mode os.FileMode) error
	Chown(name string, uid, git int) error
	Chtimes(name string, atime, mtime time.Time) error
	Create(name string) (*os.File, error)
	Lchown(name string, uid, git int) error
	Lstat(name string) (os.FileInfo, error)
	Mkdir(name string, perm os.FileMode) error
	Open(name string) (*os.File, error)
	OpenFile(name string, flag int, perm os.FileMode) (*os.File, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	Readlink(name string) (string, error)
	Remove(name string) error
	RemoveAll(name string) error
	Rename(oldpath, newpath string) error
	Stat(name string) (os.FileInfo, error)
	Symlink(oldname, newname string) error
	Truncate(name string, size int64) error
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

type infosByName []os.FileInfo

func (is infosByName) Len() int           { return len(is) }
func (is infosByName) Less(i, j int) bool { return is[i].Name() < is[j].Name() }
func (is infosByName) Swap(i, j int)      { is[i], is[j] = is[j], is[i] }

// MkdirAll is equivalent to os.MkdirAll but operates on fs.
func MkdirAll(fs MkdirStater, path string, perm os.FileMode) error {
	err := fs.Mkdir(path, perm)
	switch {
	case err == nil:
		// Mkdir was successful.
		return nil
	case os.IsExist(err):
		// path already exists, but we don't know whether it's a directory or
		// something else. We get this error if we try to create a subdirectory
		// of a non-directory, for example if the parent directory of path is a
		// file. There's a race condition here between the call to Mkdir and the
		// call to Stat but we can't avoid it because there's not enough
		// information in the returned error from Mkdir. We need to distinguish
		// between "path already exists and is already a directory" and "path
		// already exists and is not a directory". Between the call to Mkdir and
		// the call to Stat path might have changed.
		info, statErr := fs.Stat(path)
		if statErr != nil {
			return statErr
		}
		if !info.IsDir() {
			return err
		}
		return nil
	case os.IsNotExist(err):
		// Parent directory does not exist. Create the parent directory
		// recursively, then try again.
		parentDir := filepath.Dir(path)
		if parentDir == "/" || parentDir == "." {
			// We cannot create the root directory or the current directory, so
			// return the original error.
			return err
		}
		if err := MkdirAll(fs, parentDir, perm); err != nil {
			return nil
		}
		return fs.Mkdir(path, perm)
	default:
		// Some other error.
		return err
	}
}

// walk recursively walks fs from path.
func walk(fs LstatReadDirer, path string, walkFn filepath.WalkFunc, info os.FileInfo, err error) error {
	if err != nil {
		return walkFn(path, info, err)
	}
	err = walkFn(path, info, nil)
	if !info.IsDir() {
		return err
	}
	if err == filepath.SkipDir {
		return nil
	}
	infos, err := fs.ReadDir(path)
	if err != nil {
		return err
	}
	sort.Sort(infosByName(infos))
	for _, info := range infos {
		name := info.Name()
		if name == "." || name == ".." {
			continue
		}
		if err := walk(fs, filepath.Join(path, info.Name()), walkFn, info, nil); err != nil {
			return err
		}
	}
	return nil
}

// Walk is the equivalent of filepath.Walk but operates on fs. Entries are
// returned in lexicographical order.
func Walk(fs LstatReadDirer, path string, walkFn filepath.WalkFunc) error {
	info, err := fs.Lstat(path)
	return walk(fs, path, walkFn, info, err)
}
