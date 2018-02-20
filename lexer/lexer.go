package lexer

import (
	"bytes"
	"fmt"
	"strings"
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

const operatorCharacters = "+-!*/%&<>=,;#.:(){}[]|@?$"

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
	input     string           // the string being scanned.
	state     StateFn          // the next lexing function to enter
	pos       int              // current position in the input.
	start     int              // start position of this item.
	width     int              // width of last rune read from input.
	tokens    chan token.Token // channel of scanned tokens.
	lastToken token.Token      // lastToken stores the last token emitted by the lexer
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
	token := token.NewToken(t, l.input[l.start:l.pos], l.start)
	l.lastToken = token
	l.tokens <- token
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
	case '$':
		return lexGlobal
	case '\n':
		l.emit(token.NEWLINE)
		return startLexer
	case '\'':
		return lexSingleQuoteString
	case '"':
		return lexString
	case ':':
		p := l.peek()
		if p == ':' {
			l.next()
			l.emit(token.SCOPE)
			return startLexer
		}
		if isWhitespace(p) {
			l.emit(token.COLON)
			return startLexer
		}
		l.emit(token.SYMBEG)
		return startLexer
	case '.':
		l.emit(token.DOT)
		return startLexer
	case '=':
		if l.peek() == '=' {
			l.next()
			l.emit(token.EQ)
		} else if l.peek() == '>' {
			l.next()
			l.emit(token.HASHROCKET)
		} else {
			l.emit(token.ASSIGN)
		}
		return startLexer
	case '+':
		if l.peek() == '=' {
			l.next()
			l.emit(token.ADDASSIGN)
			return startLexer
		}
		l.emit(token.PLUS)
		return startLexer
	case '-':
		if l.peek() == '=' {
			l.next()
			l.emit(token.SUBASSIGN)
			return startLexer
		}
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
	case '?':
		p := l.peek()
		if isWhitespace(p) {
			l.emit(token.QMARK)
			return startLexer
		}
		if isExpressionDelimiter(p) {
			fmt.Printf("warning: invalid character syntax; use ?%c\n", r)
			l.ignore()
			return l.errorf("unexpected '?'")
		}
		return lexCharacterLiteral
	case '/':
		if l.peek() == '=' {
			l.next()
			l.emit(token.DIVASSIGN)
			return startLexer
		}
		l.emit(token.SLASH)
		return startLexer
	case '*':
		if l.peek() == '=' {
			l.next()
			l.emit(token.MULASSIGN)
			return startLexer
		}
		l.emit(token.ASTERISK)
		return startLexer
	case '%':
		if l.peek() == '=' {
			l.next()
			l.emit(token.MODASSIGN)
			return startLexer
		}
		l.emit(token.MODULO)
		return startLexer
	case '&':
		if p := l.peek(); p == '&' {
			l.next()
			l.emit(token.LOGICALAND)
			return startLexer
		}
		if p := l.peek(); isLetter(p) {
			l.emit(token.CAPTURE)
			return startLexer
		}
		l.emit(token.AND)
		return startLexer
	case '<':
		if l.peek() == '=' {
			l.next()
			if l.peek() == '>' {
				l.next()
				l.emit(token.SPACESHIP)
				return startLexer
			}
			l.emit(token.LTE)
			return startLexer
		}
		if l.peek() == '<' {
			l.next()
			l.emit(token.LSHIFT)
			return startLexer
		}
		l.emit(token.LT)
		return startLexer
	case '>':
		if l.peek() == '=' {
			l.next()
			l.emit(token.GTE)
			return startLexer
		}
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
	case '#':
		return commentLexer
	case '|':
		if l.lastToken.Type == token.DO || l.lastToken.Type == token.LBRACE {
			l.emit(token.PIPE)
			return startLexer
		}
		if p := l.peek(); p == '|' {
			l.next()
			l.emit(token.LOGICALOR)
			return startLexer
		}
		l.emit(token.PIPE)
		return startLexer
	case '@':
		l.emit(token.AT)
		return startLexer

	default:
		if isDigit(r) {
			return lexDigit
		} else if isLetter(r) {
			return lexIdentifier
		} else {
			return l.errorf("Illegal character: '%c'", r)
		}
	}
}

func lexIdentifier(l *Lexer) StateFn {
	legalIdentifierCharacters := []byte{'?', '!'}
	r := l.next()
	for {
		if unicode.IsSpace(r) || strings.ContainsRune(operatorCharacters, r) || r == eof {
			if bytes.ContainsRune(legalIdentifierCharacters, r) {
				l.next()
				break
			}
			break
		}
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

func lexCharacterLiteral(l *Lexer) StateFn {
	l.ignore()
	r := l.next()
	if isWhitespace(r) && r != '\t' && r != '\v' && r != '\f' && r != '\r' {
		return l.errorf("invalid character syntax; use ?\\s")
	}
	if r == '\\' {
		r = l.next()
	}
	if p := l.peek(); !isWhitespace(p) && !isExpressionDelimiter(p) {
		return l.errorf("unexpected '?'")
	}
	l.emit(token.STRING)
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

func lexGlobal(l *Lexer) StateFn {
	r := l.next()

	if r == '.' {
		return l.errorf("Illegal character: '%c'", r)
	}

	if isExpressionDelimiter(r) {
		return l.errorf("Illegal character: '%c'", r)
	}

	if isWhitespace(r) {
		return l.errorf("Illegal character: '%c'", r)
	}

	for !isWhitespace(r) && !isExpressionDelimiter(r) && r != '.' && r != ',' {
		r = l.next()
	}
	l.backup()
	l.emit(token.GLOBAL)
	return startLexer
}

func commentLexer(l *Lexer) StateFn {
	l.emit(token.HASH)
	r := l.next()

	for r != '\n' && r != eof {
		r = l.next()
	}
	l.backup()
	l.emit(token.STRING)
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

func isExpressionDelimiter(r rune) bool {
	return r == '\n' || r == ';' || r == eof
}
