package ast

import (
	"testing"

	"github.com/goruby/goruby/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &Variable{
					Name: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "myVar"},
						Value: "myVar",
					},
					Value: &Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
						Value: "anotherVar",
					},
				},
			},
		},
	}
	if program.String() != "myVar = anotherVar" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
