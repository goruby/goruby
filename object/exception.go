package object

import (
	"fmt"
	"reflect"
)

var (
	exceptionClass RubyClassObject = newClass(
		"Exception",
		objectClass,
		exceptionMethods,
		exceptionClassMethods,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &Exception{message: c.Name()}, nil
		},
	)
	standardErrorClass RubyClassObject = newClass(
		"StandardError",
		exceptionClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &StandardError{message: c.Name()}, nil
		},
	)
	runtimeErrorClass RubyClassObject = newClass(
		"RuntimeError",
		standardErrorClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &RuntimeError{message: c.Name()}, nil
		},
	)
	zeroDivisionErrorClass RubyClassObject = newClass(
		"ZeroDivisionError",
		standardErrorClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &ZeroDivisionError{message: c.Name()}, nil
		},
	)
	argumentErrorClass RubyClassObject = newClass(
		"ArgumentError",
		standardErrorClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &ArgumentError{message: c.Name()}, nil
		},
	)
	nameErrorClass RubyClassObject = newClass(
		"NameError",
		standardErrorClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &NameError{message: c.Name()}, nil
		},
	)
	noMethodErrorClass RubyClassObject = newClass(
		"NoMethodError",
		nameErrorClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &NoMethodError{message: c.Name()}, nil
		},
	)
	typeErrorClass RubyClassObject = newClass(
		"TypeError",
		standardErrorClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &TypeError{message: c.Name()}, nil
		},
	)
	scriptErrorClass RubyClassObject = newClass(
		"ScriptError",
		exceptionClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &ScriptError{message: c.Name()}, nil
		},
	)
	loadErrorClass RubyClassObject = newClass(
		"LoadError",
		scriptErrorClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &LoadError{message: c.Name()}, nil
		},
	)
	syntaxErrorClass RubyClassObject = newClass(
		"SyntaxError",
		scriptErrorClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &SyntaxError{message: c.Name()}, nil
		},
	)
	notImplementedErrorClass RubyClassObject = newClass(
		"NotImplementedError",
		scriptErrorClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &NotImplementedError{message: c.Name()}, nil
		},
	)
	localJumpErrorClass RubyClassObject = newClass(
		"LocalJumpError",
		standardErrorClass,
		nil,
		nil,
		func(c RubyClassObject, args ...RubyObject) (RubyObject, error) {
			return &LocalJumpError{message: c.Name()}, nil
		},
	)
)

func init() {
	classes.Set("Exception", exceptionClass)
	classes.Set("StandardError", standardErrorClass)
	classes.Set("ZeroDivisionError", zeroDivisionErrorClass)
	classes.Set("ArgumentError", argumentErrorClass)
	classes.Set("NameError", nameErrorClass)
	classes.Set("NoMethodError", noMethodErrorClass)
	classes.Set("TypeError", typeErrorClass)
	classes.Set("ScriptError", scriptErrorClass)
	classes.Set("LoadError", loadErrorClass)
	classes.Set("SyntaxError", syntaxErrorClass)
	classes.Set("NotImplementedError", notImplementedErrorClass)
}

func formatException(exception RubyObject, message string) string {
	return fmt.Sprintf("%s: %s", reflect.TypeOf(exception).Elem().Name(), message)
}

type exception interface {
	setErrorMessage(string)
	error
}

// NewException creates a new exception with the given message template and
//uses fmt.Sprintf to interpolate the args into messageinto message.
func NewException(message string, args ...interface{}) *Exception {
	return &Exception{message: fmt.Sprintf(message, args...)}
}

// Exception represents a basic exception
type Exception struct {
	message string
}

// Type returns the type of the RubyObject
func (e *Exception) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *Exception) Inspect() string { return formatException(e, e.message) }
func (e *Exception) Error() string   { return e.message }

func (e *Exception) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns exceptionClass
func (e *Exception) Class() RubyClass { return exceptionClass }

var exceptionClassMethods = map[string]RubyMethod{
	"exception": publicMethod(exceptionClassException),
}

