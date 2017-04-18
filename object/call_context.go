package object

import (
	"fmt"

	"github.com/goruby/goruby/ast"
)

// CallContext represents the context information when sending a message to an
// object.
type CallContext interface {
	// Env returns the current environment at call time
	Env() Environment
	// Eval represents an evalualtion method suitable to eval arbitrary Ruby
	// AST Nodes and transform them into a resulting Ruby object or an error.
	Eval(ast.Node, Environment) (RubyObject, error)
	// Receiver returns the Ruby object the message is sent to
	Receiver() RubyObject
}

// NewCallContext returns a new CallContext with a stubbed eval function
func NewCallContext(env Environment, receiver RubyObject) CallContext {
	return &callContext{
		env:      env,
		eval:     func(node ast.Node, env Environment) (RubyObject, error) { return nil, fmt.Errorf("No eval present") },
		receiver: receiver,
	}
}

type callContext struct {
	env      Environment
	eval     func(node ast.Node, env Environment) (RubyObject, error)
	receiver RubyObject
}

func (c *callContext) Env() Environment { return c.env }
func (c *callContext) Eval(node ast.Node, env Environment) (RubyObject, error) {
	return c.eval(node, env)
}
func (c *callContext) Receiver() RubyObject { return c.receiver }
