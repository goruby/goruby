package object

var objectClass = newMixin(NewClass("Object", basicObjectClass, objectMethods, objectClassMethods), kernelModule)

func init() {
	classes.Set("Object", objectClass)
}

// Object represents an Object in Ruby
type Object struct{}

// Inspect return ""
func (o *Object) Inspect() string { return "" }

// Type returns OBJECT_OBJ
func (o *Object) Type() Type { return OBJECT_OBJ }

// Class returns objectClass
func (o *Object) Class() RubyClass { return objectClass }

var objectClassMethods = map[string]RubyMethod{}

var objectMethods = map[string]RubyMethod{}
