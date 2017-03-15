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
