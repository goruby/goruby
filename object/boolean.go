package object

import "fmt"

var (
	BOOLEAN_EIGENCLASS RubyClass       = newEigenclass(CLASS_CLASS, nil)
	TRUE_CLASS         RubyClassObject = &TrueClass{}
	FALSE_CLASS        RubyClassObject = &FalseClass{}
	TRUE               RubyObject      = &Boolean{Value: true}
	FALSE              RubyObject      = &Boolean{Value: false}
)

type FalseClass struct{}

func (b *FalseClass) Inspect() string            { return "FalseClass" }
func (b *FalseClass) Type() ObjectType           { return BOOLEAN_OBJ }
func (b *FalseClass) Class() RubyClass           { return BOOLEAN_EIGENCLASS }
func (b *FalseClass) Methods() map[string]method { return booleanFalseMethods }
func (b *FalseClass) SuperClass() RubyClass      { return OBJECT_CLASS }

type TrueClass struct{}

func (b *TrueClass) Inspect() string            { return "TrueClass" }
func (b *TrueClass) Type() ObjectType           { return BOOLEAN_OBJ }
func (b *TrueClass) Class() RubyClass           { return BOOLEAN_EIGENCLASS }
func (b *TrueClass) Methods() map[string]method { return booleanTrueMethods }
func (b *TrueClass) SuperClass() RubyClass      { return OBJECT_CLASS }

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

var booleanTrueMethods = map[string]method{}

var booleanFalseMethods = map[string]method{}
