package object

import (
	"fmt"
)

var classClass RubyClassObject = &class{name: "Class", superClass: moduleClass, instanceMethods: NewMethodSet(classInstanceMethods)}

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
	return c.name
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
func (c *class) addMethod(name string, method RubyMethod) {
	c.instanceMethods.Set(name, method)
}

var classClassMethods = map[string]RubyMethod{}

var classInstanceMethods = map[string]RubyMethod{
	"superclass": withArity(0, publicMethod(classSuperclass)),
	"new":        publicMethod(classNew),
	"initialize": privateMethod(classInitialize),
}

type classInstance struct {
	class RubyClassObject
}

func (o *classInstance) Inspect() string  { return fmt.Sprintf("#<%s:%p>", o.class.Inspect(), o) }
func (o *classInstance) Class() RubyClass { return o.class }
func (o *classInstance) Type() Type       { return CLASS_INSTANCE_OBJ }

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

func classNew(context CallContext, args ...RubyObject) (RubyObject, error) {
	classObject := context.Receiver().(RubyClassObject)
	var instance = classInstance{class: classObject}
	Send(
		&callContext{
			receiver: &Self{&instance, classObject.Inspect()},
			env:      context.Env(),
			eval:     context.Eval,
		},
		"initialize",
		args...,
	)
	return &instance, nil
}

func classInitialize(context CallContext, args ...RubyObject) (RubyObject, error) {
	return context.Receiver(), nil
}
