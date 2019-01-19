package parser

import (
	"flag"
	"fmt"
	gotoken "go/token"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/goruby/goruby/ast"
	"github.com/goruby/goruby/token"
	"github.com/pkg/errors"
)

var parseMode Mode = ParseComments

func TestMain(m *testing.M) {
	mode := flag.String("parser.mode", "ParseComments", "parser.mode=ParseComments")
	flag.Parse()
	var ok bool
	parseMode, ok = parseModes[*mode]
	if !ok {
		fmt.Printf("Unknown parse mode %s\n", *mode)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestBlockCapture(t *testing.T) {
	tests := []struct {
		desc   string
		input  string
		result ast.Expression
		err    error
	}{
		{
			desc:  "block capture in func params as only argument",
			input: "def foo &block; end",
			result: &ast.FunctionLiteral{
				Name: &ast.Identifier{Value: "foo"},
				CapturedBlock: &ast.BlockCapture{
					Name: &ast.Identifier{Value: "block"},
				},
			},
		},
		{
			desc:  "block capture in func params as last arguments",
			input: "def foo x, y, &block; end",
			result: &ast.FunctionLiteral{
				Name: &ast.Identifier{Value: "foo"},
				Parameters: []*ast.FunctionParameter{
					&ast.FunctionParameter{Name: &ast.Identifier{Value: "x"}},
					&ast.FunctionParameter{Name: &ast.Identifier{Value: "y"}},
				},
				CapturedBlock: &ast.BlockCapture{
					Name: &ast.Identifier{Value: "block"},
				},
			},
		},
		{
			desc:   "block capture in func params not last arguments",
			input:  "def foo x, &block, y; end",
			result: nil,
			err: &unexpectedTokenError{
				expectedTokens: []token.Type{token.NEWLINE, token.SEMICOLON},
				actualToken:    token.COMMA,
			},
		},
		{
			desc:   "block capture in func params on integer",
			input:  "def foo &2; end",
			result: nil,
			err: &unexpectedTokenError{
				expectedTokens: []token.Type{token.IDENT},
				actualToken:    token.INT,
			},
		},
		{
			desc: "block capture only statement in func body",
			input: `
			def foo
				&block
			end`,
			result: &ast.FunctionLiteral{
				Name:       &ast.Identifier{Value: "foo"},
				Parameters: []*ast.FunctionParameter{},
				Body: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ExpressionStatement{
							Expression: &ast.BlockCapture{
								Name: &ast.Identifier{Value: "block"},
							},
						},
					},
				},
			},
		},
		{
			desc: "block capture as arg on call in func body",
			input: `
			def foo
				each &block
			end`,
			result: &ast.FunctionLiteral{
				Name:       &ast.Identifier{Value: "foo"},
				Parameters: []*ast.FunctionParameter{},
				Body: &ast.BlockStatement{
					Statements: []ast.Statement{
						&ast.ExpressionStatement{
							Expression: &ast.ContextCallExpression{
								Function: &ast.Identifier{Value: "each"},
								Arguments: []ast.Expression{
									&ast.BlockCapture{
										Name: &ast.Identifier{Value: "block"},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			expr, err := parseExpression(tt.input)
			compareFirstParserError(t, tt.err, err)

			if !ast.Equal(expr, tt.result) {
				t.Logf("Expected AST node to equal %v, got %v", tt.result, expr)
				t.Fail()
			}
		})
	}
}

func TestAssignment(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		leftType  reflect.Type
		rightType reflect.Type
	}{
		{
			name:      "hash index assignment",
			input:     `x[:foo] = 3`,
			leftType:  reflect.TypeOf(&ast.IndexExpression{}),
			rightType: reflect.TypeOf(&ast.IntegerLiteral{}),
		},
		{
			name:      "instance varibale",
			input:     `@x = 3`,
			leftType:  reflect.TypeOf(&ast.InstanceVariable{}),
			rightType: reflect.TypeOf(&ast.IntegerLiteral{}),
		},
		{
			name:      "local varibale",
			input:     `x = 3`,
			leftType:  reflect.TypeOf(&ast.Identifier{}),
			rightType: reflect.TypeOf(&ast.IntegerLiteral{}),
		},
		{
			name:      "method call with block on rhs",
			input:     `x = foo { |x| }`,
			leftType:  reflect.TypeOf(&ast.Identifier{}),
			rightType: reflect.TypeOf(&ast.ContextCallExpression{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			if len(program.Statements) != 1 {
				t.Fatalf(
					"program.Statements does not contain 1 statements. got=%d",
					len(program.Statements),
				)
			}
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf(
					"program.Statements[0] is not ast.ExpressionStatement. got=%T",
					program.Statements[0],
				)
			}

			assign, ok := stmt.Expression.(*ast.Assignment)
			if !ok {
				t.Fatalf(
					"stmt.Expression is not *ast.Assignment. got=%T",
					stmt.Expression,
				)
			}

			{
				actual := reflect.TypeOf(assign.Left)
				if tt.leftType != actual {
					t.Fatalf(
						"assign.Left is not %v. got=%v",
						tt.leftType,
						actual,
					)
				}
			}

			{
				actual := reflect.TypeOf(assign.Right)
				if tt.rightType != actual {
					t.Fatalf(
						"assign.Right is not %v. got=%v",
						tt.rightType,
						actual,
					)
				}
			}
		})
	}
}

func TestAssignmentOperator(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		leftType      reflect.Type
		rightOperator string
	}{
		{
			name:          "-=",
			input:         `x -= 3`,
			leftType:      reflect.TypeOf(&ast.Identifier{}),
			rightOperator: "-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			if len(program.Statements) != 1 {
				t.Logf(
					"program.Statements does not contain 1 statements. got=%d",
					len(program.Statements),
				)
				t.Logf(
					"program.Statements: %v",
					program.Statements,
				)
				t.FailNow()
			}
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf(
					"program.Statements[0] is not ast.ExpressionStatement. got=%T",
					program.Statements[0],
				)
			}

			assign, ok := stmt.Expression.(*ast.Assignment)
			if !ok {
				t.Fatalf(
					"stmt.Expression is not *ast.Assignment. got=%T",
					stmt.Expression,
				)
			}

			{
				actual := reflect.TypeOf(assign.Left)
				if tt.leftType != actual {
					t.Fatalf(
						"assign.Left is not %v. got=%v",
						tt.leftType,
						actual,
					)
				}
			}

			{
				infix, ok := assign.Right.(*ast.InfixExpression)
				if !ok {
					t.Logf("Expected right assign type to be %T, got %T", infix, assign.Right)
					t.FailNow()
				}

				if infix.Operator != tt.rightOperator {
					t.Logf(
						"Expected right assign infix operator to be %q, got %q",
						tt.rightOperator,
						infix.Operator,
					)
					t.Fail()
				}
			}
		})
	}
}

func TestVariableExpression(t *testing.T) {
	t.Run("valid variable expressions", func(t *testing.T) {
		tests := []struct {
			input              string
			expectedIdentifier string
			expectedValue      string
		}{
			{"x = 5;", "x", "5"},
			{"x = 5_0;", "x", "50"},
			{"y = true;", "y", "true"},
			{"foobar = y;", "foobar", "y"},
			{"foobar = (12 + 2 * bar) - x;", "foobar", "((12 + (2 * bar)) - x)"},
		}

		for _, tt := range tests {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			if len(program.Statements) != 1 {
				t.Fatalf(
					"program.Statements does not contain 1 statements. got=%d",
					len(program.Statements),
				)
			}
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf(
					"program.Statements[0] is not ast.ExpressionStatement. got=%T",
					program.Statements[0],
				)
			}

			variable, ok := stmt.Expression.(*ast.Assignment)
			if !ok {
				t.Fatalf(
					"stmt.Expression is not *ast.VariableAssignment. got=%T",
					stmt.Expression,
				)
			}

			if !testIdentifier(t, variable.Left, tt.expectedIdentifier) {
				return
			}

			val := variable.Right.String()

			if val != tt.expectedValue {
				t.Logf(
					"Expected variable value to equal %s, got %s\n",
					tt.expectedValue,
					val,
				)
				t.Fail()
			}
		}
	})
	t.Run("const assignment within function", func(t *testing.T) {
		tests := []struct {
			desc  string
			input string
			err   error
		}{
			{
				desc: "single const assign",
				input: `
				def foo
					Ten = 10
				end`,
				err: fmt.Errorf("dynamic constant assignment"),
			},
			{
				desc: "const assign as multiassign",
				input: `
				def foo
					x, Ten = 10, 20
				end`,
				err: fmt.Errorf("dynamic constant assignment"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {

				_, errs := parseExpression(tt.input)

				if errs == nil {
					t.Logf("Expected error, got nil")
					t.FailNow()
				}

				errors := errs.errors
				if len(errors) != 1 {
					t.Logf("Exected one error, got %d", len(errors))
					t.FailNow()
				}

				if !reflect.DeepEqual(errors[0], tt.err) {
					t.Logf("Expected error to equal\n%v\n\tgot\n%v\n", tt.err, errors[0])
					t.Fail()
				}
			})
		}
	})
}

func TestWhileExpression(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name: "with explicit do",
			input: `
			while x < y do
				x += x
			end`,
		},
		{
			name: "without explicit do",
			input: `
			while x < y
				x += x
			end`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			if len(program.Statements) != 1 {
				t.Logf(
					"program.Body does not contain %d statements. got=%d\n",
					1,
					len(program.Statements),
				)
				t.Logf("%s\n", program.Statements)
				t.FailNow()
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf(
					"program.Statements[0] is not ast.ExpressionStatement. got=%T",
					program.Statements[0],
				)
			}

			exp, ok := stmt.Expression.(*ast.LoopExpression)
			if !ok {
				t.Fatalf(
					"stmt.Expression is not %T. got=%T",
					exp,
					stmt.Expression,
				)
			}
		})
	}
}

