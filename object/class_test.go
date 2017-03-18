package object

import (
	"fmt"
	"reflect"
	"testing"
)

func TestClassInspect(t *testing.T) {
	t.Run("class Class", func(t *testing.T) {
		context := &Class{}

		actual := context.Inspect()

		expected := fmt.Sprintf("#<Class:%p>", context)

		if expected != actual {
			t.Logf("Expected Inspect to equal %q, got %q", expected, actual)
			t.Fail()
		}
	})
	t.Run("other class", func(t *testing.T) {
		context := &Class{name: "Foo"}

		actual := context.Inspect()

		expected := "Foo"

		if expected != actual {
			t.Logf("Expected Inspect to equal %q, got %q", expected, actual)
			t.Fail()
		}
	})
}

func TestClassClass(t *testing.T) {
	class := &Class{}

	context := &Class{class: class}

	actual := context.Class()

	expected := class

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Class to equal %v, got %v", expected, actual)
		t.Fail()
	}
}

func TestClassSuperClass(t *testing.T) {
	context := &Class{superClass: BASIC_OBJECT_CLASS}

	actual := context.SuperClass()

	expected := BASIC_OBJECT_CLASS

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected SuperClass to equal %+#v, got %+#v", expected, actual)
		t.Fail()
	}
}

func TestClassMethods(t *testing.T) {
	contextMethods := map[string]RubyMethod{
		"a_method": nil,
	}

	context := &Class{instanceMethods: contextMethods}

	actual := context.Methods()

	expected := contextMethods

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Methods to equal %v, got %v", expected, actual)
		t.Fail()
	}
}

func TestClassType(t *testing.T) {
	context := &Class{}

	actual := context.Type()

	expected := CLASS_OBJ

	if expected != actual {
		t.Logf("Expected Type to equal %v, got %v", expected, actual)
		t.Fail()
	}
}

func TestClassSuperclass(t *testing.T) {
	t.Run("anything else than BasicObject", func(t *testing.T) {
		context := &Class{superClass: OBJECT_CLASS}

		result := classSuperclass(context)

		_, ok := result.(*Class)
		if !ok {
			t.Logf("Expected Class object, got %T\n", result)
			t.Fail()
		}
	})
	t.Run("BasicObject", func(t *testing.T) {
		context := BASIC_OBJECT_CLASS

		result := classSuperclass(context)

		_, ok := result.(*Nil)
		if !ok {
			t.Logf("Expected Nil object, got %T\n", result)
			t.Fail()
		}
	})
	t.Run("Eigenclass", func(t *testing.T) {
		context := &Class{superClass: newEigenclass(OBJECT_CLASS, nil)}

		result := classSuperclass(context)

		_, ok := result.(*eigenclass)
		if !ok {
			t.Logf("Expected eigenClass object, got %T\n", result)
			t.Fail()
		}
	})
}
