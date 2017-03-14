package object

var BASIC_OBJECT *BasicObjectClass = newBasicObjectClass()

func newBasicObjectClass() *BasicObjectClass {
	cl := &BasicObjectClass{}
	cl.methods = &methodSet{cl, basicObjectClassMethods}
	cl.instanceMethods = &methodSet{methods: basicObjectMethods}
	return cl
}

type BasicObjectClass struct {
	methods         *methodSet
	instanceMethods *methodSet
}

func (i *BasicObjectClass) Inspect() string  { return "BasicObject" }
func (i *BasicObjectClass) Type() ObjectType { return BASIC_OBJECT_CLASS_OBJ }
func (i *BasicObjectClass) Send(name string, args ...RubyObject) RubyObject {
	return i.methods.Call(name, args...)
}

func NewBasicObject() *BasicObject {
	i := &BasicObject{}
	i.methods = BASIC_OBJECT.instanceMethods
	i.methods.context = i
	return i
}

type BasicObject struct {
	methods *methodSet
}

func (b *BasicObject) Inspect() string  { return "" }
func (b *BasicObject) Type() ObjectType { return BASIC_OBJECT_OBJ }
func (b *BasicObject) Send(name string, args ...RubyObject) RubyObject {
	return b.methods.Call(name, args...)
}

var basicObjectClassMethods = map[string]method{
	"new": func(context RubyObject, args ...RubyObject) RubyObject {
		i := &BasicObject{}
		i.methods = &methodSet{methods: basicObjectMethods}
		return i
	},
}

var basicObjectMethods = map[string]method{}
