package object

func newEigenclass(wrappedClass RubyClass, methods map[string]RubyMethod) RubyClassObject {
	return &eigenclass{methods: methods, wrappedClass: wrappedClass}
}

type eigenclass struct {
	methods      map[string]RubyMethod
	wrappedClass RubyClass
}

func (e *eigenclass) Inspect() string {
	if e.wrappedClass != nil {
		return e.wrappedClass.Inspect()
	}
	return "(singleton class)"
}
func (e *eigenclass) Type() ObjectType { return EIGENCLASS_OBJ }
func (e *eigenclass) Class() RubyClass {
	if e.wrappedClass != nil {
		return e.wrappedClass
	}
	return CLASS_CLASS
}
func (e *eigenclass) Methods() map[string]RubyMethod { return e.methods }
func (e *eigenclass) SuperClass() RubyClass {
	if e.wrappedClass != nil {
		return e.wrappedClass
	}
	return OBJECT_CLASS
}
