package object

var moduleClass RubyClassObject = &class{name: "Module", instanceMethods: moduleMethods}

func init() {
	moduleClass.(*class).superClass = objectClass
	classes.Set("Module", moduleClass)
}

func newModule(name string, methods map[string]RubyMethod) *Module {
	return &Module{name, newEigenclass(moduleClass, methods)}
}

// Module represents a module in Ruby
type Module struct {
	name  string
	class RubyClass
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

var moduleMethods = map[string]RubyMethod{
	"ancestors": withArity(0, publicMethod(moduleAncestors)),
}

func moduleAncestors(context CallContext, args ...RubyObject) (RubyObject, error) {
	class := context.Receiver().(RubyClassObject)
	var ancestors []RubyObject
	ancestors = append(ancestors, &String{class.Inspect()})

	if mixin, ok := class.(*methodSet); ok {
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

	if mixin, ok := class.(*methodSet); ok {
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
