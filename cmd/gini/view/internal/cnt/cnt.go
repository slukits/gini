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
Package cnt provides a gini view's Context and all context-features
around it.
*/
package cnt

import (
	"fmt"

	"github.com/slukits/lines"
)

const DefaultContent = "hlp/index.gnh"

type Context struct{ lines.Component }

func (c *Context) OnInit(e *lines.Env) {
	fmt.Fprint(e, DefaultContent)
}
