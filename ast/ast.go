package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/goruby/goruby/token"
)

// Node represents a node within the AST
//
// All node types implement the Node interface.
type Node interface {
	// TokenLiteral returns the literal of the node
	TokenLiteral() string
	String() string
}

// A Statement represents a statement within the AST
//
// All statement nodes implement the Statement interface.
type Statement interface {
	Node
	statementNode()
}

// An Expression represents an expression within the AST
//
// All expression nodes implement the Expression interface.
type Expression interface {
	Node
	expressionNode()
}

// literal
type literal interface {
	Node
	literalNode()
}

// A Program node is the root node within the AST.
type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// TokenLiteral returns the literal of the first statement and empty string if
// there is no statement.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// A ReturnStatement represents a return node which yields another Expression.
type ReturnStatement struct {
	Token       token.Token // the 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	return out.String()
}
func (rs *ReturnStatement) statementNode() {}

// TokenLiteral returns the 'return' token literal
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// An ExpressionStatement is a Statement wrapping an Expression
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
func (es *ExpressionStatement) statementNode() {}

// TokenLiteral returns the first token of the Expression
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// BlockStatement represents a list of statements
type BlockStatement struct {
	// the { token or the first token from the first statement
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral returns '{' or the first token from the first statement
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// GlobalAssignment represents a global assignment
type GlobalAssignment struct {
	Name  *Global
	Value Expression
}

func (v *GlobalAssignment) String() string {
	var out bytes.Buffer
	out.WriteString(v.Name.String())
	out.WriteString(" = ")
	if v.Value != nil {
		val := v.Value.String()
		hasParens := strings.HasPrefix(val, "(") && strings.HasSuffix(val, ")")
		_, isLiteral := v.Value.(literal)
		if !isLiteral && !hasParens {
			val = "(" + val + ")"
		}
		out.WriteString(val)
	}
	return out.String()
}
func (v *GlobalAssignment) expressionNode() {}

// TokenLiteral returns the literal of the Name token
func (v *GlobalAssignment) TokenLiteral() string { return v.Name.Token.Literal }

// VariableAssignment represents a variable assignment
type VariableAssignment struct {
	Name  *Identifier
	Value Expression
}

func (v *VariableAssignment) String() string {
	var out bytes.Buffer
	out.WriteString(v.Name.String())
	out.WriteString(" = ")
	if v.Value != nil {
		val := v.Value.String()
		hasParens := strings.HasPrefix(val, "(") && strings.HasSuffix(val, ")")
		_, isLiteral := v.Value.(literal)
		if !isLiteral && !hasParens {
			val = "(" + val + ")"
		}
		out.WriteString(val)
	}
	return out.String()
}
func (v *VariableAssignment) expressionNode() {}

// TokenLiteral returns the literal of the Name token
func (v *VariableAssignment) TokenLiteral() string { return v.Name.Token.Literal }

// Self represents self in the current context in the program
type Self struct {
	Token token.Token // the token.SELF token
}

func (s *Self) String() string  { return s.Token.Literal }
func (s *Self) expressionNode() {}
func (s *Self) literalNode()    {}

// TokenLiteral returns the literal of the token.SELF token
func (s *Self) TokenLiteral() string { return s.Token.Literal }

// YieldExpression represents self in the current context in the program
type YieldExpression struct {
	Token     token.Token  // the token.YIELD token
	Arguments []Expression // The arguments to yield
}

func (y *YieldExpression) String() string {
	var out bytes.Buffer
	out.WriteString(y.Token.Literal)
	if len(y.Arguments) != 0 {
		args := []string{}
		for _, a := range y.Arguments {
			args = append(args, a.String())
		}
		out.WriteString(" ")
		out.WriteString(strings.Join(args, ", "))
	}
	return out.String()
}
func (y *YieldExpression) expressionNode() {}

// TokenLiteral returns the literal of the token.YIELD token
func (y *YieldExpression) TokenLiteral() string { return y.Token.Literal }

// An Identifier represents an identifier in the program
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) String() string  { return i.Value }
func (i *Identifier) expressionNode() {}
func (i *Identifier) literalNode()    {}

// IsConstant returns true if the Identifier represents a Constant, false otherwise
func (i *Identifier) IsConstant() bool { return i.Token.Type == token.CONST }

// TokenLiteral returns the literal of the token.IDENT token
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

// Global represents a global in the AST
type Global struct {
	Token token.Token // the token.GLOBAL token
	Value string
}

func (g *Global) String() string  { return g.Value }
func (g *Global) expressionNode() {}
func (g *Global) literalNode()    {}

// TokenLiteral returns the literal of the token.GLOBAL token
func (g *Global) TokenLiteral() string { return g.Token.Literal }

// ScopedIdentifier represents a scoped Constant declaration
type ScopedIdentifier struct {
	Token token.Token // the token.SCOPE
	Outer *Identifier
	Inner Expression
}

func (i *ScopedIdentifier) String() string {
	var out bytes.Buffer
	out.WriteString(i.Outer.String())
	out.WriteString(i.Token.Literal)
	out.WriteString(i.Inner.String())
	return out.String()
}
func (i *ScopedIdentifier) expressionNode() {}
func (i *ScopedIdentifier) literalNode()    {}

// TokenLiteral returns the literal of the token.SCOPE token
func (i *ScopedIdentifier) TokenLiteral() string { return i.Token.Literal }

// IntegerLiteral represents an integer in the AST
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) literalNode()    {}

