package object

import (
	"reflect"
	"sort"
	"testing"
)

func TestObjMethods(t *testing.T) {
	contextMethods := map[string]method{
		"foo": nil,
		"bar": nil,
	}
	context := &testRubyObject{
		methods:    contextMethods,
		superClass: OBJECT_CLASS,
	}

	result := objMethods(context)

	array, ok := result.(*Array)
	if !ok {
		t.Logf("Expected array, got %T", result)
		t.FailNow()
	}

	expectedLen := len(contextMethods) + len(objectMethods) + len(basicObjectMethods)

	if len(array.Elements) != expectedLen {
		t.Logf("Expected %d items, got %d", expectedLen, len(array.Elements))
		t.Fail()
	}

	var methods []string
	for i, elem := range array.Elements {
		sym, ok := elem.(*Symbol)
		if !ok {
			t.Logf("Expected all elements to be symbols, got %T at index %d", elem, i)
			t.Fail()
		} else {
			methods = append(methods, sym.Inspect())
		}
	}

	var expectedMethods []string
	for k, _ := range contextMethods {
		expectedMethods = append(expectedMethods, ":"+k)
	}
	for k, _ := range basicObjectMethods {
		expectedMethods = append(expectedMethods, ":"+k)
	}
	for k, _ := range objectMethods {
		expectedMethods = append(expectedMethods, ":"+k)
	}

	sort.Strings(expectedMethods)
	sort.Strings(methods)

	if !reflect.DeepEqual(expectedMethods, methods) {
		t.Logf("Expected methods to equal\n%s\n\tgot\n%s\n", expectedMethods, methods)
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
