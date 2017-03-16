package object

import (
	"reflect"
	"testing"
)

func TestObjMethods(t *testing.T) {
	context := &testRubyObject{
		methods: map[string]method{
			"foo": nil,
			"bar": nil,
		},
	}

	result := objMethods(context)

	expected := NIL

	if !reflect.DeepEqual(expected, result) {
		t.Logf("Expected result to equal\n%+#v\n\tgot\n%+#v\n", expected, result)
		t.Fail()
	}
}

func TestObjectIsNil(t *testing.T) {
	result := objectIsNil(TRUE)

	boolean, ok := result.(*Boolean)
	if !ok {
		t.Logf("Expected Boolean, got %T", result)
		t.FailNow()
	}

	if boolean.Value != false {
		t.Logf("Expected false, got true")
		t.Fail()
	}
}
