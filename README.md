goruby
======

GoRuby, an implementation of Ruby written in Go

## Contribution

If anyone wants to help to get the project to the real implementation please ping me or fork it and send a pull request.

## REPL
There is a basic REPL within `cmd/girb`. It supports multiline expressions and all syntax elements the language supports yet.

To run it ad hoc run `go run cmd/girb/main.go` and exit the REPL with CTRL-D.

## Command
To run the command as one off run `go run main.go`.

## Supported features

### `goruby` Command
- [x] parse program files
- [ ] program file arguments
- [ ] Flags
  - [ ] `-0[octal]`       specify record separator (\0, if no argument)
  - [ ] `-a`              autosplit mode with -n or -p (splits $_ into $F)
  - [ ] `-c`              check syntax only
  - [ ] `-Cdirectory`     cd to directory before executing your script
  - [ ] `-d`              set debugging flags (set $DEBUG to true)
  - [x] `-e 'command'`    one line of script. Several -e's allowed. Omit [programfile]
  - [ ] `-Eex[:in]`       specify the default external and internal character encodings
  - [ ] `-Fpattern`       split() pattern for autosplit (-a)
  - [ ] `-i[extension]`   edit ARGV files in place (make backup if extension supplied)
  - [ ] `-Idirectory`     specify $LOAD_PATH directory (may be used more than once)
  - [ ] `-l`              enable line ending processing
  - [ ] `-n`              assume 'while gets(); ... end' loop around your script
  - [ ] `-p`              assume loop like -n but print line also like sed
  - [ ] `-rlibrary`       require the library before executing your script
  - [ ] `-s`              enable some switch parsing for switches after script name
  - [ ] `-S`              look for the script using PATH environment variable
  - [ ] `-T[level=1]`     turn on tainting checks
  - [ ] `-v`              print version number, then turn on verbose mode
  - [ ] `-w`              turn warnings on for your script
  - [ ] `-W[level=2]`     set warning level; 0=silence, 1=medium, 2=verbose
  - [ ] `-x[directory]`   strip off text before #!ruby line and perhaps cd to directory
  - [ ] `-h`              show this message, --help for more info

