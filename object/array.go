package object

import "strings"

var ARRAY_CLASS RubyClassObject = NewClass("Array", OBJECT_CLASS, arrayMethods, arrayClassMethods)

func init() {
	classes.Set("Array", ARRAY_CLASS)
}

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

var arrayClassMethods = map[string]RubyMethod{}

var arrayMethods = map[string]RubyMethod{}
