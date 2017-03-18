package object

var OBJECT_CLASS RubyClassObject = mixin(NewClass("Object", BASIC_OBJECT_CLASS, objectMethods, objectClassMethods), KERNEL_MODULE)

type Object struct{}

func (o *Object) Inspect() string  { return "" }
func (o *Object) Type() ObjectType { return OBJECT_OBJ }
func (o *Object) Class() RubyClass { return OBJECT_CLASS }

var objectClassMethods = map[string]RubyMethod{}

var objectMethods = map[string]RubyMethod{}
