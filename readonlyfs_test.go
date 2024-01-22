package vfs_test

import "github.com/twpayne/go-vfs/v5"

var _ vfs.FS = &vfs.ReadOnlyFS{}
