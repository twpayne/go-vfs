package vfs

import (
	"fmt"
	"os"
	"path/filepath"
)

type FS interface {
	Chmod(string, os.FileMode) error
	Lstat(string) (os.FileInfo, error)
	Mkdir(string, os.FileMode) error
	ReadDir(string) ([]os.FileInfo, error)
	ReadFile(string) ([]byte, error)
	Readlink(string) (string, error)
	Remove(string) error
	Stat(string) (os.FileInfo, error)
	Symlink(string, string) error
	WriteFile(string, []byte, os.FileMode) error
}

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

func removeAll(fs FS, path string, info os.FileInfo) error {
	if info.Mode().IsDir() {
		infos, err := fs.ReadDir(path)
		if err != nil {
			return err
		}
		for _, info := range infos {
			if err := removeAll(fs, filepath.Join(path, info.Name()), info); err != nil {
				return err
			}
		}
		return nil
	}
	return fs.Remove(path)
}

func RemoveAll(fs FS, path string) error {
	info, err := fs.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	return removeAll(fs, path, info)
}

func Walk(fs FS, path string, walkFn filepath.WalkFunc) error {
	return nil
}
