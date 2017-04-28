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
	"+":   withArity(1, publicMethod(integerAdd)),
	"-":   withArity(1, publicMethod(integerSub)),
	"<":   withArity(1, publicMethod(integerLessThan)),
	"<=":  withArity(1, publicMethod(integerLessThanOrEqual)),
	">":   withArity(1, publicMethod(integerGreaterThan)),
	">=":  withArity(1, publicMethod(integerGreaterThanOrEqual)),
	"<=>": withArity(1, publicMethod(integerSpaceShipOperator)),
	"-@":  withArity(0, publicMethod(integerUnaryMinus)),
	"+@":  withArity(0, publicMethod(integerUnaryPlus)),
}

func integerDiv(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	if right.Value == 0 {
		return nil, NewZeroDivisionError()
	}
	return NewInteger(i.Value / right.Value), nil
}

func integerMul(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value * right.Value), nil
}

func integerAdd(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value + right.Value), nil
}

func integerSub(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value - right.Value), nil
}

func integerLessThan(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	if i.Value < right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerLessThanOrEqual(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	if i.Value <= right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerGreaterThan(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	if i.Value > right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerGreaterThanOrEqual(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	if i.Value >= right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerSpaceShipOperator(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return NIL, nil
	}
	switch {
	case i.Value > right.Value:
		return NewInteger(1), nil
	case i.Value < right.Value:
		return NewInteger(-1), nil
	case i.Value == right.Value:
		return NewInteger(0), nil
	default:
		return NIL, nil
	}
}

func integerUnaryMinus(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	return NewInteger(-i.Value), nil
}

func integerUnaryPlus(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	return NewInteger(i.Value), nil
}
