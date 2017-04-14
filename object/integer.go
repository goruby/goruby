package object

import "fmt"

var integerClass RubyClassObject = newClass("Integer", objectClass, integerMethods, integerClassMethods)

func init() {
	classes.Set("Integer", integerClass)
}

// NewInteger returns a new Integer with the given value
func NewInteger(value int64) *Integer {
	return &Integer{Value: value}
}

// Integer represents an integer in Ruby
type Integer struct {
	Value int64
}

// Inspect returns the value as string
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Type returns INTEGER_OBJ
func (i *Integer) Type() Type { return INTEGER_OBJ }

// Class returns integerClass
func (i *Integer) Class() RubyClass { return integerClass }

var integerClassMethods = map[string]RubyMethod{}

var integerMethods = map[string]RubyMethod{
	"div": withArity(1, publicMethod(integerDiv)),
	"/":   withArity(1, publicMethod(integerDiv)),
	"*":   withArity(1, publicMethod(integerMul)),
}

func integerDiv(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	divisor, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	if divisor.Value == 0 {
		return nil, NewZeroDivisionError()
	}
	return NewInteger(i.Value / divisor.Value), nil
}

func integerMul(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	factor, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value * factor.Value), nil
}

func integerAdd(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	add, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value + add.Value), nil
}
