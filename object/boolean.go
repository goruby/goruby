package object

import "fmt"

var (
	TRUE_CLASS  RubyClassObject = NewClass("TrueClass", OBJECT_CLASS, booleanTrueMethods, nil)
	FALSE_CLASS RubyClassObject = NewClass("FalseClass", OBJECT_CLASS, booleanFalseMethods, nil)
	TRUE        RubyObject      = &Boolean{Value: true}
	FALSE       RubyObject      = &Boolean{Value: false}
)

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Class() RubyClass {
	if b.Value {
		return TRUE_CLASS
	}
	return FALSE_CLASS
}

var booleanTrueMethods = map[string]RubyMethod{}

var booleanFalseMethods = map[string]RubyMethod{}
