package object

var objectClass = newMixin(newClass(
	"Object",
	basicObjectClass,
	objectMethods,
	objectClassMethods,
	func(RubyClassObject, ...RubyObject) (RubyObject, error) {
		return &Object{}, nil
	},
), kernelModule)

func init() {
	classes.Set("Object", objectClass)
}

// Object represents an Object in Ruby
type Object struct {
	_ int // for uniqueness
}

// Inspect return ""
func (o *Object) Inspect() string { return "" }

// Type returns OBJECT_OBJ
func (o *Object) Type() Type { return OBJECT_OBJ }

// Class returns objectClass
func (o *Object) Class() RubyClass { return objectClass }

var objectClassMethods = map[string]RubyMethod{}

var objectMethods = map[string]RubyMethod{}
