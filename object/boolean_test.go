package object

import "testing"

func TestBoolean_hashKey(t *testing.T) {
	hello1 := &Boolean{Value: true}
	hello2 := &Boolean{Value: true}
	diff1 := &Boolean{Value: false}
	diff2 := &Boolean{Value: false}

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

func TestBooleanEq(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{FALSE},
			FALSE,
			nil,
		},
		{
			[]RubyObject{TRUE},
			TRUE,
			nil,
		},
		{
			[]RubyObject{&String{""}},
			FALSE,
			nil,
		},
		{
			[]RubyObject{&Integer{Value: 0}},
			FALSE,
			nil,
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: TRUE}

		result, err := booleanEq(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func TestBooleanNeq(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{FALSE},
			TRUE,
			nil,
		},
		{
			[]RubyObject{TRUE},
			FALSE,
			nil,
		},
		{
			[]RubyObject{&String{""}},
			TRUE,
			nil,
		},
		{
			[]RubyObject{&Integer{Value: 0}},
			TRUE,
			nil,
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: TRUE}

		result, err := booleanNeq(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}
