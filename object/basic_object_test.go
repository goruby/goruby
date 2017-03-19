package object

import (
	"reflect"
	"testing"
)

func TestBasicObjectMethodMissing(t *testing.T) {
	result := basicObjectMethodMissing(NIL, &Symbol{"foo"})

	expected := NewNoMethodError(NIL, "foo")

	if !reflect.DeepEqual(expected, result) {
		t.Logf("Expected result to equal\n%+#v\n\tgot\n%+#v\n", expected, result)
		t.Fail()
	}
}