func TestGlobalAssignment(t *testing.T) {
	input := "$foo = 3"

	program, err := parseSource(input)
	checkParserErrors(t, err)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Statements does not contain 1 statements. got=%d",
			len(program.Statements),
		)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	variable, ok := stmt.Expression.(*ast.Assignment)
	if !ok {
		t.Fatalf(
			"stmt.Expression is not %T. got=%T",
			variable,
			stmt.Expression,
		)
	}

	expectedGlobal := "$foo"

	if !testGlobal(t, variable.Left, expectedGlobal) {
		return
	}

	val := variable.Right.String()

	expectedValue := "3"

	if val != expectedValue {
		t.Logf(
			"Expected variable value to equal %s, got %s\n",
			expectedValue,
			val,
		)
		t.Fail()
	}
}

func TestParseMultiAssignment(t *testing.T) {
	tests := []struct {
		input     string
		variables []string
		values    []string
	}{
		{
			input:     "x, y, z = 3, 4, 5;",
			variables: []string{"x", "y", "z"},
			values:    []string{"3", "4", "5"},
		},
		{
			input:     "x, y = 3, 4;",
			variables: []string{"x", "y"},
			values:    []string{"3", "4"},
		},
		{
			input:     "x, y, z = 3, 4;",
			variables: []string{"x", "y", "z"},
			values:    []string{"3", "4"},
		},
		{
			input:     "x, y, z = 3;",
			variables: []string{"x", "y", "z"},
			values:    []string{"3"},
		},
		{
			input:     "x[0], @y, $z, A = 3, 4, 5, 6;",
			variables: []string{"(x[0])", "@y", "$z", "A"},
			values:    []string{"3", "4", "5", "6"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr, err := parseExpression(tt.input)
			checkParserErrors(t, err)

			assign, ok := expr.(*ast.Assignment)
			if !ok {
				t.Logf("Expected expression to be %T, got %T\n", assign, expr)
				t.FailNow()
			}

			left, ok := assign.Left.(ast.ExpressionList)
			if !ok {
				t.Logf("Expected left to be %T, got %T\n", left, assign.Left)
				t.FailNow()
			}

			actualVars := make([]string, len(left))
			for i, v := range left {
				actualVars[i] = v.String()
			}

			if !reflect.DeepEqual(tt.variables, actualVars) {
				t.Logf("Expected variable identifiers to equal %s, got %s\n", tt.variables, actualVars)
				t.Fail()
			}

			if !reflect.DeepEqual(strings.Join(tt.values, ", "), assign.Right.String()) {
				t.Logf("Expected variable values to equal %s, got %s\n", tt.values, assign.Right.String())
				t.Fail()
			}
		})
	}
}

func TestInstanceVariable(t *testing.T) {
	input := "@foo"

	program, err := parseSource(input)
	checkParserErrors(t, err)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program.Statements does not contain 1 statements. got=%d",
			len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}
	instVar, ok := stmt.Expression.(*ast.InstanceVariable)
	if !ok {
		t.Fatalf("Expression not %T. got=%T", instVar, stmt.Expression)
	}

	testLiteralExpression(t, instVar.Name, "foo")
}

func TestExceptionHandling(t *testing.T) {
	type rescue struct {
		classes   []string
		exception string
		body      string
	}
	tests := []struct {
		input   string
		body    string
		rescues []rescue
	}{
		{
			input: `
begin
end
`,
			body: "",
		},
		{
			input: `
begin
	2
end
`,
			body: "2",
		},
		{
			input: `
begin
	2
rescue
	3
end
`,
			body:    "2",
			rescues: []rescue{{body: "3"}},
		},
		{
			input: `
begin
	2
rescue Error
	3
end
`,
			body:    "2",
			rescues: []rescue{{classes: []string{"Error"}, body: "3"}},
		},
		{
			input: `
begin
	2
rescue Error, StandardError
	3
end
`,
			body:    "2",
			rescues: []rescue{{classes: []string{"Error", "StandardError"}, body: "3"}},
		},
		{
			input: `
begin
	2
rescue => e
	3
end
`,
			body:    "2",
			rescues: []rescue{{body: "3", exception: "e"}},
		},
		{
			input: `
begin
	2
rescue Error => e
	3
end
`,
			body:    "2",
			rescues: []rescue{{classes: []string{"Error"}, body: "3", exception: "e"}},
		},
	}

	for _, tt := range tests {
		program, err := parseSource(tt.input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain 1 statements. got=%d",
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}
		begin, ok := stmt.Expression.(*ast.ExceptionHandlingBlock)
		if !ok {
			t.Fatalf("Expression not %T. got=%T", begin, stmt.Expression)
		}

		body := begin.TryBody.String()
		if body != tt.body {
			t.Logf("Expected TryBody to equal\n%s\n\tgot\n%s\n", tt.body, body)
			t.Fail()
		}

		if len(begin.Rescues) != len(tt.rescues) {
			t.Logf("Expected %d rescue blocks, got %d\n", len(tt.rescues), len(begin.Rescues))
			t.Fail()
		}

		var rescues []rescue
		for _, r := range begin.Rescues {
			re := rescue{body: r.Body.String()}
			for _, ec := range r.ExceptionClasses {
				re.classes = append(re.classes, ec.String())
			}
			if r.Exception != nil {
				re.exception = r.Exception.String()
			}
			rescues = append(rescues, re)
		}

		if !reflect.DeepEqual(tt.rescues, rescues) {
			t.Logf("Expected rescues to equal\n%s\n\tgot\n%s\n", tt.rescues, rescues)
			t.Fail()
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
		{"return 3, 5, 8;", []string{"3", "5", "8"}},
	}

	for _, tt := range tests {
		program, err := parseSource(tt.input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain 1 statements. got=%d",
				len(program.Statements),
			)
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf(
				"returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral(),
			)
		}
		if !testLiteralExpression(t, returnStmt.ReturnValue, tt.expectedValue) {
			t.Fail()
		}
	}
}

func TestParseComment(t *testing.T) {
	t.Run("line comment newline", func(t *testing.T) {
		tests := []struct {
			input        string
			commentValue string
		}{
			{
				input:        "# a comment\n",
				commentValue: " a comment",
			},
			{
				input:        "# a comment",
				commentValue: " a comment",
			},
			{
				input:        "# a comment;",
				commentValue: " a comment;",
			},
		}

		for _, tt := range tests {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			if len(program.Statements) != 1 {
				t.Fatalf(
					"program has not enough statements. got=%d",
					len(program.Statements),
				)
			}

			comment, ok := program.Statements[0].(*ast.Comment)
			if !ok {
				t.Logf("Expected program.Statements[0] to be %T, got %T\n", comment, program.Statements[0])
				t.FailNow()
			}

			if comment.Value != tt.commentValue {
				t.Logf("Expected comment value to equal %q, got %q\n", tt.commentValue, comment.Value)
				t.Fail()
			}
		}
	})
	t.Run("inline comment", func(t *testing.T) {
		tests := []struct {
			input        string
			commentValue string
		}{
			{
				input:        "foo # a comment\n",
				commentValue: " a comment",
			},
			{
				input:        "foo # a comment",
				commentValue: " a comment",
			},
			{
				input:        "foo # a comment;",
				commentValue: " a comment;",
			},
		}

		for _, tt := range tests {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			if len(program.Statements) != 2 {
				t.Fatalf(
					"program has not enough statements. got=%d",
					len(program.Statements),
				)
			}

			comment, ok := program.Statements[1].(*ast.Comment)
			if !ok {
				t.Logf("Expected program.Statements[1] to be %T, got %T\n", comment, program.Statements[1])
				t.FailNow()
			}

			if comment.Value != tt.commentValue {
				t.Logf("Expected comment value to equal %q, got %q\n", tt.commentValue, comment.Value)
				t.Fail()
			}
		}
	})
}

func TestIdentifierExpression(t *testing.T) {
	t.Run("local variable", func(t *testing.T) {
		input := "foobar;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program has not enough statements. got=%d",
				len(program.Statements),
			)
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		ident, ok := stmt.Expression.(*ast.Identifier)
		if !ok {
			t.Fatalf("expression not *ast.Identifier. got=%T", stmt.Expression)
		}
		if ident.Value != "foobar" {
			t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
		}
		if ident.TokenLiteral() != "foobar" {
			t.Errorf(
				"ident.TokenLiteral not %s. got=%s", "foobar",
				ident.TokenLiteral(),
			)
		}
	})
	t.Run("constant", func(t *testing.T) {
		input := "Foobar;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program has not enough statements. got=%d",
				len(program.Statements),
			)
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		ident, ok := stmt.Expression.(*ast.Identifier)
		if !ok {
			t.Fatalf("expression not *ast.Identifier. got=%T", stmt.Expression)
		}
		if ident.Value != "Foobar" {
			t.Errorf("ident.Value not %s. got=%s", "Foobar", ident.Value)
		}
		if ident.TokenLiteral() != "Foobar" {
			t.Errorf(
				"ident.TokenLiteral not %s. got=%s", "Foobar",
				ident.TokenLiteral(),
			)
		}
	})
}

func TestGlobalExpression(t *testing.T) {
	input := "$foobar;"

	program, err := parseSource(input)
	checkParserErrors(t, err)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program has not enough statements. got=%d",
			len(program.Statements),
		)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	global, ok := stmt.Expression.(*ast.Global)
	if !ok {
		t.Fatalf("expression not *ast.Global. got=%T", stmt.Expression)
	}
	if global.Value != "$foobar" {
		t.Errorf("ident.Value not %s. got=%s", "$foobar", global.Value)
	}
	if global.TokenLiteral() != "$foobar" {
		t.Errorf(
			"global.TokenLiteral not %s. got=%s", "$foobar",
			global.TokenLiteral(),
		)
	}
}

