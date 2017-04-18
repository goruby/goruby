package object

import "fmt"

var classClass RubyClassObject = &class{name: "Class", superClass: moduleClass, instanceMethods: NewMethodSet(classMethods)}

func init() {
	classClass.(*class).class = classClass
	classes.Set("Class", classClass)
}

// NewClass returns a new Ruby Class
func NewClass(name string, superClass RubyClass, instanceMethods, classMethods map[string]RubyMethod) RubyClassObject {
	if instanceMethods == nil {
		instanceMethods = map[string]RubyMethod{}
	}
	if classMethods == nil {
		classMethods = map[string]RubyMethod{}
	}
	return &class{
		name:            name,
		superClass:      superClass,
		instanceMethods: NewMethodSet(instanceMethods),
		class:           newEigenclass(classClass, classMethods),
	}
}

// class represents a Ruby Class object
type class struct {
	name            string
	superClass      RubyClass
	class           RubyClass
	instanceMethods SettableMethodSet
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
func (c *class) Methods() MethodSet {
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
	if mixin, ok := superclass.(*mixin); ok {
		return mixin.RubyClassObject, nil
	}
	return superclass.(RubyObject), nil
}
