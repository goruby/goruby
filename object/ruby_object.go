package object

import (
	"bytes"
	"strings"

	"github.com/goruby/goruby/ast"
)

// Type represents a type of an object
type Type string

const (
	EIGENCLASS_OBJ     Type = "EIGENCLASS"
	FUNCTION_OBJ       Type = "FUNCTION"
	RETURN_VALUE_OBJ   Type = "RETURN_VALUE"
	BASIC_OBJECT_OBJ   Type = "BASIC_OBJECT"
	OBJECT_OBJ         Type = "OBJECT"
	CLASS_OBJ          Type = "CLASS"
	CLASS_INSTANCE_OBJ Type = "CLASS"
	ARRAY_OBJ          Type = "ARRAY"
	INTEGER_OBJ        Type = "INTEGER"
	STRING_OBJ         Type = "STRING"
	SYMBOL_OBJ         Type = "SYMBOL"
	BOOLEAN_OBJ        Type = "BOOLEAN"
	NIL_OBJ            Type = "NIL"
	NIL_CLASS_OBJ      Type = "NIL_CLASS"
	EXCEPTION_OBJ      Type = "EXCEPTION"
	MODULE_OBJ         Type = "MODULE"
	SELF               Type = "SELF"
)

type inspectable interface {
	Inspect() string
}

// RubyObject represents an object in Ruby
type RubyObject interface {
	inspectable
	Type() Type
	Class() RubyClass
}

// RubyClass represents a class in Ruby
type RubyClass interface {
	Methods() MethodSet
	SuperClass() RubyClass
}

// RubyClassObject represents a class object in Ruby
type RubyClassObject interface {
	RubyObject
	RubyClass
}

type extendable interface {
	addMethod(name string, method RubyMethod)
}

type extendableRubyObject interface {
	RubyObject
	extendable
}

// ReturnValue represents a wrapper object for a return statement. It is no
// real Ruby object and only used within the interpreter evaluation
type ReturnValue struct {
	Value RubyObject
}

// Type returns RETURN_VALUE_OBJ
func (rv *ReturnValue) Type() Type { return RETURN_VALUE_OBJ }

// Inspect returns the string representation of the wrapped object
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Class reurns the class of the wrapped object
func (rv *ReturnValue) Class() RubyClass { return rv.Value.Class() }

// A Function represents a user defined function. It is no real Ruby object.
type Function struct {
	Parameters       []*ast.Identifier
	Body             *ast.BlockStatement
	Env              Environment
	MethodVisibility MethodVisibility
}

// Type returns FUNCTION_OBJ
func (f *Function) Type() Type { return FUNCTION_OBJ }

// Inspect returns the function body
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// Class returns nil
func (f *Function) Class() RubyClass { return nil }

// Call implements the RubyMethod interface. It evaluates f.Body and returns its result
func (f *Function) Call(context CallContext, args ...RubyObject) (RubyObject, error) {
	if len(args) != len(f.Parameters) {
		return nil, NewWrongNumberOfArgumentsError(len(f.Parameters), len(args))
	}
	extendedEnv := f.extendFunctionEnv(args)
	evaluated, err := context.Eval(f.Body, extendedEnv)
	if err != nil {
		return nil, err
	}
	return f.unwrapReturnValue(evaluated), nil
}

// Visibility implements the RubyMethod interface. It returns f.MethodVisibility
func (f *Function) Visibility() MethodVisibility {
	return f.MethodVisibility
}

func (f *Function) extendFunctionEnv(args []RubyObject) Environment {
	env := NewEnclosedEnvironment(f.Env)
	for paramIdx, param := range f.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func (f *Function) unwrapReturnValue(obj RubyObject) RubyObject {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// Self represents the value associated to `self`. It acts as a wrapper around
// the RubyObject and is just meant to indicate that the given object is
// self in the given context.
type Self struct {
	RubyObject
	Name string
}

// Type returns SELF
func (s *Self) Type() Type { return SELF }

// Inspect returns the name of Self
func (s *Self) Inspect() string { return s.Name }

// extendedObject is a wrapper object for an object extended by methods.
type extendedObject struct {
	RubyObject
	class *eigenclass
}

func (e *extendedObject) Class() RubyClass { return e.class }
func (e *extendedObject) addMethod(name string, method RubyMethod) {
	e.class.addMethod(name, method)
}
