package evaluator

import (
	"reflect"
	"testing"

	"github.com/goruby/goruby/lexer"
	"github.com/goruby/goruby/object"
	"github.com/goruby/goruby/parser"
)

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
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
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
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
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
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
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
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
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
		evaluated := testEval(tt.input)
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
			"Exception: type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"Exception: type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"Exception: unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"Exception: unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true + false + true + false;",
			"Exception: unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"Exception: unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`"Hello" - "World"`,
			"Exception: unknown operator: STRING - STRING",
		},
		{
			"if (10 > 1); true + false; end",
			"Exception: unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1); true + false; end",
			"Exception: unknown operator: BOOLEAN + BOOLEAN",
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
			"Exception: unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"NameError: undefined local variable or method `foobar' for :Object",
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
		env.Set("self", &object.Object{})
		evaluated := testEval(tt.input, env)

		ok := IsError(evaluated)
		if !ok {
			t.Errorf(
				"no error object returned. got=%T(%+v)",
				evaluated,
				evaluated,
			)
		}

		if evaluated.Inspect() != tt.expectedMessage {
			t.Errorf(
				"wrong error message. expected=%q, got=%q",
				tt.expectedMessage,
				evaluated.Inspect(),
			)
		}
	}
}

func TestVariableAssignmentExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"a = 5; a;", 5},
		{"a = 5 * 5; a;", 25},
		{"a = 5; b = a; b;", 5},
		{"a = 5; b = a; c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	tests := []struct {
		input              string
		expectedParameters []string
		expectedBody       string
	}{
		{
			"def foo x; x + 2; end",
			[]string{"x"},
			"(x + 2)",
		},
		{
			`def foo
				2
			end`,
			[]string{},
			"2",
		},
		{
			"def foo; 2; end",
			[]string{},
			"2",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		fn, ok := evaluated.(*object.Function)
		if !ok {
			t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
		}

		if len(fn.Parameters) != len(tt.expectedParameters) {
			t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
		}

		parameters := make([]string, len(fn.Parameters))
		for i, param := range fn.Parameters {
			parameters[i] = param.String()
		}

		if !reflect.DeepEqual(tt.expectedParameters, parameters) {
			t.Fatalf("parameters are not %v. got=%v", tt.expectedParameters, parameters)
		}

		if fn.Body.String() != tt.expectedBody {
			t.Fatalf("body is not %q. got=%q", tt.expectedBody, fn.Body.String())
		}
	}
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
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
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

	evaluated := testEval(input)
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

	evaluated := testEval(input)
	sym, ok := evaluated.(*object.Symbol)
	if !ok {
		t.Fatalf("object is not Symbol. got=%T (%+v)", evaluated, evaluated)
	}

	if sym.Value != "foobar" {
		t.Errorf("Symbol has wrong value. got=%q", sym.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`puts;`, nil},
		{`puts "foo";`, nil},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, object.NewMainEnvironment())

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNilObject(t, evaluated)
		case string:
			ok := IsError(evaluated)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if evaluated.Inspect() != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, evaluated.Inspect())
			}
		}
	}
}

func TestMethodCalls(t *testing.T) {
	input := "x = 2; x.foo :bar"

	evaluated := testEval(input)

	if !IsError(evaluated) {
		t.Logf("Expected error, got %T:%s\n", evaluated, evaluated)
		t.Fail()
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
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
			nil,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
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
	evaluated := testEval(input)
	testNilObject(t, evaluated)
}

func TestRequireExpression(t *testing.T) {
	t.Run("simple require", func(t *testing.T) {
		input := `require "testfile.rb"
		x + 2
		`

		evaluated := testEval(input)

		testIntegerObject(t, evaluated, int64(7))
	})
	t.Run("simple require no filetype", func(t *testing.T) {
		input := `require "testfile"
		x + 2
		`

		evaluated := testEval(input)

		testIntegerObject(t, evaluated, int64(7))
	})
	t.Run("require return value", func(t *testing.T) {
		input := `require "testfile.rb"
		`

		evaluated := testEval(input)

		testBooleanObject(t, evaluated, true)
	})
	t.Run("require appends to $LOADED_FEATURES", func(t *testing.T) {
		input := `require "testfile.rb"
		`

		env := object.NewEnvironment()
		testEval(input, env)

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")
		if !ok {
			t.Logf("Expected env to contain $LOADED_FEATURES object")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*object.Array)
		if !ok {
			t.Logf("Expected loadedFeatures to be an array, got %T", loadedFeatures)
			t.Fail()
		}

		expectedLen := 1
		actualLen := len(arr.Elements)

		if expectedLen != actualLen {
			t.Logf("Expected array to have %d elements, got %d", expectedLen, actualLen)
			t.Fail()
		}

		expectedElemValue := "testfile.rb"
		actualElemValue := arr.Elements[0].Inspect()

		if expectedElemValue != actualElemValue {
			t.Logf("Expected elem value to equal %s, got %s", expectedElemValue, actualElemValue)
			t.Fail()
		}
	})
	t.Run("require only parses the files once", func(t *testing.T) {
		input := `require "testfile_require_once.rb"
			x
		`

		evaluated := testEval(input)

		testIntegerObject(t, evaluated, int64(7))
	})
	t.Run("require returns false if already required", func(t *testing.T) {
		input := `require "testfile"
			require "testfile"
		`

		evaluated := testEval(input)

		testBooleanObject(t, evaluated, false)
	})
	t.Run("recursive require", func(t *testing.T) {
		input := `require "testfile_recursive_require.rb"
		x + 2
		`

		evaluated := testEval(input)

		testIntegerObject(t, evaluated, int64(7))
	})
	t.Run("syntax error in file", func(t *testing.T) {
		input := `require "testfile_syntax_error.rb"
		`

		evaluated := testEval(input)

		expected := "SyntaxError: syntax error, Parsing errors:\n\texpected next token to be of type [END], got EOF instead\n"
		testExceptionObject(t, evaluated, expected)
	})
	t.Run("file not found", func(t *testing.T) {
		input := `require "this/file/does/not/exist"
		`

		evaluated := testEval(input)

		expected := "LoadError: no such file to load -- this/file/does/not/exist"
		testExceptionObject(t, evaluated, expected)
	})
}

func testExceptionObject(t *testing.T, obj object.RubyObject, errorMessage string) {
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
	if obj != object.NIL {
		t.Errorf("object is not NIL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testEval(input string, context ...object.Environment) object.RubyObject {
	env := object.NewEnvironment()
	for _, e := range context {
		env = object.NewEnclosedEnvironment(e)
	}
	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		return object.NewSyntaxError(err.Error())
	}
	return Eval(program, env)
}

func testBooleanObject(t *testing.T, obj object.RubyObject, expected bool) bool {
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
			"object has wrong value. got=%t, want=%t",
			result.Value,
			expected,
		)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.RubyObject, expected int64) bool {
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
			"object has wrong value. got=%d, want=%d",
			result.Value,
			expected,
		)
		return false
	}
	return true
}
