package object

import (
	"bytes"
	"strings"

	"github.com/goruby/goruby/ast"
)

type ObjectType string

const (
	EIGENCLASS_OBJ         = "EIGENCLASS"
	FUNCTION_OBJ           = "FUNCTION"
	RETURN_VALUE_OBJ       = "RETURN_VALUE"
	BASIC_OBJECT_OBJ       = "BASIC_OBJECT"
	BASIC_OBJECT_CLASS_OBJ = "BASIC_OBJECT_CLASS"
	OBJECT_OBJ             = "OBJECT"
	OBJECT_CLASS_OBJ       = "OBJECT_CLASS"
	ARRAY_OBJ              = "ARRAY"
	ARRAY_CLASS_OBJ        = "ARRAY_CLASS"
	INTEGER_OBJ            = "INTEGER"
	INTEGER_CLASS_OBJ      = "INTEGER_CLASS"
	STRING_OBJ             = "STRING"
	STRING_CLASS_OBJ       = "STRING_CLASS"
	SYMBOL_OBJ             = "SYMBOL"
	BOOLEAN_OBJ            = "BOOLEAN"
	BOOLEAN_CLASS_OBJ      = "BOOLEAN_CLASS"
	NIL_OBJ                = "NIL"
	ERROR_OBJ              = "ERROR"
	EXCEPTION_OBJ          = "EXCEPTION"
	EXCEPTION_CLASS_OBJ    = "EXCEPTION_CLASS"
	MODULE_OBJ             = "MODULE"
	MODULE_CLASS_OBJ       = "MODULE_CLASS"
	BUILTIN_OBJ            = "BUILTIN"
)

type RubyObject interface {
	Type() ObjectType
	Inspect() string
	Class() RubyClass
}

type RubyClass interface {
	RubyObject
	Methods() map[string]method
	SuperClass() RubyClass
}

type BuiltinFunction func(args ...RubyObject) RubyObject

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType           { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string            { return "builtin function" }
func (b *Builtin) Methods() map[string]method { return nil }
func (b *Builtin) Class() RubyClass           { return nil }

type ReturnValue struct {
	Value RubyObject
}

func (rv *ReturnValue) Type() ObjectType           { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string            { return rv.Value.Inspect() }
func (rv *ReturnValue) Methods() map[string]method { return nil }
func (rv *ReturnValue) Class() RubyClass           { return rv.Value.Class() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType           { return ERROR_OBJ }
func (e *Error) Inspect() string            { return "ERROR: " + e.Message }
func (e *Error) Methods() map[string]method { return nil }
func (e *Error) Class() RubyClass           { return nil }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
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
func (f *Function) Methods() map[string]method { return nil }
func (f *Function) Class() RubyClass           { return nil }
