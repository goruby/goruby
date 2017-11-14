package object

import (
	"fmt"
)

var basicObjectClass RubyClassObject = newClass(
	"BasicObject",
	nil,
	basicObjectMethods,
	basicObjectClassMethods,
	func(RubyClassObject, ...RubyObject) (RubyObject, error) { return &basicObject{}, nil },
)

func init() {
	classes.Set("BasicObject", basicObjectClass)
}

// basicObject represents a basicObject object in Ruby
type basicObject struct {
	_ int // for uniqueness
}

// Inspect returns empty string. BasicObjects do not have an `inspect` method.
func (b *basicObject) Inspect() string {
	fmt.Println("(Object doesn't support #inspect)")
	return ""
}

// Type returns the ObjectType of the array
func (b *basicObject) Type() Type { return BASIC_OBJECT_OBJ }

// Class returns the class of BasicObject
func (b *basicObject) Class() RubyClass { return basicObjectClass }

var basicObjectClassMethods = map[string]RubyMethod{}

var basicObjectMethods = map[string]RubyMethod{
	"initialize":     privateMethod(basicObjectInitialize),
	"method_missing": privateMethod(basicObjectMethodMissing),
}

func basicObjectMethodMissing(context CallContext, args ...RubyObject) (RubyObject, error) {
	if len(args) < 1 {
		return nil, NewWrongNumberOfArgumentsError(1, 0)
	}
	method, ok := args[0].(*Symbol)
	if !ok {
		return nil, NewImplicitConversionTypeError(method, args[0])
	}
	return nil, NewNoMethodError(context.Receiver(), method.Value)
}

func basicObjectInitialize(context CallContext, args ...RubyObject) (RubyObject, error) {
	return context.Receiver(), nil
}
