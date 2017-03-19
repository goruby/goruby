package object

import "fmt"

var INTEGER_CLASS RubyClassObject = NewClass("Integer", OBJECT_CLASS, integerMethods, integerClassMethods)

func init() {
	classes.Set("Integer", INTEGER_CLASS)
}

func NewInteger(value int64) *Integer {
	return &Integer{Value: value}
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Class() RubyClass { return INTEGER_CLASS }

var integerClassMethods = map[string]RubyMethod{}

var integerMethods = map[string]RubyMethod{
	"div": withArity(1, publicMethod(integerDiv)),
	"/":   withArity(1, publicMethod(integerDiv)),
	"*":   withArity(1, publicMethod(integerMul)),
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
