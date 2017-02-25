package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/goruby/goruby/token"
)

const (
	eof = -1
)

type LexerStateFn func(*Lexer) LexerStateFn

func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		state:  startLexer,
		tokens: make(chan token.Token, 2), // Two token sufficient.
	}
	return l
}

type Lexer struct {
	input  string           // the string being scanned.
	state  LexerStateFn     // the next lexing function to enter
	pos    int              // current position in the input.
	start  int              // start position of this item.
	width  int              // width of last rune read from input.
	tokens chan token.Token // channel of scanned tokens.
}

func (l *Lexer) NextToken() token.Token {
	for {
		select {
		case item, ok := <-l.tokens:
			if ok {
				return item
			} else {
				panic(fmt.Errorf("No items left"))
			}
		default:
			l.state = l.state(l)
			if l.state == nil {
				close(l.tokens)
			}
		}
	}
	panic("not reached")
}

// HasNext returns true if there are tokens left, false if EOF has reached
func (l *Lexer) HasNext() bool {
	return l.state != nil
}

// emit passes a token back to the client.
func (l *Lexer) emit(t token.TokenType) {
	l.tokens <- token.NewToken(t, l.input[l.start:l.pos], l.start)
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	var r rune
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// error returns an error token and terminates the scan by passing
// back a nil pointer that will be the next state, terminating l.run.
func (l *Lexer) errorf(format string, args ...interface{}) LexerStateFn {
	l.tokens <- token.NewToken(token.ILLEGAL, fmt.Sprintf(format, args...), l.start)
	return nil
}

func startLexer(l *Lexer) LexerStateFn {
	r := l.next()
	if isWhitespace(r) {
		l.ignore()
		return startLexer
	}
	switch r {
	case '\n':
		l.emit(token.NEWLINE)
		return startLexer
	case '"':
		return lexString
	case '=':
		if l.peek() == '=' {
			l.next()
			l.emit(token.EQ)
		} else {
			l.emit(token.ASSIGN)
		}
		return startLexer
	case '+':
		l.emit(token.PLUS)
		return startLexer
	case '-':
		l.emit(token.MINUS)
		return startLexer
	case '!':
		if l.peek() == '=' {
			l.next()
			l.emit(token.NOT_EQ)
		} else {
			l.emit(token.BANG)
		}
		return startLexer
	case '/':
		l.emit(token.SLASH)
		return startLexer
	case '*':
		l.emit(token.ASTERISK)
		return startLexer
	case '<':
		l.emit(token.LT)
		return startLexer
	case '>':
		l.emit(token.GT)
		return startLexer
	case '(':
		l.emit(token.LPAREN)
		return startLexer
	case ')':
		l.emit(token.RPAREN)
		return startLexer
	case '{':
		l.emit(token.LBRACE)
		return startLexer
	case '}':
		l.emit(token.RBRACE)
		return startLexer
	case ',':
		l.emit(token.COMMA)
		return startLexer
	case ';':
		l.emit(token.SEMICOLON)
		return startLexer
	case eof:
		l.emit(token.EOF)
		return startLexer
	default:
		if isLetter(r) {
			return lexIdentifier
		} else if isDigit(r) {
			return lexDigit
		} else {
			return l.errorf("Illegal character at %d: '%c'", l.start, r)
		}
	}
}

func lexIdentifier(l *Lexer) LexerStateFn {
	r := l.next()
	for isLetter(r) {
		r = l.next()
	}
	l.backup()
	literal := l.input[l.start:l.pos]
	l.emit(token.LookupIdent(literal))
	return startLexer
}

func lexDigit(l *Lexer) LexerStateFn {
	r := l.next()
	for isDigit(r) {
		r = l.next()
	}
	l.backup()
	l.emit(token.INT)
	return startLexer
}

func lexString(l *Lexer) LexerStateFn {
	l.ignore()
	r := l.next()

	for r != '"' {
		r = l.next()
	}
	l.backup()
	l.emit(token.STRING)
	l.next()
	l.ignore()
	return startLexer
}

func isWhitespace(r rune) bool {
	return unicode.IsSpace(r) && r != '\n'
}

func isLetter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}
