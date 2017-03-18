package object

type visibility int

type MethodVisibility visibility

const (
	PUBLIC_METHOD MethodVisibility = iota
	PROTECTED_METHOD
	PRIVATE_METHOD
)

type RubyMethod interface {
	Call(context RubyObject, args ...RubyObject) RubyObject
	Visibility() MethodVisibility
}

func withArity(arity int, fn RubyMethod) RubyMethod {
	return &method{
		fn: func(context RubyObject, args ...RubyObject) RubyObject {
			if len(args) != arity {
				return NewWrongNumberOfArgumentsError(arity, len(args))
			}
			return fn.Call(context, args...)
		},
		visibility: fn.Visibility(),
	}
}

func publicMethod(fn func(context RubyObject, args ...RubyObject) RubyObject) RubyMethod {
	return &method{visibility: PUBLIC_METHOD, fn: fn}
}

func protectedMethod(fn func(context RubyObject, args ...RubyObject) RubyObject) RubyMethod {
	return &method{visibility: PROTECTED_METHOD, fn: fn}
}

func privateMethod(fn func(context RubyObject, args ...RubyObject) RubyObject) RubyMethod {
	return &method{visibility: PRIVATE_METHOD, fn: fn}
}

type publicMethodX func(context RubyObject, args ...RubyObject) RubyObject

func (m publicMethodX) Call(context RubyObject, args ...RubyObject) RubyObject {
	return m(context, args...)
}
func (m publicMethodX) Visibility() MethodVisibility { return PUBLIC_METHOD }

type method struct {
	visibility MethodVisibility
	fn         func(context RubyObject, args ...RubyObject) RubyObject
}

func (m *method) Call(context RubyObject, args ...RubyObject) RubyObject {
	return m.fn(context, args...)
}
func (m *method) Visibility() MethodVisibility { return m.visibility }

func mixin(class RubyClassObject, modules ...*Module) RubyClassObject {
	return &methodSet{class, modules}
}

type methodSet struct {
	RubyClassObject
	modules []*Module
}

func (m *methodSet) Methods() map[string]RubyMethod {
	var methods = make(map[string]RubyMethod)
	for _, mod := range m.modules {
		moduleMethods := mod.Class().Methods()
		for k, v := range moduleMethods {
			methods[k] = v
		}
	}
	for k, v := range m.RubyClassObject.Methods() {
		methods[k] = v
	}
	return methods
}
