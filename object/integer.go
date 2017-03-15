package object

import "fmt"

var (
	INTEGER_EIGENCLASS RubyClass = &IntegerEigenclass{}
	INTEGER_CLASS      RubyClass = &IntegerClass{}
)

type IntegerEigenclass struct{}

func (i *IntegerEigenclass) Inspect() string  { return "Integer" }
func (i *IntegerEigenclass) Type() ObjectType { return INTEGER_CLASS_OBJ }
func (i *IntegerEigenclass) Methods() map[string]method {
	return integerClassMethods
}
func (i *IntegerEigenclass) Class() RubyClass      { return OBJECT_CLASS }
func (i *IntegerEigenclass) SuperClass() RubyClass { return BASIC_OBJECT_CLASS }

type IntegerClass struct{}

func (i *IntegerClass) Inspect() string  { return "Integer" }
func (i *IntegerClass) Type() ObjectType { return INTEGER_CLASS_OBJ }
func (i *IntegerClass) Methods() map[string]method {
	return integerMethods
}
func (i *IntegerClass) Class() RubyClass      { return INTEGER_EIGENCLASS }
func (i *IntegerClass) SuperClass() RubyClass { return OBJECT_CLASS }

func NewInteger(value int64) *Integer {
	return &Integer{Value: value}
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Class() RubyClass { return INTEGER_CLASS }

var integerClassMethods = map[string]method{}

var integerMethods = map[string]method{
	"div": withArity(1, integerDiv),
	"/":   withArity(1, integerDiv),
	"*":   withArity(1, integerMul),
}

func integerDiv(context RubyObject, args ...RubyObject) RubyObject {
	i := context.(*Integer)
	divisor, ok := args[0].(*Integer)
	if !ok {
		return NewCoercionTypeError(args[0], i)
	}
	if divisor.Value == 0 {
		return NewZeroDivisionError()
	}
	return NewInteger(i.Value / divisor.Value)
}

func integerMul(context RubyObject, args ...RubyObject) RubyObject {
	i := context.(*Integer)
	factor, ok := args[0].(*Integer)
	if !ok {
		return NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value * factor.Value)
}

func integerAdd(context RubyObject, args ...RubyObject) RubyObject {
	i := context.(*Integer)
	add, ok := args[0].(*Integer)
	if !ok {
		return NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value + add.Value)
}