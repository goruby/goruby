package ast

import (
	"bytes"
	"fmt"
	gotoken "go/token"
	"strings"

	"github.com/goruby/goruby/token"
)

// Node represents a node within the AST
//
// All node types implement the Node interface.
type Node interface {
	// Pos returns the position of first character belonging to the node
	Pos() int
	// End returns the position of first character immediately after the node
	End() int
	// TokenLiteral returns the literal of the node
	TokenLiteral() string
	// String returns a string representation of the node
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

// IsLiteral returns true if n is a literal node, false otherwise
func IsLiteral(n Node) bool {
	_, ok := n.(literal)
	return ok
}

// A Program node is the root node within the AST.
type Program struct {
	pos        int
	File       *gotoken.File
	Statements []Statement
}

// Pos returns the position of first character belonging to the node
func (p *Program) Pos() int { return p.pos }

// End returns the position of first character immediately after the node
func (p *Program) End() int {
	if len(p.Statements) == 0 {
		return p.pos
	}
	return p.Statements[len(p.Statements)-1].End()
}
func (p *Program) String() string {
	stmts := make([]string, len(p.Statements))
	for i, s := range p.Statements {
		if s != nil {
			stmts[i] = s.String()
		}
	}
	return strings.Join(stmts, "\n")
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

// Pos returns the position of first character belonging to the node
func (rs *ReturnStatement) Pos() int { return rs.Token.Pos }

// End returns the position of first character immediately after the node
func (rs *ReturnStatement) End() int { return rs.ReturnValue.End() }

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

// Pos returns the position of first character belonging to the node
func (es *ExpressionStatement) Pos() int { return es.Expression.Pos() }

// End returns the position of first character immediately after the node
func (es *ExpressionStatement) End() int { return es.Expression.End() }

// TokenLiteral returns the first token of the Expression
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// BlockStatement represents a list of statements
type BlockStatement struct {
	// the { token or the first token from the first statement
	Token      token.Token
	EndToken   token.Token // the } token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// Pos returns the position of first character belonging to the node
func (bs *BlockStatement) Pos() int { return bs.Token.Pos }

// End returns the position of first character immediately after the node
func (bs *BlockStatement) End() int { return bs.EndToken.Pos }

// TokenLiteral returns '{' or the first token from the first statement
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	for _, s := range bs.Statements {
		if s != nil {
			out.WriteString(s.String())
		}
	}
	return out.String()
}

// ExceptionHandlingBlock represents a begin/end block where exceptions are rescued
type ExceptionHandlingBlock struct {
	BeginToken token.Token
	EndToken   token.Token
	TryBody    *BlockStatement
	Rescues    []*RescueBlock
}

func (eh *ExceptionHandlingBlock) expressionNode() {}

// Pos returns the position of first character belonging to the node
func (eh *ExceptionHandlingBlock) Pos() int { return eh.BeginToken.Pos }

// End returns the position of first character immediately after the node
func (eh *ExceptionHandlingBlock) End() int { return eh.EndToken.Pos }

// TokenLiteral returns the token literal from 'begin'
func (eh *ExceptionHandlingBlock) TokenLiteral() string { return eh.BeginToken.Literal }
func (eh *ExceptionHandlingBlock) String() string {
	var out bytes.Buffer
	out.WriteString(eh.BeginToken.Literal)
	out.WriteString("\n")
	out.WriteString(eh.TryBody.String())
	out.WriteString("\n")
	for _, r := range eh.Rescues {
		out.WriteString(r.String())
	}
	out.WriteString("end")
	return out.String()
}

// A RescueBlock represents a rescue block
type RescueBlock struct {
	Token            token.Token
	ExceptionClasses []*Identifier
	Exception        *Identifier
	Body             *BlockStatement
}

func (rb *RescueBlock) expressionNode() {}

// Pos returns the position of first character belonging to the node
func (rb *RescueBlock) Pos() int { return rb.Token.Pos }

// End returns the position of first character immediately after the node
func (rb *RescueBlock) End() int { return rb.Body.End() }

// TokenLiteral returns the token literal from 'rescue'
func (rb *RescueBlock) TokenLiteral() string { return rb.Token.Literal }
func (rb *RescueBlock) String() string {
	var out bytes.Buffer
	out.WriteString(rb.Token.Literal)
	if len(rb.ExceptionClasses) != 0 {
		classes := make([]string, len(rb.ExceptionClasses))
		for i, c := range rb.ExceptionClasses {
			classes[i] = c.String()
		}
		out.WriteString(strings.Join(classes, ", "))
	}
	if rb.Exception != nil {
		out.WriteString(" => ")
		out.WriteString(rb.Exception.String())
	}
	out.WriteString("\n")
	out.WriteString(rb.Body.String())
	out.WriteString("\n")
	return out.String()
}

// Assignment represents a generic assignment
type Assignment struct {
	Token token.Token
	Left  Expression
	Right Expression
}

func (a *Assignment) String() string {
	var out bytes.Buffer
	out.WriteString(encloseInParensIfNeeded(a.Left))
	out.WriteString(" = ")
	out.WriteString(encloseInParensIfNeeded(a.Right))
	return out.String()
}
func (a *Assignment) expressionNode() {}

// Pos returns the position of first character belonging to the node
func (a *Assignment) Pos() int { return a.Left.Pos() }

// End returns the position of first character immediately after the node
func (a *Assignment) End() int { return a.Right.End() }

// TokenLiteral returns the literal of the ASSIGN token
func (a *Assignment) TokenLiteral() string { return a.Token.Literal }

// An InstanceVariable represents an instance variable in the AST
type InstanceVariable struct {
	Token token.Token
	Name  *Identifier
}

func (i *InstanceVariable) String() string {
	var out bytes.Buffer
	out.WriteString(i.Token.Literal)
	out.WriteString(i.Name.String())
	return out.String()
}
func (i *InstanceVariable) literalNode()    {}
func (i *InstanceVariable) expressionNode() {}

// Pos returns the position of first character belonging to the node
func (i *InstanceVariable) Pos() int { return i.Token.Pos }

// End returns the position of first character immediately after the node
func (i *InstanceVariable) End() int { return i.Name.End() }

// TokenLiteral returns the literal of the AT token
func (i *InstanceVariable) TokenLiteral() string { return i.Token.Literal }

// MultiAssignment represents multiple variables on the lefthand side
type MultiAssignment struct {
	Variables []*Identifier
	Values    []Expression
}

func (m *MultiAssignment) String() string {
	var out bytes.Buffer
	vars := make([]string, len(m.Variables))
	for i, v := range m.Variables {
		vars[i] = v.Value
	}
	out.WriteString(strings.Join(vars, ", "))
	out.WriteString(" = ")
	values := make([]string, len(m.Values))
	for i, v := range m.Values {
		values[i] = v.String()
	}
	out.WriteString(strings.Join(values, ", "))
	return out.String()
}
func (m *MultiAssignment) literalNode() {}

// Pos returns the position of first character belonging to the node
func (m *MultiAssignment) Pos() int { return m.Variables[0].Pos() }

// End returns the position of first character immediately after the node
func (m *MultiAssignment) End() int        { return m.Values[len(m.Values)-1].End() }
func (m *MultiAssignment) expressionNode() {}

// TokenLiteral returns the literal of the first variable token
func (m *MultiAssignment) TokenLiteral() string { return m.Variables[0].Token.Literal }

// Self represents self in the current context in the program
type Self struct {
	Token token.Token // the token.SELF token
}

func (s *Self) String() string  { return s.Token.Literal }
func (s *Self) expressionNode() {}
func (s *Self) literalNode()    {}

// Pos returns the position of first character belonging to the node
func (s *Self) Pos() int { return s.Token.Pos }

// End returns the position of first character immediately after the node
func (s *Self) End() int { return s.Token.Pos + 4 }

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

// Pos returns the position of first character belonging to the node
func (y *YieldExpression) Pos() int { return y.Token.Pos }

// End returns the position of first character immediately after the node
func (y *YieldExpression) End() int {
	if len(y.Arguments) == 0 {
		return y.Pos() + 5
	}
	return y.Arguments[len(y.Arguments)-1].End()
}

// TokenLiteral returns the literal of the token.YIELD token
func (y *YieldExpression) TokenLiteral() string { return y.Token.Literal }

// Keyword__FILE__ represents __FILE__ in the AST
type Keyword__FILE__ struct {
	Token    token.Token // the token.FILE__ token
	Filename string
}

func (f *Keyword__FILE__) String() string  { return f.Token.Literal }
func (f *Keyword__FILE__) expressionNode() {}
func (f *Keyword__FILE__) literalNode()    {}

// Pos returns the position of first character belonging to the node
func (f *Keyword__FILE__) Pos() int { return f.Token.Pos }

// End returns the position of first character immediately after the node
func (f *Keyword__FILE__) End() int { return f.Token.Pos + 8 }

// TokenLiteral returns the literal of the token.FILE__ token
func (f *Keyword__FILE__) TokenLiteral() string { return f.Token.Literal }

// An Identifier represents an identifier in the program
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) String() string  { return i.Value }
func (i *Identifier) expressionNode() {}
func (i *Identifier) literalNode()    {}

// Pos returns the position of first character belonging to the node
func (i *Identifier) Pos() int { return i.Token.Pos }

// End returns the position of first character immediately after the node
func (i *Identifier) End() int { return i.Token.Pos + len(i.Value) }

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

// Pos returns the position of first character belonging to the node
func (g *Global) Pos() int { return g.Token.Pos }

// End returns the position of first character immediately after the node
func (g *Global) End() int     { return g.Token.Pos + len(g.Value) }
func (g *Global) literalNode() {}

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

// Pos returns the position of first character belonging to the node
func (i *ScopedIdentifier) Pos() int { return i.Outer.Pos() }

// End returns the position of first character immediately after the node
func (i *ScopedIdentifier) End() int { return i.Inner.End() }

// TokenLiteral returns the literal of the token.SCOPE token
func (i *ScopedIdentifier) TokenLiteral() string { return i.Token.Literal }

// IntegerLiteral represents an integer in the AST
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) literalNode()    {}

