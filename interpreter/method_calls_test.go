package interpreter

import (
	"reflect"
	"testing"

	"github.com/goruby/goruby/object"
)

func TestMainMethodCalls(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			"method with arguments",
			`
		def add x, y
			x + y
		end

		add 7, 3
		`,
			10,
		},
		{
			"method with arguments and block",
			`
		def add x, y
			yield x + y
		end

		add 7, 3 { |sum| sum - 5 }
		`,
			5,
		},
		{
			"method with arguments and block and __BLOCK__ as local variable",
			`
		def add x, y
			__BLOCK__ = 3
			yield x + y
		end

		add 7, 3 { |sum| sum - 5 }
		`,
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := New()

			evaluated, err := i.Interpret(tt.input)
			if err != nil {
				t.Logf("Expected no error, got %T:%v", err, err)
				t.FailNow()
			}

			integer, ok := evaluated.(*object.Integer)
			if !ok {
				t.Logf("Expected evaluated value to be *object.Integer, got %T", evaluated)
				t.FailNow()
			}

			actual := integer.Value

			if tt.expected != actual {
				t.Logf("Expected result to equal %d, got %d", tt.expected, actual)
				t.Fail()
			}
		})
	}
}

func TestMethodBlockLeakage(t *testing.T) {
	t.Run("main methods", func(t *testing.T) {
		input := `
		def add x, y
		yield x + y
		end

		def sub x, y
		yield x - y
		end

		add 7, 3 { |sum| sum - 5 }

		sub 10, 4
		`
		i := New()

		_, err := i.Interpret(input)

		expected := object.NewNoBlockGivenLocalJumpError()

		if !reflect.DeepEqual(expected, err) {
			t.Logf("Expected error to equal\n%+#v\n\tgot\n%+#v\n", expected, err)
			t.Fail()
		}
	})
	t.Run("custom class methods", func(t *testing.T) {
		input := `
		class Foo
			def add x, y
			yield x + y
			end

			def sub x, y
			yield x - y
			end
		end

		foo = Foo.new

		foo.add 7, 3 { |sum| sum - 5 }

		foo.sub 10, 4
		`
		i := New()

		_, err := i.Interpret(input)

		expected := object.NewNoBlockGivenLocalJumpError()

		if !reflect.DeepEqual(expected, err) {
			t.Logf("Expected error to equal\n%+#v\n\tgot\n%+#v\n", expected, err)
			t.Fail()
		}
	})
}
