package object

import (
	"testing"
)

func TestWithArity(t *testing.T) {
	wrappedMethod := publicMethod(func(context RubyObject, args ...RubyObject) (RubyObject, error) {
		return NewInteger(1), nil
	})

	tests := []struct {
		arity     int
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			1,
			[]RubyObject{NIL},
			NewInteger(1),
			nil,
		},
		{
			1,
			[]RubyObject{NIL, NIL},
			nil,
			NewWrongNumberOfArgumentsError(1, 2),
		},
		{
			2,
			[]RubyObject{NIL},
			nil,
			NewWrongNumberOfArgumentsError(2, 1),
		},
	}

	for _, testCase := range tests {
		fn := withArity(testCase.arity, wrappedMethod)

		result, err := fn.Call(NIL, testCase.arguments...)

		checkResult(t, result, testCase.result)

		checkError(t, err, testCase.err)
	}
}
