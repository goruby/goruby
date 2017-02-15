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
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalBangOperator(t *testing.T) {
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
			testNullObject(t, evaluated)
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
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NIL {
		t.Errorf("object is not NIL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
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
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
