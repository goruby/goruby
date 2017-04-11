package object

import (
	"reflect"
	"sort"
	"testing"
)

func TestKernelMethods(t *testing.T) {
	t.Run("without superclass", func(t *testing.T) {
		contextMethods := map[string]RubyMethod{
			"foo": publicMethod(nil),
			"bar": publicMethod(nil),
		}
		context := &callContext{
			receiver: &testRubyObject{
				class: &class{
					instanceMethods: contextMethods,
					superClass:      nil,
				},
			},
		}

		result, err := kernelMethods(context)

		checkError(t, err, nil)

		array, ok := result.(*Array)
		if !ok {
			t.Logf("Expected array, got %T", result)
			t.FailNow()
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

		var expectedMethods = []string{
			":foo", ":bar",
		}

		expectedLen := len(expectedMethods)

		if len(array.Elements) != expectedLen {
			t.Logf("Expected %d items, got %d", expectedLen, len(array.Elements))
			t.Fail()
		}

		sort.Strings(expectedMethods)
		sort.Strings(methods)

		if !reflect.DeepEqual(expectedMethods, methods) {
			t.Logf("Expected methods to equal\n%s\n\tgot\n%s\n", expectedMethods, methods)
			t.Fail()
		}
	})
	t.Run("with superclass", func(t *testing.T) {
		superClassMethods := map[string]RubyMethod{
			"super_foo": publicMethod(nil),
			"super_bar": publicMethod(nil),
		}
		contextMethods := map[string]RubyMethod{
			"foo": publicMethod(nil),
			"bar": publicMethod(nil),
		}
		context := &callContext{
			receiver: &testRubyObject{
				class: &class{
					instanceMethods: contextMethods,
					superClass: &class{
						instanceMethods: superClassMethods,
						superClass:      nil,
					},
				},
			},
		}

		result, err := kernelMethods(context)

		checkError(t, err, nil)

		array, ok := result.(*Array)
		if !ok {
			t.Logf("Expected array, got %T", result)
			t.FailNow()
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

		var expectedMethods = []string{
			":foo", ":bar", ":super_foo", ":super_bar",
		}

		expectedLen := len(expectedMethods)

		if len(array.Elements) != expectedLen {
			t.Logf("Expected %d items, got %d", expectedLen, len(array.Elements))
			t.Fail()
		}

		sort.Strings(expectedMethods)
		sort.Strings(methods)

		if !reflect.DeepEqual(expectedMethods, methods) {
			t.Logf("Expected methods to equal\n%s\n\tgot\n%s\n", expectedMethods, methods)
			t.Fail()
		}
	})
	t.Run("with private methods", func(t *testing.T) {
		contextMethods := map[string]RubyMethod{
			"foo":         publicMethod(nil),
			"bar":         publicMethod(nil),
			"private_foo": privateMethod(nil),
		}
		context := &callContext{
			receiver: &testRubyObject{
				class: &class{
					instanceMethods: contextMethods,
					superClass:      nil,
				},
			},
		}

		result, err := kernelMethods(context)

		checkError(t, err, nil)

		array, ok := result.(*Array)
		if !ok {
			t.Logf("Expected array, got %T", result)
			t.FailNow()
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

		var expectedMethods = []string{
			":foo", ":bar",
		}

		expectedLen := len(expectedMethods)

		if len(array.Elements) != expectedLen {
			t.Logf("Expected %d items, got %d", expectedLen, len(array.Elements))
			t.Fail()
		}

		sort.Strings(expectedMethods)
		sort.Strings(methods)

		if !reflect.DeepEqual(expectedMethods, methods) {
			t.Logf("Expected methods to equal\n%s\n\tgot\n%s\n", expectedMethods, methods)
			t.Fail()
		}
	})
}

func TestKernelIsNil(t *testing.T) {
	context := &callContext{receiver: TRUE}
	result, err := kernelIsNil(context)

	checkError(t, err, nil)

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

func TestKernelClass(t *testing.T) {
	t.Run("regular object", func(t *testing.T) {
		context := &callContext{receiver: &Integer{1}}

		result, err := kernelClass(context)

		checkError(t, err, nil)

		cl, ok := result.(*class)
		if !ok {
			t.Logf("Expected Class, got %T", result)
			t.Fail()
		}

		expected := integerClass

		if !reflect.DeepEqual(expected, cl) {
			t.Logf("Expected class to equal %+#v, got %+#v", expected, cl)
			t.Fail()
		}
	})
	t.Run("class object", func(t *testing.T) {
		context := &callContext{receiver: stringClass}

		result, err := kernelClass(context)

		checkError(t, err, nil)

		cl, ok := result.(*class)
		if !ok {
			t.Logf("Expected Class, got %T", result)
			t.Fail()
		}

		expected := classClass

		if !reflect.DeepEqual(expected, cl) {
			t.Logf("Expected class to equal %+#v, got %+#v", expected, cl)
			t.Fail()
		}
	})
	t.Run("class class", func(t *testing.T) {
		context := &callContext{receiver: classClass}

		result, err := kernelClass(context)

		checkError(t, err, nil)

		cl, ok := result.(*class)
		if !ok {
			t.Logf("Expected Class, got %T", result)
			t.Fail()
		}

		expected := classClass

		if !reflect.DeepEqual(expected, cl) {
			t.Logf("Expected class to equal %+#v, got %+#v", expected, cl)
			t.Fail()
		}
	})
}
