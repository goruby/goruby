package object

import "fmt"

var CLASS_CLASS RubyClassObject = &Class{name: "Class", superClass: MODULE_CLASS, instanceMethods: classMethods}

func init() {
	CLASS_CLASS.(*Class).class = CLASS_CLASS
}

func NewClass(name string, superClass RubyClass, instanceMethods, classMethods map[string]RubyMethod) RubyClassObject {
	return &Class{name: name, superClass: superClass, instanceMethods: instanceMethods, class: newEigenclass(CLASS_CLASS, classMethods)}
}

type Class struct {
	name            string
	superClass      RubyClass
	class           RubyClass
	instanceMethods map[string]RubyMethod
}

func (c *Class) Inspect() string {
	if c.name != "" {
		return c.name
	}
	return fmt.Sprintf("#<Class:%p>", c)
}
func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Class() RubyClass {
	if c.class != nil {
		return c.class
	}
	return CLASS_CLASS
}
func (c *Class) SuperClass() RubyClass {
	return c.superClass
}
func (c *Class) Methods() map[string]RubyMethod {
	return c.instanceMethods
}

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
