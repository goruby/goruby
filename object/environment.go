package object

var classes = &Environment{
	store: map[string]RubyObject{},
}

func NewMainEnvironment() *Environment {
	kernelFunctions.outer = classes
	return NewEnclosedEnvironment(kernelFunctions)
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]RubyObject)
	return &Environment{store: s, outer: nil}
}

type Environment struct {
	store map[string]RubyObject
	outer *Environment
}

func (e *Environment) Get(name string) (RubyObject, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val RubyObject) RubyObject {
	e.store[name] = val
	return val
}

func (e *Environment) Enclose(outer *Environment) *Environment {
	env := e.Clone()
	env.outer = outer
	return env
}

func (e *Environment) Clone() *Environment {
	env := NewEnvironment()
	for k, v := range e.store {
		env.store[k] = v
	}
	return env
}