func TestScopedIdentifierExpression(t *testing.T) {
	input := "A::B"

	program, err := parseSource(input)
	checkParserErrors(t, err)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program has not enough statements. got=%d",
			len(program.Statements),
		)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	_, ok = stmt.Expression.(*ast.ScopedIdentifier)
	if !ok {
		t.Logf("Expected expression to be *ast.ScopedIdentifier, got %T", stmt.Expression)
		t.Fail()
	}
}

func TestSelfExpression(t *testing.T) {
	input := "self;"

	program, err := parseSource(input)
	checkParserErrors(t, err)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program has not enough statements. got=%d",
			len(program.Statements),
		)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	_, ok = stmt.Expression.(*ast.Self)
	if !ok {
		t.Fatalf("expression not *ast.Self. got=%T", stmt.Expression)
	}
}

func TestKeyword__FILE__(t *testing.T) {
	t.Run("keyword found", func(t *testing.T) {
		input := "__FILE__;"

		program, err := ParseFile(gotoken.NewFileSet(), "a_filename.rb", input, 0)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program has not enough statements. got=%d",
				len(program.Statements),
			)
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		file, ok := stmt.Expression.(*ast.Keyword__FILE__)
		if !ok {
			t.Fatalf("expression not *ast.Keyword__FILE__. got=%T", stmt.Expression)
		}

		expected := "a_filename.rb"

		if expected != file.Filename {
			t.Logf("Expected filename to equal %q, got %q\n", expected, file.Filename)
			t.Fail()
		}
	})
	t.Run("assignment to keyword", func(t *testing.T) {
		input := "__FILE__ = 42;"

		_, err := parseSource(input)

		expected := "1:9: Can't assign to __FILE__"

		parserErrors := err.errors
		if len(parserErrors) != 1 {
			t.Logf("Expected one error, got %d\n", len(parserErrors))
			t.Logf("Errors: %v\n", err)
			t.FailNow()
		}

		if expected != parserErrors[0].Error() {
			t.Logf("Expected error to equal\n%q\n\tgot\n%q\n", expected, parserErrors[0].Error())
			t.Fail()
		}

	})
}

func TestYieldExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:        "yield;",
			expectedArgs: []string{},
		},
		{
			input:        "yield 1, 2 + 3;",
			expectedArgs: []string{"1", "(2 + 3)"},
		},
		{
			input:        "yield(1, 2 + 3);",
			expectedArgs: []string{"1", "(2 + 3)"},
		},
	}

	for _, tt := range tests {
		program, err := parseSource(tt.input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program has not enough statements. got=%d",
				len(program.Statements),
			)
		}
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		yield, ok := stmt.Expression.(*ast.YieldExpression)
		if !ok {
			t.Fatalf("expression not *ast.YieldExpression. got=%T", stmt.Expression)
		}

		if len(yield.Arguments) != len(tt.expectedArgs) {
			t.Logf("Expected %d arguments, got %d", len(tt.expectedArgs), len(yield.Arguments))
			t.Fail()
		}

		actualArgs := make([]string, len(yield.Arguments))
		for i, arg := range yield.Arguments {
			actualArgs[i] = arg.String()
		}

		if !reflect.DeepEqual(tt.expectedArgs, actualArgs) {
			t.Logf("Expected arguments to equal\n%v\n\tgot\n%v\n", tt.expectedArgs, actualArgs)
			t.Fail()
		}
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	program, err := parseSource(input)
	checkParserErrors(t, err)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program has not enough statements. got=%d",
			len(program.Statements),
		)
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("expression.Value not %d. got=%d", 5, literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf(
			"expression.TokenLiteral not %s. got=%s", "5",
			literal.TokenLiteral(),
		)
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		program, err := parseSource(tt.input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf(
				"exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator,
			)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	t.Run("literal expressions", func(t *testing.T) {
		infixTests := []struct {
			input      string
			leftValue  interface{}
			operator   string
			rightValue interface{}
		}{
			{"5 + 5;", 5, "+", 5},
			{"5 - 5;", 5, "-", 5},
			{"5 * 5;", 5, "*", 5},
			{"5 / 5;", 5, "/", 5},
			{"5 % 5;", 5, "%", 5},
			{"5 > 5;", 5, ">", 5},
			{"5 < 5;", 5, "<", 5},
			{"5 >= 5;", 5, ">=", 5},
			{"5 <= 5;", 5, "<=", 5},
			{"5 == 5;", 5, "==", 5},
			{"5 != 5;", 5, "!=", 5},
			{"5 <=> 5;", 5, "<=>", 5},
			{"foobar + barfoo;", "foobar", "+", "barfoo"},
			{"foobar - barfoo;", "foobar", "-", "barfoo"},
			{"foobar * barfoo;", "foobar", "*", "barfoo"},
			{"foobar / barfoo;", "foobar", "/", "barfoo"},
			{"foobar > barfoo;", "foobar", ">", "barfoo"},
			{"foobar < barfoo;", "foobar", "<", "barfoo"},
			{"foobar == barfoo;", "foobar", "==", "barfoo"},
			{"foobar <=> barfoo;", "foobar", "<=>", "barfoo"},
			{"foobar != barfoo;", "foobar", "!=", "barfoo"},
			{"true == true", true, "==", true},
			{"true != false", true, "!=", false},
			{"false == false", false, "==", false},
			{"false || false", false, "||", false},
			{"false && false", false, "&&", false},
		}

		for _, tt := range infixTests {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			if len(program.Statements) != 1 {
				t.Fatalf(
					"program.Statements does not contain %d statements. got=%d\n",
					1,
					len(program.Statements),
				)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf(
					"program.Statements[0] is not ast.ExpressionStatement. got=%T",
					program.Statements[0],
				)
			}

			if !testInfixExpression(t, stmt.Expression, tt.leftValue,
				tt.operator, tt.rightValue) {
				return
			}
		}
	})
	t.Run("symbols expressions", func(t *testing.T) {
		input := ":bar <=> 13"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		infix, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf(
				"stmt.Expression is not %T. got=%T",
				infix,
				stmt.Expression,
			)
		}
	})
	t.Run("call expression no args", func(t *testing.T) {
		input := "foo.bar <=> 13"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		infix, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf(
				"stmt.Expression is not %T. got=%T",
				infix,
				stmt.Expression,
			)
		}
	})
	t.Run("call expression with one arg", func(t *testing.T) {
		input := "foo.bar 3 <=> 13"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		cce, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf(
				"stmt.Expression is not %T. got=%T",
				cce,
				stmt.Expression,
			)
		}
	})
	t.Run("call expression with two args", func(t *testing.T) {
		input := "foo.bar 3, 5 <=> 13"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		cce, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf(
				"stmt.Expression is not %T. got=%T",
				cce,
				stmt.Expression,
			)
		}
	})
	t.Run("complex infix with call expression with just a block", func(t *testing.T) {
		input := "1 + 21 * 8 - 3 <=> foo { |x| x }"

		expr, err := parseExpression(input)
		checkParserErrors(t, err)

		infix, ok := expr.(*ast.InfixExpression)
		if !ok {
			t.Fatalf(
				"stmt.Expression is not %T. got=%T",
				infix,
				expr,
			)
		}
	})
	t.Run("easy infix with call expression with just a block", func(t *testing.T) {
		input := "1 <=> foo { |x| x }"

		expr, err := parseExpression(input)
		checkParserErrors(t, err)

		infix, ok := expr.(*ast.InfixExpression)
		if !ok {
			t.Fatalf(
				"stmt.Expression is not %T. got=%T",
				infix,
				expr,
			)
		}
	})
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)\n((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true | true",
			"(true | true)",
		},
		{
			"true & true",
			"(true & true)",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"x = 12 * 3;",
			"x = (12 * 3)",
		},
		{
			"x = 3 + 4 * 3;",
			"x = (3 + (4 * 3))",
		},
		{
			"x = add(4) * 3;",
			"x = (add(4) * 3)",
		},
		{
			"add(x = add(4) * 3);",
			"add(x = (add(4) * 3))",
		},
		{
			"a = b = 0;",
			"a = (b = 0)",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			actual := program.String()
			if actual != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, actual)
			}
		})
	}
}

