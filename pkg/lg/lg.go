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
Package lg is a thin wrapper around the go log-package with some
convenience function especially for testing: A logger for a temp
environment is initially setup which logs to string-buffers and
*String-functions can be leveraged to see the current content of a
logger.
*/
package lg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/slukits/gini/pkg/env"
)

const (
	ErrFlags   int = log.Lmsgprefix | log.Ldate | log.Lshortfile
	Flags      int = log.Lmsgprefix | log.Ldate | log.Ltime | log.LUTC
	FileSuffix     = "log"
)

var initMutex sync.Mutex

// Logger is a convenience wrapper around the standard-library
// log-package.  It provides named loggers which are created as needed.
// The zero-Logger is ready to use and logs to os.Stderr.  EnvLogger
// also allows to retrieve the content of a logger.
type Logger struct {

	// Env is the environment this logger logs to.  Is Env is nil the
	// default logger of go's log package is used.  If Env is a
	// temporary environment WriteTempLogs determines if in-memory
	// logger are used for logging; otherwise log-files to Env's
	// Logging-directory are written.
	Env *env.Env

	// WriteTempLogs switch indicates if logs in a temp-Env(ironment)
	// should be written to the disk.
	WriteTempLogs bool

	mutex *sync.Mutex

	ll map[string]*log.Logger
	// FatalHandler is an optional callback for situations where
	// log.Fatal would be called.  It defaults to a call of log.Fatal.
	FatalHandler func(v ...interface{})

	// ErrHandler is an optional callback for log-write error-handling.
	// It defaults to a call of log.Fatal
	ErrHandler func(error)
}

func (l *Logger) initMutex() {
	initMutex.Lock()
	defer initMutex.Unlock()
	if l.mutex == nil {
		l.mutex = &sync.Mutex{}
	}
}

func (l *Logger) lock() {
	if l.mutex == nil {
		l.initMutex()
	}
	l.mutex.Lock()
}

// To logs given message msg to logger with given name.
func (l *Logger) To(name, msg string) {
	l.lock()
	defer l.mutex.Unlock()
	if l.Env == nil {
		log.Print(msg)
		return
	}
	l.logger(name).Print(msg)
}

func (l *Logger) logger(name string) *log.Logger {
	if l, ok := l.ll[name]; ok {
		return l
	}
	if l.ll == nil {
		l.ll = map[string]*log.Logger{}
	}
	lgg := l.createLogger(name)
	l.ll[name] = lgg
	return lgg
}

func (l *Logger) createLogger(name string) *log.Logger {
	ff, p := Flags, name
	if strings.Contains(strings.ToLower(name), "err") {
		ff = ErrFlags
	}
	if l.Env.IsTemp() && !l.WriteTempLogs {
		return log.New(&strings.Builder{}, p, ff)
	}
	dir := l.Env.Logging()
	fName := filepath.Join(dir, fmt.Sprintf("%s.%s", name, FileSuffix))
	file, err := os.OpenFile(
		fName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		l.Env.MkLogging()
		file, err = os.OpenFile(
			fName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(fmt.Sprintf(
				"gini: pkg: lg: can't write log-file: %v", err))
		}
	}
	return log.New(file, p, ff)
}

// Tof logs to logger with given name given format string with given
// args.
func (l *Logger) Tof(name, format string, args ...interface{}) {
	l.lock()
	defer l.mutex.Unlock()
	if l.Env == nil {
		log.Printf(format, args...)
		return
	}
	l.logger(name).Printf(format, args...)
}

// String returns the content of the logger with given name.
func (l *Logger) String(name string) string {
	l.lock()
	defer l.mutex.Unlock()
	if name == "" {
		return ""
	}
	lg, ok := l.ll[name]
	if !ok {
		return ""
	}
	buffer, ok := lg.Writer().(*strings.Builder)
	if ok {
		return buffer.String()
	}
	f, ok := lg.Writer().(*os.File)
	if !ok {
		return ""
	}
	fName, flags := f.Name(), lg.Flags()
	if err := f.Close(); err != nil {
		l.handleError(err)
	}
	bb, err := ioutil.ReadFile(fName)
	if err != nil {
		l.handleError(err)
	}
	f, err = os.OpenFile(
		fName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		l.handleError(err)
	}
	l.ll[name] = log.New(f, name, flags)
	return string(bb)
}

func (l *Logger) handleError(err error) {
	if l.ErrHandler == nil {
		l.ErrHandler = func(err error) {
			log.Fatal(err)
		}
	}
	l.ErrHandler(err)
}
