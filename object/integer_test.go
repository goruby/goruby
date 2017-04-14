package object

import (
	"reflect"
	"testing"
)

func TestIntegerDiv(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(2)},
			NewInteger(2),
			nil,
		},
		{
			[]RubyObject{&String{""}},
			nil,
			NewCoercionTypeError(&String{}, &Integer{}),
		},
		{
			[]RubyObject{NewInteger(0)},
			nil,
			NewZeroDivisionError(),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerDiv(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func TestIntegerMul(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(2)},
			NewInteger(8),
			nil,
		},
		{
			[]RubyObject{&String{""}},
			nil,
			NewCoercionTypeError(&String{}, &Integer{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerMul(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func TestIntegerAdd(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(2)},
			NewInteger(4),
			nil,
		},
		{
			[]RubyObject{&String{""}},
			nil,
			NewCoercionTypeError(&String{}, &Integer{}),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(2)}

		result, err := integerAdd(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func checkError(t *testing.T, actual, expected error) {
	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected error to equal %T:%v, got %T:%v\n", expected, expected, actual, actual)
		t.Fail()
	}
}

func checkResult(t *testing.T, actual, expected RubyObject) {
	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected result to equal %s, got %s\n", toString(expected), toString(actual))
		t.Fail()
	}
}

func toString(obj RubyObject) string {
	if obj == nil {
		return "nil"
	}
	return obj.Inspect()
}
