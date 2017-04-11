package object

type visibility int

// MethodVisibility represents the visibility of a method
type MethodVisibility visibility

const (
	// PUBLIC_METHOD declares that a method is visible from the outside of an object
	PUBLIC_METHOD MethodVisibility = iota
	// PROTECTED_METHOD declares that a method is not visible from the outside
	// of an object but to all of its decendents
	PROTECTED_METHOD
	// PRIVATE_METHOD declares that a method is not visible from the outside
	// of an object and not to all of its decendents
	PRIVATE_METHOD
)

// RubyMethod defines a Ruby method
type RubyMethod interface {
	Call(context CallContext, args ...RubyObject) (RubyObject, error)
	Visibility() MethodVisibility
}

func withArity(arity int, fn RubyMethod) RubyMethod {
	return &method{
		fn: func(context CallContext, args ...RubyObject) (RubyObject, error) {
			if len(args) != arity {
				return nil, NewWrongNumberOfArgumentsError(arity, len(args))
			}
			return fn.Call(context, args...)
		},
		visibility: fn.Visibility(),
	}
}

func publicMethod(fn func(context CallContext, args ...RubyObject) (RubyObject, error)) RubyMethod {
	return &method{visibility: PUBLIC_METHOD, fn: fn}
}

func protectedMethod(fn func(context CallContext, args ...RubyObject) (RubyObject, error)) RubyMethod {
	return &method{visibility: PROTECTED_METHOD, fn: fn}
}

func privateMethod(fn func(context CallContext, args ...RubyObject) (RubyObject, error)) RubyMethod {
	return &method{visibility: PRIVATE_METHOD, fn: fn}
}

type method struct {
	visibility MethodVisibility
	fn         func(context CallContext, args ...RubyObject) (RubyObject, error)
}

func (m *method) Call(context CallContext, args ...RubyObject) (RubyObject, error) {
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
