package object

import (
	"fmt"
	"reflect"
)

var (
	exceptionClass         RubyClassObject = newClass("Exception", objectClass, exceptionMethods, exceptionClassMethods)
	standardErrorClass     RubyClassObject = newClass("StandardError", exceptionClass, nil, nil)
	zeroDivisionErrorClass RubyClassObject = newClass("ZeroDivisionError", standardErrorClass, nil, nil)
	argumentErrorClass     RubyClassObject = newClass("ArgumentError", standardErrorClass, nil, nil)
	nameErrorClass         RubyClassObject = newClass("NameError", standardErrorClass, nil, nil)
	noMethodErrorClass     RubyClassObject = newClass("NoMethodError", nameErrorClass, nil, nil)
	typeErrorClass         RubyClassObject = newClass("TypeError", standardErrorClass, nil, nil)
)

func init() {
	classes.Set("Exception", exceptionClass)
	classes.Set("StandardError", standardErrorClass)
	classes.Set("ZeroDivisionError", zeroDivisionErrorClass)
	classes.Set("ArgumentError", argumentErrorClass)
	classes.Set("NameError", nameErrorClass)
	classes.Set("NoMethodError", noMethodErrorClass)
	classes.Set("TypeError", typeErrorClass)
}

func formatException(exception RubyObject, message string) string {
	return fmt.Sprintf("%s: %s", reflect.TypeOf(exception).Elem().Name(), message)
}

type Exception struct {
	Message string
}

func (e *Exception) Type() Type       { return EXCEPTION_OBJ }
func (e *Exception) Inspect() string  { return formatException(e, e.Message) }
func (e *Exception) Class() RubyClass { return exceptionClass }

var exceptionClassMethods = map[string]RubyMethod{}

var exceptionMethods = map[string]RubyMethod{}

func NewStandardError(message string) *StandardError {
	return &StandardError{Message: message}
}

type StandardError struct {
	Message string
}

func (e *StandardError) Type() Type       { return EXCEPTION_OBJ }
func (e *StandardError) Inspect() string  { return formatException(e, e.Message) }
func (e *StandardError) Class() RubyClass { return standardErrorClass }

func NewZeroDivisionError() *ZeroDivisionError {
	return &ZeroDivisionError{
		Message: "divided by 0",
	}
}

type ZeroDivisionError struct {
	Message string
}

func (e *ZeroDivisionError) Type() Type       { return EXCEPTION_OBJ }
func (e *ZeroDivisionError) Inspect() string  { return formatException(e, e.Message) }
func (e *ZeroDivisionError) Class() RubyClass { return zeroDivisionErrorClass }

func NewWrongNumberOfArgumentsError(expected, actual int) *ArgumentError {
	return &ArgumentError{
		Message: fmt.Sprintf(
			"wrong number of arguments (given %d, expected %d)",
			actual,
			expected,
		),
	}
}

type ArgumentError struct {
	Message string
}

func (e *ArgumentError) Type() Type       { return EXCEPTION_OBJ }
func (e *ArgumentError) Inspect() string  { return formatException(e, e.Message) }
func (e *ArgumentError) Class() RubyClass { return argumentErrorClass }

type NameError struct {
	Message string
}

func (e *NameError) Type() Type       { return EXCEPTION_OBJ }
func (e *NameError) Inspect() string  { return formatException(e, e.Message) }
func (e *NameError) Class() RubyClass { return nameErrorClass }

func NewNoMethodError(context RubyObject, method string) *NoMethodError {
	return &NoMethodError{
		Message: fmt.Sprintf(
			"undefined method `%s' for %s:%s",
			method,
			context.Inspect(),
			context.Class().(RubyObject).Inspect(),
		),
	}
}

func NewPrivateNoMethodError(context RubyObject, method string) *NoMethodError {
	return &NoMethodError{
		Message: fmt.Sprintf(
			"private method `%s' called for %s:%s",
			method,
			context.Inspect(),
			context.Class().(RubyObject).Inspect(),
		),
	}
}

type NoMethodError struct {
	Message string
}

func (e *NoMethodError) Type() Type       { return EXCEPTION_OBJ }
func (e *NoMethodError) Inspect() string  { return formatException(e, e.Message) }
func (e *NoMethodError) Class() RubyClass { return noMethodErrorClass }

func NewCoercionTypeError(expected, actual RubyObject) *TypeError {
	return &TypeError{
		Message: fmt.Sprintf(
			"%s can't be coerced into %s",
			reflect.TypeOf(actual).Elem().Name(),
			reflect.TypeOf(expected).Elem().Name(),
		),
	}
}

func NewImplicitConversionTypeError(expected, actual RubyObject) *TypeError {
	return &TypeError{
		Message: fmt.Sprintf(
			"no implicit conversion of %s into %s",
			reflect.TypeOf(actual).Elem().Name(),
			reflect.TypeOf(expected).Elem().Name(),
		),
	}
}

type TypeError struct {
	Message string
}

func (e *TypeError) Type() Type       { return EXCEPTION_OBJ }
func (e *TypeError) Inspect() string  { return formatException(e, e.Message) }
func (e *TypeError) Class() RubyClass { return typeErrorClass }
