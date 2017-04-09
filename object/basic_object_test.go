package object

import (
	"testing"
)

func TestBasicObjectMethodMissing(t *testing.T) {
	result, err := basicObjectMethodMissing(NIL, &Symbol{"foo"})

	checkResult(t, result, nil)

	expected := NewNoMethodError(NIL, "foo")

	checkError(t, err, expected)
}
