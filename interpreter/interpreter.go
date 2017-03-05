package interpreter

import (
	"fmt"

	"github.com/goruby/goruby/ast"
	"github.com/goruby/goruby/evaluator"
	"github.com/goruby/goruby/lexer"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
)

func New() *Interpreter {
	return &Interpreter{environment: object.NewMainEnvironment()}
}

type Interpreter struct {
	environment *object.Environment
}

func (i *Interpreter) Interpret(input string) (object.Object, error) {
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

func (i *Interpreter) SetEnvironment(env *object.Environment) {
	i.environment = env
}

func (i *Interpreter) parse(input string) (ast.Node, error) {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}
