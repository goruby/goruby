package object

import (
	"reflect"
	"testing"
)

func TestInteger_hashKey(t *testing.T) {
	hello1 := &Integer{Value: 1}
	hello2 := &Integer{Value: 1}
	diff1 := &Integer{Value: 3}
	diff2 := &Integer{Value: 3}

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

func TestIntegerSub(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(3)},
			NewInteger(1),
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

		result, err := integerSub(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func TestIntegerModulo(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(3)},
			NewInteger(1),
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

		result, err := integerModulo(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func TestIntegerLt(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewInteger(2)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerLt(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func TestIntegerGt(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{NewInteger(2)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerGt(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func TestIntegerEq(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{NewInteger(4)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerEq(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func TestIntegerNeq(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewInteger(4)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerNeq(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func TestIntegerGte(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{NewInteger(4)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewInteger(2)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerGte(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func TestIntegerLte(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewInteger(4)},
			TRUE,
			nil,
		},
		{
			[]RubyObject{NewInteger(2)},
			FALSE,
			nil,
		},
		{
			[]RubyObject{&String{""}},
			nil,
			NewArgumentError("comparison of Integer with String failed"),
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerLte(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func TestIntegerSpaceship(t *testing.T) {
	tests := []struct {
		arguments []RubyObject
		result    RubyObject
		err       error
	}{
		{
			[]RubyObject{NewInteger(6)},
			&Integer{Value: -1},
			nil,
		},
		{
			[]RubyObject{NewInteger(4)},
			&Integer{Value: 0},
			nil,
		},
		{
			[]RubyObject{NewInteger(2)},
			&Integer{Value: 1},
			nil,
		},
		{
			[]RubyObject{&String{""}},
			NIL,
			nil,
		},
	}

	for _, testCase := range tests {
		context := &callContext{receiver: NewInteger(4)}

		result, err := integerSpaceship(context, testCase.arguments...)

		checkError(t, err, testCase.err)

		checkResult(t, result, testCase.result)
	}
}

func checkError(t *testing.T, actual, expected error) {
	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected error to equal\n%T:%v\n\tgot\n%T:%v\n", expected, expected, actual, actual)
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
