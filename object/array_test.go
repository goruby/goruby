package object

import (
	"reflect"
	"testing"
)

func TestArrayPush(t *testing.T) {
	t.Run("one argument", func(t *testing.T) {
		array := &Array{}
		env := NewEnvironment()
		context := &callContext{
			receiver: array,
			env:      env,
		}

		result, err := arrayPush(context, &Integer{Value: 17})

		checkError(t, err, nil)

		expectedResult := &Array{Elements: []RubyObject{&Integer{Value: 17}}}

		checkResult(t, result, expectedResult)

		if !reflect.DeepEqual(expectedResult, array) {
			t.Logf("Expected array to equal\n%+#v\n\tgot\n%+#v\n", expectedResult, array)
			t.Fail()
		}
	})
	t.Run("more than one argument", func(t *testing.T) {
		array := &Array{}
		env := NewEnvironment()
		context := &callContext{
			receiver: array,
			env:      env,
		}

		result, err := arrayPush(context, &Integer{Value: 17}, NIL, TRUE, FALSE)

		checkError(t, err, nil)

		expectedResult := &Array{Elements: []RubyObject{&Integer{Value: 17}, NIL, TRUE, FALSE}}

		checkResult(t, result, expectedResult)

		if !reflect.DeepEqual(expectedResult, array) {
			t.Logf("Expected array to equal\n%+#v\n\tgot\n%+#v\n", expectedResult, array)
			t.Fail()
		}
	})
}

func TestArrayUnshift(t *testing.T) {
	t.Run("one argument", func(t *testing.T) {
		array := &Array{Elements: []RubyObject{&String{Value: "first element"}}}
		env := NewEnvironment()
		context := &callContext{
			receiver: array,
			env:      env,
		}

		result, err := arrayUnshift(context, &Integer{Value: 17})

		checkError(t, err, nil)

		expectedResult := &Array{Elements: []RubyObject{&Integer{Value: 17}, &String{Value: "first element"}}}

		checkResult(t, result, expectedResult)

		if !reflect.DeepEqual(expectedResult, array) {
			t.Logf("Expected array to equal\n%+#v\n\tgot\n%+#v\n", expectedResult, array)
			t.Fail()
		}
	})
	t.Run("more than one argument", func(t *testing.T) {
		array := &Array{Elements: []RubyObject{&String{Value: "first element"}}}
		env := NewEnvironment()
		context := &callContext{
			receiver: array,
			env:      env,
		}

		result, err := arrayUnshift(context, &Integer{Value: 17}, NIL, TRUE, FALSE)

		checkError(t, err, nil)

		expectedResult := &Array{Elements: []RubyObject{&Integer{Value: 17}, NIL, TRUE, FALSE, &String{Value: "first element"}}}

		checkResult(t, result, expectedResult)

		if !reflect.DeepEqual(expectedResult, array) {
			t.Logf("Expected array to equal\n%+#v\n\tgot\n%+#v\n", expectedResult, array)
			t.Fail()
		}
	})
}
