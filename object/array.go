package object

import "strings"

var (
	ARRAY_EIGENCLASS RubyClass = nil
	ARRAY_CLASS      RubyClass = nil
)

type ArrayEigenClass struct {
	Elements []RubyObject
}

func (a *ArrayEigenClass) Type() ObjectType      { return ARRAY_OBJ }
func (a *ArrayEigenClass) Inspect() string       { return "(Array)" }
func (a *ArrayEigenClass) Class() RubyClass      { return BASIC_OBJECT_CLASS }
func (a *ArrayEigenClass) SuperClass() RubyClass { return BASIC_OBJECT_CLASS }

type ArrayClass struct {
	Elements []RubyObject
}

func (a *ArrayClass) Type() ObjectType      { return ARRAY_OBJ }
func (a *ArrayClass) Inspect() string       { return "Array" }
func (a *ArrayClass) Class() RubyClass      { return ARRAY_EIGENCLASS }
func (a *ArrayClass) SuperClass() RubyClass { return OBJECT_CLASS }

type Array struct {
	Elements []RubyObject
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	elems := make([]string, len(a.Elements))
	for i, elem := range a.Elements {
		elems[i] = elem.Inspect()
	}
	return "[" + strings.Join(elems, ", ") + "]"
}
func (a *Array) Class() RubyClass { return ARRAY_CLASS }
