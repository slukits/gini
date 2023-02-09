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

	// FatalHandler is used if an Env-instance is unable to operate on
	// the system in the intended way which by an unset FatalHandler
	// leads to a panic.
	FatalHandler Fataler

	// Lib provides the os.* functions which may fail and are needed by
	// an Env-instance for mock ups.
	Lib Lib

	home    string
	wd      string
	conf    string
	logging string
	mutex   *sync.Mutex
	initLib bool
}

func (e *Env) lock() {
	initMutex.Lock()
	if e.mutex == nil {
		e.mutex = &sync.Mutex{}
	}
	initMutex.Unlock()
	e.mutex.Lock()
}

// lib returns given Env e's lib'er initializing it to its default if
// unset.  NOTE lib expects e to be locked, i.e. there may be race
// conditions if not.
func (e *Env) lib() Lib {
	if !e.initLib {
		e.initLib = true
		if e.Lib.Getwd == nil {
			e.Lib.Getwd = os.Getwd
		}
		if e.Lib.MkdirAll == nil {
			e.Lib.MkdirAll = os.MkdirAll
		}
		if e.Lib.UserConfigDir == nil {
			e.Lib.UserConfigDir = os.UserConfigDir
		}
		if e.Lib.UserHomeDir == nil {
			e.Lib.UserHomeDir = os.UserHomeDir
		}
		if e.Lib.Chdir == nil {
			e.Lib.Chdir = os.Chdir
		}
	}
	return e.Lib
}

// IsUser returns true if given Env e's Home directory is the user's
// home directory; false otherwise.  Note in the later case all
// Env-paths like for configuration or logging are created inside e's
// home directory.
func (e *Env) IsUser() bool {
	if e == nil {
		return false
	}
	home := e.Home()
	user, err := e.lib().UserHomeDir()
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

// SetHome allows to change an environments home directory having the
// consequence that all of the environments paths are created inside
// that home directory.  I.e. if provided a testing temp-directory we
// have a unique concurrency save environment, we are sure that the user
// space is not filled with testing data and that the testing data will
// be cleaned up after the test.
func (e *Env) SetHome(path string) *Env {
	if path == "" {
		return e
	}
	e.lock()
	defer e.mutex.Unlock()
	e.home = path
	return e
}

// Home returns the environment's home directory which defaults to the
// user's home directory but may be set differently especially for
// testing (see [Env.SetHome]).
func (e *Env) Home() string {
	e.lock()
	defer e.mutex.Unlock()
	if e.home == "" {
		home, err := e.lib().UserHomeDir()
		if err != nil {
			e.fatal("gini: Env: no home directory: %w", err)
		}
		e.home = home
	}
	return e.home
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
		wd, err := e.lib().Getwd()
		if err != nil {
			e.fatal("gini: Env: no working directory %w", err)
		}
		e.wd = wd
	}
	return e.wd
}

// ChWD changes the given environment e's and the executed binary's
// working directory.
func (e *Env) ChWD(path string) error {
	e.lock()
	defer e.mutex.Unlock()
	if err := e.lib().Chdir(path); err != nil {
		return err
	}
	e.wd = path
	return nil
}

// Conf returns the user config directory.
func (e *Env) Conf() string {
	e.lock()
	if e.conf == "" {
		e.mutex.Unlock()
		conf, err := e.lib().UserConfigDir()
		if err != nil || !e.IsUser() {
			conf = e.Home()
		}
		e.lock()
		e.conf = filepath.Join(conf, "gini/config")
	}
	defer e.mutex.Unlock()
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
	return e.lib().MkdirAll(e.Logging(), 0700)
}

// Fataler defines the interface for dealing with situation when an
// environment Env-instance can not operate in the intended way on the
// system.
type Fataler interface {
	Fatal(...interface{})
}

type Lib struct {
	UserHomeDir   func() (string, error)
	Getwd         func() (string, error)
	UserConfigDir func() (string, error)
	MkdirAll      func(path string, fm fs.FileMode) error
	Chdir         func(path string) error
}
