package token

import (
	"bytes"
	"unicode"
)

//go:generate stringer -type=Type

// Recognized token types
const (
	ILLEGAL Type = iota // An illegal/unknown character
	EOF                 // end of input

	// Identifier + literals

	IDENT
	CONST
	INT
	STRING
	SYMBOL // :symbol

	// Operators

	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	BANG     // !
	ASTERISK // *
	SLASH    // /

	LT    // <
	GT    // >
	EQ    // ==
	NOTEQ // !=

	// Delimiters

	NEWLINE // \n
	COMMA
	SEMICOLON

	DOT      // .
	COLON    // :
	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	// Keywords

	DEF
	SELF
	END
	IF
	THEN
	ELSE
	TRUE
	FALSE
	RETURN
	NIL
	MODULE
	CLASS
)

var keywords = map[string]Type{
	"def":    DEF,
	"end":    END,
	"if":     IF,
	"then":   THEN,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"nil":    NIL,
	"return": RETURN,
	"self":   SELF,
	"module": MODULE,
	"class":  CLASS,
}

// LookupIdent returns a keyword Type if ident is a keyword. If ident starts
// with an upper character it returns CONST. In any other case it returns IDENT
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	if unicode.IsUpper(bytes.Runes([]byte(ident))[0]) {
		return CONST
	}
	return IDENT
}

// A Type represents a type of a known token
type Type int

// NewToken returns a new Token associated with the given Type typ, the Literal
// literal and the Position pos
func NewToken(typ Type, literal string, pos int) Token {
	return Token{typ, literal, pos}
}

// A Token represents a known token with its literal representation
type Token struct {
	Type    Type
	Literal string
	Pos     int
}
