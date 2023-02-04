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

// IsUser returns true if given Env e's Home directory is the user's
// home directory; false otherwise.
func (e *Env) IsUser() bool {
	if e == nil {
		return false
	}
	home := e.Home()
	user, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("gini: Env: no home directory: %v", err))
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
		home, err := os.UserHomeDir()
		if err != nil {
			panic(fmt.Sprintf("gini: Env: no home directory: %v", err))
		}
		e.Root = home
	}
	return e.Root
}

// WD returns the current working directory. WD panics if the working
// directory can't be determined.  Note use ChangeWD(newWD) to change
// the current working directory.
func (e *Env) WD() string {
	e.lock()
	defer e.mutex.Unlock()
	if e.wd == "" {
		wd, err := os.Getwd()
		if err != nil {
			panic(fmt.Sprintf("gini: Env: no working directory %v", err))
		}
		e.wd = wd
	}
	return e.wd
}

// Conf returns the user config directory.
func (e *Env) Conf() string {
	e.lock()
	defer e.mutex.Unlock()
	if e.conf == "" {
		conf, err := os.UserConfigDir()
		if err != nil {
			e.mutex.Unlock()
			conf = e.Home()
			e.lock()
		}
		e.conf = filepath.Join(conf, "gini")
	}
	return e.conf
}

// Logging returns the logging directory of given environment e.
func (e *Env) Logging() string {
	e.lock()
	defer e.mutex.Unlock()
	if e.logging == "" {
		e.mutex.Unlock()
		e.logging = filepath.Join(e.Conf(), "logs")
		e.lock()
	}
	return e.logging
}

// MkLogging creates the logging directory and panics if it can't be
// created.
func (e *Env) MkLogging() {
	if err := os.MkdirAll(e.Logging(), 0700); err != nil {
		panic(fmt.Sprintf(
			"gini: pkg: env: mk-logging: can't create dir: %v", err))
	}
}
