package object

import (
	"bytes"
	"fmt"
	"unicode"
)

var classes = NewEnvironment()
var mainObj = &Object{}
var mainObject = &extendedObject{
	RubyObject:  mainObj,
	class:       newEigenclass(mainObj.Class().(RubyClassObject), map[string]RubyMethod{}),
	Environment: NewEnvironment(),
}

// NewMainEnvironment returns a new Environment populated with all Ruby classes
// and the Kernel functions
func NewMainEnvironment() Environment {
	loadPath := NewArray()
	env := classes.Clone()
	env.Set("self", &Self{RubyObject: mainObject, Name: "main"})
	env.SetGlobal("$LOADED_FEATURES", NewArray())
	env.SetGlobal("$:", loadPath)
	env.SetGlobal("$LOAD_PATH", loadPath)
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

// WithScopedLocalVariables returns an Environment which scopes the local variables
// within another env so they cannot leak out
func WithScopedLocalVariables(e Environment) Environment {
	return &localVariableGuard{
		Environment:    e,
		localVariables: &environment{store: make(map[string]RubyObject)},
	}
}

// Environment holds Ruby object referenced by strings
type Environment interface {
	// Get returns the RubyObject found for this key. If it is not found,
	// ok  will be false
	Get(key string) (object RubyObject, ok bool)
	// GetAll returns all values for the current env as a map of string to RubyObject
	GetAll() map[string]RubyObject
	// Set sets the RubyObject for the given key. If there is already an
	// object with that key it will be overridden by object
	Set(key string, object RubyObject) RubyObject
	// Unset removes the entry for the given key. It returns the removed entry
	Unset(key string) RubyObject
	// SetGlobal sets val under name at the root of the environment
	SetGlobal(name string, val RubyObject) RubyObject
	// Outer returns the parent environment
	Outer() Environment
	// Clone returns a copy of the environment. It will shallow copy its values
	//
	// Note that clone will not set its outer env, so calls to Outer will
	// return nil on cloned Environments
	Clone() Environment
}

// EnvEntryInfo describes an entry in an Environment and is returned by EnvStat
type EnvEntryInfo interface {
	Name() string
	Env() Environment
}

type envEntryInfo struct {
	name string
	env  Environment
}

func (e *envEntryInfo) Name() string     { return e.name }
func (e *envEntryInfo) Env() Environment { return e.env }

// EnvStat returns the EnvEntryInfo for obj. If obj is not found in the hierarchy
// of env, the bool will be false.
func EnvStat(env Environment, obj RubyObject) (EnvEntryInfo, bool) {
	var info envEntryInfo
	for key, value := range env.GetAll() {
		if value == obj {
			info.name = key
			info.env = env
			return &info, true
		}
	}

	outer := env.Outer()
	if outer == nil {
		return &info, false
	}

	return EnvStat(outer, obj)
}

type localVariableGuard struct {
	Environment
	localVariables *environment
}

func (l *localVariableGuard) isLocalVariable(name string) bool {
	firstChar := bytes.Runes([]byte(name))[0]
	// TODO: encapsulate this logic somewhere where the choice is first done, e.g. parser
	return unicode.IsLower(firstChar) && firstChar != '@' && name != "self"
}

func (l *localVariableGuard) Set(name string, value RubyObject) RubyObject {
	if l.isLocalVariable(name) {
		return l.localVariables.Set(name, value)
	}
	return l.Environment.Set(name, value)
}

func (l *localVariableGuard) Get(name string) (RubyObject, bool) {
	if l.isLocalVariable(name) {
		return l.localVariables.Get(name)
	}
	return l.Environment.Get(name)
}

func (l *localVariableGuard) GetAll() map[string]RubyObject {
	store := l.Environment.Clone().GetAll()
	for k, v := range l.localVariables.clone().store {
		store[k] = v
	}
	return store
}

func (l *localVariableGuard) Unset(key string) RubyObject {
	if l.isLocalVariable(key) {
		return l.localVariables.Unset(key)
	}
	return l.Environment.Unset(key)
}

func (l *localVariableGuard) Clone() Environment {
	return &localVariableGuard{
		Environment:    l.Environment.Clone(),
		localVariables: l.localVariables.clone(),
	}
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

func (e *environment) GetAll() map[string]RubyObject {
	return e.clone().store
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

func (e *environment) Unset(key string) RubyObject {
	val := e.store[key]
	delete(e.store, key)
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

func (e *environment) String() string {
	var out bytes.Buffer
	fmt.Fprintf(&out, "%v", e.store)
	return out.String()
}
