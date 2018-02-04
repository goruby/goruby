package evaluator

import (
	"go/token"
	"reflect"
	"strings"
	"testing"

	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
	"github.com/pkg/errors"
)

func TestEvalComment(t *testing.T) {
	input := "5 # five"

	evaluated, err := testEval(input)
	checkError(t, err)

	var expected int64 = 5
	testIntegerObject(t, evaluated, expected)
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"5 % 2", 1},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input, object.NewMainEnvironment())
		checkError(t, err)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"true || true", true},
		{"false || false", false},
		{"true || false", true},
		{"false || false", false},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		checkError(t, err)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		checkError(t, err)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestConditionalExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if true; 10; end", 10},
		{"if false; 10; end", nil},
		{"if 1; 10; end", 10},
		{"if 1 < 2; 10; end", 10},
		{"if 1 > 2; 10; end", nil},
		{"if 1 > 2; 10; else\n 20; end", 20},
		{"if 1 < 2; 10; else\n 20; end", 10},
		{"unless true; 10; end", nil},
		{"unless false; 10; end", 10},
		{"unless 1; 10; end", nil},
		{"unless 1 < 2; 10; end", nil},
		{"unless 1 > 2; 10; end", 10},
		{"unless 1 > 2; 10; else\n 20; end", 10},
		{"unless 1 < 2; 10; else\n 20; end", 20},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		checkError(t, err)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`if 10 > 1
			if 10 > 1
				return 10
			end
			return 1
		  end`, 10},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		checkError(t, err)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"TypeError: Integer can't be coerced into Boolean",
		},
		{
			"5 + true; 5;",
			"TypeError: Integer can't be coerced into Boolean",
		},
		{
			"-true",
			"Exception: unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"NoMethodError: undefined method `+' for true:TrueClass",
		},
		{
			"true + false + true + false;",
			"NoMethodError: undefined method `+' for true:TrueClass",
		},
		{
			"5; true + false; 5",
			"NoMethodError: undefined method `+' for true:TrueClass",
		},
		{
			`"Hello" - "World"`,
			"NoMethodError: undefined method `-' for Hello:String",
		},
		{
			"if (10 > 1); true + false; end",
			"NoMethodError: undefined method `+' for true:TrueClass",
		},
		{
			"if (10 > 1); true + false; end",
			"NoMethodError: undefined method `+' for true:TrueClass",
		},
		{
			`
if (10 > 1)
	if (10 > 1)
		return true + false;
	end
	return 1;
end
`,
			"NoMethodError: undefined method `+' for true:TrueClass",
		},
		{
			"foobar",
			"NameError: undefined local variable or method `foobar' for main:Object",
		},
		{
			"Foobar",
			"NameError: uninitialized constant Foobar",
		},
		{
			`
			def foo x, y
			end

			foo 1
			`,
			"ArgumentError: wrong number of arguments (given 1, expected 2)",
		},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		env.Set("self", &object.Self{RubyObject: &object.Object{}, Name: "main"})
		evaluated, err := testEval(tt.input, env)

		if err == nil {
			t.Errorf(
				"no error returned. got=%T(%+v)",
				evaluated,
				evaluated,
			)
		}

		actual, ok := errors.Cause(err).(object.RubyObject)
		if !ok {
			t.Logf("Error is not a RubyObject, got %T:%v\n", err, err)
			t.FailNow()
		}

		testExceptionObject(t, actual, tt.expectedMessage)
	}
}

func mustGet(obj object.RubyObject, ok bool) object.RubyObject {
	if !ok {
		panic("object not found")
	}
	return obj
}

