package parser

import (
	"fmt"
	"strconv"

	"github.com/goruby/goruby/ast"
	"github.com/goruby/goruby/lexer"
	"github.com/goruby/goruby/token"
	"github.com/pkg/errors"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	ASSIGNMENT  // x = 5
	CALL        // myFunction(X)
	CONTEXT     // foo.myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.ASSIGN:   ASSIGNMENT,
	token.LPAREN:   CALL,
	token.IDENT:    CALL,
	token.INT:      CALL,
	token.STRING:   CALL,
	token.SYMBOL:   CALL,
	token.DOT:      CONTEXT,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

var defaultExpressionTerminators = []token.TokenType{
	token.SEMICOLON,
	token.NEWLINE,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []error{},
	}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.DEF, p.parseFunctionLiteral)
	p.registerPrefix(token.SYMBOL, p.parseSymbolLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpressionWithParens)
	p.registerInfix(token.IDENT, p.parseCallExpression)
	p.registerInfix(token.INT, p.parseCallExpression)
	p.registerInfix(token.STRING, p.parseCallExpression)
	p.registerInfix(token.DOT, p.parseContextCallExpression)
	p.registerInfix(token.SYMBOL, p.parseCallExpression)
	p.registerInfix(token.ASSIGN, p.parseVariableAssignExpression)
	return p
}

