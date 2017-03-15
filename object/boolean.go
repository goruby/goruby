package object

import "fmt"

var BOOLEAN_EIGENCLASS RubyClass = &BooleanEigenclass{}
var TRUE_CLASS RubyClass = &TrueClass{}
var FALSE_CLASS RubyClass = &FalseClass{}

type BooleanEigenclass struct{}

func (b *BooleanEigenclass) Inspect() string            { return "" }
func (b *BooleanEigenclass) Type() ObjectType           { return BOOLEAN_OBJ }
func (b *BooleanEigenclass) Methods() map[string]method { return nil }
func (b *BooleanEigenclass) Class() RubyClass           { return BASIC_OBJECT_CLASS }
func (b *BooleanEigenclass) SuperClass() RubyClass      { return BASIC_OBJECT_CLASS }

type FalseClass struct{}

func (b *FalseClass) Inspect() string            { return "FalseClass" }
func (b *FalseClass) Type() ObjectType           { return BOOLEAN_OBJ }
func (b *FalseClass) Methods() map[string]method { return nil }
func (b *FalseClass) Class() RubyClass           { return BOOLEAN_EIGENCLASS }
func (b *FalseClass) SuperClass() RubyClass      { return OBJECT_CLASS }

type TrueClass struct{}

func (b *TrueClass) Inspect() string            { return "TrueClass" }
func (b *TrueClass) Type() ObjectType           { return BOOLEAN_OBJ }
func (b *TrueClass) Methods() map[string]method { return nil }
func (b *TrueClass) Class() RubyClass           { return BOOLEAN_EIGENCLASS }
func (b *TrueClass) SuperClass() RubyClass      { return OBJECT_CLASS }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string            { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType           { return BOOLEAN_OBJ }
func (b *Boolean) Methods() map[string]method { return nil }
func (b *Boolean) Class() RubyClass {
	if b.Value {
		return TRUE_CLASS
	}
	return FALSE_CLASS
}
