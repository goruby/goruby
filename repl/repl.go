package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/goruby/goruby/interpreter"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
)

const prompt = "girb:%03d> "

// Start starts the repl for the given reader in and writes its results to out.
func Start(in io.Reader, out chan<- string) {
	scanner := bufio.NewScanner(in)
	counter := 1
	interpreter := interpreter.New()
	var buffer string
	for {
		out <- fmt.Sprintf(prompt, counter)
		counter++
		scanned := scanner.Scan()
		if !scanned {
			out <- fmt.Sprintln()
			close(out)
			return
		}

		buffer += scanner.Text()
		evaluated, err := interpreter.Interpret("", buffer)
		if err != nil {
			if isEOFError(err) {
				buffer += "\n"
				continue
			}
			out <- fmt.Sprintf("%s\n", err.Error())
			buffer = ""
			continue
		}

		if evaluated != nil {
			out <- fmt.Sprintf("=> %s\n", evaluated.Inspect())
		}
		buffer = ""
	}
}

func isEOFError(err error) bool {
	syntaxError, ok := err.(*object.SyntaxError)
	if !ok {
		return false
	}
	err = syntaxError.UnderlyingError()
	return parser.IsEOFError(err)
}
