package object

var (
	MODULE_EIGENCLASS RubyClass = &ModuleEigenClass{}
	MODULE_CLASS      RubyClass = &ModuleClass{}
)

type ModuleEigenClass struct{}

func (m *ModuleEigenClass) Inspect() string            { return "Module" }
func (m *ModuleEigenClass) Type() ObjectType           { return EIGENCLASS_OBJ }
func (m *ModuleEigenClass) Methods() map[string]method { return moduleMethods }
func (m *ModuleEigenClass) Class() RubyClass           { return OBJECT_CLASS }
func (m *ModuleEigenClass) SuperClass() RubyClass      { return BASIC_OBJECT_CLASS }

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
