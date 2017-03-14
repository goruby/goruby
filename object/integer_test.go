package object

import (
	"testing"
)

func TestIntegerDiv(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
	}{
		{
			[]RubyObject{NewInteger(2)},
			NewInteger(2),
		},
		{
			[]RubyObject{&String{""}},
			NewCoercionTypeError(&String{}, &Integer{}),
		},
		{
			[]RubyObject{NewInteger(0)},
			NewZeroDivisionError(),
		},
	}

	for _, testCase := range tests {
		context := NewInteger(4)

		result := integerDiv(context, testCase.arguments...)

		if result.Inspect() != testCase.result.Inspect() {
			t.Logf(
				"Expected result to equal\n%s\n\tgot\n%s\n",
				testCase.result.Inspect(),
				result.Inspect(),
			)
			t.Fail()
		}
	}
}

func TestIntegerMul(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
	}{
		{
			[]RubyObject{NewInteger(2)},
			NewInteger(8),
		},
		{
			[]RubyObject{&String{""}},
			NewCoercionTypeError(&String{}, &Integer{}),
		},
	}

	for _, testCase := range tests {
		context := NewInteger(4)

		result := integerMul(context, testCase.arguments...)

		if result.Inspect() != testCase.result.Inspect() {
			t.Logf(
				"Expected result to equal\n%s\n\tgot\n%s\n",
				testCase.result.Inspect(),
				result.Inspect(),
			)
			t.Fail()
		}
	}
}

func TestIntegerAdd(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
	}{
		{
			[]RubyObject{NewInteger(2)},
			NewInteger(4),
		},
		{
			[]RubyObject{&String{""}},
			NewCoercionTypeError(&String{}, &Integer{}),
		},
	}

	for _, testCase := range tests {
		context := NewInteger(2)

		result := integerAdd(context, testCase.arguments...)

		if result.Inspect() != testCase.result.Inspect() {
			t.Logf(
				"Expected result to equal\n%s\n\tgot\n%s\n",
				testCase.result.Inspect(),
				result.Inspect(),
			)
			t.Fail()
		}
	}
}