// Pos returns the position of first character belonging to the node
func (il *IntegerLiteral) Pos() int { return il.Token.Pos }

// End returns the position of first character immediately after the node
func (il *IntegerLiteral) End() int { return il.Token.Pos + len(fmt.Sprintf("%d", il.Value)) }

// TokenLiteral returns the literal from the token.INT token
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return fmt.Sprintf("%d", il.Value) }

// Nil represents the 'nil' keyword
type Nil struct {
	Token token.Token
}

func (n *Nil) expressionNode() {}
func (n *Nil) literalNode()    {}

// Pos returns the position of first character belonging to the node
func (n *Nil) Pos() int { return n.Token.Pos }

// End returns the position of first character immediately after the node
func (n *Nil) End() int { return n.Token.Pos + 3 }

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

// Pos returns the position of first character belonging to the node
func (b *Boolean) Pos() int { return b.Token.Pos }

// End returns the position of first character immediately after the node
func (b *Boolean) End() int { return b.Token.Pos + len(fmt.Sprintf("%t", b.Value)) }

// TokenLiteral returns the literal from the token token.BOOLEAN
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return fmt.Sprintf("%t", b.Value) }

// StringLiteral represents a double quoted string in the AST
type StringLiteral struct {
	Token token.Token // the '"'
	Value string
}