func TestExceptionHandlingBlock(t *testing.T) {
	tests := []struct {
		input  string
		err    error
		output object.RubyObject
	}{
		{
			"begin\n2\nend",
			nil,
			&object.Integer{Value: 2},
		},
		{
			`
begin
	2
rescue Exception => e
	4
end`,
			nil,
			&object.Integer{Value: 2},
		},
		{
			`
begin
	raise Exception
end`,
			object.NewException("Exception"),
			nil,
		},
		{
			`
begin
	raise Exception
rescue Exception
	4
end`,
			nil,
			&object.Integer{Value: 4},
		},
		{
			`
begin
	raise Exception.new "foo"
rescue Exception => e
	e.to_s
end`,
			nil,
			&object.String{Value: "foo"},
		},
		{
			`
begin
	raise StandardError.new "foo"
rescue
	5
end`,
			nil,
			&object.Integer{Value: 5},
		},
		{
			`
begin
	raise StandardError.new "bar"
rescue => e
	e.to_s
end`,
			nil,
			&object.String{Value: "bar"},
		},
		{
			`
begin
	raise Exception.new "qux"
rescue
	3
end`,
			object.NewException("qux"),
			nil,
		},
		{
			`
begin
	raise StandardError.new "qux"
rescue
	3
end`,
			nil,
			&object.Integer{Value: 3},
		},
	}

	for _, tt := range tests {
		if tt.err == nil {
			t.Run("catched exception", func(t *testing.T) {
				env := object.NewMainEnvironment()
				evaluated, err := testEval(tt.input, env)
				checkError(t, err)

				if !reflect.DeepEqual(evaluated, tt.output) {
					t.Logf("Expected result to equal\n%+#v\n\tgot\n%+#v\n", tt.output, evaluated)
					t.Fail()
				}
			})
		} else {
			t.Run("uncaught exception", func(t *testing.T) {
				env := object.NewMainEnvironment()
				evaluated, err := testEval(tt.input, env)
				if evaluated != nil {
					t.Logf("expected result to be nil")
					t.Fail()
				}

				if !reflect.DeepEqual(errors.Cause(err), tt.err) {
					t.Logf("Expected err to equal\n%+#v\n\tgot\n%+#v\n", tt.err, errors.Cause(err))
					t.Fail()
				}
			})
		}
	}
}

func TestScopedIdentifierExpression(t *testing.T) {
	objectClassObject, _ := object.NewMainEnvironment().Get("Object")
	objectClass := objectClassObject.(object.RubyClassObject)
	tests := []struct {
		input           string
		expectedInspect string
		expectedClass   object.RubyClass
	}{
		{
			`
			module A
				module B
				end
			end
			A::B
			`,
			object.NewModule("B", nil).Inspect(),
			object.NewModule("B", nil).Class(),
		},
		{
			`
			module A
				class B
				end
			end
			A::B
			`,
			object.NewClass("B", objectClass, nil).Inspect(),
			object.NewClass("B", objectClass, nil).Class(),
		},
		{
			`
			class A
				class B
				end
			end
			A::B
			`,
			object.NewClass("B", objectClass, nil).Inspect(),
			object.NewClass("B", objectClass, nil).Class(),
		},
		{
			`
			class A
				module B
				end
			end
			A::B
			`,
			object.NewModule("B", nil).Inspect(),
			object.NewModule("B", nil).Class(),
		},
		{
			`
			module A
				module B
					module C
					end
				end
			end
			A::B::C
			`,
			object.NewModule("C", nil).Inspect(),
			object.NewModule("C", nil).Class(),
		},
		{
			`
			module A
				Ten = 10
			end
			A::Ten
			`,
			object.NewInteger(10).Inspect(),
			object.NewInteger(10).Class(),
		},
		{
			`
			class A
				def bar
					13
				end
			end
			A.new::bar
			`,
			object.NewInteger(13).Inspect(),
			object.NewInteger(13).Class(),
		},
	}

	for _, tt := range tests {
		env := object.NewMainEnvironment()
		evaluated, err := testEval(tt.input, env)
		checkError(t, err)

		actual := evaluated.Inspect()

		if tt.expectedInspect != actual {
			t.Logf("Expected eval return to equal\n%q\n\tgot\n%q\n", tt.expectedInspect, actual)
			t.Fail()
		}

		if !reflect.DeepEqual(tt.expectedClass, evaluated.Class()) {
			t.Logf("Expected eval return class to equal\n%+#v\n\tgot\n%+#v\n", tt.expectedClass, evaluated.Class())
			t.Fail()
		}
	}
}

func TestInstanceVariable(t *testing.T) {
	tests := []struct {
		input  string
		output object.RubyObject
	}{
		{
			input: `
class X
	@foo
end`,
			output: object.NIL,
		},
		{
			input:  "@foo",
			output: object.NIL,
		},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input, object.NewMainEnvironment())
		checkError(t, err)

		if evaluated != tt.output {
			t.Logf("Expected result to equal %v, got %v\n", tt.output, evaluated)
			t.Fail()
		}
	}
}

