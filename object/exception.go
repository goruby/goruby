package object

import (
	"fmt"
	"reflect"
)

var (
	EXCEPTION_CLASS           RubyClassObject = NewClass("Exception", OBJECT_CLASS, exceptionMethods, exceptionClassMethods)
	STANDARD_ERROR_CLASS      RubyClassObject = NewClass("StandardError", EXCEPTION_CLASS, nil, nil)
	ZERO_DIVISION_ERROR_CLASS RubyClassObject = NewClass("ZeroDivisionError", STANDARD_ERROR_CLASS, nil, nil)
	ARGUMENT_ERROR_CLASS      RubyClassObject = NewClass("ArgumentError", STANDARD_ERROR_CLASS, nil, nil)
	NAME_ERROR_CLASS          RubyClassObject = NewClass("NameError", STANDARD_ERROR_CLASS, nil, nil)
	NO_METHOD_ERROR_CLASS     RubyClassObject = NewClass("NoMethodError", NAME_ERROR_CLASS, nil, nil)
	TYPE_ERROR_CLASS          RubyClassObject = NewClass("TypeError", STANDARD_ERROR_CLASS, nil, nil)
)

func init() {
	classes.Set("Exception", EXCEPTION_CLASS)
	classes.Set("StandardError", STANDARD_ERROR_CLASS)
	classes.Set("ZeroDivisionError", ZERO_DIVISION_ERROR_CLASS)
	classes.Set("ArgumentError", ARGUMENT_ERROR_CLASS)
	classes.Set("NameError", NAME_ERROR_CLASS)
	classes.Set("NoMethodError", NO_METHOD_ERROR_CLASS)
	classes.Set("TypeError", TYPE_ERROR_CLASS)
}

func formatException(exception RubyObject, message string) string {
	return fmt.Sprintf("%s: %s", reflect.TypeOf(exception).Elem().Name(), message)
}

type Exception struct {
	Message string
}

func (e *Exception) Type() ObjectType { return EXCEPTION_OBJ }
func (e *Exception) Inspect() string  { return formatException(e, e.Message) }
func (e *Exception) Class() RubyClass { return EXCEPTION_CLASS }

var exceptionClassMethods = map[string]RubyMethod{}

var exceptionMethods = map[string]RubyMethod{}

func NewStandardError(message string) *StandardError {
	return &StandardError{Message: message}
}

type StandardError struct {
	Message string
}

func (e *StandardError) Type() ObjectType { return EXCEPTION_OBJ }
func (e *StandardError) Inspect() string  { return formatException(e, e.Message) }
func (e *StandardError) Class() RubyClass { return STANDARD_ERROR_CLASS }

func NewZeroDivisionError() *ZeroDivisionError {
	return &ZeroDivisionError{
		Message: "divided by 0",
	}
}

type ZeroDivisionError struct {
	Message string
}

func (e *ZeroDivisionError) Type() ObjectType { return EXCEPTION_OBJ }
func (e *ZeroDivisionError) Inspect() string  { return formatException(e, e.Message) }
func (e *ZeroDivisionError) Class() RubyClass { return ZERO_DIVISION_ERROR_CLASS }

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

func (e *ArgumentError) Type() ObjectType { return EXCEPTION_OBJ }
func (e *ArgumentError) Inspect() string  { return formatException(e, e.Message) }
func (e *ArgumentError) Class() RubyClass { return ARGUMENT_ERROR_CLASS }

type NameError struct {
	Message string
}

func (e *NameError) Type() ObjectType { return EXCEPTION_OBJ }
func (e *NameError) Inspect() string  { return formatException(e, e.Message) }
func (e *NameError) Class() RubyClass { return NAME_ERROR_CLASS }

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

func (e *NoMethodError) Type() ObjectType { return EXCEPTION_OBJ }
func (e *NoMethodError) Inspect() string  { return formatException(e, e.Message) }
func (e *NoMethodError) Class() RubyClass { return NO_METHOD_ERROR_CLASS }

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

func (e *TypeError) Type() ObjectType { return EXCEPTION_OBJ }
func (e *TypeError) Inspect() string  { return formatException(e, e.Message) }
func (e *TypeError) Class() RubyClass { return TYPE_ERROR_CLASS }
