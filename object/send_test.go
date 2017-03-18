package object

import (
	"reflect"
	"testing"
)

type testRubyObject struct {
	class RubyClassObject
}

func (t *testRubyObject) Type() ObjectType { return ObjectType("TEST_OBJECT") }
func (t *testRubyObject) Inspect() string  { return "TEST OBJECT" }
func (t *testRubyObject) Class() RubyClass {
	if t.class != nil {
		return t.class
	}
	return OBJECT_CLASS
}

func TestSend(t *testing.T) {
	methods := map[string]RubyMethod{
		"a_method": publicMethod(func(context RubyObject, args ...RubyObject) RubyObject {
			return TRUE
		}),
		"another_method": publicMethod(func(context RubyObject, args ...RubyObject) RubyObject {
			return FALSE
		}),
	}
	context := &testRubyObject{class: &Class{instanceMethods: methods, superClass: BASIC_OBJECT_CLASS}}

	tests := []struct {
		method         string
		expectedResult RubyObject
	}{
		{
			"a_method",
			TRUE,
		},
		{
			"another_method",
			FALSE,
		},
		{
			"unknown_method",
			NewNoMethodError(context, "unknown_method"),
		},
	}

	for _, testCase := range tests {
		result := Send(context, testCase.method)

		if !reflect.DeepEqual(result, testCase.expectedResult) {
			t.Logf("Expected result to equal\n%+#v\n\tgot\n%+#v\n", testCase.expectedResult, result)
			t.Fail()
		}
	}
}
