package object

import (
	"hash/fnv"
	"strings"
)

var arrayClass RubyClassObject = newClass(
	"Array",
	objectClass,
	arrayMethods,
	arrayClassMethods,
	func(c RubyClassObject, args ...RubyObject) (RubyObject, error) { return NewArray(args...), nil },
)

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
func (a *Array) hashKey() hashKey {
	h := fnv.New64a()
	for _, e := range a.Elements {
		h.Write(hash(e).bytes())
	}
	return hashKey{Type: a.Type(), Value: h.Sum64()}
}

var arrayClassMethods = map[string]RubyMethod{}

var arrayMethods = map[string]RubyMethod{
	"push":    publicMethod(arrayPush),
	"unshift": publicMethod(arrayUnshift),
}

func arrayPush(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	array.Elements = append(array.Elements, args...)
	return array, nil
}

func arrayUnshift(context CallContext, args ...RubyObject) (RubyObject, error) {
	array, _ := context.Receiver().(*Array)
	array.Elements = append(args, array.Elements...)
	return array, nil
}