func TestBlockExpression(t *testing.T) {
	tests := []struct {
		input             string
		expectedArguments []*ast.Identifier
		expectedBody      string
	}{
		{
			"method { x }",
			nil,
			"x",
		},
		{
			"method { |x| x }",
			[]*ast.Identifier{&ast.Identifier{Value: "x"}},
			"x",
		},
		{
			"method do; x; end",
			nil,
			"x",
		},
		{
			`
			method do
				x
			end`,
			nil,
			"x",
		},
		{
			"method do |x| x; end",
			[]*ast.Identifier{&ast.Identifier{Value: "x"}},
			"x",
		},
		{
			`method do |x|
				x
			end`,
			[]*ast.Identifier{&ast.Identifier{Value: "x"}},
			"x",
		},
	}

	for _, tt := range tests {
		program, err := parseSource(tt.input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program has not enough statements. got=%d",
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		call, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("exp not *ast.ContextCallExpression. got=%T", stmt.Expression)
		}

		block := call.Block
		if block == nil {
			t.Logf("Expected block not to be nil")
			t.FailNow()
		}

		if len(block.Parameters) != len(tt.expectedArguments) {
			t.Logf("Expected %d parameters, got %d", len(tt.expectedArguments), len(block.Parameters))
			t.Fail()
		}

		for i, arg := range block.Parameters {
			expected := tt.expectedArguments[i]
			expectedArg := expected.String()
			actualArg := arg.String()

			if expectedArg != actualArg {
				t.Logf(
					"Expected block argument %d to equal\n%s\n\tgot\n%s\n",
					i,
					expectedArg,
					actualArg,
				)
				t.Fail()
			}
		}

		body := block.Body.String()
		expectedBody := tt.expectedBody
		if expectedBody != body {
			t.Logf("Expected body to equal\n%s\n\tgot\n%s\n", expectedBody, body)
			t.Fail()
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		program, err := parseSource(tt.input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program has not enough statements. got=%d",
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf(
				"boolean.Value not %t. got=%t",
				tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func TestNilExpression(t *testing.T) {
	input := "nil;"

	program, err := parseSource(input)
	checkParserErrors(t, err)

	if len(program.Statements) != 1 {
		t.Fatalf(
			"program has not enough statements. got=%d",
			len(program.Statements),
		)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf(
			"program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0],
		)
	}

	if _, ok := stmt.Expression.(*ast.Nil); !ok {
		t.Fatalf("exp not *ast.Nil. got=%T", stmt.Expression)
	}
}

func TestConditionalExpression(t *testing.T) {
	t.Run("with operator expression", func(t *testing.T) {
		tests := []struct {
			input                         string
			expectedConditionLeft         string
			expectedConditionOperator     string
			expectedConditionRight        string
			expectedConsequenceExpression string
		}{
			{`if x < y
			x
			end`, "x", "<", "y", "x"},
			{`if x < y then
			x
			end`, "x", "<", "y", "x"},
			{`if x < y; x
			end`, "x", "<", "y", "x"},
			{`if x < y
			if x == 3
			y
			end
			x
			end`, "x", "<", "y", "if(x == 3) y endx"},
			{`if x < y
			x = Object x
			end`, "x", "<", "y", "x = (Object(x))"},
			{"x 3 if x < y", "x", "<", "y", "x(3)"},
			{"x.add 3 if x < y", "x", "<", "y", "x.add(3)"},
			{"yield 3 if x < y", "x", "<", "y", "yield 3"},
			{"yield self if x < y", "x", "<", "y", "yield self"},
			{`unless x < y
			x
			end`, "x", "<", "y", "x"},
			{`unless x < y then
			x
			end`, "x", "<", "y", "x"},
			{`unless x < y; x
			end`, "x", "<", "y", "x"},
			{`unless x < y
			if x == 3
			y
			end
			x
			end`, "x", "<", "y", "if(x == 3) y endx"},
			{`unless x < y
			x = Object x
			end`, "x", "<", "y", "x = (Object(x))"},
			{"x = 3 if x < y", "x", "<", "y", "x = 3"},
			{"@x = 3 if x < y", "x", "<", "y", "@x = 3"},
			{"x = 3 unless x < y", "x", "<", "y", "x = 3"},
			{"@x = 3 unless x < y", "x", "<", "y", "@x = 3"},
			{"x 3 unless x < y", "x", "<", "y", "x(3)"},
			{"x.add 3 unless x < y", "x", "<", "y", "x.add(3)"},
			{"yield 3 unless x < y", "x", "<", "y", "yield 3"},
			{"yield self unless x < y", "x", "<", "y", "yield self"},
		}

		for _, tt := range tests {
			t.Run("expression "+tt.input, func(t *testing.T) {
				program, err := parseSource(tt.input)
				checkParserErrors(t, err)

				if len(program.Statements) != 1 {
					t.Fatalf(
						"program.Body does not contain %d statements. got=%d\n",
						1,
						len(program.Statements),
					)
				}

				stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf(
						"program.Statements[0] is not ast.ExpressionStatement. got=%T",
						program.Statements[0],
					)
				}

				exp, ok := stmt.Expression.(*ast.ConditionalExpression)
				if !ok {
					t.Fatalf(
						"stmt.Expression is not %T. got=%T",
						exp,
						stmt.Expression,
					)
				}

				if !testInfixExpression(
					t,
					exp.Condition,
					tt.expectedConditionLeft,
					tt.expectedConditionOperator,
					tt.expectedConditionRight,
				) {
					return
				}

				consequenceBody := ""
				for _, stmt := range exp.Consequence.Statements {
					consequence, ok := stmt.(*ast.ExpressionStatement)
					if !ok {
						t.Fatalf(
							"Statements[0] is not ast.ExpressionStatement. got=%T",
							exp.Consequence.Statements[0],
						)
					}

					consequenceBody += consequence.Expression.String()
				}

				if consequenceBody != tt.expectedConsequenceExpression {
					t.Logf(
						"Expected consequence to equal %q, got %q\n",
						tt.expectedConsequenceExpression,
						consequenceBody,
					)
					t.Fail()
				}

				if exp.Alternative != nil {
					t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
				}
			})
		}
	})
	t.Run("with method call expression", func(t *testing.T) {
		tests := []struct {
			input       string
			condContext string
			condMethod  string
			condArg     string
			consequence string
		}{
			{`unless x.exist? :y
			x
			end`, "x", "exist?", "y", "x"},
			{`unless x.exist? :y
			x = Object x
			end`, "x", "exist?", "y", "x = Object x"},
			{`unless x.exist? :y
			x
			end`, "x", "exist?", "y", "x"},
			{`unless x.exist? :y
			x = Object x
			end`, "x", "exist?", "y", "x = Object x"},
		}

		for _, tt := range tests {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			if len(program.Statements) != 1 {
				t.Fatalf(
					"program.Body does not contain %d statements. got=%d\n",
					1,
					len(program.Statements),
				)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf(
					"program.Statements[0] is not ast.ExpressionStatement. got=%T",
					program.Statements[0],
				)
			}

			exp, ok := stmt.Expression.(*ast.ConditionalExpression)
			if !ok {
				t.Fatalf(
					"stmt.Expression is not %T. got=%T",
					exp,
					stmt.Expression,
				)
			}

			call, ok := exp.Condition.(*ast.ContextCallExpression)
			if !ok {
				t.Fatalf(
					"exp.Condition is not %T. got=%T",
					call,
					exp.Condition,
				)
			}

			if call.Function.String() != tt.condMethod {
				t.Logf(
					"Expected condition call method to equal %q, got %q\n",
					tt.condMethod,
					call.Function.String(),
				)
			}

			args := []string{}
			for _, a := range call.Arguments {
				args = append(args, a.String())
			}
			if strings.Join(args, " ") != tt.condArg {
				t.Logf(
					"Expected condition call args to equal %q, got %q\n",
					tt.condArg,
					strings.Join(args, " "),
				)
			}

			if call.Context.String() != tt.condContext {
				t.Logf(
					"Expected condition call context to equal %q, got %q\n",
					tt.condContext,
					call.Context.String(),
				)
			}

			consequenceBody := ""
			for _, stmt := range exp.Consequence.Statements {
				consequence, ok := stmt.(*ast.ExpressionStatement)
				if !ok {
					t.Fatalf(
						"Statements[0] is not ast.ExpressionStatement. got=%T",
						exp.Consequence.Statements[0],
					)
				}

				consequenceBody += consequence.Expression.String()
			}

			if consequenceBody != tt.consequence {
				t.Logf(
					"Expected consequence to equal %q, got %q\n",
					tt.consequence,
					consequenceBody,
				)
			}

			if exp.Alternative != nil {
				t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
			}
		}
	})
}

func TestConditionalExpressionWithAlternative(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		condition   [3]string
		consequence string
		alternative string
	}{
		{
			"regular if else",
			`
			if x < y
				x
			else
				y
			end`,
			[3]string{"x", "<", "y"},
			"x",
			"y",
		},
		{
			"tenary if",
			"x < y ? x : y;",
			[3]string{"x", "<", "y"},
			"x",
			"y",
		},
		{
			"tenary if with symbol as consequence",
			"x < y ? :x : y;",
			[3]string{"x", "<", "y"},
			":x",
			"y",
		},
		{
			"tenary if with symbol as alternative",
			"x < y ? x : :y;",
			[3]string{"x", "<", "y"},
			"x",
			":y",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Body does not contain %d statements. got=%d\n",
					1, len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
					program.Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.ConditionalExpression)
			if !ok {
				t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
			}

			if !testInfixExpression(t, exp.Condition, tt.condition[0], tt.condition[1], tt.condition[2]) {
				return
			}

			if len(exp.Consequence.Statements) != 1 {
				t.Errorf("consequence is not 1 statements. got=%d\n",
					len(exp.Consequence.Statements))
			}

			consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
					exp.Consequence.Statements[0])
			}

			if !testLiteralExpression(t, consequence.Expression, tt.consequence) {
				return
			}

			if len(exp.Alternative.Statements) != 1 {
				t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
					len(exp.Alternative.Statements))
			}

			alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
					exp.Alternative.Statements[0])
			}

			if !testLiteralExpression(t, alternative.Expression, tt.alternative) {
				return
			}
		})
	}
	t.Run("tenary if with call as consequence", func(t *testing.T) {
		tt := struct {
			input       string
			condition   [3]string
			consequence string
			alternative string
		}{
			"x < y ? x.foo : y;",
			[3]string{"x", "<", "y"},
			"x.foo()",
			"y",
		}
		program, err := parseSource(tt.input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Body does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ConditionalExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
		}

		if !testInfixExpression(t, exp.Condition, tt.condition[0], tt.condition[1], tt.condition[2]) {
			return
		}

		if len(exp.Consequence.Statements) != 1 {
			t.Errorf("consequence is not 1 statements. got=%d\n",
				len(exp.Consequence.Statements))
		}

		consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
				exp.Consequence.Statements[0])
		}

		if consequence.String() != tt.consequence {
			t.Logf("Expected consequence to equal %s, got %s", tt.consequence, consequence.String())
			t.Fail()
		}

		if len(exp.Alternative.Statements) != 1 {
			t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
				len(exp.Alternative.Statements))
		}

		alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
				exp.Alternative.Statements[0])
		}

		if !testLiteralExpression(t, alternative.Expression, tt.alternative) {
			return
		}
	})
}

