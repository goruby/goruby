package object

type function func(args ...RubyObject) RubyObject

type methodSet struct {
	context RubyObject
	methods map[string]function
}

func (m *methodSet) Define(name string, fn function) *Symbol {
	m.methods[name] = fn
	return &Symbol{name}
}