func TestAssignment(t *testing.T) {
	t.Run("assign to hash", func(t *testing.T) {
		tests := []struct {
			name     string
			input    string
			expected interface{}
		}{
			{
				"return value for anonymous hash",
				`{:foo => 3}[:foo] = 5`,
				5,
			},
			{
				"return value for hash variable",
				`x = {:foo => 3}; x[:foo] = 5`,
				5,
			},
			{
				"variable value after hash assignment",
				`x = {:foo => 3}; x[:foo] = 5; x`,
				map[string]string{":foo": "5"},
			},
		}

		for _, tt := range tests {
			evaluated, err := testEval(tt.input)
			checkError(t, err)

			testObject(t, evaluated, tt.expected)
		}
	})
	t.Run("assign to local variable", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{
				`foo = 5`,
				5,
			},
			{
				`foo = 5; x = foo; x = 3; x`,
				3,
			},
			{"a = 5; a;", 5},
			{"a = 5 * 5; a;", 25},
			{"a = 5; b = a; b;", 5},
			{"a = 5; b = a; c = a + b + 5; c;", 15},
		}

		for _, tt := range tests {
			evaluated, err := testEval(tt.input, object.NewMainEnvironment())
			checkError(t, err)

			testIntegerObject(t, evaluated, tt.expected)
		}
	})
	t.Run("assign more than one value", func(t *testing.T) {
		tests := []struct {
			input    string
			expected interface{}
		}{
			{
				`foo = 5, 4`,
				[]string{"5", "4"},
			},
			{
				`foo = 5, 4; foo`,
				[]string{"5", "4"},
			},
			{
				`@foo = 5, 6; @foo`,
				[]string{"5", "6"},
			},
		}

		for _, tt := range tests {
			evaluated, err := testEval(tt.input, object.NewMainEnvironment())
			checkError(t, err)

			testObject(t, evaluated, tt.expected)
		}
	})
	t.Run("assign to InstanceVariable", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{
				`@foo = 5`,
				5,
			},
			{
				`@foo = 5; x = @foo; x = 3; x`,
				3,
			},
		}

		for _, tt := range tests {
			evaluated, err := testEval(tt.input, object.NewMainEnvironment())
			checkError(t, err)

			testIntegerObject(t, evaluated, tt.expected)
		}
	})
	t.Run("assign to array", func(t *testing.T) {
		tests := []struct {
			input    string
			size     int
			elements []object.RubyObject
		}{
			{
				`x = [3]; x[0] = 5; x`,
				1,
				[]object.RubyObject{&object.Integer{Value: 5}},
			},
			{
				`x = []; x[0] = 5; x`,
				1,
				[]object.RubyObject{&object.Integer{Value: 5}},
			},
			{
				`x = [3]; x[3] = 5; x`,
				4,
				[]object.RubyObject{&object.Integer{Value: 3}, object.NIL, object.NIL, &object.Integer{Value: 5}},
			},
		}

		for _, tt := range tests {
			evaluated, err := testEval(tt.input)
			checkError(t, err)

			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Logf("Expected to eval to array, got %T\n", evaluated)
				t.FailNow()
			}

			if len(array.Elements) != tt.size {
				t.Logf("Expected array size to equal %d, got %d\n", tt.size, len(array.Elements))
				t.Fail()
			}

			if !reflect.DeepEqual(array.Elements, tt.elements) {
				t.Logf("Expected elements to equal\n%s\n\tgot\n%s\n", tt.elements, array.Elements)
				t.Fail()
			}
		}
	})
	t.Run("assign operator on local variable", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{
				`foo = 2; foo += 5`,
				7,
			},
			{
				`foo = 5; foo -= 3; foo`,
				2,
			},
		}

		for _, tt := range tests {
			evaluated, err := testEval(tt.input, object.NewMainEnvironment())
			checkError(t, err)

			testIntegerObject(t, evaluated, tt.expected)
		}
	})
}

func TestMultiAssignment(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output *object.Array
	}{
		{
			name:  "evenly distributed sides",
			input: "x, y, z = 1, 2, 3; [x, y, z]",
			output: &object.Array{Elements: []object.RubyObject{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
			}},
		},
		{
			name:  "value side one less",
			input: "x, y, z = 1, 2; [x, y, z]",
			output: &object.Array{Elements: []object.RubyObject{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				object.NIL,
			}},
		},
		{
			name:  "value side two less",
			input: "x, y, z = 1; [x, y, z]",
			output: &object.Array{Elements: []object.RubyObject{
				&object.Integer{Value: 1},
				object.NIL,
				object.NIL,
			}},
		},
		{
			name:  "lhs with array index and instance var",
			input: "x = []; x[0], y, @z = 1, 2, 3; [x[0], y, @z]",
			output: &object.Array{Elements: []object.RubyObject{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
			}},
		},
		{
			name:  "lhs with global and const",
			input: "$x, Y = 1, 2; [$x, Y]",
			output: &object.Array{Elements: []object.RubyObject{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated, err := testEval(tt.input, object.NewMainEnvironment())
			checkError(t, err)

			if !reflect.DeepEqual(tt.output, evaluated) {
				t.Logf("Expected result to equal\n%s\n\tgot\n%s\n", tt.output.Inspect(), evaluated.Inspect())
				t.Fail()
			}
		})
	}
}

