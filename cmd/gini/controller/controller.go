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
Package controller maps from a view reported user input to queries or
modifications of the model whose results define the view-changes
expressing the response.
*/
package controller

import (
	"github.com/slukits/gini/cmd/gini/view"
	"github.com/slukits/gini/pkg/lg"
	"github.com/slukits/lines"
)

type Controller struct{}

// New creates an ui by initializing a new View v and using init's
// UIFactory followed by waiting (blocking) on v for user-input to
// process until a quit request is received which terminates executed
// gini-instance.  Note the init-instance argument allows to inject an
// ui-factory which creates a lines terminal fixture for testing.
func New(init Init) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		init.Log.Fatalf("GINI: controller: panic: %v", err)
	}()
	ll := init.UIFactory()(&view.View{})
	ll.WaitForQuit()
}

type Init struct {
	Lines func(lines.Componenter) *lines.Lines
	Log   lg.Logger
}

func (i *Init) UIFactory() func(lines.Componenter) *lines.Lines {
	if i.Lines == nil {
		return lines.Term
	}
	return i.Lines
}
