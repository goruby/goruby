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

// MethodSet represents a set of methods
type MethodSet interface {
	// Get returns the method found for name. The boolean will return true if
	// a method was found, false otherwise
	Get(name string) (RubyMethod, bool)
	// GetAll returns a map of name to methods representing the MethodSet.
	GetAll() map[string]RubyMethod
}

// SettableMethodSet represents a MethodSet which can be mutated by setting
// methods on it.
type SettableMethodSet interface {
	MethodSet
	// Set will set method to key name. If there was a method prior defined
	// under name it will be overridden.
	Set(name string, method RubyMethod)
}

// NewMethodSet returns a new method set populated with the given methods
func NewMethodSet(methods map[string]RubyMethod) SettableMethodSet {
	return &methodSet{methods: methods}
}

type methodSet struct {
	methods map[string]RubyMethod
}

func (m *methodSet) GetAll() map[string]RubyMethod {
	methods := make(map[string]RubyMethod)
	for k, v := range m.methods {
		methods[k] = v
	}
	return methods
}

func (m *methodSet) Get(name string) (RubyMethod, bool) {
	method, ok := m.methods[name]
	return method, ok
}

func (m *methodSet) Set(name string, method RubyMethod) {
	m.methods[name] = method
}
