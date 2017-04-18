package object

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewClass(t *testing.T) {
	t.Run("full constructor use", func(t *testing.T) {
		instanceMethods := map[string]RubyMethod{
			"instance_foo": publicMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
				return TRUE, nil
			}),
		}
		classMethods := map[string]RubyMethod{
			"class_foo": publicMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
				return TRUE, nil
			}),
		}
		superClass := arrayClass

		classObject := NewClass("Abc", superClass, instanceMethods, classMethods)

		class, ok := classObject.(*class)
		if !ok {
			t.Logf("Expected returned object to be a class, got %T", classObject)
			t.Fail()
		}

		expectedInstanceMethodSet := NewMethodSet(instanceMethods)
		if !reflect.DeepEqual(expectedInstanceMethodSet, class.instanceMethods) {
			t.Logf(
				"Expected class.instanceMethods to equal\n%+#v\n\tgot\n%+#v\n",
				expectedInstanceMethodSet,
				class.instanceMethods,
			)
			t.Fail()
		}

		if !reflect.DeepEqual(superClass, class.superClass) {
			t.Logf(
				"Expected class.superClass to equal\n%+#v\n\tgot\n%+#v\n",
				superClass,
				class.superClass,
			)
			t.Fail()
		}

		expectedClassMethodSet := NewMethodSet(classMethods)
		actualClassMethods := class.class.Methods()
		if !reflect.DeepEqual(expectedClassMethodSet, actualClassMethods) {
			t.Logf(
				"Expected class.class.Methods to equal\n%+#v\n\tgot\n%+#v\n",
				expectedClassMethodSet,
				actualClassMethods,
			)
			t.Fail()
		}
	})
	t.Run("missing instanceMethods", func(t *testing.T) {
		classMethods := map[string]RubyMethod{
			"class_foo": publicMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
				return TRUE, nil
			}),
		}
		superClass := arrayClass

		classObject := NewClass("Abc", superClass, nil, classMethods)

		class, ok := classObject.(*class)
		if !ok {
			t.Logf("Expected returned object to be a class, got %T", classObject)
			t.Fail()
		}

		expectedInstanceMethodSet := NewMethodSet(map[string]RubyMethod{})
		if !reflect.DeepEqual(expectedInstanceMethodSet, class.instanceMethods) {
			t.Logf(
				"Expected class.instanceMethods to equal\n%+#v\n\tgot\n%+#v\n",
				expectedInstanceMethodSet,
				class.instanceMethods,
			)
			t.Fail()
		}
	})
	t.Run("missing instanceMethods", func(t *testing.T) {
		instanceMethods := map[string]RubyMethod{
			"class_foo": publicMethod(func(context CallContext, args ...RubyObject) (RubyObject, error) {
				return TRUE, nil
			}),
		}
		superClass := arrayClass

		classObject := NewClass("Abc", superClass, instanceMethods, nil)

		class, ok := classObject.(*class)
		if !ok {
			t.Logf("Expected returned object to be a class, got %T", classObject)
			t.Fail()
		}

		expectedClassMethodSet := NewMethodSet(map[string]RubyMethod{})

		actualClassMethods := class.class.Methods()
		if !reflect.DeepEqual(expectedClassMethodSet, actualClassMethods) {
			t.Logf(
				"Expected class.class.Methods to equal\n%+#v\n\tgot\n%+#v\n",
				expectedClassMethodSet,
				actualClassMethods,
			)
			t.Fail()
		}
	})
}

func TestClassInspect(t *testing.T) {
	t.Run("class Class", func(t *testing.T) {
		context := &class{}

		actual := context.Inspect()

		expected := fmt.Sprintf("#<Class:%p>", context)

		if expected != actual {
			t.Logf("Expected Inspect to equal %q, got %q", expected, actual)
			t.Fail()
		}
	})
	t.Run("other class", func(t *testing.T) {
		context := &class{name: "Foo"}

		actual := context.Inspect()

		expected := "Foo"

		if expected != actual {
			t.Logf("Expected Inspect to equal %q, got %q", expected, actual)
			t.Fail()
		}
	})
}

func TestClassClass(t *testing.T) {
	clazz := &class{}

	context := &class{class: clazz}

	actual := context.Class()

	expected := clazz

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Class to equal %v, got %v", expected, actual)
		t.Fail()
	}
}

func TestClassSuperClass(t *testing.T) {
	context := &class{superClass: basicObjectClass}

	actual := context.SuperClass()

	expected := basicObjectClass

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected SuperClass to equal %+#v, got %+#v", expected, actual)
		t.Fail()
	}
}

func TestClassMethods(t *testing.T) {
	contextMethods := map[string]RubyMethod{
		"a_method": nil,
	}

	context := &class{instanceMethods: NewMethodSet(contextMethods)}

	actual := context.Methods()

	expected := NewMethodSet(contextMethods)

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Methods to equal %v, got %v", expected, actual)
		t.Fail()
	}
}

func TestClassType(t *testing.T) {
	context := &class{}

	actual := context.Type()

	expected := CLASS_OBJ

	if expected != actual {
		t.Logf("Expected Type to equal %v, got %v", expected, actual)
		t.Fail()
	}
}

func TestClassSuperclass(t *testing.T) {
	t.Run("anything else than BasicObject", func(t *testing.T) {
		context := &callContext{receiver: &class{superClass: objectClass}}

		result, err := classSuperclass(context)
		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		_, ok := result.(*class)
		if !ok {
			t.Logf("Expected Class object, got %T\n", result)
			t.Fail()
		}
	})
	t.Run("BasicObject", func(t *testing.T) {
		context := &callContext{receiver: basicObjectClass}

		result, err := classSuperclass(context)
		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		_, ok := result.(*nilObject)
		if !ok {
			t.Logf("Expected Nil object, got %T\n", result)
			t.Fail()
		}
	})
	t.Run("Eigenclass", func(t *testing.T) {
		context := &callContext{receiver: &class{superClass: newEigenclass(objectClass, nil)}}

		result, err := classSuperclass(context)
		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		_, ok := result.(*eigenclass)
		if !ok {
			t.Logf("Expected eigenClass object, got %T\n", result)
			t.Fail()
		}
	})
}