var exceptionMethods = map[string]RubyMethod{
	"initialize": privateMethod(exceptionInitialize),
	"exception":  publicMethod(exceptionException),
	"to_s":       withArity(0, publicMethod(exceptionToS)),
}

func exceptionInitialize(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	if self, ok := receiver.(*Self); ok {
		receiver = self.RubyObject
	}
	var message string
	message = receiver.Class().Name()
	if len(args) == 1 {
		msg, err := stringify(args[0])
		if err != nil {
			return nil, err
		}
		message = msg
	}
	if exception, ok := receiver.(exception); ok {
		exception.setErrorMessage(message)
	}
	return receiver, nil
}

func exceptionClassException(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	var message string
	class, ok := receiver.(RubyClass)
	if ok {
		receiver, _ = class.New()
	}
	if class == nil {
		class = receiver.Class()
	}
	message = class.Name()
	if len(args) == 1 {
		msg, err := stringify(args[0])
		if err != nil {
			return nil, err
		}
		message = msg
	}
	if exception, ok := receiver.(exception); ok {
		msg := exception.Error()
		if msg != message {
			exception.setErrorMessage(message)
		}
	}
	return receiver, nil
}

func exceptionException(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	if len(args) == 0 {
		return receiver, nil
	}
	var oldMessage string
	if err, ok := receiver.(error); ok {
		oldMessage = err.Error()
	}
	message, err := stringify(args[0])
	if err != nil {
		return nil, err
	}

	if oldMessage != message {
		class := receiver.Class()
		exc, err := class.New()
		if err != nil {
			return nil, err
		}
		if exception, ok := exc.(exception); ok {
			exception.setErrorMessage(message)
		}
		return exc, nil
	}
	return receiver, nil
}

func exceptionToS(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	if err, ok := receiver.(exception); ok {
		return &String{Value: err.Error()}, nil
	}
	return nil, nil
}

// NewStandardError returns a StandardError with the given message
func NewStandardError(message string) *StandardError {
	return &StandardError{message: message}
}

// StandardError is the default class for rescue blocks
type StandardError struct {
	message string
}

// Type returns EXCEPTION_OBJ
func (e *StandardError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *StandardError) Inspect() string { return formatException(e, e.message) }
func (e *StandardError) Error() string   { return e.message }