func (sl *StringLiteral) expressionNode() {}
func (sl *StringLiteral) literalNode()    {}

// Pos returns the position of first character belonging to the node
func (sl *StringLiteral) Pos() int { return sl.Token.Pos }

// End returns the position of first character immediately after the node
func (sl *StringLiteral) End() int { return sl.Token.Pos + len(sl.Value) }

// TokenLiteral returns the literal from token token.STRING
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Value }

// Comment represents a double quoted string in the AST
type Comment struct {
	Token token.Token // the #
	Value string
}

func (c *Comment) statementNode() {}
func (c *Comment) literalNode()   {}

// Pos returns the position of first character belonging to the node
func (c *Comment) Pos() int { return c.Token.Pos }

// End returns the position of first character immediately after the node
func (c *Comment) End() int { return c.Token.Pos + len(c.Value) }

// TokenLiteral returns the literal from token token.STRING
func (c *Comment) TokenLiteral() string { return c.Token.Literal }
func (c *Comment) String() string       { return c.Value }

// SymbolLiteral represents a symbol within the AST
type SymbolLiteral struct {
	Token token.Token // the ':'
	Value Expression
}

func (s *SymbolLiteral) expressionNode() {}
func (s *SymbolLiteral) literalNode()    {}

// Pos returns the position of first character belonging to the node
func (s *SymbolLiteral) Pos() int { return s.Token.Pos }

