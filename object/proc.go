package object

import (
	"bytes"
	"strings"

	"github.com/goruby/goruby/ast"
)

var procClass RubyClassObject = NewClass(
	"Proc", objectClass, procMethods, procClassMethods,
)

func init() {
	classes.Set("Proc", procClass)
}

// A Proc represents a user defined block of code.
type Proc struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        Environment
}

// Type returns proc_OBJ
func (p *Proc) Type() Type { return "" }

// Inspect returns the proc body
func (p *Proc) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range p.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("do |")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString("| \n")
	out.WriteString(p.Body.String())
	out.WriteString("\nend")
	return out.String()
}

// Class returns procClass
func (p *Proc) Class() RubyClass { return procClass }

// Call implements the RubyMethod interface. It evaluates p.Body and returns its result
func (p *Proc) Call(context CallContext, args ...RubyObject) (RubyObject, error) {
	if len(args) != len(p.Parameters) {
		return nil, NewWrongNumberOfArgumentsError(len(p.Parameters), len(args))
	}
	extendedEnv := p.extendProcEnv(args)
	evaluated, err := context.Eval(p.Body, extendedEnv)
	if err != nil {
		return nil, err
	}
	return evaluated, nil
}

func (p *Proc) extendProcEnv(args []RubyObject) Environment {
	env := NewEnclosedEnvironment(p.Env)
	for paramIdx, param := range p.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

var procClassMethods = map[string]RubyMethod{}

var procMethods = map[string]RubyMethod{}