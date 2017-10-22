package object

import "fmt"

var integerClass RubyClassObject = newClass(
	"Integer", objectClass, integerMethods, integerClassMethods, notInstantiatable,
)

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

func (i *Integer) hashKey() hashKey {
	return hashKey{Type: i.Type(), Value: uint64(i.Value)}
}

var integerClassMethods = map[string]RubyMethod{}

var integerMethods = map[string]RubyMethod{
	"div": withArity(1, publicMethod(integerDiv)),
	"/":   withArity(1, publicMethod(integerDiv)),
	"*":   withArity(1, publicMethod(integerMul)),
	"+":   withArity(1, publicMethod(integerAdd)),
	"-":   withArity(1, publicMethod(integerSub)),
	"%":   withArity(1, publicMethod(integerModulo)),
	"<":   withArity(1, publicMethod(integerLt)),
	">":   withArity(1, publicMethod(integerGt)),
	"==":  withArity(1, publicMethod(integerEq)),
	"!=":  withArity(1, publicMethod(integerNeq)),
	">=":  withArity(1, publicMethod(integerGte)),
	"<=":  withArity(1, publicMethod(integerLte)),
	"<=>": withArity(1, publicMethod(integerSpaceship)),
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

func integerSub(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	sub, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value - sub.Value), nil
}

func integerModulo(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	mod, ok := args[0].(*Integer)
	if !ok {
		return nil, NewCoercionTypeError(args[0], i)
	}
	return NewInteger(i.Value % mod.Value), nil
}

func integerLt(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value < right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerGt(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value > right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerEq(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value == right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerNeq(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value != right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerSpaceship(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return NIL, nil
	}
	switch {
	case i.Value > right.Value:
		return &Integer{Value: 1}, nil
	case i.Value < right.Value:
		return &Integer{Value: -1}, nil
	case i.Value == right.Value:
		return &Integer{Value: 0}, nil
	default:
		panic("not reachable")
	}
}

func integerGte(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value >= right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}

func integerLte(context CallContext, args ...RubyObject) (RubyObject, error) {
	i := context.Receiver().(*Integer)
	right, ok := args[0].(*Integer)
	if !ok {
		return nil, NewArgumentError(
			"comparison of Integer with %s failed",
			args[0].Class().(RubyObject).Inspect(),
		)
	}
	if i.Value <= right.Value {
		return TRUE, nil
	}
	return FALSE, nil
}