// End returns the position of first character immediately after the node
func (s *SymbolLiteral) End() int { return s.Value.End() }

// TokenLiteral returns the literal from token token.SYMBOL
func (s *SymbolLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *SymbolLiteral) String() string       { return ":" + s.Value.String() }

// ConditionalExpression represents an if expression within the AST
type ConditionalExpression struct {
	Token       token.Token // The 'if' or 'unless' token
	EndToken    token.Token // The 'end' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

// IsNegated indicates if the condition uses unless, i.e. is negated
func (ce *ConditionalExpression) IsNegated() bool {
	return ce.Token.Type == token.UNLESS
}

func (ce *ConditionalExpression) expressionNode() {}

// Pos returns the position of first character belonging to the node
func (ce *ConditionalExpression) Pos() int {
	if ce.EndToken.Type == token.ILLEGAL {
		return ce.Consequence.Pos()
	}
	return ce.Token.Pos
}

// End returns the position of first character immediately after the node
func (ce *ConditionalExpression) End() int {
	if ce.EndToken.Type == token.ILLEGAL {
		return ce.Consequence.Pos()
	}
	return ce.EndToken.Pos
}

// TokenLiteral returns the literal from token token.IF or token.UNLESS
func (ce *ConditionalExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *ConditionalExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Token.Literal)
	out.WriteString(ce.Condition.String())
	out.WriteString(" ")
	out.WriteString(ce.Consequence.String())
	if ce.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ce.Alternative.String())
	}
	out.WriteString(" end")
	return out.String()
}

// A LoopExpression represents a loop
type LoopExpression struct {
	Token     token.Token // while
	EndToken  token.Token // end
	Condition Expression
	Block     *BlockStatement
}

func (ce *LoopExpression) expressionNode() {}

// Pos returns the position of first character belonging to the node
func (ce *LoopExpression) Pos() int {
	return ce.Token.Pos
}

// End returns the position of first character immediately after the node
func (ce *LoopExpression) End() int {
	return ce.EndToken.Pos
}

// TokenLiteral returns the literal from token token.WHILE
func (ce *LoopExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *LoopExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Token.Literal)
	out.WriteString(ce.Condition.String())
	out.WriteString(" do ")
	out.WriteString(ce.Block.String())
	out.WriteString(" end")
	return out.String()
}

// ExpressionList represents a list of expressions within the AST divided by commas
type ExpressionList []Expression

func (el ExpressionList) expressionNode() {}
func (el ExpressionList) literalNode()    {}

// Pos returns the position of first character from the first expression
func (el ExpressionList) Pos() int {
	if len(el) == 0 {
		return 0
	}
	return el[0].End()
}

// End returns End of the last element
func (el ExpressionList) End() int {
	if len(el) == 0 {
		return 0
	}
	return el[len(el)-1].End()
}

// TokenLiteral returns the literal of the first element
func (el ExpressionList) TokenLiteral() string {
	if len(el) == 0 {
		return ""
	}
	return el[0].TokenLiteral()
}
func (el ExpressionList) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range el {
		elements = append(elements, e.String())
	}
	out.WriteString(strings.Join(elements, ", "))
	return out.String()
}

