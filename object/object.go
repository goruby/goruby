package object

var (
	OBJECT_EIGENCLASS RubyClass       = newEigenClass(OBJECT_CLASS, objectClassMethods)
	OBJECT_CLASS      RubyClassObject = &ObjectClass{}
)

type ObjectClass struct{}

func (o *ObjectClass) Inspect() string            { return "Object" }
func (o *ObjectClass) Type() ObjectType           { return OBJECT_CLASS_OBJ }
func (o *ObjectClass) Methods() map[string]method { return objectMethods }
func (o *ObjectClass) Class() RubyClass           { return OBJECT_EIGENCLASS }
func (o *ObjectClass) SuperClass() RubyClass      { return BASIC_OBJECT_CLASS }

type Object struct{}

func (o *Object) Inspect() string  { return "" }
func (o *Object) Type() ObjectType { return OBJECT_OBJ }
func (o *Object) Class() RubyClass { return OBJECT_CLASS }

var objectClassMethods = map[string]method{}

var objectMethods = map[string]method{
	"nil?":    withArity(0, objectIsNil),
	"methods": withArity(0, objMethods),
}

func objMethods(context RubyObject, args ...RubyObject) RubyObject {
	var methodSymbols []RubyObject
	class := context.Class()
	for class != nil {
		methods := class.Methods()
		for meth, _ := range methods {
			methodSymbols = append(methodSymbols, &Symbol{meth})
		}
		class = class.SuperClass()
	}

	return &Array{Elements: methodSymbols}
}

func objectIsNil(context RubyObject, args ...RubyObject) RubyObject {
	return FALSE
}
