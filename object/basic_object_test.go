package object

import (
	"testing"
)

func TestBasicObjectMethodMissing(t *testing.T) {
	context := &callContext{receiver: NIL}
	result, err := basicObjectMethodMissing(context, &Symbol{"foo"})

	checkResult(t, result, nil)

	expected := NewNoMethodError(NIL, "foo")

	checkError(t, err, expected)
}

func TestBasicObjectNegate(t *testing.T) {
	t.Run("on basic object", func(t *testing.T) {
		context := &callContext{receiver: &basicObject{}}

		result, err := basicObjectNegate(context)

		checkError(t, err, nil)

		expected := FALSE

		checkResult(t, result, expected)
	})
	t.Run("on nil", func(t *testing.T) {
		context := &callContext{receiver: NIL}

		result, err := basicObjectNegate(context)

		checkError(t, err, nil)

		expected := TRUE

		checkResult(t, result, expected)
	})
	t.Run("on false", func(t *testing.T) {
		context := &callContext{receiver: FALSE}

		result, err := basicObjectNegate(context)

		checkError(t, err, nil)

		expected := TRUE

		checkResult(t, result, expected)
	})
}