// ArrayLiteral represents an Array literal within the AST
type ArrayLiteral struct {
	Token    token.Token // the '['
	Rbracket token.Token // the ']'
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}
func (al *ArrayLiteral) literalNode()    {}

// Pos returns the position of first character belonging to the node
func (al *ArrayLiteral) Pos() int { return al.Token.Pos }

// End returns the position of first character immediately after the node
func (al *ArrayLiteral) End() int {
	return al.Rbracket.Pos
}

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
	Token  token.Token // the '{'
	Rbrace token.Token // the '{'
	Map    map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}
func (hl *HashLiteral) literalNode()    {}

// Pos returns the position of the left brace
func (hl *HashLiteral) Pos() int { return hl.Token.Pos }

// End returns the position of the right brace
func (hl *HashLiteral) End() int { return hl.Rbrace.Pos }

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

// A BlockCapture represents a function scoped variable capturing a block
type BlockCapture struct {
	Token token.Token // the `&`
	Name  *Identifier
}

func (b *BlockCapture) expressionNode() {}
func (b *BlockCapture) literalNode()    {}

// Pos returns the position of the ampersand
func (b *BlockCapture) Pos() int { return b.Token.Pos }

// End returns the position of the last character of Name
func (b *BlockCapture) End() int { return b.Name.End() }
func (b *BlockCapture) String() string {
	return "&" + b.Name.Value
}

// TokenLiteral returns the literal of the token
func (b *BlockCapture) TokenLiteral() string { return b.Token.Literal }

// A FunctionLiteral represents a function definition in the AST
type FunctionLiteral struct {
	Token         token.Token // The 'def' token
	EndToken      token.Token // the 'end' token
	Receiver      *Identifier
	Name          *Identifier
	Parameters    []*FunctionParameter
	CapturedBlock *BlockCapture
	Body          *BlockStatement
	Rescues       []*RescueBlock
}

func (fl *FunctionLiteral) expressionNode() {}
func (fl *FunctionLiteral) literalNode()    {}

// Pos returns the position of the `def` keyword
func (fl *FunctionLiteral) Pos() int { return fl.Token.Pos }

// End returns the position of the `end` keyword
func (fl *FunctionLiteral) End() int { return fl.EndToken.Pos }

// TokenLiteral returns the literal from token.DEF
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	if fl.CapturedBlock != nil {
		params = append(params, fl.CapturedBlock.String())
	}
	out.WriteString("def ")
	if fl.Receiver != nil {
		out.WriteString(fl.Receiver.String())
		out.WriteString(".")
	}
	out.WriteString(fl.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	if fl.Body != nil {
		out.WriteString(fl.Body.String())
	}
	for _, r := range fl.Rescues {
		out.WriteString(r.String())
	}
	out.WriteString(" end")
	return out.String()
}

// A FunctionParameter represents a parameter in a function literal
type FunctionParameter struct {
	Name    *Identifier
	Default Expression
}

func (f *FunctionParameter) expressionNode() {}

// Pos returns the position of first character belonging to the node
func (f *FunctionParameter) Pos() int { return f.Name.Pos() }

// End returns the position of the default end if it exists, otherwise the end position of Name
func (f *FunctionParameter) End() int {
	if f.Default != nil {
		return f.Default.End()
	}
	return f.Name.End()
}

// TokenLiteral returns the token of the parameter name
func (f *FunctionParameter) TokenLiteral() string { return f.Name.TokenLiteral() }
func (f *FunctionParameter) String() string {
	var out bytes.Buffer
	out.WriteString(f.Name.String())
	if f.Default != nil {
		out.WriteString(" = ")
		out.WriteString(encloseInParensIfNeeded(f.Default))
	}
	return out.String()
}

// An IndexExpression represents an array or hash access in the AST
type IndexExpression struct {
	Token  token.Token // The [ token
	Left   Expression
	Index  Expression
	Length Expression
}

func (ie *IndexExpression) expressionNode() {}

// Pos returns the position of first character belonging to the node
func (ie *IndexExpression) Pos() int { return ie.Token.Pos }