type Parser struct {
	l      *lexer.Lexer
	errors []error

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) peekError(t ...token.TokenType) {
	err := &unexpectedTokenError{
		expectedTokens: t,
		actualToken:    p.peekToken.Type,
	}
	p.errors = append(p.errors, err)
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	if len(p.errors) != 0 {
		return program, NewErrors("Parsing errors", p.Errors()...)
	}
	return program, nil
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.EOF:
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.NEWLINE},
			actualToken:    token.EOF,
		}
		p.errors = append(p.errors, err)

		return nil
	case token.NEWLINE:
		return nil
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	if p.currentTokenOneOf(token.NEWLINE, token.SEMICOLON) {
		p.nextToken()
		return stmt
	}

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if !p.acceptOneOf(token.NEWLINE, token.SEMICOLON) {
		return nil
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekTokenOneOf(token.SEMICOLON, token.NEWLINE) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Errorf("no prefix parse function for type %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseExpression(precedence int, expressionTerminators ...token.TokenType) ast.Expression {
	terminators := defaultExpressionTerminators
	if expressionTerminators != nil {
		terminators = append(terminators, expressionTerminators...)
	}
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()
	for !p.peekTokenOneOf(terminators...) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseVariableAssignExpression(variable ast.Expression) ast.Expression {
	ident, ok := variable.(*ast.Identifier)
	if !ok {
		msg := fmt.Errorf("could not parse variable assignment: expected identifier, got token '%T'", variable)
		p.errors = append(p.errors, msg)
		return nil
	}
	variableExp := &ast.Variable{
		Name: ident,
	}
	p.nextToken()
	variableExp.Value = p.parseExpression(LOWEST)
	return variableExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Errorf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseSymbolLiteral() ast.Expression {
	return &ast.SymbolLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}

	array.Elements = p.parseExpressionList(token.RBRACKET)

	return array
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.currentTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.accept(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.THEN) {
		p.accept(token.THEN)
	}
	if !p.peekTokenOneOf(token.NEWLINE, token.SEMICOLON) {
		msg := fmt.Sprintf(
			"could not parse if expression: unexpected token %s: '%s'",
			p.peekToken.Type,
			p.peekToken.Literal,
		)
		err := errors.Wrap(
			&unexpectedTokenError{
				expectedTokens: []token.TokenType{token.NEWLINE, token.SEMICOLON},
				actualToken:    p.peekToken.Type,
			},
			msg,
		)
		p.errors = append(p.errors, err)
		return nil
	}
	p.acceptOneOf(token.NEWLINE, token.SEMICOLON)
	consequence := p.parseBlockStatement(token.ELSE)
	expression.Consequence = consequence
	if p.peekTokenIs(token.ELSE) {
		p.accept(token.ELSE)
		p.accept(token.NEWLINE)
		expression.Alternative = p.parseBlockStatement()
	}
	p.accept(token.END)
	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.accept(token.IDENT) {
		return nil
	}
	lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	lit.Parameters = p.parseFunctionParameters()

	if !p.acceptOneOf(token.NEWLINE, token.SEMICOLON) {
		return nil
	}
	lit.Body = p.parseBlockStatement()
	if !p.accept(token.END) {
		return nil
	}
	return lit
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	if p.peekTokenIs(token.LPAREN) {
		p.accept(token.LPAREN)
	}

	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.accept(token.RPAREN)
		return identifiers
	}

	if p.peekTokenOneOf(token.NEWLINE, token.SEMICOLON) {
		return identifiers
	}

	p.accept(token.IDENT)

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.accept(token.COMMA)
		p.accept(token.IDENT)
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if p.peekTokenIs(token.RPAREN) {
		p.accept(token.RPAREN)
	}

	return identifiers
}

func (p *Parser) parseBlockStatement(t ...token.TokenType) *ast.BlockStatement {
	terminatorTokens := append(
		[]token.TokenType{
			token.END,
			token.EOF,
		},
		t...,
	)
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	for !p.peekTokenOneOf(terminatorTokens...) {
		p.nextToken()
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	return block
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	ident, ok := function.(*ast.Identifier)
	if !ok {
		msg := fmt.Errorf("could not parse call expression: expected identifier, got token '%T'", function)
		p.errors = append(p.errors, msg)
		return nil
	}
	exp := &ast.CallExpression{Token: ident.Token, Function: ident}
	args := []ast.Expression{}
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.consume(token.COMMA)
		args = append(args, p.parseExpression(LOWEST))
	}
	exp.Arguments = args
	return exp
}

func (p *Parser) parseContextCallExpression(context ast.Expression) ast.Expression {
	contextCallExpression := &ast.ContextCallExpression{Token: p.curToken, Context: context}
	p.nextToken()
	expr := p.parseExpression(LOWEST)
	var callExp *ast.CallExpression
	switch expr := expr.(type) {
	case *ast.Identifier:
		callExp = &ast.CallExpression{Token: expr.Token, Function: expr, Arguments: []ast.Expression{}}
	case *ast.CallExpression:
		callExp = expr
	default:
		msg := fmt.Errorf("could not parse context call expression: expected Identifier or CallExpression, got token '%T'", expr)
		p.errors = append(p.errors, msg)
		return nil
	}
	contextCallExpression.Call = callExp
	return contextCallExpression
}

func (p *Parser) parseCallExpressionWithParens(function ast.Expression) ast.Expression {
	ident, ok := function.(*ast.Identifier)
	if !ok {
		msg := fmt.Errorf("could not parse call expression: expected identifier, got token '%T'", function)
		p.errors = append(p.errors, msg)
		return nil
	}
	exp := &ast.CallExpression{Token: p.curToken, Function: ident}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}
	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.consume(token.COMMA)
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.accept(end) {
		return nil
	}

	return list
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentTokenOneOf(types ...token.TokenType) bool {
	for _, typ := range types {
		if p.curToken.Type == typ {
			return true
		}
	}
	return false
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenOneOf(types ...token.TokenType) bool {
	for _, typ := range types {
		if p.peekToken.Type == typ {
			return true
		}
	}
	return false
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// accept moves to the next Token
// if it's from the valid set.
func (p *Parser) accept(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

// acceptOneOf moves to the next Token
// if it's from the valid set.
func (p *Parser) acceptOneOf(t ...token.TokenType) bool {
	if p.peekTokenOneOf(t...) {
		p.nextToken()
		return true
	} else {
		p.peekError(t...)
		return false
	}
}

// consume consumes the next token
// if it's from the valid set.
func (p *Parser) consume(t token.TokenType) bool {
	isRightToken := p.accept(t)
	if isRightToken {
		p.nextToken()
	}
	return isRightToken
}
