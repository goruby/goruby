package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"unicode"
)

var moduleClass RubyClassObject = &class{
	name:            "Module",
	instanceMethods: NewMethodSet(moduleMethods),
	builder: func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
		return &Module{}, nil
	},
}

func init() {
	moduleClass.(*class).superClass = objectClass
	classes.Set("Module", moduleClass)
}

// NewModule returns a new module with the given name and adds methods to its method set.
func NewModule(name string, outerEnv Environment) *Module {
	methods := make(map[string]RubyMethod)
	return newModule(name, methods, outerEnv)
}

// newModule returns a new module with the given name and adds methods to its method set.
func newModule(name string, methods map[string]RubyMethod, outerEnv Environment) *Module {
	if methods == nil {
		methods = make(map[string]RubyMethod)
	}
	return &Module{
		name:        name,
		class:       newEigenclass(moduleClass, methods),
		Environment: NewEnclosedEnvironment(outerEnv),
	}
}

// Module represents a module in Ruby
type Module struct {
	name  string
	class *eigenclass
	Environment
}

// Inspect returns the name of the module
func (m *Module) Inspect() string { return m.name }

// Type returns MODULE_OBJ
func (m *Module) Type() Type { return MODULE_OBJ }

// Class returns the set class or moduleClass
func (m *Module) Class() RubyClass {
	if m.class != nil {
		return m.class
	}
	return moduleClass
}
func (m *Module) hashKey() hashKey {
	h := fnv.New64a()
	h.Write([]byte(m.name))
	return hashKey{Type: m.Type(), Value: h.Sum64()}
}

func (m *Module) addMethod(name string, method RubyMethod) {
	m.class.addMethod(name, method)
}

var moduleMethods = map[string]RubyMethod{
	"ancestors":                  withArity(0, publicMethod(moduleAncestors)),
	"included_modules":           withArity(0, publicMethod(moduleIncludedModules)),
	"instance_methods":           publicMethod(modulePublicInstanceMethods),
	"public_instance_methods":    publicMethod(modulePublicInstanceMethods),
	"protected_instance_methods": publicMethod(moduleProtectedInstanceMethods),
	"private_instance_methods":   publicMethod(modulePrivateInstanceMethods),
	"include":                    publicMethod(moduleInclude),
	"append_features":            withArity(1, privateMethod(moduleAppendFeatures)),
	"to_s":                       withArity(0, publicMethod(moduleToS)),
	"inspect":                    withArity(0, publicMethod(moduleToS)),
}

func moduleToS(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	if self, ok := receiver.(*Self); ok {
		receiver = self.RubyObject
	}
	if class, ok := receiver.(RubyClass); ok && class.Name() != "" {
		return &String{Value: class.Name()}, nil
	}
	val := fmt.Sprintf("#<%s:%p>", receiver.Class().Name(), receiver)
	return &String{Value: val}, nil
}

func moduleAncestors(context CallContext, args ...RubyObject) (RubyObject, error) {
	class := context.Receiver().(RubyClassObject)
	var ancestors []RubyObject
	ancestors = append(ancestors, &String{class.Inspect()})

	if mixin, ok := class.(*mixin); ok {
		for _, m := range mixin.modules {
			ancestors = append(ancestors, &String{m.name})
		}
	}
	superClass := class.SuperClass()
	if superClass != nil {
		callContext := NewCallContext(context.Env(), superClass.(RubyObject))
		superAncestors, err := moduleAncestors(callContext)
		if err != nil {
			return nil, err
		}
		ancestors = append(ancestors, superAncestors.(*Array).Elements...)
	}
	return &Array{ancestors}, nil
}

func moduleIncludedModules(context CallContext, args ...RubyObject) (RubyObject, error) {
	class := context.Receiver().(RubyClassObject)
	var includedModules []RubyObject

	if mixin, ok := class.(*mixin); ok {
		for _, m := range mixin.modules {
			includedModules = append(includedModules, &String{m.name})
		}
	}

	superClass := class.SuperClass()
	if superClass != nil {
		callContext := NewCallContext(context.Env(), superClass.(RubyObject))
		superIncludedModules, err := moduleIncludedModules(callContext)
		if err != nil {
			return nil, err
		}
		includedModules = append(includedModules, superIncludedModules.(*Array).Elements...)
	}

	return &Array{includedModules}, nil
}

func modulePublicInstanceMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showSuperClassInstanceMethods := true
	if len(args) == 1 {
		boolean, ok := args[0].(*Boolean)
		if !ok {
			boolean = TRUE.(*Boolean)
		}
		showSuperClassInstanceMethods = boolean.Value
	}
	class := context.Receiver().(RubyClass)

	return getMethods(class, PUBLIC_METHOD, showSuperClassInstanceMethods), nil
}

func moduleProtectedInstanceMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showSuperClassInstanceMethods := true
	if len(args) == 1 {
		boolean, ok := args[0].(*Boolean)
		if !ok {
			boolean = TRUE.(*Boolean)
		}
		showSuperClassInstanceMethods = boolean.Value
	}
	class := context.Receiver().(RubyClass)

	return getMethods(class, PROTECTED_METHOD, showSuperClassInstanceMethods), nil
}

func modulePrivateInstanceMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showSuperClassInstanceMethods := true
	if len(args) == 1 {
		boolean, ok := args[0].(*Boolean)
		if !ok {
			boolean = TRUE.(*Boolean)
		}
		showSuperClassInstanceMethods = boolean.Value
	}
	class := context.Receiver().(RubyClass)

	return getMethods(class, PRIVATE_METHOD, showSuperClassInstanceMethods), nil
}

func moduleInclude(context CallContext, args ...RubyObject) (RubyObject, error) {
	if len(args) == 0 {
		return nil, NewWrongNumberOfArgumentsError(1, 0)
	}
	modules := []RubyObject{}
	for _, a := range args {
		mod, ok := a.(*Module)
		if !ok {
			return nil, NewWrongArgumentTypeError(&Module{}, a)
		}
		modules = append([]RubyObject{mod}, modules...)
	}
	for _, m := range modules {
		_, err := Send(context, "append_features", m)
		if err != nil {
			return nil, err
		}
	}
	return context.Receiver(), nil
}

func moduleAppendFeatures(context CallContext, args ...RubyObject) (RubyObject, error) {
	module, ok := args[0].(*Module)
	if !ok {
		return nil, NewWrongArgumentTypeError(module, args[0])
	}
	self, ok := context.Receiver().(*Self)
	if !ok {
		return nil, NewPrivateNoMethodError(context.Receiver(), "append_features")
	}
	receiverEnv, ok := self.RubyObject.(Environment)
	if !ok {
		return nil, NewNoMethodError(context.Receiver(), "append_features")
	}
	envs := []Environment{receiverEnv}
	env := receiverEnv.Outer()
	for env != nil {
		envs = append([]Environment{env}, envs...)
		env = env.Outer()
	}
	for _, e := range envs {
		for k, v := range e.GetAll() {
			firstChar := bytes.Runes([]byte(k))[0]
			if unicode.IsUpper(firstChar) || firstChar == '@' {
				module.Set(k, v)
			}
		}
	}
	for k, v := range self.RubyObject.Class().Methods().GetAll() {
		module.addMethod(k, v)
	}
	return module, nil
}
