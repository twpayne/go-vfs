// +build !windows

package vfst

import (
	"syscall"
	"testing"

	"github.com/twpayne/go-vfs"
)

// TestSysNlink returns a PathTest that verifies that the the path's
// Sys().(*syscall.Stat_t).Nlink is equal to wantNlink. If path's Sys() cannot
// be converted to a *syscall.Stat_t, it does nothing.
func TestSysNlink(wantNlink int) PathTest {
	return func(t *testing.T, fs vfs.FS, path string) {
		info, err := fs.Lstat(path)
		if err != nil {
			t.Errorf("fs.Lstat(%q) == %+v, %v, want !<nil>, <nil>", path, info, err)
			return
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok && int(stat.Nlink) != wantNlink {
			t.Errorf("fs.Lstat(%q).Sys().(*syscall.Stat_t).Nlink == %d, want %d", path, stat.Nlink, wantNlink)
		}
	}
}