func TestGlobalAssignmentExpression(t *testing.T) {
	t.Run("assignments", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{"$a = 5; $a;", 5},
			{"$a = 5 * 5; $a;", 25},
			{"$a = 5; $b = $a; $b;", 5},
			{"$a = 5; $b = $a; $c = $a + $b + 5; $c;", 15},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				evaluated, err := testEval(tt.input)
				checkError(t, err)
				testIntegerObject(t, evaluated, tt.expected)
			})
		}
	})
	t.Run("set as global", func(t *testing.T) {
		input := "$Foo = 3"

		outer := object.NewEnvironment()
		env := object.NewEnclosedEnvironment(outer)
		_, err := testEval(input, env)
		checkError(t, err)

		_, ok := outer.Get("$Foo")
		if !ok {
			t.Logf("Expected $FOO to be set in outer env, was not")
			t.Fail()
		}

		_, ok = env.Clone().Get("$Foo")
		if ok {
			t.Logf("Expected $FOO not to be set in inner env")
			t.Fail()
		}
	})
}

func TestModuleObject(t *testing.T) {
	t.Run("module definition", func(t *testing.T) {
		tests := []struct {
			input           string
			expectedName    string
			expectedMethods map[string]string
			expectedReturn  object.RubyObject
		}{
			{
				`module Foo
				def a
				"foo"
				end
				end`,
				"Foo",
				map[string]string{"a": "fn() {\nfoo\n}"},
				&object.Symbol{Value: "a"},
			},
			{
				`module Foo
				3
				end`,
				"Foo",
				map[string]string{},
				&object.Integer{Value: 3},
			},
			{
				`module Foo
				end`,
				"Foo",
				map[string]string{},
				object.NIL,
			},
		}

		for _, tt := range tests {
			env := object.NewEnvironment()
			env.Set("self", &object.Self{RubyObject: &object.Object{}, Name: "main"})
			evaluated, err := testEval(tt.input, env)
			checkError(t, err)

			if !reflect.DeepEqual(evaluated, tt.expectedReturn) {
				t.Logf("Expected return object to equal\n%+#v\n\tgot\n%+#v\n", tt.expectedReturn, evaluated)
				t.Fail()
			}

			module, ok := env.Get(tt.expectedName)
			if !ok {
				t.Logf("Expected module to exist in env")
				t.Logf("Env: %+#v\n", env)
				t.FailNow()
			}

			actualMethods := make(map[string]string)

			methods := module.Class().Methods().GetAll()
			for name, method := range methods {
				if function, ok := method.(*object.Function); ok {
					actualMethods[name] = function.String()
				}
			}

			if !reflect.DeepEqual(tt.expectedMethods, actualMethods) {
				t.Logf(
					"Expected module methods to equal\n%+#v\n\tgot\n%+#v\n",
					tt.expectedMethods,
					actualMethods,
				)
				t.Fail()
			}
		}
	})
	t.Run("self after module definition", func(t *testing.T) {
		input := `
		module Foo
			def a
				"foo"
			end
		end
		`

		main := &object.Self{RubyObject: &object.Object{}, Name: "main"}
		env := object.NewEnvironment()
		env.Set("self", main)
		_, err := testEval(input, env)
		checkError(t, err)

		self, ok := env.Get("self")
		if !ok {
			t.Logf("Expected self in the env")
			t.FailNow()
		}

		if !reflect.DeepEqual(main, self) {
			t.Logf(
				"Expected self to equal\n%+#v\n\tgot\n%+#v\n",
				main,
				self,
			)
			t.Fail()
		}
	})
	t.Run("module as open classes", func(t *testing.T) {
		input :=
			`module Foo
				def a
					"foo"
				end
			end
			module Foo
				def b
					"bar"
				end
			end
			`
		env := object.NewEnvironment()
		env.Set("self", &object.Self{RubyObject: &object.Object{}, Name: "main"})
		_, err := testEval(input, env)
		checkError(t, err)

		module, ok := env.Get("Foo")
		if !ok {
			t.Logf("Expected module to exist in env")
			t.Logf("Env: %+#v\n", env)
			t.FailNow()
		}

		actualMethods := make(map[string]string)

		methods := module.Class().Methods().GetAll()
		for name, method := range methods {
			if function, ok := method.(*object.Function); ok {
				actualMethods[name] = function.String()
			}
		}

		expectedMethods := map[string]string{
			"a": "fn() {\nfoo\n}",
			"b": "fn() {\nbar\n}",
		}

		if !reflect.DeepEqual(expectedMethods, actualMethods) {
			t.Logf(
				"Expected module methods to equal\n%+#v\n\tgot\n%+#v\n",
				expectedMethods,
				actualMethods,
			)
			t.Fail()
		}
	})
}

