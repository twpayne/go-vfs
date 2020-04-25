package vfs

import "github.com/bmatcuk/doublestar"

var (
	_ FS            = &ReadOnlyFS{}
	_ doublestar.OS = &ReadOnlyFS{}
)
