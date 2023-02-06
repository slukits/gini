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
Package env provides access to a GINI instance's environment and is used
as identifier for environment specific type-instances of various types
making concurrency save parallel testing easier.
*/
package env

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var initMutex = sync.Mutex{}

// Env provides concurrency save access to the runtime environment of a
// GINI instance.  The zero-value is ready to use.
type Env struct {

	// Root is an environments root directory which defaults to the
	// users home directory.  Is root set to a different directory but
	// the users home directory all other directories like Caching,
	// Config or Logging are created inside Root instead of the user's
	// caching directory etc.  E.g. for testing it is handy to set Root
	// to a temporary directory.
	Root string

	// FatalHandler is used if an Env-instance is unable to operate on
	// the system in the intended way which by an unset FatalHandler
	// leads to a panic.
	FatalHandler Fataler

	// OS provides the os.* functions which may fail and are needed by
	// an Env-instance.  OS defaults to an instance of an internal type
	// mapping OSer methods to the corresponding os.* calls.  This
	// indirection allows to easily mock up error situations.
	OS OSer

	wd      string
	conf    string
	logging string
	mutex   *sync.Mutex
}

func (e *Env) initMutex() {
	initMutex.Lock()
	defer initMutex.Unlock()
	if e.mutex != nil {
		return
	}
	e.mutex = &sync.Mutex{}
}

func (e *Env) lock() {
	if e.mutex == nil {
		e.initMutex()
	}
	e.mutex.Lock()
}

// os returns given Env e's os'er initializing it to its default if
// unset.  NOTE os expects e to be locked, i.e. there may be race
// conditions if not.
func (e *Env) os() OSer {
	if e.OS == nil {
		e.OS = &oser{}
	}
	return e.OS
}

// IsUser returns true if given Env e's Home directory is the user's
// home directory; false otherwise.
func (e *Env) IsUser() bool {
	if e == nil {
		return false
	}
	home := e.Home()
	user, err := e.os().UserHomeDir()
	if err != nil {
		e.fatal("gini: Env: no home directory: %w", err)
	}
	return home == user
}

// IsTemp return true if given Env e's Home is prefixed by the systems
// temp directory; false otherwise.
func (e *Env) IsTemp() bool {
	if e == nil {
		return false
	}
	return strings.HasPrefix(e.Home(), os.TempDir())
}

// Home returns the content of Root respectively sets Root to its
// default which is the user's home directory.
func (e *Env) Home() string {
	e.lock()
	defer e.mutex.Unlock()
	if e.Root == "" {
		home, err := e.os().UserHomeDir()
		if err != nil {
			e.fatal("gini: Env: no home directory: %w", err)
		}
		e.Root = home
	}
	return e.Root
}

func (e *Env) fatal(msg string, err error) {
	if e.FatalHandler == nil {
		panic(fmt.Errorf(msg, err))
	}
	e.FatalHandler.Fatal(fmt.Errorf(msg, err))
	panic("env: handling fatal: expected execution to stop")
}

// WD returns the current working directory. WD panics if the working
// directory can't be determined.  Note use ChangeWD(newWD) to change
// the current working directory.
func (e *Env) WD() string {
	e.lock()
	defer e.mutex.Unlock()
	if e.wd == "" {
		wd, err := e.os().Getwd()
		if err != nil {
			e.fatal("gini: Env: no working directory %w", err)
		}
		e.wd = wd
	}
	return e.wd
}

// Conf returns the user config directory.
func (e *Env) Conf() string {
	e.lock()
	if e.conf == "" {
		conf, err := e.os().UserConfigDir()
		if err != nil {
			e.mutex.Unlock()
			conf = e.Home()
			e.lock()
		}
		e.conf = filepath.Join(conf, "gini/config")
	}
	e.mutex.Unlock()
	return e.conf
}

// Logging returns the logging directory of given environment e.
func (e *Env) Logging() string {
	e.lock()
	defer e.mutex.Unlock()
	if e.logging == "" {
		e.mutex.Unlock()
		logging := filepath.Join(e.Conf(), "logs")
		e.lock()
		e.logging = logging
	}
	return e.logging
}

// MkLogging creates the logging directory and errors if MkdirAll errors.
func (e *Env) MkLogging() error {
	return e.os().MkdirAll(e.Logging(), 0700)
}

// Fataler defines the interface for dealing with situation when an
// environment Env-instance can not operate in the intended way on the
// system.
type Fataler interface {
	Fatal(...interface{})
}

// Oser defines the interface of os-operations which may fail and are
// needed by an environment Env-instance to provide its features.
type OSer interface {
	UserHomeDir() (string, error)
	Getwd() (string, error)
	UserConfigDir() (string, error)
	MkdirAll(path string, _ fs.FileMode) error
}

type oser struct{}

func (oser) UserHomeDir() (string, error) { return os.UserHomeDir() }
func (oser) Getwd() (string, error)       { return os.Getwd() }
func (oser) UserConfigDir() (string, error) {
	return os.UserConfigDir()
}
func (oser) MkdirAll(path string, fm fs.FileMode) error {
	return os.MkdirAll(path, fm)
}
