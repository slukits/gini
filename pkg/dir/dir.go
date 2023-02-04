/*
Copyright 2022 - present Stephan Lukits. All rights reserved.
Use of this source code is governed by the GNU GPLv3 that can
be found in the LICENSE file.

This file is part of GINI.

GINI is free software: you can redistribute it and/or modify it
under the terms of the GNU General Public License as published
by the Free Software Foundation, either version 3 of the License,
or (at your option) any later version.

GINI is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with GINI. If not, see <https://www.gnu.org/licenses/#GPL>.
*/

/*
Package dir provides features around the file system which exceed
provided features by the os-package.
*/
package dir

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/slukits/gini/pkg/env"
)

// Dir provides file-system operation compounded of os-file-system
// operations.  The zero-type is ready to use.
type Dir struct {

	// Env instance used for concurrency save logging; Env defaults to
	// &env.Env{}
	Env *env.Env

	// Path of an Dir instance defaulting to Env's working directory.
	Path string
}

func (d *Dir) env() *env.Env {
	if d.Env == nil {
		d.Env = &env.Env{}
	}
	return d.Env
}

func (d *Dir) path() string {
	if d.Path == "" {
		d.Path = d.env().WD()
	}
	return d.Path
}

// WalkUp returns a closure which goes with each call one directory up
// starting at given directory d.
func (d *Dir) WalkUp() (up func() (*Dir, bool)) {
	p := d.path()
	return func() (*Dir, bool) {
		if p == filepath.Dir(p) {
			return nil, false
		}
		p = filepath.Dir(p)
		return &Dir{Env: d.Env, Path: p}, true
	}
}

// Repo walks given Dir d's path up until a directory is found
// containing a .git directory and returns it along with a true value.
// If no such directory is found nil and false is returned.
func (d *Dir) Repo() (*Dir, bool) {
	up, _d, next := d.WalkUp(), d, true
	for next && !_d.Contains(".git") {
		_d, next = up()
	}
	if _d == nil {
		return nil, false
	}
	return _d, true
}

// Caller returns the directory of the file containing Caller calling
// function.  Caller panics if caller can't be determined.
func (d *Dir) Caller() *Dir {
	_, f, _, ok := runtime.Caller(1)
	if !ok {
		panic("gini: pkg: dir: Caller: cant determine caller stack")
	}
	return &Dir{
		Env:  d.Env,
		Path: filepath.Dir(f),
	}
}

// Contains returns true if given directory d contains a directory with
// given name dirName;  otherwise false is returned.
func (d *Dir) Contains(dirName string) bool {
	ee, err := os.ReadDir(d.path())
	if err != nil {
		return false
	}
	for _, e := range ee {
		if !e.IsDir() {
			continue
		}
		if e.Name() != dirName {
			continue
		}
		return true
	}
	return false
}

// FileContains returns ture if given file fl in given directory d
// contains given bytes bb; otherwise false is returned.
func (d *Dir) FileContains(fl string, bb []byte) bool {
	fbb, err := ioutil.ReadFile(filepath.Join(d.path(), fl))
	if err != nil {
		return false
	}
	return bytes.Contains(fbb, bb)
}

// String returns given directories d path.
func (d *Dir) String() string {
	return d.path()
}
