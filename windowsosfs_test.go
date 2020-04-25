// +build windows

package vfs

import "github.com/bmatcuk/doublestar"

var (
	_ FS            = WindowsOSFS{}
	_ doublestar.OS = WindowsOSFS{}
)
