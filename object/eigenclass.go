package object

func newEigenclass(wrappedClass RubyClassObject, methods map[string]RubyMethod) *eigenclass {
	return &eigenclass{methods: methods, wrappedClass: wrappedClass}
}

type eigenclass struct {
	methods      map[string]RubyMethod
	wrappedClass RubyClassObject
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
func (e *eigenclass) Methods() map[string]RubyMethod { return e.methods }
func (e *eigenclass) SuperClass() RubyClass {
	if e.wrappedClass != nil {
		return e.wrappedClass
	}
	return objectClass
}
func (e *eigenclass) addMethod(name string, method RubyMethod) {
	e.methods[name] = method
}
