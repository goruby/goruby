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
