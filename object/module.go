package object

var moduleClass RubyClassObject = &class{name: "Module", instanceMethods: NewMethodSet(moduleMethods)}

func init() {
	moduleClass.(*class).superClass = objectClass
	classes.Set("Module", moduleClass)
}

// NewModule returns a new module with the given name and adds methods to its method set.
func NewModule(name string, methods map[string]RubyMethod) *Module {
	if methods == nil {
		methods = make(map[string]RubyMethod)
	}
	return &Module{
		name:  name,
		class: newEigenclass(moduleClass, methods),
	}
}

// Module represents a module in Ruby
type Module struct {
	name  string
	class *eigenclass
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
