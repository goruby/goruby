package object

import (
	"reflect"
	"testing"
)

type testRubyObject struct {
	methods map[string]method
}

func (t *testRubyObject) Type() ObjectType           { return ObjectType("TEST_OBJECT") }
func (t *testRubyObject) Inspect() string            { return "TEST OBJECT" }
func (t *testRubyObject) Class() RubyClass           { return t }
func (t *testRubyObject) SuperClass() RubyClass      { return BASIC_OBJECT_CLASS }
func (t *testRubyObject) Methods() map[string]method { return t.methods }

func TestSend(t *testing.T) {
	methods := map[string]method{
		"a_method": func(context RubyObject, args ...RubyObject) RubyObject {
			return TRUE
		},
		"another_method": func(context RubyObject, args ...RubyObject) RubyObject {
			return FALSE
		},
	}
	context := &testRubyObject{methods}

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
