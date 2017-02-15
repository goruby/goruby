package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"github.com/goruby/goruby/ast"
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
		fmt.Fprintf(out, PROMPT, counter)
		counter++
		scanned := scanner.Scan()
		if !scanned {
			fmt.Fprintln(out)
			return
		}

		line := scanner.Text()
		node, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(out, "%s", err.Error())
		}

		evaluated := evaluator.Eval(node, env)
		if evaluated != nil {
			fmt.Fprintf(out, "=> %s\n", evaluated.Inspect())
		}
	}
}

func parseLine(line string) (ast.Node, error) {
	l := lexer.New(line)
	p := parser.New(l)
	var err error
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		err = mergeParserErrors(p.Errors())
	}
	return program, err
}

func mergeParserErrors(errors []error) error {
	var buf bytes.Buffer
	printParserErrors(&buf, errors)
	return fmt.Errorf(buf.String())
}

func printParserErrors(out io.Writer, errors []error) {
	fmt.Println("Parser errors: ")
	for _, err := range errors {
		fmt.Fprintf(out, "\t%s\n", err.Error())
	}
}
