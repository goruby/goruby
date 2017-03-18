package object

var (
	MODULE_EIGENCLASS RubyClass       = newEigenclass(CLASS_CLASS, moduleMethods)
	MODULE_CLASS      RubyClassObject = &ModuleClass{}
)

type ModuleClass struct{}

func (m *ModuleClass) Inspect() string                { return "Module" }
func (m *ModuleClass) Type() ObjectType               { return MODULE_CLASS_OBJ }
func (m *ModuleClass) Class() RubyClass               { return MODULE_EIGENCLASS }
func (m *ModuleClass) Methods() map[string]RubyMethod { return moduleMethods }
func (m *ModuleClass) SuperClass() RubyClass          { return OBJECT_CLASS }

func newModule(name string, class RubyClass) *Module {
	return &Module{name, class}
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
