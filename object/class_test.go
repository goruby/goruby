package object

import "testing"

func TestClassSuperclass(t *testing.T) {
	t.Run("anything else than BasicObject", func(t *testing.T) {
		context := &testRubyObject{superClass: OBJECT_CLASS}

		result := classSuperclass(context)

		_, ok := result.(*ObjectClass)
		if !ok {
			t.Logf("Expected ObjectClass object, got %T\n", result)
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
		context := &testRubyObject{superClass: newEigenClass(OBJECT_CLASS, nil)}

		result := classSuperclass(context)

		_, ok := result.(*eigenClass)
		if !ok {
			t.Logf("Expected eigenClass object, got %T\n", result)
			t.Fail()
		}
	})
}
