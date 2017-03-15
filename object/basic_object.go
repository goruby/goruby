package object

var BASIC_OBJECT_EIGENCLASS *BasicObjectEigenclass = &BasicObjectEigenclass{}
var BASIC_OBJECT_CLASS *BasicObjectClass = &BasicObjectClass{}

type BasicObjectEigenclass struct{}

func (b *BasicObjectEigenclass) Inspect() string  { return "BasicObject" }
func (b *BasicObjectEigenclass) Type() ObjectType { return EIGENCLASS_OBJ }
func (b *BasicObjectEigenclass) Methods() map[string]method {
	return basicObjectClassMethods
}
func (b *BasicObjectEigenclass) Class() RubyClass      { return BASIC_OBJECT_CLASS }
func (b *BasicObjectEigenclass) SuperClass() RubyClass { return nil }

type BasicObjectClass struct{}

func (b *BasicObjectClass) Inspect() string  { return "BasicObject" }
func (b *BasicObjectClass) Type() ObjectType { return BASIC_OBJECT_CLASS_OBJ }
func (b *BasicObjectClass) Methods() map[string]method {
	return basicObjectMethods
}
func (b *BasicObjectClass) Class() RubyClass      { return BASIC_OBJECT_EIGENCLASS }
func (b *BasicObjectClass) SuperClass() RubyClass { return nil }

type BasicObject struct{}

func (b *BasicObject) Inspect() string  { return "" }
func (b *BasicObject) Type() ObjectType { return BASIC_OBJECT_OBJ }
func (b *BasicObject) Class() RubyClass { return BASIC_OBJECT_CLASS }

var basicObjectClassMethods = map[string]method{
	"new": func(context RubyObject, args ...RubyObject) RubyObject {
		return &BasicObject{}
	},
}

var basicObjectMethods = map[string]method{
	"method_missing": basicObjectMethodMissing,
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
