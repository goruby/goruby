package object

var classes = NewEnvironment()

// NewMainEnvironment returns a new Environment populated with all Ruby classes
// and the Kernel functions
func NewMainEnvironment() Environment {
	env := kernelFunctions.Clone()
	env.Set("self", &Self{&Object{}})
	env.SetGlobal("$LOADED_FEATURES", NewArray())
	return env
}

// NewEnclosedEnvironment returns an Environment wrapped by outer
func NewEnclosedEnvironment(outer Environment) Environment {
	s := make(map[string]RubyObject)
	env := &environment{store: s, outer: nil}
	env.outer = outer
	return env
}

// NewEnvironment returns a new Environment ready to use
func NewEnvironment() Environment {
	s := make(map[string]RubyObject)
	return &environment{store: s, outer: nil}
}

// Environment holds Ruby object referenced by strings
type Environment interface {
	// Get returns the RubyObject found for this key. If it is not found,
	// ok  will be false
	Get(key string) (object RubyObject, ok bool)
	// Set sets the RubyObject for the given key. If there is already an
	// object with that key it will be overridden by object
	Set(key string, object RubyObject) RubyObject
	// SetGlobal sets val under name at the root of the environment
	SetGlobal(name string, val RubyObject) RubyObject
	// Outer returns the parent environment
	Outer() Environment
	// Clone returns a copy of the environment. It will shallow copy its values
	//
	// Note that clone will also not set its outer env, so calls to Outer will
	// return nil on cloned Environments
	Clone() Environment
}

type environment struct {
	store map[string]RubyObject
	outer Environment
}

// Get returns the RubyObject found for this key. If it is not found,
// ok  will be false
func (e *environment) Get(name string) (RubyObject, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set sets the RubyObject for the given key. If there is already an
// object with that key it will be overridden by object
func (e *environment) Set(name string, val RubyObject) RubyObject {
	e.store[name] = val
	return val
}

// SetGlobal sets val under name at the root of the environment
func (e *environment) SetGlobal(name string, val RubyObject) RubyObject {
	var env Environment = e
	for env.Outer() != nil {
		env = env.Outer()
	}
	env.Set(name, val)
	return val
}

// Outer returns the parent environment
func (e *environment) Outer() Environment {
	return e.outer
}

// Enclose encloses the environment and returns a new one wrapped by outer
func (e *environment) Enclose(outer Environment) Environment {
	env := e.clone()
	env.outer = outer
	return env
}

func (e *environment) Clone() Environment {
	return e.clone()
}

func (e *environment) clone() *environment {
	s := make(map[string]RubyObject)
	env := &environment{store: s, outer: nil}
	for k, v := range e.store {
		env.store[k] = v
	}
	return env
}