func TestFunctionLiteralParsing(t *testing.T) {
	type funcParam struct {
		name         string
		defaultValue interface{}
	}
	tests := []struct {
		desc          string
		input         string
		receiver      string
		name          string
		parameters    []funcParam
		bodyStatement string
	}{
		{
			"with parens",
			`def foo(x, y)
			  x + y
          end`,
			"",
			"foo",
			[]funcParam{
				{name: "x", defaultValue: nil},
				{name: "y", defaultValue: nil},
			},
			"(x + y)",
		},
		{
			"without parens",
			`def bar x, y
          x + y
          end`,
			"",
			"bar",
			[]funcParam{
				{name: "x", defaultValue: nil},
				{name: "y", defaultValue: nil},
			},
			"(x + y)",
		},
		{
			"without arguments",
			`def qux
          x + y
          end`,
			"",
			"qux",
			[]funcParam{},
			"(x + y)",
		},
		{
			"expression separator semicolon no arguments",
			"def qux; x + y; end",
			"",
			"qux",
			[]funcParam{},
			"(x + y)",
		},
		{
			"expression separator semicolon two arguments",
			"def foo x, y; x + y; end",
			"",
			"foo",
			[]funcParam{
				{name: "x", defaultValue: nil},
				{name: "y", defaultValue: nil},
			},
			"(x + y)",
		},
		{
			"expression separator semicolon with parens and two arguments",
			"def foo(x, y); x + y; end",
			"",
			"foo",
			[]funcParam{
				{name: "x", defaultValue: nil},
				{name: "y", defaultValue: nil},
			},
			"(x + y)",
		},
		{
			"upcase function name",
			`def Qux
          x + y
          end
          `,
			"",
			"Qux",
			[]funcParam{},
			"(x + y)",
		},
		{
			"two arguments with defaults without parens",
			`def foo x = 2, y = 3
          x + y
          end
          `,
			"",
			"foo",
			[]funcParam{
				{name: "x", defaultValue: 2},
				{name: "y", defaultValue: 3},
			},
			"(x + y)",
		},
		{
			"operator as function name",
			`def <=>
          x + y
          end
          `,
			"",
			"<=>",
			[]funcParam{},
			"(x + y)",
		},
		{
			"function on local variable context",
			`def a.qux
          x + y
          end`,
			"a",
			"qux",
			[]funcParam{},
			"(x + y)",
		},
		{
			"function on const context",
			`def A.qux
          x + y
          end`,
			"A",
			"qux",
			[]funcParam{},
			"(x + y)",
		},
		{
			"upcase function on const context",
			`def A.Qux
          x + y
          end`,
			"A",
			"Qux",
			[]funcParam{},
			"(x + y)",
		},
		{
			"function on self context",
			`def self.qux
          x + y
          end`,
			"self",
			"qux",
			[]funcParam{},
			"(x + y)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			if len(program.Statements) != 1 {
				t.Fatalf(
					"program.Body does not contain %d statements. got=%d\n",
					1,
					len(program.Statements),
				)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf(
					"program.Statements[0] is not ast.ExpressionStatement. got=%T",
					program.Statements[0],
				)
			}

			function, ok := stmt.Expression.(*ast.FunctionLiteral)
			if !ok {
				t.Fatalf(
					"stmt.Expression is not ast.FunctionLiteral. got=%T",
					stmt.Expression,
				)
			}

			receiver := ""
			if function.Receiver != nil {
				receiver = function.Receiver.Value
			}
			if receiver != tt.receiver {
				t.Logf("function receiver wrong, want %q, got %q", tt.receiver, receiver)
				t.Fail()
			}

			functionName := function.Name.Value
			if functionName != tt.name {
				t.Logf("function name wrong, want %q, got %q", tt.name, functionName)
				t.Fail()
			}

			if len(function.Parameters) != len(tt.parameters) {
				t.Fatalf(
					"function literal parameters wrong. want %d, got=%d\n",
					len(tt.parameters),
					len(function.Parameters),
				)
			}

			for i, param := range function.Parameters {
				testLiteralExpression(t, param.Name, tt.parameters[i].name)
				testLiteralExpression(t, param.Default, tt.parameters[i].defaultValue)
			}

			if len(function.Body.Statements) != 1 {
				t.Fatalf(
					"function.Body.Statements has not 1 statements. got=%d\n",
					len(function.Body.Statements),
				)
			}

			bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf(
					"function body stmt is not ast.ExpressionStatement. got=%T",
					function.Body.Statements[0],
				)
			}

			statement := bodyStmt.String()
			if statement != tt.bodyStatement {
				t.Logf(
					"Expected body statement to equal\n%q\n\tgot\n%q\n",
					tt.bodyStatement,
					statement,
				)
				t.Fail()
			}
		})
	}
	t.Run("test function rescue block", func(t *testing.T) {
		input := `
			def foo
				3
			rescue Exception => e
				puts e
			end
		`

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Body does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf(
				"stmt.Expression is not %T. got=%T",
				function,
				stmt.Expression,
			)
		}
	})
}

func TestBlockExpressionParsing(t *testing.T) {
	tests := []struct {
		input         string
		parameters    []string
		bodyStatement string
	}{
		{
			`method do |x, y|
          x + y
          end`,
			[]string{"x", "y"},
			"(x + y)",
		},
		{
			`method do
          x + y
          end`,
			[]string{},
			"(x + y)",
		},
		{
			"method do ; x + y; end",
			[]string{},
			"(x + y)",
		},
		{
			"method do |x, y|; x + y; end",
			[]string{"x", "y"},
			"(x + y)",
		},
		{
			"method do |x, y|; x + y; end",
			[]string{"x", "y"},
			"(x + y)",
		},
		{
			`method { |x, y|
			  x + y
			  }`,
			[]string{"x", "y"},
			"(x + y)",
		},
		{
			`method {
          x + y
          }`,
			[]string{},
			"(x + y)",
		},
		{
			"method { x + y; }",
			[]string{},
			"(x + y)",
		},
		{
			"method { |x, y|; x + y; }",
			[]string{"x", "y"},
			"(x + y)",
		},
		{
			"method { |x, y|; x + y; }",
			[]string{"x", "y"},
			"(x + y)",
		},
		{
			"method { |x, y|; x.add y }",
			[]string{"x", "y"},
			"x.add(y)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			if len(program.Statements) != 1 {
				t.Fatalf(
					"program.Body does not contain %d statements. got=%d\n",
					1,
					len(program.Statements),
				)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Logf(
					"program.Statements[0] is not ast.ExpressionStatement. got=%T",
					program.Statements[0],
				)
				t.Log(program.Statements)
				t.FailNow()
			}

			call, ok := stmt.Expression.(*ast.ContextCallExpression)
			if !ok {
				t.Logf(
					"stmt.Expression is not *ast.ContextCallExpression. got=%T",
					stmt.Expression,
				)
				t.Fail()
			}

			block := call.Block

			if block == nil {
				t.Logf("Expected block not to be nil")
				t.FailNow()
			}

			if len(block.Parameters) != len(tt.parameters) {
				t.Fatalf(
					"block literal parameters wrong. want %d, got=%d\n",
					len(tt.parameters),
					len(block.Parameters),
				)
			}

			for i, param := range block.Parameters {
				testLiteralExpression(t, param.Name, tt.parameters[i])
			}

			if len(block.Body.Statements) != 1 {
				t.Fatalf(
					"block.Body.Statements has not 1 statements. got=%d\n",
					len(block.Body.Statements),
				)
			}

			bodyStmt, ok := block.Body.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf(
					"block body stmt is not ast.ExpressionStatement. got=%T",
					block.Body.Statements[0],
				)
			}

			statement := bodyStmt.String()
			if statement != tt.bodyStatement {
				t.Logf(
					"Expected body statement to equal\n%q\n\tgot\n%q\n",
					tt.bodyStatement,
					statement,
				)
				t.Fail()
			}
		})
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	type funcParam struct {
		name         string
		defaultValue interface{}
	}
	tests := []struct {
		desc           string
		input          string
		expectedParams []funcParam
	}{
		{
			desc:           "no params with parens",
			input:          "def fn(); end",
			expectedParams: []funcParam{},
		},
		{
			desc:           "one param with parens",
			input:          "def fn(x); end",
			expectedParams: []funcParam{{name: "x"}},
		},
		{
			desc:           "multiple params with parens",
			input:          "def fn(x, y, z); end",
			expectedParams: []funcParam{{name: "x"}, {name: "y"}, {name: "z"}},
		},
		{
			desc:           "multiple params first two defaults with parens",
			input:          "def fn(x = 3, y = 18, z); end",
			expectedParams: []funcParam{{name: "x", defaultValue: 3}, {name: "y", defaultValue: 18}, {name: "z"}},
		},
		{
			desc:           "multiple params middle default with parens",
			input:          "def fn(x, y = 18, z); end",
			expectedParams: []funcParam{{name: "x"}, {name: "y", defaultValue: 18}, {name: "z"}},
		},
		{
			desc:           "multiple params last default with parens",
			input:          "def fn(x, y, z = 1); end",
			expectedParams: []funcParam{{name: "x"}, {name: "y"}, {name: "z", defaultValue: 1}},
		},
		{
			desc:           "multiple params last array splat with parens",
			input:          "def fn(x, y, *z); end",
			expectedParams: []funcParam{{name: "x"}, {name: "y"}, {name: "z"}},
		},
		{
			desc:           "one param array splat with parens",
			input:          "def fn(*x); end",
			expectedParams: []funcParam{{name: "x"}},
		},
		{
			desc:           "multiple params last block capture with parens",
			input:          "def fn(x, y, &z); end",
			expectedParams: []funcParam{{name: "x"}, {name: "y"}},
		},
		{
			desc:           "one param block capture with parens",
			input:          "def fn(&x); end",
			expectedParams: []funcParam{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			function := stmt.Expression.(*ast.FunctionLiteral)

			if len(function.Parameters) != len(tt.expectedParams) {
				t.Errorf(
					"length parameters wrong. want %d, got=%d\n",
					len(tt.expectedParams),
					len(function.Parameters),
				)
			}

			for i, ident := range tt.expectedParams {
				testLiteralExpression(t, function.Parameters[i].Name, ident.name)
				testLiteralExpression(t, function.Parameters[i].Default, ident.defaultValue)
			}
		})
	}
}

