package object

import (
	"bytes"
	"strings"

	"github.com/goruby/goruby/ast"
)

// Type represents a type of an object
type Type string

const (
	EIGENCLASS_OBJ         Type = "EIGENCLASS"
	FUNCTION_OBJ           Type = "FUNCTION"
	RETURN_VALUE_OBJ       Type = "RETURN_VALUE"
	BASIC_OBJECT_OBJ       Type = "BASIC_OBJECT"
	BASIC_OBJECT_CLASS_OBJ Type = "BASIC_OBJECT_CLASS"
	OBJECT_OBJ             Type = "OBJECT"
	OBJECT_CLASS_OBJ       Type = "OBJECT_CLASS"
	CLASS_OBJ              Type = "CLASS"
	CLASS_CLASS_OBJ        Type = "CLASS_CLASS"
	ARRAY_OBJ              Type = "ARRAY"
	ARRAY_CLASS_OBJ        Type = "ARRAY_CLASS"
	INTEGER_OBJ            Type = "INTEGER"
	INTEGER_CLASS_OBJ      Type = "INTEGER_CLASS"
	STRING_OBJ             Type = "STRING"
	STRING_CLASS_OBJ       Type = "STRING_CLASS"
	SYMBOL_OBJ             Type = "SYMBOL"
	BOOLEAN_OBJ            Type = "BOOLEAN"
	BOOLEAN_CLASS_OBJ      Type = "BOOLEAN_CLASS"
	NIL_OBJ                Type = "NIL"
	NIL_CLASS_OBJ          Type = "NIL_CLASS"
	ERROR_OBJ              Type = "ERROR"
	EXCEPTION_OBJ          Type = "EXCEPTION"
	EXCEPTION_CLASS_OBJ    Type = "EXCEPTION_CLASS"
	MODULE_OBJ             Type = "MODULE"
	MODULE_CLASS_OBJ       Type = "MODULE_CLASS"
	BUILTIN_OBJ            Type = "BUILTIN"
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
	Methods() map[string]RubyMethod
	SuperClass() RubyClass
}

// RubyClassObject represents a class object in Ruby
type RubyClassObject interface {
	RubyObject
	RubyClass
}

// A BuiltinFunction represents a function
type BuiltinFunction func(args ...RubyObject) RubyObject

// Builtin represents a builtin within the interpreter. It holds a function
// which can be called directly. It is no real Ruby object.
//
// Ruby does not have any builtin functions as everything is bound to an object.
//
// This object will go away soon. Don't depend on it.
type Builtin struct {
	Fn BuiltinFunction
}

// Type returns BUILTIN_OBJ
func (b *Builtin) Type() Type { return BUILTIN_OBJ }

// Inspect returns 'buitin function'
func (b *Builtin) Inspect() string { return "builtin function" }

// Class returns nil
func (b *Builtin) Class() RubyClass { return nil }

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

// Error represents the old way of marking something within the evaluation
// phase as error. It is obsolete and should be replaced by a real exception.
//
// This object will go away soon. Don't depend on it.
type Error struct {
	Message string
}

// Type returns ERROR_OBJ
func (e *Error) Type() Type { return ERROR_OBJ }

// Inspect returns the wrapped message preprended with 'ERROR: '
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Class returns nil
func (e *Error) Class() RubyClass { return nil }

// A Function represents a user defined function. It is no real Ruby object.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        Environment
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
