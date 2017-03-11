package interpreter

import (
	"fmt"

	"github.com/goruby/goruby/ast"
	"github.com/goruby/goruby/evaluator"
	"github.com/goruby/goruby/lexer"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
)

type Interpreter interface {
	Interpret(string) (object.RubyObject, error)
	SetEnvironment(*object.Environment)
}

func New() Interpreter {
	return &interpreter{environment: object.NewMainEnvironment()}
}

type interpreter struct {
	environment *object.Environment
}

func (i *interpreter) Interpret(input string) (object.RubyObject, error) {
	node, err := i.parse(input)
	if err != nil {
		return nil, err
	}
	evaluated := evaluator.Eval(node, i.environment)
	if evaluator.IsError(evaluated) {
		return nil, fmt.Errorf(evaluated.Inspect())
	}
	return evaluated, nil
}

func (i *interpreter) SetEnvironment(env *object.Environment) {
	i.environment = env
}

func (i *interpreter) parse(input string) (ast.Node, error) {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
