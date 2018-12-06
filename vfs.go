package vfs

import (
	"fmt"
	"os"
	"path/filepath"
)

// An FS is an abstraction over commonly-used functions in the os and ioutil
// packages.
type FS interface {
	Chmod(name string, mode os.FileMode) error
	Chown(name string, uid, git int) error
	Lstat(name string) (os.FileInfo, error)
	Mkdir(name string, perm os.FileMode) error
	Open(name string) (*os.File, error)
	ReadDir(dirname string) ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	Readlink(name string) (string, error)
	Remove(name string) error
	RemoveAll(name string) error
	Rename(oldpath, newpath string) error
	Stat(name string) (os.FileInfo, error)
	Symlink(oldname, newname string) error
	WriteFile(filename string, data []byte, perm os.FileMode) error
}

// MkdirAll is equivalent to os.MkdirAll but operates on fs.
func MkdirAll(fs FS, path string, perm os.FileMode) error {
	if parentDir := filepath.Dir(path); parentDir != "." {
		info, err := fs.Stat(parentDir)
		if err != nil && os.IsNotExist(err) {
			if err := MkdirAll(fs, parentDir, perm); err != nil {
				return err
			}
		} else if err != nil {
			return err
		} else if err == nil && !info.Mode().IsDir() {
			return fmt.Errorf("%s: not a directory", parentDir)
		}
	}
	info, err := fs.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil && info.Mode().IsDir() {
		return nil
	}
	return fs.Mkdir(path, perm)
}

// walk recursively walks fs from path.
func walk(fs FS, path string, walkFn filepath.WalkFunc, info os.FileInfo, err error) error {
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

// Walk is the equivalent of filepath.Walk but operates on fs.
func Walk(fs FS, path string, walkFn filepath.WalkFunc) error {
	info, err := fs.Lstat(path)
	return walk(fs, path, walkFn, info, err)
}