// End returns the position of the last character belonging to the node
func (ie *IndexExpression) End() int { return ie.Index.End() }

// TokenLiteral returns the literal from token.LBRACKET
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	if ie.Length != nil {
		out.WriteString(", ")
		out.WriteString(ie.Index.String())
	}
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

// Pos returns the position of first character belonging to the node
func (ce *ContextCallExpression) Pos() int {
	if ce.Context != nil {
		return ce.Context.Pos()
	}
	return ce.Function.Pos()
}

// End returns the end position of the block if it exists. If not, it returns
// the end position of the last argument if any. Otherwise it returns the end
// of the function identifier
func (ce *ContextCallExpression) End() int {
	if ce.Block != nil {
		return ce.Block.End()
	}
	if len(ce.Arguments) == 0 {
		return ce.Function.End()
	}
	return ce.Arguments[len(ce.Arguments)-1].End()
}

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
	Token      token.Token          // token.DO or token.LBRACE
	EndToken   token.Token          // token.END or token.RBRACE
	Parameters []*FunctionParameter // the block parameters
	Body       *BlockStatement      // the block body
}

func (b *BlockExpression) expressionNode() {}

// Pos returns the position of first character belonging to the node
func (b *BlockExpression) Pos() int { return b.Token.Pos }

// End returns the position of the end token
func (b *BlockExpression) End() int { return b.EndToken.Pos }

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
	Token    token.Token // The module keyword
	EndToken token.Token // The end token
	Name     *Identifier // The module name, will always be a const
	Body     *BlockStatement
}

func (m *ModuleExpression) expressionNode() {}

// Pos returns the position of first character belonging to the node
func (m *ModuleExpression) Pos() int { return m.Token.Pos }

// End returns the position of the `end` token
func (m *ModuleExpression) End() int { return m.EndToken.Pos }

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
	Token      token.Token // The class keyword
	EndToken   token.Token // The end token
	Name       *Identifier // The class name, will always be a const
	SuperClass *Identifier // The superclass, if any
	Body       *BlockStatement
}

func (m *ClassExpression) expressionNode() {}

// Pos returns the position of first character belonging to the node
func (m *ClassExpression) Pos() int { return m.Token.Pos }

// End returns the position of the `end` token
func (m *ClassExpression) End() int { return m.EndToken.Pos }

// TokenLiteral returns the literal from token.CLASS
func (m *ClassExpression) TokenLiteral() string { return m.Token.Literal }
func (m *ClassExpression) String() string {
	var out bytes.Buffer
	out.WriteString(m.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(m.Name.String())
	if m.SuperClass != nil {
		out.WriteString(" ")
		out.WriteString("<")
		out.WriteString(" ")
		out.WriteString(m.SuperClass.String())
	}
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

// Pos returns the position of first character belonging to the node
func (pe *PrefixExpression) Pos() int { return pe.Token.Pos }

// End returns the end of the right expression
func (pe *PrefixExpression) End() int { return pe.Right.End() }

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

// MustEvaluateRight returns true if it is mandatory to evaluate the right side
// of the operator, false otherwise
func (oe *InfixExpression) MustEvaluateRight() bool {
	return oe.Token.Type != token.LOGICALOR
}

// IsControlExpression returns true if the infix is used for control flow,
// false otherwise
func (oe *InfixExpression) IsControlExpression() bool {
	return oe.Token.Type == token.LOGICALOR || oe.Token.Type == token.LOGICALAND
}

func (oe *InfixExpression) expressionNode() {}

// Pos returns the position of first character belonging to the left node
func (oe *InfixExpression) Pos() int { return oe.Left.Pos() }

// End returns the position of last character belonging to the right node
func (oe *InfixExpression) End() int { return oe.Right.End() }

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

func encloseInParensIfNeeded(expr Expression) string {
	val := expr.String()
	hasParens := strings.HasPrefix(val, "(") && strings.HasSuffix(val, ")")
	_, isLiteral := expr.(literal)
	if !isLiteral && !hasParens {
		val = "(" + val + ")"
	}
	return val
}
