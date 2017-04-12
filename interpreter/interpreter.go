package interpreter

import (
	"github.com/goruby/goruby/ast"
	"github.com/goruby/goruby/evaluator"
	"github.com/goruby/goruby/lexer"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
)

// Interpreter defines the methods of an interpreter
type Interpreter interface {
	Interpret(string) (object.RubyObject, error)
	SetEnvironment(object.Environment)
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
	node, err := i.parse(input)
	if err != nil {
		return nil, object.NewSyntaxError(err.Error())
	}
	return evaluator.Eval(node, i.environment)
}

func (i *interpreter) SetEnvironment(env object.Environment) {
	i.environment = env
}

func (i *interpreter) parse(input string) (ast.Node, error) {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
