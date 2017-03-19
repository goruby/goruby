package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/goruby/goruby/interpreter"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
)

const PROMPT = "girb:%03d> "

func Start(in io.Reader, out chan<- string) {
	scanner := bufio.NewScanner(in)
	counter := 1
	env := object.NewMainEnvironment()
	interpreter := interpreter.New()
	interpreter.SetEnvironment(env)
	var buffer string
	for {
		out <- fmt.Sprintf(PROMPT, counter)
		counter++
		scanned := scanner.Scan()
		if !scanned {
			out <- fmt.Sprintln()
			close(out)
			return
		}

		buffer += scanner.Text()
		evaluated, err := interpreter.Interpret(buffer)
		if err != nil {
			if parser.IsEOFError(err) {
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
