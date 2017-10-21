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
