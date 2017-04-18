package lexer

import (
	"bytes"
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/goruby/goruby/token"
)

const (
	eof = -1
)

// LexStartFn represents the entrypoint the Lexer uses to start processing the
// input.
var LexStartFn = startLexer

// StateFn represents a function which is capable of lexing parts of the
// input. It returns another StateFn to proceed with.
//
// Typically a state function would get called from LexStartFn and should
// return LexStartFn to go back to the decision loop. It also could return
// another non start state function if the partial input to parse is abiguous.
type StateFn func(*Lexer) StateFn

// New returns a Lexer instance ready to process the given input.
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		state:  startLexer,
		tokens: make(chan token.Token, 2), // Two token sufficient.
	}
	return l
}

// Lexer is the engine to process input and emit Tokens
type Lexer struct {
	input  string           // the string being scanned.
	state  StateFn          // the next lexing function to enter
	pos    int              // current position in the input.
	start  int              // start position of this item.
	width  int              // width of last rune read from input.
	tokens chan token.Token // channel of scanned tokens.
}

// NextToken will return the next token processed from the lexer.
//
// Callers should make sure to call Lexer.HasNext before calling this method
// as it will panic if it is called after token.EOF is returned.
func (l *Lexer) NextToken() token.Token {
	for {
		select {
		case item, ok := <-l.tokens:
			if ok {
				return item
			}
			panic(fmt.Errorf("No items left"))
		default:
			l.state = l.state(l)
			if l.state == nil {
				close(l.tokens)
			}
		}
	}
}

// HasNext returns true if there are tokens left, false if EOF has reached
func (l *Lexer) HasNext() bool {
	return l.state != nil
}

// emit passes a token back to the client.
func (l *Lexer) emit(t token.Type) {
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
func (l *Lexer) errorf(format string, args ...interface{}) StateFn {
	l.tokens <- token.NewToken(token.ILLEGAL, fmt.Sprintf(format, args...), l.start)
	return nil
}

func startLexer(l *Lexer) StateFn {
	r := l.next()
	if isWhitespace(r) {
		l.ignore()
		return startLexer
	}
	switch r {
	case '\n':
		l.emit(token.NEWLINE)
		return startLexer
	case '\'':
		return lexSingleQuoteString
	case '"':
		return lexString
	case ':':
		return lexSymbol
	case '.':
		l.emit(token.DOT)
		return startLexer
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
			l.emit(token.NOTEQ)
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
	case '[':
		l.emit(token.LBRACKET)
		return startLexer
	case ']':
		l.emit(token.RBRACKET)
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

func lexIdentifier(l *Lexer) StateFn {
	legalIdentifierCharacters := []byte{'?', '!'}
	r := l.next()
	for isLetter(r) || isDigit(r) || bytes.ContainsRune(legalIdentifierCharacters, r) {
		r = l.next()
	}
	l.backup()
	literal := l.input[l.start:l.pos]
	l.emit(token.LookupIdent(literal))
	return startLexer
}

func lexDigit(l *Lexer) StateFn {
	r := l.next()
	for isDigitOrUnderscore(r) {
		r = l.next()
	}
	l.backup()
	l.emit(token.INT)
	return startLexer
}

func lexSingleQuoteString(l *Lexer) StateFn {
	l.ignore()
	r := l.next()

	for r != '\'' {
		r = l.next()
	}
	l.backup()
	l.emit(token.STRING)
	l.next()
	l.ignore()
	return startLexer
}

func lexString(l *Lexer) StateFn {
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

func lexSymbol(l *Lexer) StateFn {
	l.ignore()
	r := l.next()

	for isLetter(r) || isDigit(r) {
		r = l.next()
	}
	l.backup()
	l.emit(token.SYMBOL)
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

func isDigitOrUnderscore(r rune) bool {
	return isDigit(r) || r == '_'
}
