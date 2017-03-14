package object

import (
	"reflect"
	"testing"
)

func TestMethodSetDefine(t *testing.T) {
	set := methodSet{context: NewInteger(3), methods: make(map[string]method)}

	sym := set.Define("foo", nil)

	expectedSymbol := &Symbol{"foo"}

	if !reflect.DeepEqual(expectedSymbol, sym) {
		t.Logf("Expected symbol to equal %s, got %s\n", expectedSymbol, sym)
		t.Fail()
	}

	expectedMethodsLength := 1
	actualMethodsLength := len(set.methods)
	if actualMethodsLength != expectedMethodsLength {
		t.Logf("Expected methods to have %d items, got %d\n", expectedMethodsLength, actualMethodsLength)
		t.Fail()
	}
}

func TestWithArity(t *testing.T) {
	wrappedMethod := func(context RubyObject, args ...RubyObject) RubyObject {
		return NewInteger(1)
	}

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

		result := fn(NIL, testCase.arguments...)

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