func TestBlockParameterParsing(t *testing.T) {
	type funcParam struct {
		name         string
		defaultValue interface{}
	}
	tests := []struct {
		desc           string
		input          string
		expectedParams []funcParam
	}{
		{
			desc:           "empty brace block",
			input:          "method {}",
			expectedParams: []funcParam{},
		},
		{
			desc:           "empty brace block params",
			input:          "method { || }",
			expectedParams: []funcParam{},
		},
		{
			desc:           "one brace block param",
			input:          "method { |x| }",
			expectedParams: []funcParam{{name: "x"}},
		},
		{
			desc:           "multiple brace block params",
			input:          "method { |x, y, z| }",
			expectedParams: []funcParam{{name: "x"}, {name: "y"}, {name: "z"}},
		},
		{
			desc:           "empty do block",
			input:          "method do; end",
			expectedParams: []funcParam{},
		},
		{
			desc:           "empty do block params",
			input:          "method do ||; end",
			expectedParams: []funcParam{},
		},
		{
			desc:           "one do block param",
			input:          "method do |x|; end",
			expectedParams: []funcParam{{name: "x"}},
		},
		{
			desc:           "multiple do block params",
			input:          "method do |x, y, z|; end",
			expectedParams: []funcParam{{name: "x"}, {name: "y"}, {name: "z"}},
		},
		{
			desc:           "multiple brace block params with defaults",
			input:          "method { |x = 3, y = 2, z| }",
			expectedParams: []funcParam{{name: "x", defaultValue: 3}, {name: "y", defaultValue: 2}, {name: "z"}},
		},
		{
			desc:           "multiple do block params starting defaults",
			input:          "method do |x = 1, y = 8, z|; end",
			expectedParams: []funcParam{{name: "x", defaultValue: 1}, {name: "y", defaultValue: 8}, {name: "z"}},
		},
		{
			desc:           "multiple brace block params with middle default",
			input:          "method { |x, y = 2, z| }",
			expectedParams: []funcParam{{name: "x"}, {name: "y", defaultValue: 2}, {name: "z"}},
		},
		{
			desc:           "multiple do block params with middle default",
			input:          "method do |x, y = 8, z|; end",
			expectedParams: []funcParam{{name: "x"}, {name: "y", defaultValue: 8}, {name: "z"}},
		},
		{
			desc:           "multiple brace block params last defaults",
			input:          "method { |x, y, z = 2| }",
			expectedParams: []funcParam{{name: "x"}, {name: "y"}, {name: "z", defaultValue: 2}},
		},
		{
			desc:           "multiple do block params last defaults",
			input:          "method do |x, y, z = 4|; end",
			expectedParams: []funcParam{{name: "x"}, {name: "y"}, {name: "z", defaultValue: 4}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			call, ok := stmt.Expression.(*ast.ContextCallExpression)
			if !ok {
				t.Logf(
					"stmt.Expression is not *ast.ContextCallExpression. got=%T",
					stmt.Expression,
				)
				t.Fail()
			}

			block := call.Block

			if block == nil {
				t.Logf("Expected block not to be nil")
				t.FailNow()
			}

			if len(block.Parameters) != len(tt.expectedParams) {
				t.Errorf(
					"length parameters wrong. want %d, got=%d\n",
					len(tt.expectedParams),
					len(block.Parameters),
				)
			}

			for i, ident := range tt.expectedParams {
				testLiteralExpression(t, block.Parameters[i].Name, ident.name)
				testLiteralExpression(t, block.Parameters[i].Default, ident.defaultValue)
			}
		})
	}
}

func TestCallExpressionParsing(t *testing.T) {
	testCases := []struct {
		desc        string
		input       string
		context     string
		funcName    string
		arguments   []interface{}
		hasBlock    bool
		blockParams []string
	}{
		{
			desc:     "with parens",
			input:    "add(1, 2 * 3, 4 + 5);",
			funcName: "add",
			arguments: []interface{}{
				1, infix{2, "*", 3}, infix{4, "+", 5},
			},
		},
		{
			desc:     "without parens",
			input:    "add 1, 2 * 3, 4 + 5;",
			funcName: "add",
			arguments: []interface{}{
				1, infix{2, "*", 3}, infix{4, "+", 5},
			},
		},
		{
			desc:     "with parens and brace block",
			input:    "add(1, 2 * 3, 4 + 5) { |x| x };",
			funcName: "add",
			arguments: []interface{}{
				1, infix{2, "*", 3}, infix{4, "+", 5},
			},
			hasBlock:    true,
			blockParams: []string{"x"},
		},
		{
			desc:     "with parens and do block",
			input:    "add(1, 2 * 3, 4 + 5) do |x| x; end;",
			funcName: "add",
			arguments: []interface{}{
				1, infix{2, "*", 3}, infix{4, "+", 5},
			},
			hasBlock:    true,
			blockParams: []string{"x"},
		},
		{
			desc:     "without parens with block",
			input:    "add 1, 2 * 3, 4 + 5 { |x| x };",
			funcName: "add",
			arguments: []interface{}{
				1, infix{2, "*", 3}, infix{4, "+", 5},
			},
			hasBlock:    true,
			blockParams: []string{"x"},
		},
		{
			desc:        "without parens without args with block",
			input:       "add { |x| x };",
			funcName:    "add",
			hasBlock:    true,
			blockParams: []string{"x"},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			expr, err := parseExpression(tt.input)
			checkParserErrors(t, err)

			call, ok := expr.(*ast.ContextCallExpression)
			if !ok {
				t.Fatalf(
					"expression is not %T, got=%T",
					call,
					expr,
				)
			}

			if !testIdentifier(t, call.Function, tt.funcName) {
				return
			}

			if tt.context != "" && !testIdentifier(t, call.Context, tt.context) {
				return
			}

			if len(call.Arguments) != len(tt.arguments) {
				t.Logf(
					"wrong length of arguments. want %d, got=%d",
					len(tt.arguments),
					len(call.Arguments),
				)
				t.Fail()
				for len(call.Arguments) > len(tt.arguments) {
					tt.arguments = append(tt.arguments, "<unexpected>")
				}
			}

			for i, arg := range call.Arguments {
				t.Logf("argument %d", i+1)
				testExpression(t, arg, tt.arguments[i])
			}

			if tt.hasBlock {
				if call.Block == nil {
					t.Logf("Expected function block not to be nil")
					t.FailNow()
				}

				if len(call.Block.Parameters) != len(tt.blockParams) {
					t.Logf(
						"wrong length of block parameters. want %d, got=%d",
						len(tt.blockParams),
						len(call.Block.Parameters),
					)
					t.Fail()
					for len(call.Block.Parameters) > len(tt.blockParams) {
						tt.blockParams = append(tt.blockParams, "<unexpected>")
					}
				}

				for i, param := range call.Block.Parameters {
					expected := tt.blockParams[i]
					actual := param.String()

					if expected != actual {
						t.Logf(
							"Expected block param %d to equal\n%s\n\tgot\n%s\n",
							i,
							expected,
							actual,
						)
						t.Fail()
					}
				}
			}
		})
	}
	t.Run("with parens and do block", func(t *testing.T) {
	})
	t.Run("without parens with block", func(t *testing.T) {
	})
	t.Run("without parens", func(t *testing.T) {
	})
	t.Run("without parens without args with block", func(t *testing.T) {
	})
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
		{
			input:         "add 1;",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         `add "foo";`,
			expectedIdent: "add",
			expectedArgs:  []string{"foo"},
		},
		{
			input:         `add :foo;`,
			expectedIdent: "add",
			expectedArgs:  []string{":foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			program, err := parseSource(tt.input)
			checkParserErrors(t, err)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			exp, ok := stmt.Expression.(*ast.ContextCallExpression)
			if !ok {
				t.Fatalf(
					"stmt.Expression is not ast.ContextCallExpression. got=%T",
					stmt.Expression,
				)
			}

			if !testIdentifier(t, exp.Function, tt.expectedIdent) {
				return
			}

			if len(exp.Arguments) != len(tt.expectedArgs) {
				t.Fatalf("wrong number of arguments. want=%d, got=%d",
					len(tt.expectedArgs), len(exp.Arguments))
			}

			for i, arg := range tt.expectedArgs {
				if exp.Arguments[i].String() != arg {
					t.Errorf("argument %d wrong. want=%q, got=%q", i,
						arg, exp.Arguments[i].String())
				}
			}
		})
	}
}

