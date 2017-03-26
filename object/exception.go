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

// Exception represents a basic exception
type Exception struct {
	Message string
}

// Type returns the type of the RubyObject
func (e *Exception) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *Exception) Inspect() string { return formatException(e, e.Message) }

// Class returns exceptionClass
func (e *Exception) Class() RubyClass { return exceptionClass }

var exceptionClassMethods = map[string]RubyMethod{}

var exceptionMethods = map[string]RubyMethod{}

// NewStandardError returns a StandardError with the given message
func NewStandardError(message string) *StandardError {
	return &StandardError{Message: message}
}

// StandardError is the default class for rescue blocks
type StandardError struct {
	Message string
}

// Type returns EXCEPTION_OBJ
func (e *StandardError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *StandardError) Inspect() string { return formatException(e, e.Message) }

// Class returns standardErrorClass
func (e *StandardError) Class() RubyClass { return standardErrorClass }

// NewZeroDivisionError returns a new ZeroDivisionError with the default message
func NewZeroDivisionError() *ZeroDivisionError {
	return &ZeroDivisionError{
		Message: "divided by 0",
	}
}

// ZeroDivisionError represents an arithmethic error when dividing through 0
type ZeroDivisionError struct {
	Message string
}

// Type returns EXCEPTION_OBJ
func (e *ZeroDivisionError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *ZeroDivisionError) Inspect() string { return formatException(e, e.Message) }

// Class returns zeroDivisionErrorClass
func (e *ZeroDivisionError) Class() RubyClass { return zeroDivisionErrorClass }

// NewWrongNumberOfArgumentsError returns an ArgumentError populated with the default message
func NewWrongNumberOfArgumentsError(expected, actual int) *ArgumentError {
	return &ArgumentError{
		Message: fmt.Sprintf(
			"wrong number of arguments (given %d, expected %d)",
			actual,
			expected,
		),
	}
}

// ArgumentError represents an error in method call arguments
type ArgumentError struct {
	Message string
}

// Type returns EXCEPTION_OBJ
func (e *ArgumentError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *ArgumentError) Inspect() string { return formatException(e, e.Message) }

// Class returns argumentErrorClass
func (e *ArgumentError) Class() RubyClass { return argumentErrorClass }

// A NameError represents an error accessing an identifier unknown to the environment
type NameError struct {
	Message string
}

// Type returns EXCEPTION_OBJ
func (e *NameError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *NameError) Inspect() string { return formatException(e, e.Message) }

// Class returns nameErrorClass
func (e *NameError) Class() RubyClass { return nameErrorClass }

// NewNoMethodError returns a NoMethodError with the default message for undefined methods
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

// NewPrivateNoMethodError returns a NoMethodError with the default message for private methods
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

// NoMethodError represents an error finding a fitting method on an object
type NoMethodError struct {
	Message string
}

// Type returns EXCEPTION_OBJ
func (e *NoMethodError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *NoMethodError) Inspect() string { return formatException(e, e.Message) }

// Class returns noMethodErrorClass
func (e *NoMethodError) Class() RubyClass { return noMethodErrorClass }

// NewCoercionTypeError returns a TypeError with the default message for coercing errors
func NewCoercionTypeError(expected, actual RubyObject) *TypeError {
	return &TypeError{
		Message: fmt.Sprintf(
			"%s can't be coerced into %s",
			reflect.TypeOf(actual).Elem().Name(),
			reflect.TypeOf(expected).Elem().Name(),
		),
	}
}

// NewImplicitConversionTypeError returns a TypeError with the default message for impossible implicit conversions
func NewImplicitConversionTypeError(expected, actual RubyObject) *TypeError {
	return &TypeError{
		Message: fmt.Sprintf(
			"no implicit conversion of %s into %s",
			reflect.TypeOf(actual).Elem().Name(),
			reflect.TypeOf(expected).Elem().Name(),
		),
	}
}

// TypeError represents an error when the given type does not fit in the given context
type TypeError struct {
	Message string
}

// Type returns EXCEPTION_OBJ
func (e *TypeError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *TypeError) Inspect() string { return formatException(e, e.Message) }

// Class returns typeErrorClass
func (e *TypeError) Class() RubyClass { return typeErrorClass }