### `girb` Command
- [ ] parse program files
- [ ] program file arguments
- [ ] Flags
  - [ ] `-f`		    Suppress read of ~/.irbrc
  - [ ] `-m`		    Bc mode (load mathn, fraction or matrix are available)
  - [ ] `-d`                Set $DEBUG to true (same as `ruby -d')
  - [ ] `-r load-module`    Same as `ruby -r'
  - [ ] `-I path`           Specify $LOAD_PATH directory
  - [ ] `-U`                Same as `ruby -U`
  - [ ] `-E enc`            Same as `ruby -E`
  - [ ] `-w`                Same as `ruby -w`
  - [ ] `-W[level=2]`       Same as `ruby -W`
  - [ ] `--context-mode n`  Set n[0-3] to method to create Binding Object,
                    when new workspace was created
  - [ ] `--echo`            Show result(default)
  - [ ] `--noecho`          Don't show result
  - [ ] `--inspect`	    Use `inspect' for output (default except for bc mode)
  - [ ] `--noinspect`	    Don't use inspect for output
  - [ ] `--readline`        Use Readline extension module
  - [ ] `--noreadline`	    Don't use Readline extension module
  - [ ] `--prompt prompt-mode`/`--prompt-mode prompt-mode`
		    Switch prompt mode. Pre-defined prompt modes are
		    `default', `simple', `xmp' and `inf-ruby'
  - [ ] `--inf-ruby-mode`   Use prompt appropriate for inf-ruby-mode on emacs.
		    Suppresses --readline.
  - [ ] `--sample-book-mode`/`--simple-prompt`
                    Simple prompt mode
  - [ ] `--noprompt`        No prompt mode
  - [ ] `--single-irb`      Share self with sub-irb.
  - [ ] `--tracer`          Display trace for each execution of commands.
  - [ ] `--back-trace-limit n`
		    Display backtrace top n and tail n. The default
		    value is 16.
  - [ ] `--irb_debug n`	    Set internal debug level to n (not for popular use)
  - [ ] `--verbose`         Show details
  - [ ] `--noverbose`       Don't show details
  - [ ] `-v`, `--version`	  Print the version of irb
  - [ ] `-h`, `--help`      Print help
  - [ ] `--`                Separate options of irb from the list of command-line args


### Supported language feature
- [ ] everything is an object
	- [x] allow method calls on everything
	- [ ] operators are method calls
- [ ] full UTF8 support
	- [ ] Unicode identifier
	- [ ] Unicode symbols
- [x] functions
	- [x] with parens
	- [x] without parens
	- [x] return keyword
	- [ ] default values for parameters
	- [ ] keyword arguments
	- [ ] block arguments
- [x] function calls
	- [x] with parens
	- [x] without parens
- [ ] conditionals
	- [x] if
	- [x] if/else
	- [ ] if/elif/else
	- [ ] tenary `? : `
	- [ ] unless
	- [ ] case
	- [ ] `||`
	- [ ] `&&`
- [ ] control flow
	- [ ] for loop
	- [ ] while loop
	- [ ] until loop
	- [ ] break
	- [ ] next
	- [ ] redo
	- [ ] flip flop
- [ ] numbers
	- [ ] integers
		- [x] integer arithmetics
		- [x] integers `1234`
		- [ ] integers with underscores `1_234`
		- [ ] decimal numbers `0d170`, `0D170`
		- [ ] octal numbers `0252`, `0o252`, `0O252`
		- [ ] hexadecimal numbers `0xaa`, `0xAa`, `0xAA`, `0Xaa`, `0XAa`, `0XaA`
		- [ ] binary numbers `0b10101010`, `0B10101010`
	- [ ] floats
		- [ ] float arithmetics
		- [ ] `12.34`
		- [ ] `1234e-2`
		- [ ] `1.234E1`
- [x] booleans
- [ ] strings
	- [x] double quoted
	- [ ] single quoted
	- [ ] `%q{}`
	- [ ] `%Q{}`
	- [ ] heredoc
		- [ ] without indentation (`<<EOF`)
		- [ ] indented (`<<-EOF`)
		- [ ] “squiggly” heredoc `<<~`
		- [ ] quoted heredoc
			- [ ] single quotes `<<-'HEREDOC'`
 			- [ ] double quotes `<<-"HEREDOC"`
 			- [ ] backticks <<-\`HEREDOC\`"
	- [ ] escaped characters
		- [ ] `\a` bell, ASCII 07h (BEL)
		- [ ] 	`\b` backspace, ASCII 08h (BS)
		- [ ] 	`\t` horizontal tab, ASCII 09h (TAB)
		- [ ] 	`\n` newline (line feed), ASCII 0Ah (LF)
		- [ ] 	`\v` vertical tab, ASCII 0Bh (VT)
		- [ ] 	`\f` form feed, ASCII 0Ch (FF)
		- [ ] 	`\r` carriage return, ASCII 0Dh (CR)
		- [ ] 	`\e` escape, ASCII 1Bh (ESC)
		- [ ] 	`\s` space, ASCII 20h (SPC)
		- [ ] 	`\\` backslash, \
		- [ ] 	`\nnn` octal bit pattern, where nnn is 1-3 octal digits ([0-7])
		- [ ] 	`\xnn` hexadecimal bit pattern, where nn is 1-2 hexadecimal digits ([0-9a-fA-F])
		- [ ] `\unnnn` Unicode character, where nnnn is exactly 4 hexadecimal digits ([0-9a-fA-F])
		- [ ] `\u{nnnn ...}` Unicode character(s), where each nnnn is 1-6 hexadecimal digits ([0-9a-fA-F])
		- [ ] `\cx` or `\C-x` control character, where x is an ASCII printable character
		- [ ] `\M-x` meta character, where x is an ASCII printable character
		- [ ] `\M-\C-x` meta control character, where x is an ASCII printable character
		- [ ] `\M-\cx` same as above
		- [ ] `\c\M-x` same as above
		- [ ] `\c?` or `\C-?` delete, ASCII 7Fh (DEL)
	- [ ] interpolation `#{}`
	- [ ] automatic concatenation
- [ ] arrays
	- [x] array literal `[1,2]`
	- [x] array indexing `arr[2]`
	- [ ] splat
	- [ ] array decomposition
	- [ ] implicit array assignment
	- [ ] array of strings `%w{}`
	- [ ] array of symbols `%i{}`
- [ ] nil
- [ ] hashes
- [ ] symbols
	- [x] `:symbol`
	- [ ] `:"symbol"`
	- [ ] `:"symbol"` with interpolation
	- [ ] `:'symbol'`
	- [ ] `%s{symbol}`
	- [ ] singleton symbols
- [ ] regexp
	- [ ] `/regex/`
	- [ ] `%r{regex}`
- [ ] ranges
	- [ ] `..` inclusive
	- [ ] `...` exclusive
- [ ] procs `->`
- [ ] variables
	- [x] variable assignments
	- [ ] globals
- [ ] operators
	- [x] `+`
	- [x] `-`
	- [x] `/`
	- [x] `*`
	- [x] `!`
	- [x] `<`
	- [x] `>`
	- [ ] `**` (pow)
	- [ ] `%` (modulus)
	- [ ] `&` (AND)
	- [ ] `^` (XOR)
	- [ ] `>>` (right shift)
	- [ ] `<<` (left shift, append)
	- [x] `==` (equal)
	- [x] `!=` (not equal)
	- [ ] `===` (case equality)
	- [ ] `=~` (pattern match)
	- [ ] `!~` (does not match)
	- [ ] `<=>` (comparison or spaceship operator)
	- [ ] `<=` (less or equal)
	- [ ] `>=` (greater or equal)
	- [ ] assignment operators
		- [ ] `+=`
		- [ ] `-=`
		- [ ] `/=`
		- [ ] `*=`
		- [ ] `%=`
		- [ ] `**=`
		- [ ] `&=`
		- [ ] `|=`
		- [ ] `^=`
		- [ ] `<<=`
		- [ ] `>>=`
		- [ ] `||=`
		- [ ] `&&=`
- [ ] function blocks (procs)
- [ ] constants
- [ ] scope operator `::`
- [ ] classes
	- [x] class objects
	- [x] class Class
	- [ ] instance variables
	- [ ] class variables
	- [x] class methods
	- [x] instance methods
	- [ ] method overrides
	- [ ] private
	- [ ] protected
	- [ ] public
	- [x] inheritance
	- [ ] constructors
	- [x] new
	- [ ] `self`
	- [ ] singleton classes (also known as the metaclass or eigenclass) `class << self`
	- [ ] assigment methods
	- [ ] self defined classes
	- [ ] self defined classes with inheritance
- [x] modules
- [ ] object main

