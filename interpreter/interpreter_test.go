package interpreter

import (
	"testing"

	"github.com/goruby/goruby/object"
	"github.com/pkg/errors"
)

func TestInterpreterInterpret(t *testing.T) {
	t.Run("return proper result", func(t *testing.T) {
		input := `
			def foo
				3
			end

			x = 5

			def add x, y
				x + y
			end

			add foo, x
			`
		i := New()

		out, err := i.Interpret("", input)
		if err != nil {
			panic(err)
		}

		res, ok := out.(*object.Integer)
		if !ok {
			t.Logf("Expected *object.Integer, got %T\n", out)
			t.Fail()
		}

		if res.Value != 8 {
			t.Logf("Expected result to equal 8, got %d\n", res.Value)
			t.Fail()
		}
	})
}

func TestModuleInEnv(t *testing.T) {
	input := `
		module Foo
			def foo
				3
			end
		end
		Foo
	`
	interpreter := New()

	evaluated, err := interpreter.Interpret("", input)
	if err != nil {
		t.Logf("Expected no error, got %T:%v", err, err)
		t.Fail()
	}

	_, ok := evaluated.(*object.Module)
	if !ok {
		t.Logf("Expected evaluated return to be a object.Module, got %T", evaluated)
		t.Fail()
	}
}

func TestInterpretModules(t *testing.T) {
	t.Run("evaluate module object", func(t *testing.T) {
		input := `
		module Foo
			def foo
				3
			end
		end
		Foo
		`

		i := New()

		evaluated, err := i.Interpret("", input)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.FailNow()
		}

		module, ok := evaluated.(*object.Module)
		if !ok {
			t.Logf("Expected module object, got %T", evaluated)
			t.FailNow()
		}

		_, ok = module.Class().Methods().Get("foo")
		if !ok {
			t.Logf("Expected class object to have method foo")
			t.Fail()
		}
	})
	t.Run("self after module definition", func(t *testing.T) {
		input := `
		module Foo
			def foo
				3
			end
		end
		self
		`

		i := New()

		evaluated, err := i.Interpret("", input)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.FailNow()
		}

		self, ok := evaluated.(*object.Self)
		if !ok {
			t.Logf("Expected self, got %T", evaluated)
			t.FailNow()
		}

		if self.Name != "main" {
			t.Logf("Expected main object, got %s", self.Name)
			t.Fail()
		}
	})
	t.Run("module definitions do not overwrite but add", func(t *testing.T) {
		input := `
		module Foo
			def foo
				3
			end
		end

		module Foo
			def bar
				"bar"
			end
		end

		Foo
		`

		i := New()

		evaluated, err := i.Interpret("", input)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.FailNow()
		}

		module, ok := evaluated.(*object.Module)
		if !ok {
			t.Logf("Expected module object, got %T", evaluated)
			t.FailNow()
		}

		_, ok = module.Class().Methods().Get("foo")
		if !ok {
			t.Logf("Expected module object to have method foo")
			t.Fail()
		}

		_, ok = module.Class().Methods().Get("bar")
		if !ok {
			t.Logf("Expected module object to have method bar")
			t.Fail()
		}
	})
	t.Run("nested module definitions", func(t *testing.T) {
		input := `
		module Foo
			def foo
				3
			end

			module Bar
				def bar
					"bar"
				end
			end
			Bar
		end
		`

		i := New()

		evaluated, err := i.Interpret("", input)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.FailNow()
		}

		module, ok := evaluated.(*object.Module)
		if !ok {
			t.Logf("Expected module object, got %T", evaluated)
			t.FailNow()
		}

		if module.Inspect() != "Bar" {
			t.Logf("Expected module object to stringify to 'Bar', got %q", module.Inspect())
			t.Fail()
		}

		_, ok = module.Class().Methods().Get("bar")
		if !ok {
			t.Logf("Expected module object to have method bar")
			t.Fail()
		}
	})
	t.Run("nested module definitions are scoped to the inner module", func(t *testing.T) {
		input := `
		module Foo
			module Bar
			end
		end
		Bar
		`

		i := New()

		_, err := i.Interpret("", input)
		if err == nil {
			t.Logf("Expected error, got nil")
			t.FailNow()
		}

		expected := object.NewUninitializedConstantNameError("Bar")

		if errors.Cause(err).Error() != expected.Error() {
			t.Logf("Expected error to equal %T:%v, got %T:%v", expected, expected.Error(), err, err)
			t.Fail()
		}
	})
}
func TestInterpretClasses(t *testing.T) {
	t.Run("evaluate class object", func(t *testing.T) {
		input := `
		class Foo
			def foo
				3
			end
		end
		Foo
		`

		i := New()

		evaluated, err := i.Interpret("", input)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.FailNow()
		}

		classObject, ok := evaluated.(object.RubyClassObject)
		if !ok {
			t.Logf("Expected class object, got %T", evaluated)
			t.FailNow()
		}

		_, ok = classObject.Methods().Get("foo")
		if !ok {
			t.Logf("Expected class object to have method foo")
			t.Fail()
		}
	})
	t.Run("self after class definition", func(t *testing.T) {
		input := `
		class Foo
			def foo
				3
			end
		end
		self
		`

		i := New()

		evaluated, err := i.Interpret("", input)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.FailNow()
		}

		self, ok := evaluated.(*object.Self)
		if !ok {
			t.Logf("Expected self, got %T", evaluated)
			t.FailNow()
		}

		if self.Name != "main" {
			t.Logf("Expected main object, got %s", self.Name)
			t.Fail()
		}
	})
	t.Run("class definitions do not overwrite but add", func(t *testing.T) {
		input := `
		class Foo
			def foo
				3
			end
		end

		class Foo
			def bar
				"bar"
			end
		end

		Foo
		`

		i := New()

		evaluated, err := i.Interpret("", input)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.FailNow()
		}

		classObject, ok := evaluated.(object.RubyClassObject)
		if !ok {
			t.Logf("Expected class object, got %T", evaluated)
			t.FailNow()
		}

		_, ok = classObject.Methods().Get("foo")
		if !ok {
			t.Logf("Expected class object to have method foo")
			t.Fail()
		}

		_, ok = classObject.Methods().Get("bar")
		if !ok {
			t.Logf("Expected class object to have method bar")
			t.Fail()
		}
	})
	t.Run("nested class definitions", func(t *testing.T) {
		input := `
		class Foo
			def foo
				3
			end

			class Bar
				def bar
					"bar"
				end
			end
			Bar
		end
		`

		i := New()

		evaluated, err := i.Interpret("", input)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.FailNow()
		}

		classObject, ok := evaluated.(object.RubyClassObject)
		if !ok {
			t.Logf("Expected class object, got %T", evaluated)
			t.FailNow()
		}

		if classObject.Inspect() != "Bar" {
			t.Logf("Expected class object to stringify to 'Bar', got %q", classObject.Inspect())
			t.Fail()
		}

		_, ok = classObject.Methods().Get("bar")
		if !ok {
			t.Logf("Expected class object to have method bar")
			t.Fail()
		}
	})
	t.Run("nested class definitions are scoped to the inner class", func(t *testing.T) {
		input := `
		class Foo
			class Bar
			end
		end
		Bar
		`

		i := New()

		_, err := i.Interpret("", input)
		if err == nil {
			t.Logf("Expected error, got nil")
			t.FailNow()
		}

		expected := object.NewUninitializedConstantNameError("Bar")

		if errors.Cause(err).Error() != expected.Error() {
			t.Logf("Expected error to equal %T:%v, got %T:%v", expected, expected.Error(), err, err)
			t.Fail()
		}
	})
}