func TestContextCallExpression(t *testing.T) {
	t.Run("context call with multiple args with parens", func(t *testing.T) {
		input := "foo.add(1, 2 * 3, 4 + 5);"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Context, "foo") {
			return
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 3 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}

		testLiteralExpression(t, exp.Arguments[0], 1)
		testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
		testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
	})
	t.Run("context call with multiple args with parens and block", func(t *testing.T) {
		input := "foo.add(1, 2 * 3, 4 + 5) { |x|x.to_s };"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Logf("Input: %s\n", input)
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Context, "foo") {
			return
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 3 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}

		testLiteralExpression(t, exp.Arguments[0], 1)
		testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
		testInfixExpression(t, exp.Arguments[2], 4, "+", 5)

		if exp.Block == nil {
			t.Logf("Expected block not to be nil")
			t.Fail()
		}
	})
	t.Run("context call with multiple args no parens", func(t *testing.T) {
		input := "foo.add 1, 2 * 3, 4 + 5;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Context, "foo") {
			return
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 3 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}

		testLiteralExpression(t, exp.Arguments[0], 1)
		testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
		testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
	})
	t.Run("context call with multiple args no parens with block", func(t *testing.T) {
		input := "foo.add 1, 2 * 3, 4 + 5 { |x| x };"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Context, "foo") {
			return
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 3 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}

		testLiteralExpression(t, exp.Arguments[0], 1)
		testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
		testInfixExpression(t, exp.Arguments[2], 4, "+", 5)

		if exp.Block == nil {
			t.Logf("Expected block not to be nil")
			t.Fail()
		}
	})
	t.Run("context call with no args", func(t *testing.T) {
		input := "foo.add;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Context, "foo") {
			return
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}
	})
	t.Run("context call on self with no args", func(t *testing.T) {
		input := "self.add;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if _, ok := exp.Context.(*ast.Self); !ok {
			t.Logf("exp.Context is not ast.Self, got=%T", exp.Context)
			t.Fail()
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}
	})
	t.Run("context call on self with no dot", func(t *testing.T) {
		input := "self add;"

		_, err := parseSource(input)

		if err == nil {
			t.Logf("Expected parser error, got nil")
			t.FailNow()
		}

		errs := err.errors
		cause := errors.Cause(errs[0])

		unexpectErr, ok := cause.(*unexpectedTokenError)
		if !ok {
			t.Logf("Expected err to be %T, got %T\n", unexpectErr, cause)
			t.FailNow()
		}

		{
			expected := []token.Type{token.NEWLINE, token.SEMICOLON, token.DOT, token.EOF}
			actual := unexpectErr.expectedTokens
			if !reflect.DeepEqual(expected, actual) {
				t.Logf("Expected error to equal\n%+#v\n\tgot\n%+#v\n", expected, actual)
				t.Fail()
			}
		}

		{
			expected := token.IDENT
			actual := unexpectErr.actualToken
			if !reflect.DeepEqual(expected, actual) {
				t.Logf("Expected error to equal\n%+#v\n\tgot\n%+#v\n", expected, actual)
				t.Fail()
			}
		}
	})
	t.Run("context call on nonident with no dot", func(t *testing.T) {
		input := "1 add;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Context, 1) {
			return
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}
	})
	t.Run("context call on nonident with dot", func(t *testing.T) {
		input := "1.add"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Context, 1) {
			return
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}
	})
	t.Run("context call on nonident with no dot multiargs", func(t *testing.T) {
		input := "1 add 1"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf(
				"stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0],
			)
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Context, 1) {
			return
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 1 {
			t.Fatalf(
				"wrong length of arguments. got=%d",
				len(exp.Arguments),
			)
		}

		if !testIntegerLiteral(t, exp.Arguments[0], 1) {
			return
		}
	})
	t.Run("context call on ident with no dot", func(t *testing.T) {
		input := "foo add;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, "foo") {
			return
		}

		if len(exp.Arguments) != 1 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}

		if !testIdentifier(t, exp.Arguments[0], "add") {
			return
		}
	})
	t.Run("context call on const with no dot", func(t *testing.T) {
		input := "Integer add;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, "Integer") {
			return
		}

		if len(exp.Arguments) != 1 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}

		if !testIdentifier(t, exp.Arguments[0], "add") {
			return
		}
	})
	t.Run("context call on ident with no dot Const as arg", func(t *testing.T) {
		input := "add Integer;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 1 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}

		if !testIdentifier(t, exp.Arguments[0], "Integer") {
			return
		}
	})
	t.Run("chained context call with dot without parens", func(t *testing.T) {
		input := "foo.add.bar;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		context, ok := exp.Context.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf(
				"expr.Context is not ast.ContextCallExpression. got=%T",
				exp.Context,
			)
		}

		if !testIdentifier(t, context.Context, "foo") {
			return
		}

		if !testIdentifier(t, context.Function, "add") {
			return
		}

		if len(context.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(context.Arguments))
		}

		if !testIdentifier(t, exp.Function, "bar") {
			return
		}

		if len(exp.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}
	})
	t.Run("chained context call with dot without parens", func(t *testing.T) {
		input := "1.add.bar;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		context, ok := exp.Context.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf(
				"expr.Context is not ast.ContextCallExpression. got=%T",
				exp.Context,
			)
		}

		if !testIntegerLiteral(t, context.Context, 1) {
			return
		}

		if !testIdentifier(t, context.Function, "add") {
			return
		}

		if len(context.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(context.Arguments))
		}

		if !testIdentifier(t, exp.Function, "bar") {
			return
		}

		if len(exp.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}
	})
	t.Run("chained context call with dot with parens", func(t *testing.T) {
		input := "foo.add().bar();"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		context, ok := exp.Context.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf(
				"expr.Context is not ast.ContextCallExpression. got=%T",
				exp.Context,
			)
		}

		if !testIdentifier(t, context.Function, "add") {
			return
		}

		if len(context.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(context.Arguments))
		}

		if !testIdentifier(t, exp.Function, "bar") {
			return
		}

		if len(exp.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}
	})
	t.Run("scope as context call", func(t *testing.T) {
		input := "foo.add::bar;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		context, ok := exp.Context.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf(
				"expr.Context is not ast.ContextCallExpression. got=%T",
				exp.Context,
			)
		}

		if !testIdentifier(t, context.Function, "add") {
			return
		}

		if len(context.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(context.Arguments))
		}

		if !testIdentifier(t, exp.Function, "bar") {
			return
		}

		if len(exp.Arguments) != 0 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}
	})
	t.Run("allow `class` as method name", func(t *testing.T) {
		input := "foo.class;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, expr.Context, "foo") {
			return
		}

		if !testIdentifier(t, expr.Function, "class") {
			return
		}
	})
	t.Run("allow operators as method name", func(t *testing.T) {
		input := "foo.<=>;"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		if len(program.Statements) != 1 {
			t.Fatalf(
				"program.Statements does not contain %d statements. got=%d\n",
				1,
				len(program.Statements),
			)
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.ContextCallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.ContextCallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, expr.Context, "foo") {
			return
		}

		if !testIdentifier(t, expr.Function, "<=>") {
			return
		}
	})
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	program, err := parseSource(input)
	checkParserErrors(t, err)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestSymbolExpression(t *testing.T) {
	tests := []struct {
		input string
		value string
	}{
		{
			`:symbol;`,
			"symbol",
		},
		{
			`:"symbol";`,
			"symbol",
		},
		{
			`:'symbol';`,
			"symbol",
		},
	}

	for _, tt := range tests {
		program, err := parseSource(tt.input)
		checkParserErrors(t, err)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		literal, ok := stmt.Expression.(*ast.SymbolLiteral)
		if !ok {
			t.Fatalf("exp not *ast.SymbolLiteral. got=%T", stmt.Expression)
		}

		if literal.Value.String() != tt.value {
			t.Errorf("literal.Value not %q. got=%q", tt.value, literal.Value)
		}
	}
}

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3, {'foo'=>2}]"
	program, err := parseSource(input)
	checkParserErrors(t, err)

	if len(program.Statements) != 1 {
		t.Logf("Expected only one statement, got %d\n", len(program.Statements))
		t.Fail()
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not ast.ArrayLiteral. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 4 {
		t.Fatalf("len(array.Elements) not 4. got=%d", len(array.Elements))
	}
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
	testHashLiteral(t, array.Elements[3], map[string]string{"foo": "2"})
}

func TestParsingIndexExpressions(t *testing.T) {
	t.Run("one arg as index", func(t *testing.T) {
		input := "myArray[1 + 1]"
		program, err := parseSource(input)
		checkParserErrors(t, err)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		indexExp, ok := stmt.Expression.(*ast.IndexExpression)
		if !ok {
			t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
		}

		if !testIdentifier(t, indexExp.Left, "myArray") {
			return
		}

		if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
			return
		}
	})
	t.Run("two args as index", func(t *testing.T) {
		t.Run("integers", func(t *testing.T) {
			input := "myArray[1, 1]"
			program, err := parseSource(input)
			checkParserErrors(t, err)

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			indexExp, ok := stmt.Expression.(*ast.IndexExpression)
			if !ok {
				t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
			}

			if !testIdentifier(t, indexExp.Left, "myArray") {
				return
			}

			if !testIntegerLiteral(t, indexExp.Index, 1) {
				return
			}

			if !testIntegerLiteral(t, indexExp.Length, 1) {
				return
			}
		})
		t.Run("method calls as index", func(t *testing.T) {
			input := "myArray[foo.bar, 1]"
			program, err := parseSource(input)
			checkParserErrors(t, err)

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			indexExp, ok := stmt.Expression.(*ast.IndexExpression)
			if !ok {
				t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
			}

			if !testIdentifier(t, indexExp.Left, "myArray") {
				return
			}

			index := indexExp.Index.String()
			if index != "foo.bar()" {
				t.Logf("Expected index arg to equal %s, got %s", "foo.bar()", index)
				t.Fail()
			}

			if !testIntegerLiteral(t, indexExp.Length, 1) {
				return
			}
		})
		t.Run("method calls as length", func(t *testing.T) {
			input := "myArray[1, foo.bar]"
			program, err := parseSource(input)
			checkParserErrors(t, err)

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			indexExp, ok := stmt.Expression.(*ast.IndexExpression)
			if !ok {
				t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
			}

			if !testIdentifier(t, indexExp.Left, "myArray") {
				return
			}

			if !testIntegerLiteral(t, indexExp.Index, 1) {
				return
			}

			length := indexExp.Length.String()
			if length != "foo.bar()" {
				t.Logf("Expected length arg to equal %s, got %s", "foo.bar()", length)
				t.Fail()
			}
		})
	})
}

