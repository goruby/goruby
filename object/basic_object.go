package object

var basicObjectClass RubyClassObject = newClass("BasicObject", nil, basicObjectMethods, basicObjectClassMethods)

func init() {
	classes.Set("BasicObject", basicObjectClass)
}

// basicObject represents a basicObject object in Ruby
type basicObject struct{}

// Inspect returns empty string. BasicObjects do not have an `inspect` method.
func (b *basicObject) Inspect() string { return "" }

// Type returns the ObjectType of the array
func (b *basicObject) Type() Type { return BASIC_OBJECT_OBJ }

// Class returns the class of BasicObject
func (b *basicObject) Class() RubyClass { return basicObjectClass }

var basicObjectClassMethods = map[string]RubyMethod{
	"new": publicMethod(func(context RubyObject, args ...RubyObject) RubyObject {
		return &basicObject{}
	}),
}

var basicObjectMethods = map[string]RubyMethod{
	"method_missing": privateMethod(basicObjectMethodMissing),
}

func basicObjectMethodMissing(context RubyObject, args ...RubyObject) RubyObject {
	if len(args) < 1 {
		// TODO: can we protect against this
		panic("wrong number of call arguments for method_missing")
	}
	method, ok := args[0].(*Symbol)
	if !ok {
		// TODO: can we protect against this?
		panic("wrong call argument for method_missing")
	}
	return NewNoMethodError(context, method.Value)
}
