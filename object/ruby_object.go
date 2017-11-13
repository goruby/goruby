package object

import (
	"bytes"
	"fmt"
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
	HASH_OBJ           Type = "HASH"
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
	New(args ...RubyObject) (RubyObject, error)
	Name() string
}

// RubyClassObject represents a class object in Ruby
type RubyClassObject interface {
	RubyObject
	RubyClass
}

type hashable interface {
	hashKey() hashKey
}

type extendable interface {
	addMethod(name string, method RubyMethod)
}

type extendableRubyObject interface {
	RubyObject
	extendable
}

type objects []RubyObject

func (o objects) String() string {
	out := []string{}
	for _, v := range o {
		out = append(out, v.Inspect())
	}
	return fmt.Sprintf("%s", out)
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

type functionParameters []*FunctionParameter

func (f functionParameters) defaultParamCount() int {
	count := 0
	for _, p := range f {
		if p.Default != nil {
			count++
		}
	}
	return count
}
func (f functionParameters) mandatoryParams() []*FunctionParameter {
	params := make([]*FunctionParameter, 0)
	for _, p := range f {
		if p.Default == nil {
			params = append(params, p)
		}
	}
	return params
}
func (f functionParameters) optionalParams() []*FunctionParameter {
	params := make([]*FunctionParameter, 0)
	for _, p := range f {
		if p.Default != nil {
			params = append(params, p)
		}
	}
	return params
}
func (f functionParameters) separateDefaultParams() ([]*FunctionParameter, []*FunctionParameter) {
	mandatory, defaults := make([]*FunctionParameter, 0), make([]*FunctionParameter, 0)
	for _, p := range f {
		if p.Default != nil {
			defaults = append(defaults, p)
		} else {
			mandatory = append(mandatory, p)
		}
	}
	return mandatory, defaults
}

// FunctionParameter represents a parameter within a function
type FunctionParameter struct {
	Name    string
	Default RubyObject
}

func (f *FunctionParameter) String() string {
	var out bytes.Buffer
	out.WriteString(f.Name)
	if f.Default != nil {
		out.WriteString(" = ")
		out.WriteString(f.Default.Inspect())
	}
	return out.String()
}

// A Function represents a user defined function. It is no real Ruby object.
type Function struct {
	Parameters       []*FunctionParameter
	Body             *ast.BlockStatement
	Env              Environment
	MethodVisibility MethodVisibility
}

// String returns the function literal
func (f *Function) String() string {
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

// Call implements the RubyMethod interface. It evaluates f.Body and returns its result
func (f *Function) Call(context CallContext, args ...RubyObject) (RubyObject, error) {
	block, arguments, _ := extractBlockFromArgs(args)
	defaultParams := functionParameters(f.Parameters).defaultParamCount()
	if len(arguments) < len(f.Parameters)-defaultParams || len(arguments) > len(f.Parameters) {
		return nil, NewWrongNumberOfArgumentsError(len(f.Parameters), len(arguments))
	}
	params, err := f.populateParameters(arguments)
	if err != nil {
		return nil, err
	}
	contextSelf, _ := context.Env().Get("self")
	contextSelfObject := contextSelf.(*Self)
	extendedEnv := f.extendFunctionEnv(contextSelfObject, params, block)
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

func (f *Function) populateParameters(args []RubyObject) (map[string]RubyObject, error) {
	if len(args) > len(f.Parameters) {
		return nil, NewWrongNumberOfArgumentsError(len(f.Parameters), len(args))
	}
	params := make(map[string]RubyObject)

	mandatory, defaults := functionParameters(f.Parameters).separateDefaultParams()

	if len(args) < len(mandatory)-len(defaults) || len(args) > len(f.Parameters) {
		return nil, NewWrongNumberOfArgumentsError(len(f.Parameters), len(args))
	}

	if len(args) == len(f.Parameters) {
		for paramIdx, param := range f.Parameters {
			params[param.Name] = args[paramIdx]
		}
		return params, nil
	}

	parameters := append(mandatory, defaults...)

	for paramIdx, param := range parameters {
		if paramIdx >= len(args) {
			params[param.Name] = param.Default
			continue
		}
		params[param.Name] = args[paramIdx]
	}
	return params, nil
}

func (f *Function) extendFunctionEnv(context *Self, params map[string]RubyObject, block *Proc) Environment {
	// encapsulate the block within a new self, but with the same object
	funcSelf := &Self{RubyObject: context.RubyObject, Name: context.Name, Block: block}
	env := NewEnclosedEnvironment(f.Env)
	env.Set("self", funcSelf)
	for k, v := range params {
		env.Set(k, v)
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
	RubyObject        // The encapsuled object acting as self
	Block      *Proc  // the block given to the current execution binding
	Name       string // The name of self in this context
}

// Type returns SELF
func (s *Self) Type() Type { return SELF }

// Inspect returns the name of Self
func (s *Self) Inspect() string { return s.Name }

// extendedObject is a wrapper object for an object extended by methods.
type extendedObject struct {
	RubyObject
	class *eigenclass
	Environment
}

func (e *extendedObject) Class() RubyClass { return e.class }
func (e *extendedObject) addMethod(name string, method RubyMethod) {
	e.class.addMethod(name, method)
}
