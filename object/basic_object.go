package object

var BASIC_OBJECT_CLASS RubyClassObject = NewClass("BasicObject", nil, basicObjectMethods, basicObjectClassMethods)

func init() {
	classes.Set("BasicObject", BASIC_OBJECT_CLASS)
}

type BasicObject struct{}

func (b *BasicObject) Inspect() string  { return "" }
func (b *BasicObject) Type() ObjectType { return BASIC_OBJECT_OBJ }
func (b *BasicObject) Class() RubyClass { return BASIC_OBJECT_CLASS }

var basicObjectClassMethods = map[string]RubyMethod{
	"new": publicMethod(func(context RubyObject, args ...RubyObject) RubyObject {
		return &BasicObject{}
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
