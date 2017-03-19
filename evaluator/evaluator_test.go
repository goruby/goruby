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
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"true + false + true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			"if (10 > 1); true + false; end",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1); true + false; end",
			"unknown operator: BOOLEAN + BOOLEAN",
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
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
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
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf(
				"no error object returned. got=%T(%+v)",
				evaluated,
				evaluated,
			)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf(
				"wrong error message. expected=%q, got=%q",
				tt.expectedMessage,
				errObj.Message,
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
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		}
	}
}

func testNilObject(t *testing.T, obj object.Object) bool {
	if obj != object.NIL {
		t.Errorf("object is not NIL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testEval(input string, context ...*object.Environment) object.Object {
	env := object.NewEnvironment()
	for _, e := range context {
		env = object.NewEnclosedEnvironment(e)
	}
	l := lexer.New(input)
	p := parser.New(l)
	program, err := p.ParseProgram()
	if err != nil {
		return newError(err.Error())
	}
	return Eval(program, env)
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
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
			"object has wrong value. got=%d, want=%d",
			result.Value,
			expected,
		)
		return false
	}
	return true
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
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
