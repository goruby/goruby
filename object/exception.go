package object

import (
	"fmt"
	"reflect"
)

var (
	EXCEPTION_EIGENCLASS RubyClass       = newEigenclass(CLASS_CLASS, exceptionClassMethods)
	EXCEPTION_CLASS      RubyClassObject = &ExceptionClass{}
)

type ExceptionClass struct{}

func (e *ExceptionClass) Type() ObjectType           { return EXCEPTION_OBJ }
func (e *ExceptionClass) Inspect() string            { return "Exception" }
func (e *ExceptionClass) Class() RubyClass           { return EXCEPTION_EIGENCLASS }
func (e *ExceptionClass) Methods() map[string]method { return exceptionMethods }
func (e *ExceptionClass) SuperClass() RubyClass      { return OBJECT_CLASS }

type Exception struct {
	exception interface{}
	Message   string
}

func (e *Exception) Type() ObjectType { return EXCEPTION_OBJ }
func (e *Exception) Inspect() string {
	return fmt.Sprintf("%s: %s", reflect.TypeOf(e.exception).Elem().Name(), e.Message)
}
func (e *Exception) Class() RubyClass { return EXCEPTION_CLASS }

var exceptionClassMethods = map[string]method{}

var exceptionMethods = map[string]method{}

func NewStandardError(message string) *StandardError {
	e := &StandardError{Exception{Message: message}}
	e.exception = e
	return e
}

type StandardError struct {
	Exception
}

func NewZeroDivisionError() *ZeroDivisionError {
	e := &ZeroDivisionError{
		StandardError{
			Exception{
				Message: "divided by 0",
			},
		},
	}
	e.exception = e
	return e
}

type ZeroDivisionError struct {
	StandardError
}

func NewWrongNumberOfArgumentsError(expected, actual int) *ArgumentError {
	e := &ArgumentError{
		StandardError{
			Exception{
				Message: fmt.Sprintf(
					"wrong number of arguments (given %d, expected %d)",
					actual,
					expected,
				),
			},
		},
	}
	e.exception = e
	return e
}

type ArgumentError struct {
	StandardError
}

type NameError struct {
	StandardError
}

func NewNoMethodError(context RubyObject, method string) *NoMethodError {
	e := &NoMethodError{
		NameError{
			StandardError{
				Exception{
					Message: fmt.Sprintf(
						"undefined method `%s' for %s:%s",
						method,
						context.Inspect(),
						context.Class().(RubyObject).Inspect(),
					),
				},
			},
		},
	}
	e.exception = e
	return e
}

type NoMethodError struct {
	NameError
}

func NewCoercionTypeError(expected, actual RubyObject) *TypeError {
	e := &TypeError{
		StandardError{
			Exception{
				Message: fmt.Sprintf(
					"%s can't be coerced into %s",
					reflect.TypeOf(actual).Elem().Name(),
					reflect.TypeOf(expected).Elem().Name(),
				),
			},
		},
	}
	e.exception = e
	return e
}

func NewImplicitConversionTypeError(expected, actual RubyObject) *TypeError {
	e := &TypeError{
		StandardError{
			Exception{
				Message: fmt.Sprintf(
					"no implicit conversion of %s into %s",
					reflect.TypeOf(actual).Elem().Name(),
					reflect.TypeOf(expected).Elem().Name(),
				),
			},
		},
	}
	e.exception = e
	return e
}

type TypeError struct {
	StandardError
}
