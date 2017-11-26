package repl

import (
	"fmt"
	"io"
	"strings"

	"github.com/goruby/goruby/interpreter"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
	"github.com/pkg/errors"
)

// Input defines the input interface for the repl
type Input interface {
	// Readline returns the next line of input. When it returns io.EOF, the
	// repl exits
	Readline() (string, error)
}

// Prompt represents a way to set the prompt
type Prompt interface {
	// SetPrompt is called after every evaluation
	SetPrompt(string)
}

// The PromptFunc type is an adapter to allow the use of ordinary functions as
// prompt. If f is a function with the appropriate signature, PromptFunc(f) is
// a Prompt that calls f.
type PromptFunc func(string)

// SetPrompt calls f(prompt).
func (f PromptFunc) SetPrompt(prompt string) { f(prompt) }

// Repl defines the interface to the Repl
type Repl interface {
	// Start starts the repl. If an error occurs, e.g. from Input, it will be returned.
	Start() error
}

// New returns a repl
func New(input Input, output io.Writer, prompt Prompt) Repl {
	return &repl{
		input:  input,
		output: output,
		prompt: prompt,
		interpreter: &bufferedInterpreter{
			interpreter: interpreter.New(),
		},
	}
}

type repl struct {
	input       Input
	output      io.Writer
	prompt      Prompt
	interpreter *bufferedInterpreter
}

func (r *repl) Start() error {
	p := &prompt{}
	r.prompt.SetPrompt(p.prompt())
	for {
		// Read a line
		line, err := r.input.Readline()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		r.interpreter.interpretLine(line, r.output)
		r.prompt.SetPrompt(p.prompt())
	}
	return nil
}

// bufferedInterpreter provides a buffered evaluator
type bufferedInterpreter struct {
	interpreter interpreter.Interpreter
	buffer      string
	linecount   int
}

// interpretLine evaluates the next line from the input and writes the results to out
func (b *bufferedInterpreter) interpretLine(line string, out io.Writer) {
	b.buffer += line
	evaluated, err := b.interpreter.Interpret("(irb)", b.preserveLineCount(b.buffer))
	if err != nil {
		if isEOFError(err) {
			b.buffer += "\n"
			return
		}
		fmt.Fprintf(out, "%s\n", errors.Cause(err).Error())
		b.buffer = ""
		return
	}

	if evaluated != nil {
		fmt.Fprintf(out, "=> %s\n", evaluated.Inspect())
	}
	b.buffer = ""
}

func (b *bufferedInterpreter) preserveLineCount(input string) string {
	return strings.Repeat("\n", b.linecount) + input
}

func (b *bufferedInterpreter) getLines(input string) {
	b.linecount = b.linecount + len(strings.Split(input, "\n"))
}

func isEOFError(err error) bool {
	syntaxError, ok := err.(*object.SyntaxError)
	if !ok {
		return false
	}
	err = syntaxError.UnderlyingError()
	return parser.IsEOFError(err)
}

const promptTemplate = "girb:%03d> "

type prompt struct {
	counter int
}

func (p *prompt) prompt() string {
	p.counter++
	return fmt.Sprintf(promptTemplate, p.counter)
}
