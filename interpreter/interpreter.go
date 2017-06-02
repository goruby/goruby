package interpreter

import (
	"go/token"
	"log"
	"os"

	"github.com/goruby/goruby/evaluator"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
)

// Interpreter defines the methods of an interpreter
type Interpreter interface {
	Interpret(filename string, input interface{}) (object.RubyObject, error)
}

// New returns an Interpreter ready to use and with the environment set to
// object.NewMainEnvironment()
func New() Interpreter {
	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Cannot get working directory: %s\n", err)
	}
	env := object.NewMainEnvironment()
	loadPath, _ := env.Get("$:")
	loadPathArr := loadPath.(*object.Array)
	loadPathArr.Elements = append(loadPathArr.Elements, &object.String{Value: cwd})
	env.SetGlobal("$:", loadPathArr)
	return &interpreter{environment: env}
}

type interpreter struct {
	environment object.Environment
}

func (i *interpreter) Interpret(filename string, input interface{}) (object.RubyObject, error) {
	node, err := parser.ParseFile(token.NewFileSet(), filename, input, 0)
	if err != nil {
		return nil, object.NewSyntaxError(err)
	}
	return evaluator.Eval(node, i.environment)
}
