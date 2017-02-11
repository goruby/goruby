package parser

import (
	"testing"

	"github.com/goruby/goruby/ast"
	"github.com/goruby/goruby/lexer"
)

func TestVariableStatements(t *testing.T) {
	input := `
x = 5
y = 10
foobar = 838383
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}
	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVariableStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVariableStatement(t *testing.T, s ast.Statement, name string) bool {
	variableStmt, ok := s.(*ast.VariableStatement)
	if !ok {
		t.Errorf("s not *ast.VariableStatement. got=%T", s)
		return false
	}
	if variableStmt.Name.Value != name {
		t.Errorf("variableStmt.Name.Value not '%s'. got=%s", name, variableStmt.Name.Value)
		return false
	}
	if variableStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, variableStmt.Name)
		return false
	}
	return true
}
