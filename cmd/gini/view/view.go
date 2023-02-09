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
Package view provides gini's user interface.  A gini-ui is made up of a
context bar and columns while a column consists of splits.  A split is
a component like an Editor.  For example the initial view of gini
executed inside its repository shows in one column and one split the
introductory help file.  This split is an Editor instance allowing to
edit the help file.  Further more on top in the first line the context
bar for displayed editable help files is shown:

+--------------------------------------------------------------------+
| ln: 0  cl: 0  mod:cmd  hlp/index.gnh        [messages]           • |
+-------------+-+------------------------------------+-+-------------+
|  remaining  |g|   only initial column and split    |g|  remaining  |
|   screen    |a|   showing the introductory help    |a|   screen    |
|    area     |p|   page within an view-Editor.      |p|    area     |
+-------------+-+------------------------------------+-+-------------+
*/
package view

import (
	"github.com/slukits/gini/cmd/gini/view/internal/cnt"
	"github.com/slukits/lines"
)

type View struct {
	lines.Component
	lines.Stacking
}

func (v *View) OnInit(e *lines.Env) {
	v.CC = append(v.CC, &cnt.Context{})
}

func (v *View) Context() lines.Componenter {
	return v.CC[0]
}
