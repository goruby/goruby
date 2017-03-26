package object

import (
	"reflect"
	"testing"
)

type testRubyObject struct {
	class RubyClassObject
}

func (t *testRubyObject) Type() Type      { return Type("TEST_OBJECT") }
func (t *testRubyObject) Inspect() string { return "TEST OBJECT" }
func (t *testRubyObject) Class() RubyClass {
	if t.class != nil {
		return t.class
	}
	return objectClass
}

func TestSend(t *testing.T) {
	superMethods := map[string]RubyMethod{
		"a_super_method": publicMethod(func(context RubyObject, args ...RubyObject) RubyObject {
			return TRUE
		}),
		"a_private_super_method": privateMethod(func(context RubyObject, args ...RubyObject) RubyObject {
			return FALSE
		}),
	}
	methods := map[string]RubyMethod{
		"a_method": publicMethod(func(context RubyObject, args ...RubyObject) RubyObject {
			return TRUE
		}),
		"another_method": publicMethod(func(context RubyObject, args ...RubyObject) RubyObject {
			return FALSE
		}),
		"a_private_method": privateMethod(func(context RubyObject, args ...RubyObject) RubyObject {
			return FALSE
		}),
	}
	context := &testRubyObject{
		class: &class{
			name:            "base class",
			instanceMethods: methods,
			superClass: &class{
				name:            "super class",
				instanceMethods: superMethods,
				superClass:      basicObjectClass,
			},
		},
	}

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
			"a_super_method",
			TRUE,
		},
		{
			"a_private_method",
			NewPrivateNoMethodError(context, "a_private_method"),
		},
		{
			"a_private_super_method",
			NewPrivateNoMethodError(context, "a_private_super_method"),
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
