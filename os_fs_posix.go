// +build !windows

package vfs

import "os"

// LstatIfPossible calls Lstat if it is available, Stat otherwise.
func (osfs) LstatIfPossible(name string) (os.FileInfo, bool, error) {
	info, err := os.Lstat(name)
	return info, true, err
}
