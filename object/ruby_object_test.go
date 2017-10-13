package object

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/goruby/goruby/ast"
)

func mustCall(obj RubyObject, err error) RubyObject {
	if err != nil {
		panic(err)
	}
	return obj
}

func TestFunctionCall(t *testing.T) {
	t.Run("calls CallContext#Eval with its Body", func(t *testing.T) {
		functionBody := &ast.BlockStatement{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Expression: &ast.IntegerLiteral{Value: 13},
				},
			},
		}

		function := &Function{
			Body: functionBody,
		}

		var actualEvalNode ast.Node
		context := &callContext{
			env: NewMainEnvironment(),
			eval: func(node ast.Node, env Environment) (RubyObject, error) {
				actualEvalNode = node
				return nil, nil
			},
		}

		mustCall(function.Call(context))

		var expected ast.Node = functionBody
		if !reflect.DeepEqual(expected, actualEvalNode) {
			t.Logf("Expected Eval argument to equal\n%v\n\tgot\n%v\n", expected, actualEvalNode)
			t.Fail()
		}
	})
	t.Run("returns any error returned by CallContext#Eval", func(t *testing.T) {
		evalErr := fmt.Errorf("An error")

		context := &callContext{
			env:  NewMainEnvironment(),
			eval: func(ast.Node, Environment) (RubyObject, error) { return nil, evalErr },
		}

		function := &Function{
			Parameters: []*FunctionParameter{},
		}

		_, err := function.Call(context)

		if !reflect.DeepEqual(evalErr, err) {
			t.Logf("Expected error to equal\n%v\n\tgot\n%v\n", evalErr, err)
			t.Fail()
		}
	})
	t.Run("uses the function env as env for CallContext#Eval", func(t *testing.T) {
		contextEnv := NewEnvironment()
		contextEnv.Set("self", &Self{RubyObject: &Integer{Value: 42}, Name: "context self"})
		contextEnv.Set("bar", &String{Value: "not reachable in Eval"})
		var evalEnv Environment
		context := &callContext{
			env: contextEnv,
			eval: func(node ast.Node, env Environment) (RubyObject, error) {
				evalEnv = env
				return nil, nil
			},
		}

		functionEnv := NewEnvironment()
		functionEnv.Set("foo", &Symbol{Value: "bar"})
		function := &Function{
			Parameters: []*FunctionParameter{},
			Env:        functionEnv,
		}

		mustCall(function.Call(context))

		{
			expected := &Symbol{Value: "bar"}
			actual, ok := evalEnv.Get("foo")

			if !ok {
				t.Logf("Expected key 'foo' to be in Eval env")
				t.FailNow()
			}

			if !reflect.DeepEqual(expected, actual) {
				t.Logf("Expected 'foo' to equal\n%v\n\tgot\n%v\n", expected, actual)
				t.Fail()
			}
		}

		_, ok := evalEnv.Get("bar")
		if ok {
			t.Logf("Expected key 'bar' not to be in Eval env")
			t.Fail()
		}

		{
			expected := &Self{RubyObject: &Integer{Value: 42}, Name: "context self"}
			actual, _ := evalEnv.Get("self")
			if !reflect.DeepEqual(expected, actual) {
				t.Logf("Expected Eval env self to equal\n%+v\n\tgot\n%+v\n", expected, actual)
				t.Fail()
			}
			if expected == actual {
				t.Logf("Expected Eval env self to be a new Self object\n")
				t.Fail()
			}
		}
	})
	t.Run("puts the Call args into the env for CallContext#Eval", func(t *testing.T) {
		contextEnv := NewEnvironment()
		contextEnv.Set("self", &Self{RubyObject: &Integer{Value: 42}, Name: "context self"})
		var evalEnv Environment
		context := &callContext{
			env: contextEnv,
			eval: func(node ast.Node, env Environment) (RubyObject, error) {
				evalEnv = env
				return nil, nil
			},
		}

		t.Run("without default params", func(t *testing.T) {
			function := &Function{
				Parameters: []*FunctionParameter{
					&FunctionParameter{Name: "foo"},
					&FunctionParameter{Name: "bar"},
				},
			}

			mustCall(function.Call(context, &Integer{Value: 300}, &Symbol{Value: "sym"}))

			{
				expected := &Integer{Value: 300}
				actual, ok := evalEnv.Get("foo")

				if !ok {
					t.Logf("Expected function parameter %q to be in Eval env", "foo")
					t.FailNow()
				}

				if !reflect.DeepEqual(expected, actual) {
					t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, actual)
					t.Fail()
				}
			}
			{
				expected := &Symbol{Value: "sym"}
				actual, ok := evalEnv.Get("bar")

				if !ok {
					t.Logf("Expected function parameter %q to be in Eval env", "bar")
					t.FailNow()
				}

				if !reflect.DeepEqual(expected, actual) {
					t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, actual)
					t.Fail()
				}
			}
		})
		t.Run("with default params", func(t *testing.T) {
			t.Skip()
			function := &Function{
				Parameters: []*FunctionParameter{
					&FunctionParameter{Name: "foo", Default: &Integer{Value: 12}},
					&FunctionParameter{Name: "bar"},
					&FunctionParameter{Name: "qux"},
				},
			}

			mustCall(function.Call(context, &Integer{Value: 300}, &Symbol{Value: "sym"}))

			{
				expected := &Integer{Value: 12}
				actual, ok := evalEnv.Get("foo")

				if !ok {
					t.Logf("Expected function parameter %q to be in Eval env", "foo")
					t.FailNow()
				}

				if !reflect.DeepEqual(expected, actual) {
					t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, actual)
					t.Fail()
				}
			}
			{
				expected := &Integer{Value: 300}
				actual, ok := evalEnv.Get("bar")

				if !ok {
					t.Logf("Expected function parameter %q to be in Eval env", "bar")
					t.FailNow()
				}

				if !reflect.DeepEqual(expected, actual) {
					t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, actual)
					t.Fail()
				}
			}
			{
				expected := &Symbol{Value: "sym"}
				actual, ok := evalEnv.Get("qux")

				if !ok {
					t.Logf("Expected function parameter %q to be in Eval env", "qux")
					t.FailNow()
				}

				if !reflect.DeepEqual(expected, actual) {
					t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, actual)
					t.Fail()
				}
			}
		})
	})
	t.Run("returns the object returned by CallContext#Eval", func(t *testing.T) {
		t.Run("vanilla object", func(t *testing.T) {
			context := &callContext{
				env:  NewMainEnvironment(),
				eval: func(ast.Node, Environment) (RubyObject, error) { return &Integer{Value: 8}, nil },
			}

			function := &Function{}

			result, _ := function.Call(context)

			expected := &Integer{Value: 8}

			if !reflect.DeepEqual(expected, result) {
				t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, result)
				t.Fail()
			}
		})
		t.Run("wrapped into a return value", func(t *testing.T) {
			context := &callContext{
				env:  NewMainEnvironment(),
				eval: func(ast.Node, Environment) (RubyObject, error) { return &ReturnValue{Value: &Integer{Value: 8}}, nil },
			}

			function := &Function{}

			result, _ := function.Call(context)

			expected := &Integer{Value: 8}

			if !reflect.DeepEqual(expected, result) {
				t.Logf("Expected result to equal\n%v\n\tgot\n%v\n", expected, result)
				t.Fail()
			}
		})
	})
	t.Run("extracts the Call block from the rest of the arguments", func(t *testing.T) {
		contextEnv := NewEnvironment()
		contextEnv.Set("self", &Self{RubyObject: &Integer{Value: 42}, Name: "context self"})

		var evalEnv Environment
		context := &callContext{
			env: contextEnv,
			eval: func(node ast.Node, env Environment) (RubyObject, error) {
				evalEnv = env
				return nil, nil
			},
		}

		function := &Function{
			Parameters: []*FunctionParameter{
				&FunctionParameter{Name: "x"},
			},
		}

		mustCall(function.Call(context, &String{Value: "the x value"}, &Proc{}))

		actual := len(evalEnv.GetAll())
		expected := 2 // `self` and `x`

		if !reflect.DeepEqual(expected, actual) {
			t.Logf("Expected Eval env to have %d items, got %d\n", expected, actual)
			t.Fail()
		}
	})
	t.Run("propagates the Call block to CallContext#Eval", func(t *testing.T) {
		contextEnv := NewEnvironment()
		contextEnv.Set("self", &Self{RubyObject: &Integer{Value: 42}, Name: "context self"})

		var evalEnv Environment
		context := &callContext{
			env: contextEnv,
			eval: func(node ast.Node, env Environment) (RubyObject, error) {
				evalEnv = env
				return nil, nil
			},
		}

		function := &Function{
			Parameters: []*FunctionParameter{},
		}

		mustCall(function.Call(context, &Proc{ArgumentCountMandatory: true}))

		expected := &Proc{ArgumentCountMandatory: true}
		envSelf, _ := evalEnv.Get("self")
		actual := envSelf.(*Self).Block
		if !reflect.DeepEqual(expected, actual) {
			t.Logf("Expected block to equal\n%+v\n\tgot\n%+v\n", expected, actual)
			t.Fail()
		}
	})
	t.Run("validates that the arguments match the function parameters", func(t *testing.T) {
		context := &callContext{
			env:  NewMainEnvironment(),
			eval: func(ast.Node, Environment) (RubyObject, error) { return nil, nil },
		}

		function := &Function{
			Parameters: []*FunctionParameter{},
		}

		t.Run("without block argument", func(t *testing.T) {
			expected := NewWrongNumberOfArgumentsError(0, 1)

			_, err := function.Call(context, &String{Value: "foo"})

			if !reflect.DeepEqual(expected, err) {
				t.Logf("Expected error to equal\n%v\n\tgot\n%v\n", expected, err)
				t.Fail()
			}
		})

		t.Run("with block argument", func(t *testing.T) {
			_, err := function.Call(context, &Proc{})

			if err != nil {
				t.Logf("Expected no error, got %T:%v\n", err, err)
				t.Fail()
			}
		})

		t.Run("with default arguments", func(t *testing.T) {
			function.Parameters = []*FunctionParameter{
				&FunctionParameter{Name: "x", Default: TRUE},
				&FunctionParameter{Name: "y"},
			}

			_, err := function.Call(context, &Integer{Value: 8})

			if err != nil {
				t.Logf("Expected no error, got %T:%v\n", err, err)
				t.Fail()
			}
		})
	})
}
