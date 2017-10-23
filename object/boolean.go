package object

import "fmt"

var (
	trueClass  RubyClassObject = newClass("TrueClass", objectClass, booleanTrueMethods, nil, notInstantiatable)
	falseClass RubyClassObject = newClass("FalseClass", objectClass, booleanFalseMethods, nil, notInstantiatable)
	// TRUE represents the singleton object for the Boolean true
	TRUE RubyObject = &Boolean{Value: true}
	// FALSE represents the singleton object for the Boolean false
	FALSE RubyObject = &Boolean{Value: false}
)

func init() {
	classes.Set("TrueClass", trueClass)
	classes.Set("FalseClass", falseClass)
}

// Boolean represents a Boolean object in Ruby
type Boolean struct {
	Value bool
}

// Inspect returns the string representation of the boolean
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Type returns the object type of the boolean
func (b *Boolean) Type() Type { return BOOLEAN_OBJ }

// Class returns the TrueClass for true, and the FalseClass otherwise
func (b *Boolean) Class() RubyClass {
	if b.Value {
		return trueClass
	}
	return falseClass
}

func (b *Boolean) hashKey() hashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return hashKey{Type: b.Type(), Value: value}
}

var booleanTrueMethods = map[string]RubyMethod{}

var booleanFalseMethods = map[string]RubyMethod{}
