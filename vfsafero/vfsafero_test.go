package vfsafero

import "github.com/spf13/afero"

var (
	_ afero.Fs = &AferoFS{}
)
