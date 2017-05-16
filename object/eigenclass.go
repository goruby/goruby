package object

func newEigenclass(wrappedClass RubyClass, methods map[string]RubyMethod) *eigenclass {
	return &eigenclass{methods: NewMethodSet(methods), wrappedClass: wrappedClass}
}

type eigenclass struct {
	methods      SettableMethodSet
	wrappedClass RubyClass
}

func (e *eigenclass) Inspect() string {
	if e.wrappedClass != nil {
		return e.wrappedClass.(RubyClassObject).Inspect()
	}
	return "(singleton class)"
}
func (e *eigenclass) Type() Type { return EIGENCLASS_OBJ }
func (e *eigenclass) Class() RubyClass {
	if e.wrappedClass != nil {
		return e.wrappedClass
	}
	return classClass
}
func (e *eigenclass) Methods() MethodSet { return e.methods }
func (e *eigenclass) SuperClass() RubyClass {
	if e.wrappedClass != nil {
		return e.wrappedClass
	}
	return objectClass
}
func (e *eigenclass) New() RubyObject { return e.wrappedClass.New() }
func (e *eigenclass) addMethod(name string, method RubyMethod) {
	e.methods.Set(name, method)
}