func (e *StandardError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns standardErrorClass
func (e *StandardError) Class() RubyClass { return standardErrorClass }

// NewRuntimeError returns a new RuntimeError with the formatted message
func NewRuntimeError(format string, args ...interface{}) *RuntimeError {
	return &RuntimeError{
		message: fmt.Sprintf(format, args...),
	}
}

// RuntimeError is a generic error class raised when an invalid operation is attempted
type RuntimeError struct {
	message string
}

// Type returns EXCEPTION_OBJ
func (e *RuntimeError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *RuntimeError) Inspect() string { return formatException(e, e.message) }
func (e *RuntimeError) Error() string   { return e.message }

func (e *RuntimeError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns runtimeErrorClass
func (e *RuntimeError) Class() RubyClass { return runtimeErrorClass }

// NewZeroDivisionError returns a new ZeroDivisionError with the default message
func NewZeroDivisionError() *ZeroDivisionError {
	return &ZeroDivisionError{
		message: "divided by 0",
	}
}

// ZeroDivisionError represents an arithmethic error when dividing through 0
type ZeroDivisionError struct {
	message string
}

// Type returns EXCEPTION_OBJ
func (e *ZeroDivisionError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *ZeroDivisionError) Inspect() string { return formatException(e, e.message) }
func (e *ZeroDivisionError) Error() string   { return e.message }

func (e *ZeroDivisionError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns zeroDivisionErrorClass
func (e *ZeroDivisionError) Class() RubyClass { return zeroDivisionErrorClass }

// NewWrongNumberOfArgumentsError returns an ArgumentError populated with the default message
func NewWrongNumberOfArgumentsError(expected, actual int) *ArgumentError {
	return &ArgumentError{
		message: fmt.Sprintf(
			"wrong number of arguments (given %d, expected %d)",
			actual,
			expected,
		),
	}
}

// NewArgumentError creates an ArgumentError. It has the same API as fmt.Errorf
func NewArgumentError(format string, args ...interface{}) *ArgumentError {
	return &ArgumentError{
		message: fmt.Sprintf(format, args...),
	}
}

// ArgumentError represents an error in method call arguments
type ArgumentError struct {
	message string
}

// Type returns EXCEPTION_OBJ
func (e *ArgumentError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *ArgumentError) Inspect() string { return formatException(e, e.message) }
func (e *ArgumentError) Error() string   { return e.message }

func (e *ArgumentError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns argumentErrorClass
func (e *ArgumentError) Class() RubyClass { return argumentErrorClass }

// NewUninitializedConstantNameError returns a NameError with the default message for uninitialized constants
func NewUninitializedConstantNameError(name string) *NameError {
	return &NameError{
		message: fmt.Sprintf(
			"uninitialized constant %s",
			name,
		),
	}
}

// NewUndefinedLocalVariableOrMethodNameError returns a NameError with the default message for undefined names
func NewUndefinedLocalVariableOrMethodNameError(context RubyObject, name string) *NameError {
	return &NameError{
		message: fmt.Sprintf(
			"undefined local variable or method `%s' for %s:%s",
			name,
			context.Inspect(),
			context.Class().(RubyObject).Inspect(),
		),
	}
}

// A NameError represents an error accessing an identifier unknown to the environment
type NameError struct {
	message string
}

// Type returns EXCEPTION_OBJ
func (e *NameError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *NameError) Inspect() string { return formatException(e, e.message) }
func (e *NameError) Error() string   { return e.message }

func (e *NameError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns nameErrorClass
func (e *NameError) Class() RubyClass { return nameErrorClass }

// NewNoMethodError returns a NoMethodError with the default message for undefined methods
func NewNoMethodError(context RubyObject, method string) *NoMethodError {
	return &NoMethodError{
		message: fmt.Sprintf(
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
		message: fmt.Sprintf(
			"private method `%s' called for %s:%s",
			method,
			context.Inspect(),
			context.Class().(RubyObject).Inspect(),
		),
	}
}

// NoMethodError represents an error finding a fitting method on an object
type NoMethodError struct {
	message string
}

// Type returns EXCEPTION_OBJ
func (e *NoMethodError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *NoMethodError) Inspect() string { return formatException(e, e.message) }
func (e *NoMethodError) Error() string   { return e.message }

func (e *NoMethodError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns noMethodErrorClass
func (e *NoMethodError) Class() RubyClass { return noMethodErrorClass }

// NewWrongArgumentTypeError returns a TypeError with the default message for wrong arugument type errors
func NewWrongArgumentTypeError(expected, actual RubyObject) *TypeError {
	return &TypeError{
		message: fmt.Sprintf(
			"wrong argument type %s (expected %s)",
			reflect.TypeOf(actual).Elem().Name(),
			reflect.TypeOf(expected).Elem().Name(),
		),
	}
}

// NewCoercionTypeError returns a TypeError with the default message for coercing errors
func NewCoercionTypeError(expected, actual RubyObject) *TypeError {
	return &TypeError{
		message: fmt.Sprintf(
			"%s can't be coerced into %s",
			reflect.TypeOf(actual).Elem().Name(),
			reflect.TypeOf(expected).Elem().Name(),
		),
	}
}

// NewImplicitConversionTypeError returns a TypeError with the default message for impossible implicit conversions
func NewImplicitConversionTypeError(expected, actual RubyObject) *TypeError {
	return &TypeError{
		message: fmt.Sprintf(
			"no implicit conversion of %s into %s",
			reflect.TypeOf(actual).Elem().Name(),
			reflect.TypeOf(expected).Elem().Name(),
		),
	}
}

// NewTypeError returns a TypeError with the provided message
func NewTypeError(message string) *TypeError {
	return &TypeError{message: message}
}

// TypeError represents an error when the given type does not fit in the given context
type TypeError struct {
	message string
}

// Type returns EXCEPTION_OBJ
func (e *TypeError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *TypeError) Inspect() string { return formatException(e, e.message) }
func (e *TypeError) Error() string   { return e.message }

func (e *TypeError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns typeErrorClass
func (e *TypeError) Class() RubyClass { return typeErrorClass }

// NewScriptError returns a new script error with the provided message
func NewScriptError(format string, args ...interface{}) *ScriptError {
	return &ScriptError{message: fmt.Sprintf(format, args...)}
}

// ScriptError represetns an error in the loaded script
type ScriptError struct {
	message string
}

// Type returns EXCEPTION_OBJ
func (e *ScriptError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *ScriptError) Inspect() string { return formatException(e, e.message) }
func (e *ScriptError) Error() string   { return e.message }

func (e *ScriptError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns scriptErrorClass
func (e *ScriptError) Class() RubyClass { return scriptErrorClass }

// NewNoSuchFileLoadError returns a new LoadError with the default message
func NewNoSuchFileLoadError(filepath string) *LoadError {
	return &LoadError{
		message: fmt.Sprintf(
			"cannot load such file -- %s",
			filepath,
		),
	}
}

// LoadError represents an error while loading another file
type LoadError struct {
	message string
}

// Type returns EXCEPTION_OBJ
func (e *LoadError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *LoadError) Inspect() string { return formatException(e, e.message) }
func (e *LoadError) Error() string   { return e.message }

func (e *LoadError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns loadErrorClass
func (e *LoadError) Class() RubyClass { return loadErrorClass }

// NewSyntaxError returns a new SyntaxError with the default message
func NewSyntaxError(syntaxError error) *SyntaxError {
	return &SyntaxError{
		message: fmt.Sprintf(
			"syntax error, %s",
			syntaxError.Error(),
		),
		err: syntaxError,
	}
}

// SyntaxError represents a syntax error in the ruby scripts
type SyntaxError struct {
	err     error
	message string
}

// Type returns EXCEPTION_OBJ
func (e *SyntaxError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *SyntaxError) Inspect() string { return formatException(e, e.message) }
func (e *SyntaxError) Error() string   { return e.message }

func (e *SyntaxError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns syntaxErrorClass
func (e *SyntaxError) Class() RubyClass { return syntaxErrorClass }

// UnderlyingError returns the parser error wrapped by SyntaxError
func (e *SyntaxError) UnderlyingError() error { return e.err }

// NewNotImplementedError returns a NotImplementedError with the provided message
func NewNotImplementedError(format string, args ...interface{}) *NotImplementedError {
	return &NotImplementedError{message: fmt.Sprintf(format, args...)}
}

// NotImplementedError represents an error for a not implemented feature on a given platform
type NotImplementedError struct {
	message string
}

// Type returns EXCEPTION_OBJ
func (e *NotImplementedError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *NotImplementedError) Inspect() string { return formatException(e, e.message) }
func (e *NotImplementedError) Error() string   { return e.message }

func (e *NotImplementedError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns notImplementedErrorClass
func (e *NotImplementedError) Class() RubyClass { return notImplementedErrorClass }

// NewNoBlockGivenLocalJumpError returns a LocalJumpError with the default message for missing blocks
func NewNoBlockGivenLocalJumpError() *LocalJumpError {
	return &LocalJumpError{message: "no block given (yield)"}
}

// LocalJumpError represents an error for a not supported jump
type LocalJumpError struct {
	message string
}

// Type returns EXCEPTION_OBJ
func (e *LocalJumpError) Type() Type { return EXCEPTION_OBJ }

// Inspect returns a string starting with the exception class name, followed by the message
func (e *LocalJumpError) Inspect() string { return formatException(e, e.message) }
func (e *LocalJumpError) Error() string   { return e.message }

func (e *LocalJumpError) setErrorMessage(msg string) {
	e.message = msg
}

// Class returns notImplementedErrorClass
func (e *LocalJumpError) Class() RubyClass { return notImplementedErrorClass }
