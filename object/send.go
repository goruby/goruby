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

// Send sends message method with args to context and returns its result
func Send(context CallContext, method string, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	class := receiver.Class()

	// search for the method in the ancestry tree
	for class != nil {
		fn, ok := class.Methods()[method]
		if !ok {
			class = class.SuperClass()
			continue
		}

		if fn.Visibility() == PRIVATE_METHOD && receiver.Type() != SELF {
			return nil, NewPrivateNoMethodError(receiver, method)
		}

		return fn.Call(context, args...)
	}

	methodMissingArgs := append(
		[]RubyObject{&Symbol{method}},
		args...,
	)

	return methodMissing(context, methodMissingArgs...)
}

// AddMethod adds a method to a given object. It returns the object with the modified method set
func AddMethod(context RubyObject, methodName string, method *Function) RubyObject {
	objectToExtend := context
	self, contextIsSelf := context.(*Self)
	if contextIsSelf {
		objectToExtend = self.RubyObject
	}
	extended, ok := objectToExtend.(*extendedObject)
	if !ok {
		extended = &extendedObject{
			RubyObject: context,
			class:      newEigenclass(context.Class(), map[string]RubyMethod{}),
		}
	}
	extended.addMethod(methodName, method)
	if contextIsSelf {
		self.RubyObject = extended
		return self
	}
	return extended
}

func methodMissing(context CallContext, args ...RubyObject) (RubyObject, error) {
	class := context.Receiver().Class()

	// search for method_missing in the ancestry tree
	for class != nil {
		fn, ok := class.Methods()["method_missing"]
		if !ok {
			class = class.SuperClass()
			continue
		}
		return fn.Call(context, args...)
	}
	return nil, NewNoMethodError(context.Receiver(), args[0].(*Symbol).Value)
}
