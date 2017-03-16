package object

func newEigenClass(super RubyClass, methods map[string]method) RubyClass {
	return &eigenClass{methods: methods, superClass: super}
}

type eigenClass struct {
	methods    map[string]method
	superClass RubyClass
}

func (e *eigenClass) Inspect() string            { return "" }
func (e *eigenClass) Type() ObjectType           { return EIGENCLASS_OBJ }
func (e *eigenClass) Class() RubyClass           { return OBJECT_CLASS }
func (e *eigenClass) Methods() map[string]method { return e.methods }
func (e *eigenClass) SuperClass() RubyClass      { return e.superClass }
