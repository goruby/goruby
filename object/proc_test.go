package object

import (
	"reflect"
	"testing"

	"github.com/goruby/goruby/ast"
)

func TestExtractBlockFromArgs(t *testing.T) {
	t.Run("args empty", func(t *testing.T) {
		args := []RubyObject{}

		block, remaining, ok := extractBlockFromArgs(args)

		if ok {
			t.Logf("Expected no block found")
			t.Fail()
		}

		if block != nil {
			t.Logf("Expected block to be nil, got %+#v\n", block)
			t.Fail()
		}

		if len(remaining) != 0 {
			t.Logf("Expected remaining args to have length %d, got %d", 0, len(remaining))
			t.Fail()
		}
	})
	t.Run("args with block only", func(t *testing.T) {
		args := []RubyObject{&Proc{}}

		block, remaining, ok := extractBlockFromArgs(args)

		if !ok {
			t.Logf("Expected block found")
			t.Fail()
		}

		if block == nil {
			t.Logf("Expected block not to be nil")
			t.Fail()
		}

		if len(remaining) != 0 {
			t.Logf("Expected remaining args to have length %d, got %d", 0, len(remaining))
			t.Fail()
		}
	})
	t.Run("args with nil and block", func(t *testing.T) {
		args := []RubyObject{NIL, &Proc{}}

		block, remaining, ok := extractBlockFromArgs(args)

		if !ok {
			t.Logf("Expected block found")
			t.Fail()
		}

		if block == nil {
			t.Logf("Expected block not to be nil")
			t.Fail()
		}

		if len(remaining) != 1 {
			t.Logf("Expected remaining args to have length %d, got %d", 1, len(remaining))
			t.Fail()
		}

		expected := []RubyObject{NIL}

		if !reflect.DeepEqual(expected, remaining) {
			t.Logf("Expected remaining args to equal\n%+#v\n\tgot\n%+#v\n", expected, remaining)
			t.Fail()
		}
	})
	t.Run("args with nil and block but block not at the end", func(t *testing.T) {
		args := []RubyObject{&Proc{}, NIL}

		block, remaining, ok := extractBlockFromArgs(args)

		if ok {
			t.Logf("Expected no block found")
			t.Fail()
		}

		if block != nil {
			t.Logf("Expected block to be nil, got %+#v\n", block)
			t.Fail()
		}

		if len(remaining) != 2 {
			t.Logf("Expected remaining args to have length %d, got %d", 2, len(remaining))
			t.Fail()
		}

		expected := []RubyObject{&Proc{}, NIL}

		if !reflect.DeepEqual(expected, remaining) {
			t.Logf("Expected remaining args to equal\n%+#v\n\tgot\n%+#v\n", expected, remaining)
			t.Fail()
		}
	})
}

func TestProcCall(t *testing.T) {
	t.Run("argument count not mandatory", func(t *testing.T) {
		proc := &Proc{
			Parameters: []*ast.FunctionParameter{&ast.FunctionParameter{Name: &ast.Identifier{Value: "a"}}},
			Body:       &ast.BlockStatement{Statements: []ast.Statement{}},
			Env:        NewEnvironment(),
			ArgumentCountMandatory: false,
		}
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}
		context := &callContext{
			receiver: NIL,
			env:      NewEnvironment(),
			eval:     eval,
		}

		_, err := proc.Call(context)

		checkError(t, err, nil)
	})
	t.Run("argument count mandatory", func(t *testing.T) {
		proc := &Proc{
			Parameters: []*ast.FunctionParameter{&ast.FunctionParameter{Name: &ast.Identifier{Value: "a"}}},
			Body:       &ast.BlockStatement{Statements: []ast.Statement{}},
			Env:        NewEnvironment(),
			ArgumentCountMandatory: true,
		}
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}
		context := &callContext{
			receiver: NIL,
			env:      NewEnvironment(),
			eval:     eval,
		}

		_, err := proc.Call(context)

		expected := NewWrongNumberOfArgumentsError(1, 0)

		checkError(t, err, expected)
	})
}
