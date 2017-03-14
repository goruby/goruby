package object

func withArity(arity int, fn method) method {
	return func(context RubyObject, args ...RubyObject) RubyObject {
		if len(args) != arity {
			return NewWrongNumberOfArgumentsError(arity, len(args))
		}
		return fn(context, args...)
	}
}

type method func(context RubyObject, args ...RubyObject) RubyObject

type methodSet struct {
	context RubyObject
	methods map[string]method
}

func (m *methodSet) SetContext(context RubyObject) *methodSet {
	return &methodSet{context: context, methods: m.methods}
}

func (m *methodSet) Define(name string, fn method) *Symbol {
	m.methods[name] = fn
	return &Symbol{name}
}

func (m *methodSet) Call(name string, args ...RubyObject) RubyObject {
	method, ok := m.methods[name]
	if !ok {
		return NewNoMethodError(m.context, name)
	}
	return method(m.context, args...)
}
