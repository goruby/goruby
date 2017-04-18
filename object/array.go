package object

import "strings"

var arrayClass RubyClassObject = NewClass("Array", objectClass, arrayMethods, arrayClassMethods)

func init() {
	classes.Set("Array", arrayClass)
}

// NewArray returns a new array populated with elements.
func NewArray(elements ...RubyObject) *Array {
	arr := &Array{Elements: make([]RubyObject, len(elements))}
	for i, elem := range elements {
		arr.Elements[i] = elem
	}
	return arr
}

// An Array represents a Ruby Array
type Array struct {
	Elements []RubyObject
}

// Type returns the ObjectType of the array
func (a *Array) Type() Type { return ARRAY_OBJ }

// Inspect returns all elements within the array, divided by comma and
// surrounded by brackets
func (a *Array) Inspect() string {
	elems := make([]string, len(a.Elements))
	for i, elem := range a.Elements {
		elems[i] = elem.Inspect()
	}
	return "[" + strings.Join(elems, ", ") + "]"
}

// Class returns the class of the Array
func (a *Array) Class() RubyClass { return arrayClass }

var arrayClassMethods = map[string]RubyMethod{}

var arrayMethods = map[string]RubyMethod{}
