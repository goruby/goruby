package object

import (
	"reflect"
	"testing"
)

func TestWithArity(t *testing.T) {
	wrappedMethod := publicMethod(func(context RubyObject, args ...RubyObject) RubyObject {
		return NewInteger(1)
	})

	tests := []struct {
		arity     int
		arguments []RubyObject
		result    RubyObject
	}{
		{
			1,
			[]RubyObject{NIL},
			NewInteger(1),
		},
		{
			1,
			[]RubyObject{NIL, NIL},
			NewWrongNumberOfArgumentsError(1, 2),
		},
		{
			2,
			[]RubyObject{NIL},
			NewWrongNumberOfArgumentsError(2, 1),
		},
	}

	for _, testCase := range tests {
		fn := withArity(testCase.arity, wrappedMethod)

		result := fn.Call(NIL, testCase.arguments...)

		if !reflect.DeepEqual(result, testCase.result) {
			t.Logf(
				"Expected result to equal\n%s\n\tgot\n%s\n",
				testCase.result.Inspect(),
				result.Inspect(),
			)
			t.Fail()
		}
	}
}
