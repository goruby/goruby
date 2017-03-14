package object

import "fmt"

var INTEGER *IntegerClass = newIntegerClass()

func newIntegerClass() *IntegerClass {
	cl := &IntegerClass{}
	cl.methods = &methodSet{cl, integerClassMethods}
	cl.instanceMethods = &methodSet{methods: integerMethods}
	return cl
}

type IntegerClass struct {
	methods         *methodSet
	instanceMethods *methodSet
}

func (i *IntegerClass) Inspect() string  { return "Integer" }
func (i *IntegerClass) Type() ObjectType { return INTEGER_CLASS_OBJ }
func (i *IntegerClass) Send(name string, args ...RubyObject) RubyObject {
	return i.methods.Call(name, args...)
}

func NewInteger(value int64) *Integer {
	i := &Integer{Value: value}
	i.methods = INTEGER.instanceMethods
	i.methods.context = i
	return i
}

type Integer struct {
	Value   int64
	methods *methodSet
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Send(name string, args ...RubyObject) RubyObject {
	return i.methods.Call(name, args...)
}

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
	result := &Integer{Value: i.Value / divisor.Value}
	result.methods = i.methods.SetContext(result)
	return result
}

func integerMul(context RubyObject, args ...RubyObject) RubyObject {
	i := context.(*Integer)
	factor, ok := args[0].(*Integer)
	if !ok {
		return NewCoercionTypeError(args[0], i)
	}
	result := &Integer{Value: i.Value * factor.Value}
	result.methods = i.methods.SetContext(result)
	return result
}

func integerAdd(context RubyObject, args ...RubyObject) RubyObject {
	i := context.(*Integer)
	add, ok := args[0].(*Integer)
	if !ok {
		return NewCoercionTypeError(args[0], i)
	}
	result := &Integer{Value: i.Value + add.Value}
	result.methods = i.methods.SetContext(result)
	return result
}
