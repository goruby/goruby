package object

var OBJECT_CLASS RubyClassObject = mixin(NewClass("Object", BASIC_OBJECT_CLASS, objectMethods, objectClassMethods), KERNEL_MODULE)

type Object struct{}

func (o *Object) Inspect() string  { return "" }
func (o *Object) Type() ObjectType { return OBJECT_OBJ }
func (o *Object) Class() RubyClass { return OBJECT_CLASS }

var objectClassMethods = map[string]RubyMethod{}

var objectMethods = map[string]RubyMethod{
	"nil?":    withArity(0, publicMethod(objectIsNil)),
	"methods": withArity(0, publicMethod(objMethods)),
	"class":   withArity(0, publicMethod(objectClass)),
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

func objectClass(context RubyObject, args ...RubyObject) RubyObject {
	class := context.Class()
	if eigenClass, ok := class.(*eigenclass); ok {
		class = eigenClass.Class()
	}
	classObj := class.(RubyClassObject)
	return classObj
}