// TokenLiteral returns the literal from the token.INT token
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// Nil represents the 'nil' keyword
type Nil struct {
	Token token.Token
}

func (n *Nil) expressionNode() {}
func (n *Nil) literalNode()    {}

// TokenLiteral returns the literal from the token token.NIL
func (n *Nil) TokenLiteral() string { return n.Token.Literal }
func (n *Nil) String() string       { return "nil" }

// Boolean represents a boolean in the AST
type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) literalNode()    {}

// TokenLiteral returns the literal from the token token.BOOLEAN
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

// StringLiteral represents a double quoted string in the AST
type StringLiteral struct {
	Token token.Token // the '"'
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) literalNode()    {}

// TokenLiteral returns the literal from token token.STRING
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Value }

// SymbolLiteral represents a symbol within the AST
type SymbolLiteral struct {
	Token token.Token // the ':'
	Value Expression
}

func (s *SymbolLiteral) expressionNode() {}
func (s *SymbolLiteral) literalNode()    {}

// TokenLiteral returns the literal from token token.SYMBOL
func (s *SymbolLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *SymbolLiteral) String() string       { return ":" + s.Value.String() }

// IfExpression represents an if expression within the AST
type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral returns the literal from token token.IF
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}
	out.WriteString(" end")
	return out.String()
}

// ArrayLiteral represents an Array literal within the AST
type ArrayLiteral struct {
	Token    token.Token // the '['
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) literalNode()    {}

// TokenLiteral returns the literal of the token token.LBRACKET
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// HashLiteral represents an Hash literal within the AST
type HashLiteral struct {
	Token token.Token // the '{'
	Map   map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}
func (hl *HashLiteral) literalNode()    {}

// TokenLiteral returns the literal of the token token.LBRACE
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for key, val := range hl.Map {
		elements = append(elements, fmt.Sprintf("%q => %q", key.String(), val.String()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("}")
	return out.String()
}

// A FunctionLiteral represents a function definition in the AST
type FunctionLiteral struct {
	Token      token.Token // The 'def' token
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) literalNode()    {}

// TokenLiteral returns the literal from token.DEF
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fl.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())
	out.WriteString(" end")
	return out.String()
}

// An IndexExpression represents an array or hash access in the AST
type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral returns the literal from token.LBRACKET
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// A ContextCallExpression represents a method call on a given Context
type ContextCallExpression struct {
	Token     token.Token      // The '.' token
	Context   Expression       // The lefthandside expression
	Function  *Identifier      // The function to call
	Arguments []Expression     // The function arguments
	Block     *BlockExpression // The function block
}

func (ce *ContextCallExpression) expressionNode() {}

// TokenLiteral returns the literal from token.DOT
func (ce *ContextCallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *ContextCallExpression) String() string {
	var out bytes.Buffer
	if ce.Context != nil {
		out.WriteString(ce.Context.String())
		out.WriteString(".")
	}
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	if ce.Block != nil {
		out.WriteString("\n")
		out.WriteString(ce.Block.String())
	}
	return out.String()
}

// A BlockExpression represents a Ruby block
type BlockExpression struct {
	Token      token.Token     // token.DO or token.LBRACE
	Parameters []*Identifier   // the block parameters
	Body       *BlockStatement // the block body
}

func (b *BlockExpression) expressionNode() {}

// TokenLiteral returns the literal from the Token
func (b *BlockExpression) TokenLiteral() string { return b.Token.Literal }

// String returns a string representation of the block statement
func (b *BlockExpression) String() string {
	var out bytes.Buffer
	out.WriteString(b.Token.Literal)
	if len(b.Parameters) != 0 {
		args := []string{}
		for _, a := range b.Parameters {
			args = append(args, a.String())
		}
		out.WriteString("|")
		out.WriteString(strings.Join(args, ", "))
		out.WriteString("|")
		out.WriteString("\n")
	}
	out.WriteString(b.Body.String())
	out.WriteString("\n")
	if b.Token.Type == token.LBRACE {
		out.WriteString("}")
	} else {
		out.WriteString("end")
	}
	return out.String()
}

// ModuleExpression represents a module definition
type ModuleExpression struct {
	Token token.Token // The module keyword
	Name  *Identifier // The module name, will always be a const
	Body  *BlockStatement
}

func (m *ModuleExpression) expressionNode() {}

// TokenLiteral returns the literal from token.MODULE
func (m *ModuleExpression) TokenLiteral() string { return m.Token.Literal }
func (m *ModuleExpression) String() string {
	var out bytes.Buffer
	out.WriteString(m.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(m.Name.String())
	out.WriteString("\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n")
	out.WriteString("end")
	return out.String()
}

// ClassExpression represents a module definition
type ClassExpression struct {
	Token token.Token // The class keyword
	Name  *Identifier // The class name, will always be a const
	Body  *BlockStatement
}

func (m *ClassExpression) expressionNode() {}

// TokenLiteral returns the literal from token.CLASS
func (m *ClassExpression) TokenLiteral() string { return m.Token.Literal }
func (m *ClassExpression) String() string {
	var out bytes.Buffer
	out.WriteString(m.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(m.Name.String())
	out.WriteString("\n")
	out.WriteString(m.Body.String())
	out.WriteString("\n")
	out.WriteString(" end")
	return out.String()
}

// PrefixExpression represents a prefix operator
type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral returns the literal from the prefix operator token
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// An InfixExpression represents an infix operator in the AST
type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (oe *InfixExpression) expressionNode() {}

// TokenLiteral returns the literal from the infix operator token
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")
	return out.String()
}
