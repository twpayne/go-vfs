package vfs

import "github.com/bmatcuk/doublestar"

var (
	_ FS            = OSFS
	_ doublestar.OS = OSFS
)
