# `go-vfs`

[![Build Status](https://travis-ci.org/twpayne/go-vfs.svg?branch=master)](https://travis-ci.org/twpayne/go-vfs)
[![GoDoc](https://godoc.org/github.com/twpayne/go-vfs?status.svg)](https://godoc.org/github.com/twpayne/go-vfs)
[![Report Card](https://goreportcard.com/badge/github.com/twpayne/go-vfs)](https://goreportcard.com/report/github.com/twpayne/go-vfs)

Package `go-vfs` provides an abstraction of the `os` and `ioutil` packages that is easy to test.

## Key features

 * File system abstraction layer for commonly-used `os` and `ioutil` functions
   from the standard library.

 * Powerful testing framework, `vfstest`. For a quick tour of `vfstest`'s
   features, see [the examples in the
documentation](https://godoc.org/github.com/twpayne/go-vfs/vfstest#pkg-examples).

## Motivation

`go-vfs` was inspired by
[`github.com/spf13/afero`](https://github.com/spf13/afero) and
[`github.com/twpayne/aferot`](https://github.com/twpayne/aferot). So, why not
use these?

 * `afero` has several critical bugs in its in-memory mock filesystem
   implementation `MemMapFs`, to the point that it is unusable for non-trivial
test cases.  `vfs` does not attempt to implent an in-memory mock filesystem,
and instead only provides a thin layer around the standard libary's `os` and
`ioutil` packages, and as such should have fewer bugs.

 * `afero` does not support creating or reading symbolic links, and its
   `LstatIfPossible` interface is clumsy to use. `vfs` provides out-of-the-box
support for symbolic links.

 * `afero` has been effectively abandoned by its author, and a "friendly fork"
   ([`github.com/absfs/afero`](https://github.com/absfs/afero)) has not seen
much activity. `vfs`, by providing much less functionality than `afero`, should
be smaller and easier to maintain.

## License

The MIT License (MIT)

Copyright (c) 2018 Tom Payne

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
