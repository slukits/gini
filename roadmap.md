# Road map

version 0.1
- minimal context-bar implementation
- initial program-start showing introductory help page
- editing help pages having the following commands available:
  * Space  activate context bar
    * Space leave context bar if not input-box in context bar
      otherwise insert space into context-bar's input-box
    * Esc leave context bar discard/undo all user input associated with
      the last context.
    * Enter leave context bar applying/keeping effect of user input
    * jk in input box (e.g. context-bar) switches to mini command mode
      * mini command mode supports i/h/l/H/L/c/C/d/D/p/P whereas 
	H/L/c/C/d/D expect a letter or $^ defining until which letter should be
	moved/copied/deleted.
  * h/j/k/l/cursor-keys  left/down/up/right movement
  * H/L left/right movement until given rune
  * J/K down/up movement for given number of lines
  * t go til context-bar:
      'search text' ␣  word end  next word  line end  including
  * T go backward til context-bar:
      'search text' ␣  word beginning previous word line beginning
  * $^0 special movement keys for line end/first non-blank/line begin
  * S save all modified files
  * Backspace delete rune left of cursor or connect lines
  * Delete delete rune under cursor or connect lines
  * Enter add new line after current line
  * Shift+Enter add new line before current line
  * c copy til:
      'search text' ␣  word end  next word  line end  including
  * C copy backwards til:
      'search text' ␣  word begin  prev word  line begin  including
  * d delete til:
      'search text' ␣  word end  next word  line end  including
  * D delete backwards til:
      'search text' ␣  word begin  prev word  line begin  including
  * p past last copy/delete at current position
  * P context-bar
      'search text overwrite'  ␣  copy-ring delete-ring including
  * g grep current file; context-bar:
      'search text' ␣  under  re  replace  save  ring
  * G grep current directory; context-bar: see g
  * CtrlG grep current project directory recursively; context-bar: see g
  * * + grep-save hotkey: grep for saved expression
  * Insert switch to insert-mode if in command-mode
	switch to overwrite-mode if in insert-mode
	switch to insert-mode if in overwrite-mode
  * Esc/jk  switch to command mode if in insert- or overwrite-mode
  * s save/touch modified file
  * S save all modified files
  * r record keyboard macro
  * e+* edit something, e.g.:
    * e+g edit grep
    * e+r edit keyboard macro recording
    * e+p edit current file-type parser
    * e+l edit current file-type lexer
  * a do last command again
  * A do last command from command ring again
  * u undo last modification
  * U redo last undo
  * m do last movement again
  * M move to position in the movement ring
  Note Backspace/Delete/Enter/Shift+Enter/cursor-keys work in insert and
  command mode the same.


version 0.2
- minimal help-context implementation
- minimal settings-context implementation


version 0.3
- spell-checking lexer
- comandline spell-checker integration

version 0.4
- basic file management
  * open files
  * recently opend files
  * files of current directory
  * navigate "project directories"

version 0.5
- basic git integration
  * initialize a git repository
  * show current state
  * stage changs
  * commit staged changes
  * log of a file's commit-history

version 0.6
- parsing help pages
- parser based display formatting of help pages
- parser based help-text navigation by indices

version 0.7
- create lexer/parser for commmand output and connect it with the ui
- integration of command-line linter and compiler
- integration of command-line refactoring

version 0.8
- completion server
  * path completion
  * code completion based on tag-files

From here on GINI should be useable to program GINI

version 0.9
- improve lexer and parser to parse (most of) the go standard library
  * scopes/name-spaces
  * typing
  * imports
- parser based code-completion
- parser based code-navigation
- parser based refactoring

version 1.0
- add graphical user interface (next to terminal user interface)
- improve lexer and parser to parse zig, python, skala and java

