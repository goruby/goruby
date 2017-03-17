package object

import "strings"

var (
	ARRAY_EIGENCLASS RubyClass       = newEigenclass(CLASS_CLASS, arrayClassMethods)
	ARRAY_CLASS      RubyClassObject = &ArrayClass{}
)

type ArrayClass struct{}

func (a *ArrayClass) Type() ObjectType           { return ARRAY_OBJ }
func (a *ArrayClass) Inspect() string            { return "Array" }
func (a *ArrayClass) Class() RubyClass           { return ARRAY_EIGENCLASS }
func (a *ArrayClass) Methods() map[string]method { return arrayMethods }
func (a *ArrayClass) SuperClass() RubyClass      { return OBJECT_CLASS }

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

var arrayClassMethods = map[string]method{}

var arrayMethods = map[string]method{}
