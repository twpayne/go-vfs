package vfs

import "github.com/bmatcuk/doublestar"

var (
	_ FS            = &PathFS{}
	_ doublestar.OS = &PathFS{}
)
