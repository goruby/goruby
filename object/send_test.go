package object

import (
	"reflect"
	"testing"

	"github.com/goruby/goruby/ast"
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
		"a_super_method": publicMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
			return TRUE, nil
		}),
		"a_private_super_method": privateMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
			return FALSE, nil
		}),
	}
	methods := map[string]RubyMethod{
		"a_method": publicMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
			return TRUE, nil
		}),
		"another_method": publicMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
			return FALSE, nil
		}),
		"a_private_method": privateMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
			return FALSE, nil
		}),
	}
	t.Run("normal object as context", func(t *testing.T) {
		context := &callContext{
			receiver: &testRubyObject{
				class: &class{
					name:            "base class",
					instanceMethods: methods,
					superClass: &class{
						name:            "super class",
						instanceMethods: superMethods,
						superClass:      basicObjectClass,
					},
				},
			},
		}

		tests := []struct {
			method         string
			expectedResult RubyObject
			expectedError  error
		}{
			{
				"a_method",
				TRUE,
				nil,
			},
			{
				"another_method",
				FALSE,
				nil,
			},
			{
				"a_super_method",
				TRUE,
				nil,
			},
			{
				"a_private_method",
				nil,
				NewPrivateNoMethodError(context.receiver, "a_private_method"),
			},
			{
				"a_private_super_method",
				nil,
				NewPrivateNoMethodError(context.receiver, "a_private_super_method"),
			},
			{
				"unknown_method",
				nil,
				NewNoMethodError(context.receiver, "unknown_method"),
			},
		}

		for _, testCase := range tests {
			result, err := Send(context, testCase.method)

			checkError(t, err, testCase.expectedError)

			checkResult(t, result, testCase.expectedResult)
		}
	})
	t.Run("self as context", func(t *testing.T) {
		context := &callContext{
			receiver: &Self{
				&testRubyObject{
					class: &class{
						name:            "base class",
						instanceMethods: methods,
						superClass: &class{
							name:            "super class",
							instanceMethods: superMethods,
							superClass:      basicObjectClass,
						},
					},
				},
			},
		}

		tests := []struct {
			method         string
			expectedResult RubyObject
			expectedError  error
		}{
			{
				"a_method",
				TRUE,
				nil,
			},
			{
				"another_method",
				FALSE,
				nil,
			},
			{
				"a_super_method",
				TRUE,
				nil,
			},
			{
				"a_private_method",
				FALSE,
				nil,
			},
			{
				"a_private_super_method",
				FALSE,
				nil,
			},
			{
				"unknown_method",
				nil,
				NewNoMethodError(context.receiver, "unknown_method"),
			},
		}

		for _, testCase := range tests {
			result, err := Send(context, testCase.method)

			if !reflect.DeepEqual(err, testCase.expectedError) {
				t.Logf("Expected err to equal\n%+#v\n\tgot\n%+#v\n", testCase.expectedError, err)
				t.Fail()
			}

			if !reflect.DeepEqual(result, testCase.expectedResult) {
				t.Logf("Expected result to equal\n%+#v\n\tgot\n%+#v\n", testCase.expectedResult, result)
				t.Fail()
			}
		}
	})
}

func TestAddMethod(t *testing.T) {
	t.Run("vanilla object", func(t *testing.T) {
		context := &testRubyObject{
			class: &class{
				name:            "base class",
				instanceMethods: map[string]RubyMethod{},
				superClass:      objectClass,
			},
		}

		fn := &Function{
			Parameters: []*ast.Identifier{
				&ast.Identifier{Value: "x"},
			},
			Env:  &environment{store: map[string]RubyObject{}},
			Body: nil,
		}

		newContext := AddMethod(context, "foo", fn)

		_, ok := newContext.Class().Methods()["foo"]
		if !ok {
			t.Logf("Expected object to have method foo")
			t.Fail()
		}
	})
	t.Run("extended object", func(t *testing.T) {
		context :=
			&extendedObject{
				RubyObject: &testRubyObject{
					class: &class{
						name:            "base class",
						instanceMethods: map[string]RubyMethod{},
						superClass:      objectClass,
					},
				},
				class: newEigenclass(objectClass, map[string]RubyMethod{
					"bar": publicMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
						return NIL, nil
					}),
				}),
			}

		fn := &Function{
			Parameters: []*ast.Identifier{
				&ast.Identifier{Value: "x"},
			},
			Env:  &environment{store: map[string]RubyObject{}},
			Body: nil,
		}

		newContext := AddMethod(context, "foo", fn)

		_, ok := newContext.Class().Methods()["foo"]
		if !ok {
			t.Logf("Expected object to have method foo")
			t.Fail()
		}

		_, ok = newContext.Class().Methods()["bar"]
		if !ok {
			t.Logf("Expected object to have method bar")
			t.Fail()
		}
	})
	t.Run("vanilla self object", func(t *testing.T) {
		context := &Self{
			RubyObject: &testRubyObject{
				class: &class{
					name:            "base class",
					instanceMethods: map[string]RubyMethod{},
					superClass:      objectClass,
				},
			},
		}

		fn := &Function{
			Parameters: []*ast.Identifier{
				&ast.Identifier{Value: "x"},
			},
			Env:  &environment{store: map[string]RubyObject{}},
			Body: nil,
		}

		newContext := AddMethod(context, "foo", fn)

		_, ok := newContext.Class().Methods()["foo"]
		if !ok {
			t.Logf("Expected object to have method foo")
			t.Fail()
		}

		returnedSelf, ok := newContext.(*Self)
		if !ok {
			t.Logf("Expected returned object to be self, got %T", newContext)
			t.Fail()
		}

		returnPointer := reflect.ValueOf(returnedSelf).Pointer()
		contextPointer := reflect.ValueOf(context).Pointer()

		if returnPointer != contextPointer {
			t.Logf("Expected input and return context to be the same")
			t.Fail()
		}
	})
	t.Run("extended self object", func(t *testing.T) {
		context := &Self{
			RubyObject: &extendedObject{
				RubyObject: &testRubyObject{
					class: &class{
						name:            "base class",
						instanceMethods: map[string]RubyMethod{},
						superClass:      objectClass,
					},
				},
				class: newEigenclass(objectClass, map[string]RubyMethod{
					"bar": publicMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
						return NIL, nil
					}),
				}),
			},
		}

		fn := &Function{
			Parameters: []*ast.Identifier{
				&ast.Identifier{Value: "x"},
			},
			Env:  &environment{store: map[string]RubyObject{}},
			Body: nil,
		}

		newContext := AddMethod(context, "foo", fn)

		_, ok := newContext.Class().Methods()["foo"]
		if !ok {
			t.Logf("Expected object to have method foo")
			t.Fail()
		}

		_, ok = newContext.Class().Methods()["bar"]
		if !ok {
			t.Logf("Expected object to have method bar")
			t.Fail()
		}

		returnedSelf, ok := newContext.(*Self)
		if !ok {
			t.Logf("Expected returned object to be self, got %T", newContext)
			t.Fail()
		}

		returnPointer := reflect.ValueOf(returnedSelf).Pointer()
		contextPointer := reflect.ValueOf(context).Pointer()

		if returnPointer != contextPointer {
			t.Logf("Expected input and return context to be the same")
			t.Fail()
		}
	})
}