func TestClassObject(t *testing.T) {
	tests := []struct {
		input              string
		expectedName       string
		expectedSuperclass string
		expectedMethods    map[string]string
		expectedReturn     object.RubyObject
	}{
		{
			`class Foo
				def a
					"foo"
				end
			end`,
			"Foo",
			"Object",
			map[string]string{"a": "fn() {\nfoo\n}"},
			&object.Symbol{Value: "a"},
		},
		{
			`class Foo
				3
			end`,
			"Foo",
			"Object",
			map[string]string{},
			&object.Integer{Value: 3},
		},
		{
			`class Foo
			end`,
			"Foo",
			"Object",
			map[string]string{},
			object.NIL,
		},
		{
			`class Foo < BasicObject
			end`,
			"Foo",
			"BasicObject",
			map[string]string{},
			object.NIL,
		},
	}

	for _, tt := range tests {
		env := object.NewMainEnvironment()
		evaluated, err := testEval(tt.input, env)
		checkError(t, err)

		if !reflect.DeepEqual(evaluated, tt.expectedReturn) {
			t.Logf("Expected return object to equal\n%+#v\n\tgot\n%+#v\n", tt.expectedReturn, evaluated)
			t.Fail()
		}

		class, ok := env.Get(tt.expectedName)
		if !ok {
			t.Logf("Expected class to exist in env")
			t.Logf("Env: %+#v\n", env)
			t.FailNow()
		}

		classClass, ok := class.(object.RubyClassObject)
		if !ok {
			t.Logf("Expected class to be a object.RubyClassObject, got %T", classClass)
			t.FailNow()
		}

		superClass := classClass.SuperClass().(object.RubyClassObject)

		if superClass.Inspect() != tt.expectedSuperclass {
			t.Logf("Expected superclass %q, got %q\n", tt.expectedSuperclass, superClass.Inspect())
			t.Fail()
		}

		actualMethods := make(map[string]string)

		methods := classClass.Methods().GetAll()
		for name, method := range methods {
			if function, ok := method.(*object.Function); ok {
				actualMethods[name] = function.String()
			}
		}

		if !reflect.DeepEqual(tt.expectedMethods, actualMethods) {
			t.Logf(
				"Expected class methods to equal\n%+#v\n\tgot\n%+#v\n",
				tt.expectedMethods,
				actualMethods,
			)
			t.Fail()
		}
	}
}

