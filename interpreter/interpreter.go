package interpreter

import (
	"go/token"

	"github.com/goruby/goruby/evaluator"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
)

// Interpreter defines the methods of an interpreter
type Interpreter interface {
	Interpret(string) (object.RubyObject, error)
}

// New returns an Interpreter ready to use and with the environment set to
// object.NewMainEnvironment()
func New() Interpreter {
	return &interpreter{environment: object.NewMainEnvironment()}
}

type interpreter struct {
	environment object.Environment
}

func (i *interpreter) Interpret(input string) (object.RubyObject, error) {
	node, err := parser.ParseFile(token.NewFileSet(), "", input)
	if err != nil {
		return nil, object.NewSyntaxError(err)
	}
	return evaluator.Eval(node, i.environment)
}
