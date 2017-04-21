package interpreter

import (
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