func TestFunctionObject(t *testing.T) {
	type funcParam struct {
		name         string
		defaultValue object.RubyObject
	}
	t.Run("methods without receiver", func(t *testing.T) {
		tests := []struct {
			input              string
			expectedParameters []funcParam
			expectedBody       string
		}{
			{
				"def foo x; x + 2; end",
				[]funcParam{{name: "x"}},
				"(x + 2)",
			},
			{
				`def foo
				2
				end`,
				[]funcParam{},
				"2",
			},
			{
				"def foo; 2; end",
				[]funcParam{},
				"2",
			},
			{
				"def foo x = 4; 2; end",
				[]funcParam{{name: "x", defaultValue: &object.Integer{Value: 4}}},
				"2",
			},
		}

		for _, tt := range tests {
			env := object.NewEnvironment()
			env.Set("self", &object.Self{RubyObject: &object.Object{}, Name: "main"})
			evaluated, err := testEval(tt.input, env)
			checkError(t, err)
			sym, ok := evaluated.(*object.Symbol)
			if !ok {
				t.Fatalf("object is not Symbol. got=%T (%+v)", evaluated, evaluated)
			}
			if sym.Value != "foo" {
				t.Logf("Expected returned symbol to have value %q, got %q", "foo", sym.Value)
				t.Fail()
			}

			self, _ := env.Get("self")
			method, ok := self.Class().Methods().Get("foo")
			if !ok {
				t.Logf("Expected function to be added to self")
				t.Fail()
			}
			fn, ok := method.(*object.Function)
			if !ok {
				t.Logf("self method is not Function, got=%T (%+v)", method, method)
				t.Fail()
			}

			if len(fn.Parameters) != len(tt.expectedParameters) {
				t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
			}

			for i, param := range fn.Parameters {
				testParam := tt.expectedParameters[i]
				if testParam.name != param.Name {
					t.Logf("Expected parameter %d to have name %q, got %q\n", i+1, testParam.name, param.Name)
					t.Fail()
				}
				if !reflect.DeepEqual(testParam.defaultValue, param.Default) {
					t.Logf("Expected parameter %d to have default %v, got %v\n", i+1, testParam.defaultValue, param.Default)
					t.Fail()
				}
			}

			if fn.Body.String() != tt.expectedBody {
				t.Fatalf("body is not %q. got=%q", tt.expectedBody, fn.Body.String())
			}
		}
	})
	t.Run("methods with variable receiver", func(t *testing.T) {
		env := object.NewEnvironment()
		env.Set("self", &object.Self{RubyObject: &object.Object{}, Name: "main"})
		input := `a = "foo"
def a.truth
	42
end
`

		_, err := testEval(input, env)
		checkError(t, err)

		a, ok := env.Get("a")
		if !ok {
			t.Logf("Expected env to have 'a'")
			t.FailNow()
		}

		method, ok := a.Class().Methods().Get("truth")
		if !ok {
			t.Logf("Expected function to be added to 'a'")
			t.Fail()
		}
		fn, ok := method.(*object.Function)
		if !ok {
			t.Logf("method is not %T, got=%T (%+v)", fn, method, method)
			t.Fail()
		}
	})
	t.Run("methods with const receiver", func(t *testing.T) {
		env := object.NewMainEnvironment()
		env.Set("self", &object.Self{RubyObject: &object.Object{}, Name: "main"})
		input := `class A
end

def A.truth
	42
end
`

		_, err := testEval(input, env)
		checkError(t, err)

		A, ok := env.Get("A")
		if !ok {
			t.Logf("Expected env to have 'A'")
			t.FailNow()
		}

		method, ok := A.Class().Methods().Get("truth")
		if !ok {
			t.Logf("Expected function to be added to 'A'")
			t.Fail()
		}
		fn, ok := method.(*object.Function)
		if !ok {
			t.Logf("method is not %T, got=%T (%+v)", fn, method, method)
			t.Fail()
		}
	})
	t.Run("methods with self in main", func(t *testing.T) {
		env := object.NewEnvironment()
		env.Set("self", &object.Self{RubyObject: &object.Object{}, Name: "main"})
		input := `
def self.truth
	42
end
`

		_, err := testEval(input, env)
		checkError(t, err)

		self, _ := env.Get("self")

		method, ok := self.Class().Methods().Get("truth")
		if !ok {
			t.Logf("Expected function to be added to 'main'")
			t.Fail()
		}
		fn, ok := method.(*object.Function)
		if !ok {
			t.Logf("method is not %T, got=%T (%+v)", fn, method, method)
			t.Fail()
		}
	})
	t.Run("methods with self in class", func(t *testing.T) {
		env := object.NewMainEnvironment()
		input := `
class A
	def self.truth
		42
	end
end
`

		_, err := testEval(input, env)
		checkError(t, err)

		A, ok := env.Get("A")
		if !ok {
			t.Logf("Expected env to have 'A'")
			t.FailNow()
		}

		method, ok := A.Class().Methods().Get("truth")
		if !ok {
			t.Logf("Expected function to be added to 'A'")
			t.Fail()
		}
		fn, ok := method.(*object.Function)
		if !ok {
			t.Logf("method is not %T, got=%T (%+v)", fn, method, method)
			t.Fail()
		}
	})
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"def identity x; x; end; identity(5);", 5},
		{"def identity x; return x; end; identity(5);", 5},
		{"def double x; x * 2; end; double(5);", 10},
		{"def add x, y; x + y; end; add(5, 5);", 10},
		{"def add x, y; x + y; end; add(5 + 5, add(5, 5));", 20},
		{"def double x; x * 2; end; double 5;", 10},
		{"def identity x; x; end; identity 5;", 5},
		{"def foo; 3; end; foo;", 3},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		env.Set("self", &object.Self{RubyObject: &object.Object{}, Name: "main"})
		evaluated, err := testEval(tt.input, env)
		checkError(t, err)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestGlobalLiteral(t *testing.T) {
	input := `$foo = 'bar'; $foo`

	evaluated, err := testEval(input)
	checkError(t, err)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "bar" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated, err := testEval(input)
	checkError(t, err)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated, err := testEval(input)
	checkError(t, err)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestSymbolLiteral(t *testing.T) {
	input := `:foobar;`

	evaluated, err := testEval(input)
	checkError(t, err)
	sym, ok := evaluated.(*object.Symbol)
	if !ok {
		t.Fatalf("object is not Symbol. got=%T (%+v)", evaluated, evaluated)
	}

	if sym.Value != "foobar" {
		t.Errorf("Symbol has wrong value. got=%q", sym.Value)
	}
}

func TestMethodCalls(t *testing.T) {
	input := "x = 2; x.foo :bar"

	evaluated, err := testEval(input)

	if err == nil {
		t.Logf("Expected error, got %T:%s\n", evaluated, evaluated)
		t.Fail()
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated, err := testEval(input)
	checkError(t, err)

	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}
	if len(result.Elements) != 3 {
		t.Fatalf(
			"array has wrong num of elements. got=%d",
			len(result.Elements),
		)
	}
	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"myArray = [1, 2, 3]; i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			3,
		},
		{
			"[1, 2, 3][-2]",
			2,
		},
		{
			"[1, 2, 3][-3]",
			1,
		},
		{
			"[1, 2, 3][-4]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		checkError(t, err)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}

func TestNilExpression(t *testing.T) {
	input := "nil"
	evaluated, err := testEval(input)
	checkError(t, err)
	testNilObject(t, evaluated)
}

func TestSelfExpression(t *testing.T) {
	input := "self"

	env := object.NewMainEnvironment()
	env.Set("self", &object.Self{RubyObject: &object.Integer{Value: 3}, Name: "3"})
	evaluated, err := testEval(input, env)
	checkError(t, err)

	self, ok := evaluated.(*object.Self)
	if !ok {
		t.Logf("Expected evaluated object to be object.Self, got=%T", evaluated)
		t.Fail()
	}

	testIntegerObject(t, self.RubyObject, 3)
}

func TestHashLiteral(t *testing.T) {
	input := `{"foo" => 42, :bar => 2, true => false, nil => true, 2 => 2}`

	env := object.NewMainEnvironment()
	evaluated, err := testEval(input, env)
	checkError(t, err)

	hash, ok := evaluated.(*object.Hash)
	if !ok {
		t.Logf("Expected evaluated object to be *object.Hash, got=%T", evaluated)
		t.FailNow()
	}

	expected := map[string]object.RubyObject{
		"foo":  &object.Integer{Value: 42},
		":bar": &object.Integer{Value: 2},
		"true": object.FALSE,
		"nil":  object.TRUE,
		"2":    &object.Integer{Value: 2},
	}

	actual := make(map[string]object.RubyObject)
	for k, v := range hash.Map() {
		actual[k.Inspect()] = v
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected hash to equal\n%s\n\tgot\n%s\n", expected, actual)
		t.Fail()
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"{'foo' => 1, 'bar' => 2, 'qux' => 3}['foo']",
			1,
		},
		{
			"{'foo' => 1, 'bar' => 2, 'qux' => 3}['bar']",
			2,
		},
		{
			"{'foo' => 1, 'bar' => 2, 'qux' => 3}['qux']",
			3,
		},
		{
			"i = 'foo'; {'foo'=>1}[i];",
			1,
		},
		{
			"{1=>1, 2=>2, 3=>3}[1 + 1];",
			2,
		},
		{
			"myHash = {1=>1, 2=>2, 3=>3}; myHash[2];",
			2,
		},
		{
			"myHash = {0=>1, 1=>2, 2=>3}; myHash[0] + myHash[1] + myHash[2];",
			6,
		},
		{
			"myHash = {0=>1, 1=>2, 2=>3}; i = myHash[0]; myHash[i]",
			2,
		},
		{
			"{0=>1, 1=>2, 2=>3}[3]",
			nil,
		},
		{
			"{0=>1, 1=>2, 2=>3}[-1]",
			nil,
		},
		{
			"{:foo => 1, :bar => 2, :qux => 3}[:qux]",
			3,
		},
		{
			"{'foo' =>1, true => 2, false => 3}[true]",
			2,
		},
		{
			"{nil =>1, :qux => 2, 3=>3}[nil]",
			1,
		},
	}

	for _, tt := range tests {
		evaluated, err := testEval(tt.input)
		checkError(t, err)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNilObject(t, evaluated)
		}
	}
}

