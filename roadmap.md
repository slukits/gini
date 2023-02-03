# Road map

[version 0.1](rm01.md)
- initial program-start showing introductory help page
- minimal context-bar implementation for initial help page
- minimal editor implementation for editing initial help page


[version 0.2](rm02.md)
- minimal (global) help-context implementation
- minimal (global) settings-context implementation
- extend editor commands by "find-movement" for movement, copying and
  deleting.

version 0.3
- spell-checking lexer
- command line spell-checker integration
- extend editor by GINI's idea of a grep-command

version 0.4
- system integration
  * command-line linter
  * command-line compiler
  * command-line testing
  * command-line refactoring
- support for go to develope GINI with GINI

version 0.5
- basic file management
  * open files
  * recently opened files
  * files of current directory
  * navigate "project directories"
- basic git integration

version 0.6
- parsing help pages
- parser based display formatting of help pages
- parser based help-text navigation by indices

version 0.7
- create lexer/parser for command output and connect it with the ui
- integration of command-line linter and compiler
- integration of command-line refactoring

version 0.8
- completion server
  * path completion
  * code completion based on tag-files

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

