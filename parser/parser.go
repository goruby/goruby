package parser

import (
	"github.com/goruby/goruby/ast"
	"github.com/goruby/goruby/lexer"
	"github.com/goruby/goruby/token"
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.IDENT:
		if p.peekToken.Type == token.ASSIGN {
			return p.parseVariableStatement()
		} else {
			return nil
		}
	default:
		return nil
	}
}

func (p *Parser) parseVariableStatement() *ast.VariableStatement {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if !p.accept(token.ASSIGN) {
		return nil
	}
	// TODO: We're skipping the expressions until we
	// encounter a newline
	for !p.currentTokenIs(token.NEWLINE) {
		p.nextToken()
	}
	return &ast.VariableStatement{Name: ident}
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// accept consumes the next Token
// if it's from the valid set.
func (p *Parser) accept(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		return false
	}
}
