package object

import "fmt"

var classClass RubyClassObject = &class{name: "Class", superClass: moduleClass, instanceMethods: classMethods}

func init() {
	classClass.(*class).class = classClass
	classes.Set("Class", classClass)
}

// newClass returns a new Ruby Class
func newClass(name string, superClass RubyClass, instanceMethods, classMethods map[string]RubyMethod) *class {
	return &class{name: name, superClass: superClass, instanceMethods: instanceMethods, class: newEigenclass(classClass, classMethods)}
}

// class represents a Ruby Class object
type class struct {
	name            string
	superClass      RubyClass
	class           RubyClass
	instanceMethods map[string]RubyMethod
}

func (c *class) Inspect() string {
	if c.name != "" {
		return c.name
	}
	return fmt.Sprintf("#<Class:%p>", c)
}
func (c *class) Type() Type { return CLASS_OBJ }
func (c *class) Class() RubyClass {
	if c.class != nil {
		return c.class
	}
	return classClass
}
func (c *class) SuperClass() RubyClass {
	return c.superClass
}
func (c *class) Methods() map[string]RubyMethod {
	return c.instanceMethods
}

var classClassMethods = map[string]RubyMethod{}

var classMethods = map[string]RubyMethod{
	"superclass": withArity(0, publicMethod(classSuperclass)),
}

func classSuperclass(context CallContext, args ...RubyObject) (RubyObject, error) {
	class := context.Receiver().(RubyClass)
	superclass := class.SuperClass()
	if superclass == nil {
		return NIL, nil
	}
	if mixin, ok := superclass.(*methodSet); ok {
		return mixin.RubyClassObject, nil
	}
	return superclass.(RubyObject), nil
}
