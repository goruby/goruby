package object

import (
	"bytes"
	"strings"

	"github.com/goruby/goruby/ast"
)

type ObjectType string

const (
	EIGENCLASS_OBJ         ObjectType = "EIGENCLASS"
	FUNCTION_OBJ           ObjectType = "FUNCTION"
	RETURN_VALUE_OBJ       ObjectType = "RETURN_VALUE"
	BASIC_OBJECT_OBJ       ObjectType = "BASIC_OBJECT"
	BASIC_OBJECT_CLASS_OBJ ObjectType = "BASIC_OBJECT_CLASS"
	OBJECT_OBJ             ObjectType = "OBJECT"
	OBJECT_CLASS_OBJ       ObjectType = "OBJECT_CLASS"
	CLASS_OBJ              ObjectType = "CLASS"
	CLASS_CLASS_OBJ        ObjectType = "CLASS_CLASS"
	ARRAY_OBJ              ObjectType = "ARRAY"
	ARRAY_CLASS_OBJ        ObjectType = "ARRAY_CLASS"
	INTEGER_OBJ            ObjectType = "INTEGER"
	INTEGER_CLASS_OBJ      ObjectType = "INTEGER_CLASS"
	STRING_OBJ             ObjectType = "STRING"
	STRING_CLASS_OBJ       ObjectType = "STRING_CLASS"
	SYMBOL_OBJ             ObjectType = "SYMBOL"
	BOOLEAN_OBJ            ObjectType = "BOOLEAN"
	BOOLEAN_CLASS_OBJ      ObjectType = "BOOLEAN_CLASS"
	NIL_OBJ                ObjectType = "NIL"
	NIL_CLASS_OBJ          ObjectType = "NIL_CLASS"
	ERROR_OBJ              ObjectType = "ERROR"
	EXCEPTION_OBJ          ObjectType = "EXCEPTION"
	EXCEPTION_CLASS_OBJ    ObjectType = "EXCEPTION_CLASS"
	MODULE_OBJ             ObjectType = "MODULE"
	MODULE_CLASS_OBJ       ObjectType = "MODULE_CLASS"
	BUILTIN_OBJ            ObjectType = "BUILTIN"
)

type inspectable interface {
	Inspect() string
}

type RubyObject interface {
	inspectable
	Type() ObjectType
	Class() RubyClass
}

type RubyClass interface {
	Methods() map[string]RubyMethod
	SuperClass() RubyClass
}

type RubyClassObject interface {
	RubyObject
	RubyClass
}

type BuiltinFunction func(args ...RubyObject) RubyObject

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) Class() RubyClass { return nil }

type ReturnValue struct {
	Value RubyObject
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Class() RubyClass { return rv.Value.Class() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) Class() RubyClass { return nil }

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
func (f *Function) Class() RubyClass { return nil }
