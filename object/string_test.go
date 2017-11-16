package object

import "testing"

func TestString_hashKey(t *testing.T) {
	hello1 := &String{Value: "Hello World"}
	hello2 := &String{Value: "Hello World"}
	diff1 := &String{Value: "My name is johnny"}
	diff2 := &String{Value: "My name is johnny"}

	if hello1.hashKey() != hello2.hashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.hashKey() != diff2.hashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.hashKey() == diff1.hashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func Test_stringify(t *testing.T) {
	t.Run("object with regular `to_s`", func(t *testing.T) {
		obj := &Symbol{Value: "sym"}

		res, err := stringify(obj)

		checkError(t, err, nil)

		if res != "sym" {
			t.Logf("Expected stringify to return 'sym', got %q\n", res)
			t.Fail()
		}
	})
	t.Run("object with `to_s` returning not string", func(t *testing.T) {
		toS := func(CallContext, ...RubyObject) (RubyObject, error) {
			return &Integer{Value: 42}, nil
		}
		obj := &extendedObject{
			RubyObject: &Object{},
			class: newEigenclass(objectClass, map[string]RubyMethod{
				"to_s": publicMethod(toS),
			}),
			Environment: NewEnvironment(),
		}

		_, err := stringify(obj)

		checkError(t, err, NewTypeError(
			"can't convert Object to String (Object#to_s gives Integer)",
		))
	})
	t.Run("object without `to_s`", func(t *testing.T) {
		obj := &basicObject{}

		_, err := stringify(obj)

		checkError(t, err, NewTypeError("can't convert BasicObject into String"))
	})
}

func TestStringAdd(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{&String{Value: " bar"}},
			&String{Value: "foo bar"},
			nil,
		},
		{
			[]RubyObject{&Integer{Value: 3}},
			nil,
			NewImplicitConversionTypeError(&String{}, &Integer{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: &String{Value: "foo"}}

		result, err := stringAdd(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}
