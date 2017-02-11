package ast

import (
	"github.com/goruby/goruby/token"
)

type Node interface {
	TokenLiteral() string
}
type Statement interface {
	Node
	statementNode()
}
type Expression interface {
	Node
	expressionNode()
}
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type VariableStatement struct {
	Name  *Identifier
	Value Expression
}

func (v *VariableStatement) statementNode()       {}
func (v *VariableStatement) TokenLiteral() string { return v.Name.Token.Literal }

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
