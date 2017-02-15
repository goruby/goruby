package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/goruby/goruby/evaluator"
	"github.com/goruby/goruby/lexer"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
)

const PROMPT = "girb:%03d> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	counter := 1
	env := object.NewEnvironment()
	for {
		fmt.Printf(PROMPT, counter)
		counter++
		scanned := scanner.Scan()
		if !scanned {
			fmt.Println()
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			fmt.Fprintf(out, "=> %s\n", evaluated.Inspect())
		}
	}
}

func printParserErrors(out io.Writer, errors []error) {
	fmt.Println("Parser errors: ")
	for _, err := range errors {
		fmt.Fprintf(out, "\t%s\n", err.Error())
	}
}