func TestParsingModuleExpressions(t *testing.T) {
	input := "module A\n3\nend\n"

	program, err := parseSource(input)
	checkParserErrors(t, err)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	_, ok = stmt.Expression.(*ast.ModuleExpression)
	if !ok {
		t.Fatalf("exp not *ast.ModuleExpression. got=%T", stmt.Expression)
	}
}

func TestParsingClassExpressions(t *testing.T) {
	t.Run("basic class", func(t *testing.T) {
		input := "class A\n3\nend\n"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		class, ok := stmt.Expression.(*ast.ClassExpression)
		if !ok {
			t.Fatalf("exp not *ast.ClassExpression. got=%T", stmt.Expression)
		}

		className := "A"
		if className != class.Name.String() {
			t.Logf("Expected class name to equal %q, got %q\n", className, class.Name.String())
			t.Fail()
		}
	})
	t.Run("class with superclass", func(t *testing.T) {
		input := "class A < B\n3\nend\n"

		program, err := parseSource(input)
		checkParserErrors(t, err)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		class, ok := stmt.Expression.(*ast.ClassExpression)
		if !ok {
			t.Fatalf("exp not *ast.ClassExpression. got=%T", stmt.Expression)
		}

		className := "A"
		if className != class.Name.Value {
			t.Logf("Expected class name to equal %q, got %q\n", className, class.Name.Value)
			t.Fail()
		}

		superclassName := "B"
		if superclassName != class.SuperClass.Value {
			t.Logf("Expected superclass name to equal %q, got %q\n", superclassName, class.SuperClass.Value)
			t.Fail()
		}
	})
	t.Run("downcase class", func(t *testing.T) {
		t.Skip("evaluate error")
		input := "class a\n3\nend\n"

		_, err := parseSource(input)
		checkParserErrors(t, err)
	})
}

func TestParseHash(t *testing.T) {
	tests := []struct {
		input   string
		hashMap map[string]string
	}{
		{
			input:   `{"foo" => 42}`,
			hashMap: map[string]string{"foo": "42"},
		},
		{
			input:   `{"foo" => 42, "bar" => "baz"}`,
			hashMap: map[string]string{"foo": "42", "bar": "baz"},
		},
	}

	for _, tt := range tests {
		program, err := parseSource(tt.input)
		checkParserErrors(t, err)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Logf("Expected first statement to be *ast.ExpressionStatement, got %T\n", stmt)
			t.FailNow()
		}

		testHashLiteral(t, stmt.Expression, tt.hashMap)
	}
}

func testExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	t.Helper()
	if inf, ok := expected.(infix); ok {
		return testInfixExpression(t, exp, inf.left, inf.operator, inf.right)
	}
	return testLiteralExpression(t, exp, expected)
}

type infix struct {
	left     interface{}
	operator string
	right    interface{}
}

func testInfixExpression(
	t *testing.T,
	exp ast.Expression,
	left interface{},
	operator string,
	right interface{},
) bool {
	t.Helper()
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	t.Helper()
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		if strings.HasPrefix(v, ":") {
			return testSymbol(t, exp, strings.TrimPrefix(v, ":"))
		}
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	case map[string]string:
		return testHashLiteral(t, exp, v)
	case []string:
		return testArrayLiteral(t, exp, v)
	case nil:
		return true
	}
	t.Errorf("type of expression not handled. got=%T", exp)
	return false
}

func testStringLiteral(t *testing.T, sl ast.Expression, value string) bool {
	t.Helper()
	str, ok := sl.(*ast.StringLiteral)
	if !ok {
		t.Errorf("expression not *ast.StringLiteral. got=%T", sl)
		return false
	}

	if str.Value != value {
		t.Errorf("str.Value not %s. got=%s", value, str.Value)
		return false
	}

	if str.TokenLiteral() != value {
		t.Errorf(
			"integer.TokenLiteral not %s. got=%s", value,
			str.TokenLiteral(),
		)
		return false
	}

	return true
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	t.Helper()
	if prefix, ok := il.(*ast.PrefixExpression); ok {
		if _, ok := prefix.Right.(*ast.IntegerLiteral); !ok {
			t.Errorf("expression not *ast.IntegerLiteral. got=%T", il)
			return false
		}
		if !strings.ContainsAny(prefix.Operator, "+-") {
			t.Errorf("unsupported prefix: %q", prefix.Operator)
			return false
		}
		prefixedInt := fmt.Sprintf("%s%s", prefix.Operator, prefix.Right.String())
		i, err := strconv.ParseInt(prefixedInt, 10, 64)
		if err != nil {
			t.Errorf("could not parse prefix: %v", err)
			return false
		}
		il = &ast.IntegerLiteral{Value: i}
	}
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expression not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integer.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf(
			"integer.TokenLiteral not %d. got=%s", value,
			integ.TokenLiteral(),
		)
		return false
	}

	return true
}

func testGlobal(t *testing.T, exp ast.Expression, value string) bool {
	t.Helper()
	global, ok := exp.(*ast.Global)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if global.Value != value {
		t.Errorf("global.Value not %s. got=%s", value, global.Value)
		return false
	}

	if global.TokenLiteral() != value {
		t.Errorf("global.TokenLiteral not %s. got=%s", value,
			global.TokenLiteral())
		return false
	}

	return true
}

func testSymbol(t *testing.T, exp ast.Expression, value string) bool {
	t.Helper()
	symbol, ok := exp.(*ast.SymbolLiteral)
	if !ok {
		t.Errorf("exp not %T. got=%T", symbol, exp)
		return false
	}

	if symbol.Value.String() != value {
		t.Errorf("symbol.Value not %s. got=%s", value, symbol.Value)
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	t.Helper()
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	t.Helper()
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func testArrayLiteral(t *testing.T, expr ast.Expression, value []string) bool {
	t.Helper()
	array, ok := expr.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expr not *ast.ArrayLiteral. got=%T", expr)
		return false
	}

	if len(array.Elements) != len(value) {
		t.Fatalf("len(array.Elements) not %d. got=%d", len(value), len(array.Elements))
	}

	arr := make([]string, len(array.Elements))
	for i, v := range array.Elements {
		arr[i] = v.String()
	}
	if !reflect.DeepEqual(arr, value) {
		t.Logf("Expected array to equal\n%q\n\tgot\n%q\n", value, array)
		return false
	}
	return true
}

func testHashLiteral(t *testing.T, expr ast.Expression, value map[string]string) bool {
	t.Helper()
	hash, ok := expr.(*ast.HashLiteral)
	if !ok {
		t.Errorf("expr not *ast.HashLiteral. got=%T", expr)
		return false
	}
	hashMap := make(map[string]string)
	for k, v := range hash.Map {
		hashMap[k.String()] = v.String()
	}

	if !reflect.DeepEqual(hashMap, value) {
		t.Logf("Expected hash to equal\n%q\n\tgot\n%q\n", value, hashMap)
		return false
	}
	return true
}

func parseSource(src string, modes ...Mode) (*ast.Program, *Errors) {
	mode := parseMode
	for _, m := range modes {
		mode = mode | m
	}
	prog, err := ParseFile(gotoken.NewFileSet(), "", src, mode)
	var parserErrors *Errors
	if err != nil {
		parserErrors = err.(*Errors)
	}
	return prog, parserErrors
}

func parseExpression(src string, modes ...Mode) (ast.Expression, *Errors) {
	mode := parseMode
	for _, m := range modes {
		mode = mode | m
	}
	expr, err := ParseExprFrom(gotoken.NewFileSet(), "", src, mode)
	var parserErrors *Errors
	if err != nil {
		parserErrors = err.(*Errors)
	}
	return expr, parserErrors
}

func compareFirstParserError(t *testing.T, expected, actual error) {
	t.Helper()
	if expected == nil && actual == nil {
		return
	}
	parserErrors, ok := actual.(*Errors)
	if parserErrors == nil && expected == nil {
		return
	}
	if !ok {
		t.Logf("Unexpected parser error: %T:%v\n", actual, actual)
		t.FailNow()
	}
	if expected == nil && parserErrors != nil {
		t.Logf("Expected no error, got %T:%v", actual, actual)
		t.FailNow()
	}
	firstErr := parserErrors.errors[0]
	err := firstErr.Error()
	firstSpace := strings.Index(err, " ")
	err = err[firstSpace+1:]
	if err != expected.Error() {
		t.Logf("Expected first parser error to equal %v, got %v", expected, firstErr)
		t.FailNow()
	}
}

func checkParserErrors(t *testing.T, err error, withStack ...bool) {
	t.Helper()
	if err == nil {
		return
	}
	printStack := false
	if len(withStack) != 0 {
		printStack = withStack[0]
	}
	parserErrors, ok := err.(*Errors)
	if parserErrors == nil {
		return
	}
	if !ok {
		t.Logf("Unexpected parser error: %T:%v\n", err, err)
		t.FailNow()
	}

	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	t.Errorf("parser has %d errors", len(parserErrors.errors))
	for _, e := range parserErrors.errors {
		t.Errorf("%v", e)
		if stackErr, ok := e.(stackTracer); ok && printStack {
			st := stackErr.StackTrace()
			fmt.Printf("Error stack:%+v\n", st[0:2]) // top two frames
		}

	}
	t.FailNow()
}
