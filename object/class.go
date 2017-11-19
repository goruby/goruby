package object

import (
	"fmt"
	"hash/fnv"
)

var classClass RubyClassObject = &class{
	name:            "Class",
	superClass:      moduleClass,
	instanceMethods: NewMethodSet(classInstanceMethods),
	builder:         defaultBuilder,
}

var notInstantiatable = func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
	return nil, NewNoMethodError(c, "new")
}
var defaultBuilder = func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
	return &classInstance{class: c}, nil
}

func init() {
	classClass.(*class).class = newEigenclass(moduleClass.Class(), classClassMethods)
	classes.Set("Class", classClass)
}

// NewClass returns a new Ruby Class
func NewClass(name string, superClass RubyClass, env Environment) RubyClassObject {
	instanceMethods := map[string]RubyMethod{}
	classMethods := map[string]RubyMethod{}
	return newClassWithEnv(name, superClass, instanceMethods, classMethods, defaultBuilder, env)
}

// newClass returns a new Ruby Class
func newClass(
	name string,
	superClass RubyClass,
	instanceMethods,
	classMethods map[string]RubyMethod,
	builder func(RubyClassObject, ...RubyObject) (RubyObject, error),
) *class {
	return newClassWithEnv(
		name,
		superClass,
		instanceMethods,
		classMethods,
		builder,
		nil,
	)
}

// newClass returns a new Ruby Class
func newClassWithEnv(
	name string,
	superClass RubyClass,
	instanceMethods,
	classMethods map[string]RubyMethod,
	builder func(RubyClassObject, ...RubyObject) (RubyObject, error),
	env Environment,
) *class {
	var superclassClass RubyClass = classClass
	if superClass != nil {
		superclassClass = superClass.(RubyClassObject).Class()
	}
	return &class{
		name:            name,
		superClass:      superClass,
		instanceMethods: NewMethodSet(instanceMethods),
		class:           newEigenclass(superclassClass, classMethods),
		builder:         builder,
		Environment:     NewEnclosedEnvironment(env),
	}
}

// class represents a Ruby Class object
type class struct {
	name            string
	superClass      RubyClass
	class           RubyClass
	instanceMethods SettableMethodSet
	builder         func(RubyClassObject, ...RubyObject) (RubyObject, error)
	Environment
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
func (c *class) hashKey() hashKey {
	h := fnv.New64a()
	h.Write([]byte(c.name))
	return hashKey{Type: c.Type(), Value: h.Sum64()}
}
func (c *class) addMethod(name string, method RubyMethod) {
	c.instanceMethods.Set(name, method)
}
func (c *class) New(args ...RubyObject) (RubyObject, error) {
	return c.builder(c)
}
func (c *class) Name() string { return c.name }

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
	instance, err := classObject.New(args...)
	if err != nil {
		return nil, err
	}
	self := &Self{RubyObject: instance, Name: classObject.Inspect()}
	callContext := &callContext{
		receiver: self,
		env:      context.Env(),
		eval:     context.Eval,
	}
	_, err = Send(callContext, "initialize", args...)
	if err != nil {
		return nil, err
	}

	return self.RubyObject, nil
}

func classInitialize(context CallContext, args ...RubyObject) (RubyObject, error) {
	return context.Receiver(), nil
}