func TestKeyword__File__(t *testing.T) {
	input := "__FILE__"

	env := object.NewEnvironment()
	program, err := parser.ParseFile(token.NewFileSet(), "some_file.rb", input, 0)
	checkError(t, err)
	evaluated, err := Eval(program, env)
	checkError(t, err)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Logf("Expected evaluated to be *object.String, got %T\n", evaluated)
		t.FailNow()
	}

	expected := "some_file.rb"

	if expected != str.Value {
		t.Logf("Expected __FILE__ to equal %q, got %q\n", expected, str.Value)
		t.Fail()
	}
}

func testExceptionObject(t *testing.T, obj object.RubyObject, errorMessage string) {
	t.Helper()
	if !IsError(obj) {
		t.Logf("Expected error or exception, got %T", obj)
		t.Fail()
	}

	actual := obj.Inspect()

	if errorMessage != actual {
		t.Logf("Expected obj to stringify to %q, got %q", errorMessage, actual)
		t.Fail()
	}
}

func testNilObject(t *testing.T, obj object.RubyObject) bool {
	t.Helper()
	if obj != object.NIL {
		t.Errorf("object is not NIL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testEval(input string, context ...object.Environment) (object.RubyObject, error) {
	env := object.NewEnvironment()
	if len(context) > 0 {
		env = context[0]
	}
	program, err := parser.ParseFile(token.NewFileSet(), "", input, parser.ParseComments)
	if err != nil {
		return nil, object.NewSyntaxError(err)
	}
	return Eval(program, env)
}

func checkError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}
}

