package object

var (
	OBJECT_EIGENCLASS RubyClass = &ObjectEigenclass{}
	OBJECT_CLASS      RubyClass = &ObjectClass{}
)

type ObjectEigenclass struct{}

func (o *ObjectEigenclass) Inspect() string  { return "" }
func (o *ObjectEigenclass) Type() ObjectType { return EIGENCLASS_OBJ }
func (o *ObjectEigenclass) Methods() map[string]method {
	return objectClassMethods
}
func (o *ObjectEigenclass) Class() RubyClass      { return OBJECT_CLASS }
func (o *ObjectEigenclass) SuperClass() RubyClass { return OBJECT_CLASS }

type ObjectClass struct{}

func (o *ObjectClass) Inspect() string  { return "Object" }
func (o *ObjectClass) Type() ObjectType { return OBJECT_CLASS_OBJ }
func (o *ObjectClass) Methods() map[string]method {
	return objectMethods
}
func (o *ObjectClass) Class() RubyClass      { return OBJECT_EIGENCLASS }
func (o *ObjectClass) SuperClass() RubyClass { return BASIC_OBJECT_CLASS }

type Object struct{}

func (o *Object) Inspect() string  { return "" }
func (o *Object) Type() ObjectType { return OBJECT_OBJ }
func (o *Object) Methods() map[string]method {
	return nil
}
func (o *Object) Class() RubyClass { return OBJECT_CLASS }

var objectClassMethods = map[string]method{}

var objectMethods = map[string]method{
	"nil?": objectIsNil,
}

func objMethods(context RubyObject, args ...RubyObject) RubyObject {
	return NIL
}

func objectIsNil(context RubyObject, args ...RubyObject) RubyObject {
	return FALSE
}
