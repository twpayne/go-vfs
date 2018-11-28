package vfst

import (
	"testing"

	"github.com/twpayne/go-vfs"
)

// TestSysNlink returns a PathTest that verifies that the the path's
// Sys().(*syscall.Stat_t).Nlink is equal to wantNlink. If path's Sys() cannot
// be converted to a *syscall.Stat_t, it does nothing.
func TestSysNlink(wantNlink int) PathTest {
	return func(*testing.T, vfs.FS, string) {
	}
}
