package vfsafero

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/twpayne/go-vfs/vfst"
)

var (
	_ afero.Fs = &AferoFS{}
)

func TestAferoFS(t *testing.T) {
	fs, cleanup, err := vfst.NewTestFS(map[string]interface{}{
		"/home/user/.bashrc": "# contents of .bashrc\n",
	})
	defer cleanup()
	if err != nil {
		t.Fatal(err)
	}

	aferoFS := NewAferoFS(fs)
	afero.WriteFile(aferoFS, "/home/user/foo", []byte("bar"), 0666)

	vfst.RunTests(t, fs, "",
		vfst.TestPath("/home/user/foo",
			vfst.TestContentsString("bar"),
		),
	)
}
