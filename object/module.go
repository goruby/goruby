package object

var MODULE_CLASS RubyClassObject = &Class{name: "Module", instanceMethods: moduleMethods}

func init() {
	MODULE_CLASS.(*Class).superClass = OBJECT_CLASS
}

func newModule(name string, methods map[string]RubyMethod) *Module {
	return &Module{name, newEigenclass(MODULE_CLASS, methods)}
}

type Module struct {
	name  string
	class RubyClass
}

func (m *Module) Inspect() string  { return m.name }
func (m *Module) Type() ObjectType { return MODULE_OBJ }
func (m *Module) Class() RubyClass {
	if m.class != nil {
		return m.class
	}
	return MODULE_CLASS
}

var moduleMethods = map[string]RubyMethod{
	"ancestors": withArity(0, publicMethod(moduleAncestors)),
}

func moduleAncestors(context RubyObject, args ...RubyObject) RubyObject {
	class := context.(RubyClassObject)
	var ancestors []RubyObject
	ancestors = append(ancestors, &String{class.Inspect()})

	if mixin, ok := class.(*methodSet); ok {
		for _, m := range mixin.modules {
			ancestors = append(ancestors, &String{m.name})
		}
	}
	superClass := class.SuperClass()
	if superClass != nil {
		superAncestors := moduleAncestors(superClass.(RubyObject))
		ancestors = append(ancestors, superAncestors.(*Array).Elements...)
	}
	return &Array{ancestors}
}

func moduleIncludedModules(context RubyObject, args ...RubyObject) RubyObject {
	class := context.(RubyClassObject)
	var includedModules []RubyObject

	if mixin, ok := class.(*methodSet); ok {
		for _, m := range mixin.modules {
			includedModules = append(includedModules, &String{m.name})
		}
	}

	superClass := class.SuperClass()
	if superClass != nil {
		superIncludedModules := moduleIncludedModules(superClass.(RubyObject))
		includedModules = append(includedModules, superIncludedModules.(*Array).Elements...)
	}

	return &Array{includedModules}
}
