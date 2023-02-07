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
convenience features especially for testing:
  - A Logger for a temp environment logs to in-memory string-buffers
  - A Logger's String-functions can be leveraged to retrieve the logging
    messages associated with a logging name
  - A Logger's Fatal- and Error-Handlers allow to replace log.Fatal
    calls which would end the program execution with a function that for
    example ends on the execution of a specific test evaluating the
    fatal situation.
*/
package lg

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/slukits/gini/pkg/env"
)

const (

	// ErrFlags are the logging flags for error messages which are used
	// automatically if a Logger-name contains "err"
	ErrFlags int = log.Lmsgprefix | log.Ldate | log.Lshortfile

	// Flags are the default logging flags logging messages.
	Flags int = log.Lmsgprefix | log.Ldate | log.Ltime | log.LUTC

	// FileSuffix of text files created by a Logger instance.
	FileSuffix = "log"

	// ERR may be used as the identifier for an error logger
	ERR = "err"

	// INF may be used as the identifier for an info logger
	INF = "inf"

	// SEC may be used as the identifier for an security logger
	SEC = "sec"
)

var initMutex sync.Mutex

// Logger is a convenience wrapper around the standard-library
// log-package.  It provides named logs which are created as needed.
// The zero-Logger is ready to use and logs to the log-package's default
// logger.  A Logger also allows to retrieve the content of any named
// log by using the String method.  Finally the Lib-property allows for
// easy mockups of std-lib-calls.  NOTE Lib is initialized to its
// defaults only once.  While it is possible to mock-up Fatal Logger
// expects that after the call of Fatal the execution ends, i.e. it will
// panic in the line after calling Fatal.
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

	// Lib allows to mock up used std lib functions which may fail or
	// terminate execution.
	Lib Lib

	mutex   *sync.Mutex
	initLib bool

	ll map[string]*log.Logger
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

// lib returns library function needed by given Logger l set to their
// defaults if unset.  Note lib expects l.mutex to be locked.
func (l *Logger) lib() Lib {
	if !l.initLib {
		l.initLib = true
		if l.Lib.OpenFile == nil {
			l.Lib.OpenFile = os.OpenFile
		}
		if l.Lib.ReadFile == nil {
			l.Lib.ReadFile = ioutil.ReadFile
		}
		if l.Lib.Fatal == nil {
			l.Lib.Fatal = log.Fatal
		}
		if l.Lib.Println == nil {
			l.Lib.Println = log.Println
		}
	}
	return l.Lib
}

// LibDefaults returns given Logger l's Lib property having all unset
// function properties set to their defaults.
func (l *Logger) LibDefaults() Lib { return l.lib() }

// Writer returns the writer for the log associated with given name
// which is either a file or a string builder.
func (l *Logger) Writer(name string) io.Writer {
	l.lock()
	defer l.mutex.Unlock()
	lg := l.ll[name]
	if lg == nil {
		return nil
	}
	return lg.Writer()
}

func (l *Logger) Flags(name string) int {
	l.lock()
	defer l.mutex.Unlock()
	lg := l.ll[name]
	if lg == nil {
		return 0
	}
	return lg.Flags()
}

// InMemory returns true if log with given name is in memory.
func (l *Logger) InMemory(name string) bool {
	if l == nil || l.ll == nil {
		return false
	}
	lgg, ok := l.ll[name]
	if !ok {
		return ok
	}
	_, ok = lgg.Writer().(*strings.Builder)
	return ok
}

// To logs given message msg to given Logger l to the log with given
// name.  Note if l has no environment Env specified defining the
// Logger's environment DefaultHandler is used to do the logging which
// defaults to log.Println.
func (l *Logger) To(name, msg string) {
	l.lock()
	defer l.mutex.Unlock()
	if l.Env == nil {
		l.handleDefault(name, msg)
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
	ff, p := Flags, name+": "
	if strings.Contains(strings.ToLower(name), "err") {
		ff = ErrFlags
	}
	if l.Env.IsTemp() && !l.WriteTempLogs {
		return log.New(&strings.Builder{}, p, ff)
	}
	dir := l.Env.Logging()
	fName := filepath.Join(dir, fmt.Sprintf("%s.%s", name, FileSuffix))
	file, err := l.lib().OpenFile(
		fName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		if err := l.Env.MkLogging(); err != nil {
			l.handleError("gini: pkg: lg: create dir: %v", err)
		}
		file, err = l.lib().OpenFile(
			fName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			l.handleError("gini: pkg: lg: open log-file: %v", err)
		}
	}
	return log.New(file, p, ff)
}

// Tof logs given values vv formatted with given format specification
// format with given Logger l to the log with given name.  Note if l has
// no environment Env specified defining the Logger's environment
// DefaultHandler is used to do the logging which defaults to
// log.Println(fmt.Sprintf(format, vv...)).
func (l *Logger) Tof(name, format string, vv ...interface{}) {
	l.lock()
	defer l.mutex.Unlock()
	if l.Env == nil {
		l.handleDefault(name, fmt.Sprintf(format, vv...))
		return
	}
	l.logger(name).Printf(format, vv...)
}

// String returns the content of the logger with given name.
func (l *Logger) String(name string) string {
	if l == nil {
		return ""
	}
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
	f := lg.Writer().(*os.File)
	fName, flags := f.Name(), lg.Flags()
	if err := f.Close(); err != nil {
		l.handleError("gini: pkg: lg: String: close log-file: %v", err)
	}
	bb, err := l.lib().ReadFile(fName)
	if err != nil {
		l.handleError("gini: pkg: lg: String: read log-file: %v", err)
	}
	f, err = l.lib().OpenFile(
		fName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		l.handleError("gini: pkg: lg: String: reopen log-file: %v", err)
	}
	l.ll[name] = log.New(f, name, flags)
	return string(bb)
}

func (l *Logger) handleError(format string, err error) {
	l.lib().Fatal(fmt.Sprintf(format, err))
	panic("gini: pkg: lg: handling error: expected " +
		"execution to stop")
}

func (l *Logger) handleDefault(name, msg string) {
	msg = fmt.Sprintf("%s: %s", name, msg)
	l.lib().Println(msg)
}

// Fatal calls given Logger l's FatalHandler which defaults to log.Fatal
// to given values vv
func (l *Logger) Fatal(vv ...interface{}) {
	l.lock()
	defer l.mutex.Unlock()
	l.fatal(vv...)
}

// Fatalf logs given values vv formatted according to given format
// specifier format with fmt.Sprintf to given Logger l's FatalHandler.
func (l *Logger) Fatalf(format string, vv ...interface{}) {
	l.lock()
	defer l.mutex.Unlock()
	l.fatal(fmt.Sprintf(format, vv...))
}

func (l *Logger) fatal(vv ...interface{}) {
	l.lib().Fatal(vv...)
	panic("gini: pkg: lg: handling error: expected " +
		"execution to stop")
}

// Lib provides standard library function which may fail or stop
// execution to make them easily mockable.
type Lib struct {

	// OpenFile defaults to os.OpenFile
	OpenFile func(name string, flags int, perm fs.FileMode) (
		*os.File, error)

	// ReadFile defaults to ioutil.ReadFile
	ReadFile func(name string) ([]byte, error)

	// Fatal defaults log.Fatal and is used by default for
	// [Logger.Fatal], [Logger.Fatalf] and for reporting log file
	// filesystem errors.
	Fatal func(vv ...interface{})

	// Println default logger is log.Println and is used if there is no
	// environment Env set.
	Println func(vv ...interface{})
}
