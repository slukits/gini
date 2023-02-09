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
Package main provides the gini-command.  GINI Is Not an IDE but lets you
build one to create, modify and process parsable text.

gini's architecture is following a typical MVC-pattern (model, view,
controller).  We over-package to accomplish shy and strongly
encapsulated code.  NOTE the MVC-pattern description in Wikipedia at the
time of this writing is plain wrong.  The whole purpose of the MVC
pattern is that Model and View can be developed in isolation of any
other project code only needing to care about their respective API.
I.e. letting the view access the model is an anti-pattern without any
benefit for the software architecture.  Hence in GINI neither the view
nor the model access any other part of the GINI command implementation.
Instead the controller leverages their APIs to receive, to process and
to respond to user input.
*/
package main

import "github.com/slukits/gini/cmd/gini/controller"

func main() { controller.New(controller.Init{}) }