func testObject(t *testing.T, exp object.RubyObject, expected interface{}) bool {
	t.Helper()
	switch v := expected.(type) {
	case int:
		return testIntegerObject(t, exp, int64(v))
	case int64:
		return testIntegerObject(t, exp, v)
	case string:
		if strings.HasPrefix(v, ":") {
			return testSymbolObject(t, exp, strings.TrimPrefix(v, ":"))
		}
		return testStringObject(t, exp, v)
	case bool:
		return testBooleanObject(t, exp, v)
	case map[string]string:
		return testHashObject(t, exp, v)
	case []string:
		return testArrayObject(t, exp, v)
	case nil:
		return true
	}
	t.Errorf("type of object not handled. got=%T", exp)
	return false
}

func testBooleanObject(t *testing.T, obj object.RubyObject, expected bool) bool {
	t.Helper()
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf(
			"object is not Boolean. got=%T (%+v)",
			obj,
			obj,
		)
		return false
	}
	if result.Value != expected {
		t.Errorf(
			"object has wrong value. got=%v, want=%v",
			result.Value,
			expected,
		)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.RubyObject, expected int64) bool {
	t.Helper()
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf(
			"object is not Integer. got=%T (%+v)",
			obj,
			obj,
		)
		return false
	}
	if result.Value != expected {
		t.Errorf(
			"object has wrong value. got=%v, want=%v",
			result.Value,
			expected,
		)
		return false
	}
	return true
}

func testSymbolObject(t *testing.T, obj object.RubyObject, expected string) bool {
	t.Helper()
	result, ok := obj.(*object.Symbol)
	if !ok {
		t.Errorf(
			"object is not Symbol. got=%T (%+v)",
			obj,
			obj,
		)
		return false
	}
	if result.Value != expected {
		t.Errorf(
			"object has wrong value. got=%v, want=%v",
			result.Value,
			expected,
		)
		return false
	}
	return true
}

func testStringObject(t *testing.T, obj object.RubyObject, expected string) bool {
	t.Helper()
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf(
			"object is not String. got=%T (%+v)",
			obj,
			obj,
		)
		return false
	}
	if result.Value != expected {
		t.Errorf(
			"object has wrong value. got=%v, want=%v",
			result.Value,
			expected,
		)
		return false
	}
	return true
}

func testHashObject(t *testing.T, obj object.RubyObject, expected map[string]string) bool {
	t.Helper()
	result, ok := obj.(*object.Hash)
	if !ok {
		t.Errorf(
			"object is not Hash. got=%T (%+v)",
			obj,
			obj,
		)
		return false
	}
	hashMap := make(map[string]string)
	for k, v := range result.Map() {
		hashMap[k.Inspect()] = v.Inspect()
	}
	if !reflect.DeepEqual(hashMap, expected) {
		t.Errorf(
			"object has wrong value. got=%v, want=%v",
			hashMap,
			expected,
		)
		return false
	}
	return true
}

func testArrayObject(t *testing.T, obj object.RubyObject, expected []string) bool {
	t.Helper()
	result, ok := obj.(*object.Array)
	if !ok {
		t.Errorf(
			"object is not Array. got=%T (%+v)",
			obj,
			obj,
		)
		return false
	}
	array := make([]string, len(result.Elements))
	for i, v := range result.Elements {
		array[i] = v.Inspect()
	}
	if !reflect.DeepEqual(array, expected) {
		t.Errorf(
			"object has wrong value. got=%v, want=%v",
			array,
			expected,
		)
		return false
	}
	return true
}
