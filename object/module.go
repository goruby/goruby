package object

var (
	MODULE_EIGENCLASS RubyClass = newEigenClass(OBJECT_CLASS, moduleMethods)
	MODULE_CLASS      RubyClass = &ModuleClass{}
)

type ModuleClass struct{}

func (m *ModuleClass) Inspect() string            { return "Module" }
func (m *ModuleClass) Type() ObjectType           { return MODULE_CLASS_OBJ }
func (m *ModuleClass) Methods() map[string]method { return moduleMethods }
func (m *ModuleClass) Class() RubyClass           { return MODULE_EIGENCLASS }
func (m *ModuleClass) SuperClass() RubyClass      { return OBJECT_CLASS }

type Module struct{}

func (m *Module) Inspect() string  { return "Module" }
func (m *Module) Type() ObjectType { return MODULE_OBJ }
func (m *Module) Class() RubyClass { return MODULE_CLASS }

var moduleClassMethods = map[string]method{}
var moduleMethods = map[string]method{}
