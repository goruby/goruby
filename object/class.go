package object

import "fmt"

var (
	CLASS_EIGENCLASS RubyClass       = newEigenclass(CLASS_CLASS, classClassMethods)
	CLASS_CLASS      RubyClassObject = &ClassClass{}
)

type ClassClass struct{}

func (c *ClassClass) Inspect() string                { return "Class" }
func (c *ClassClass) Type() ObjectType               { return CLASS_OBJ }
func (c *ClassClass) Class() RubyClass               { return CLASS_EIGENCLASS }
func (c *ClassClass) Methods() map[string]RubyMethod { return classMethods }
func (c *ClassClass) SuperClass() RubyClass          { return MODULE_CLASS }

type Class struct{}

func (c *Class) Inspect() string  { return fmt.Sprintf("#<Class:%p>", c) }
func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Class() RubyClass { return CLASS_CLASS }

var classClassMethods = map[string]RubyMethod{}

var classMethods = map[string]RubyMethod{
	"superclass": withArity(0, publicMethod(classSuperclass)),
}

func classSuperclass(context RubyObject, args ...RubyObject) RubyObject {
	class := context.(RubyClass)
	superclass := class.SuperClass()
	if superclass == nil {
		return NIL
	}
	if mixin, ok := superclass.(*methodSet); ok {
		return mixin.RubyClassObject
	}
	return superclass.(RubyObject)
}
