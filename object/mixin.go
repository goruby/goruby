package object

func newMixin(class RubyClassObject, modules ...*Module) *mixin {
	return &mixin{class, modules}
}

type mixin struct {
	RubyClassObject
	modules []*Module
}

func (m *mixin) Methods() MethodSet {
	var methods = make(map[string]RubyMethod)
	for _, mod := range m.modules {
		moduleMethods := mod.Class().Methods()
		for k, v := range moduleMethods.GetAll() {
			methods[k] = v
		}
	}
	for k, v := range m.RubyClassObject.Methods().GetAll() {
		methods[k] = v
	}
	return NewMethodSet(methods)
}
