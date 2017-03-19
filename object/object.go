package object

var objectClass RubyClassObject = mixin(newClass("Object", basicObjectClass, objectMethods, objectClassMethods), kernelModule)

func init() {
	classes.Set("Object", objectClass)
}

type Object struct{}

func (o *Object) Inspect() string  { return "" }
func (o *Object) Type() Type       { return OBJECT_OBJ }
func (o *Object) Class() RubyClass { return objectClass }

var objectClassMethods = map[string]RubyMethod{}

var objectMethods = map[string]RubyMethod{}
