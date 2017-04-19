package interpreter

import (
	"reflect"
	"testing"

	"github.com/goruby/goruby/object"
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

		out, err := i.Interpret(input)
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
	t.Run("return proper result with changed env", func(t *testing.T) {
		input := `
			def foo
				3
			end

			def add x, y
				x + y
			end

			add foo, x
			`
		env := object.NewMainEnvironment()
		env.Set("x", &object.Integer{Value: 3})
		i := New()
		i.SetEnvironment(env)

		out, err := i.Interpret(input)
		if err != nil {
			panic(err)
		}

		res, ok := out.(*object.Integer)
		if !ok {
			t.Logf("Expected *object.Integer, got %T\n", out)
			t.Fail()
		}

		if res.Value != 6 {
			t.Logf("Expected result to equal 8, got %d\n", res.Value)
			t.Fail()
		}
	})
}

func TestSelfAfterModuleDefinition(t *testing.T) {
	input := `
		module Foo
		end

		self
	`
	interpreter := New()

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected no panic, got %T:%v\n", r, r)
			t.Fail()
		}
	}()
	evaluated, err := interpreter.Interpret(input)

	if err != nil {
		t.Logf("Expected no error, got %T:%v", err, err)
		t.Fail()
	}

	expected := &object.Self{&object.Object{}, "main"}
	if !reflect.DeepEqual(expected, evaluated) {
		t.Logf("Expected self to equal\n%+#v\n\tgot\n%+#v\n", expected, evaluated)
		t.Fail()
	}
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

	evaluated, err := interpreter.Interpret(input)
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
