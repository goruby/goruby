package token

import (
	"bytes"
	"strconv"
	"unicode"
)

//go:generate stringer -type=Type

// Recognized token types
const (
	ILLEGAL Type = iota // An illegal/unknown character
	EOF                 // end of input

	// Identifier + literals
	literal_beg
	IDENT
	CONST
	GLOBAL
	INT
	STRING
	literal_end

	// Operators
	operator_beg
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	BANG     // !
	ASTERISK // *
	SLASH    // /
	MODULO   // %

	LT        // <
	LTE       // <=
	GT        // >
	GTE       // >=
	EQ        // ==
	NOTEQ     // !=
	SPACESHIP // <=>
	operator_end

	HASHROCKET // =>

	// Delimiters

	NEWLINE // \n
	COMMA
	SEMICOLON
	HASH // #

	DOT      // .
	COLON    // :
	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]
	PIPE     // |

	SCOPE // ::
	AT    // @

	// Keywords
	keyword_beg
	DEF
	SELF
	END
	IF
	THEN
	ELSE
	UNLESS
	TRUE
	FALSE
	RETURN
	NIL
	MODULE
	CLASS
	DO
	YIELD
	BEGIN
	RESCUE
	KEYWORD__FILE__
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	CONST:  "CONST",
	GLOBAL: "GLOBAL",
	INT:    "INT",
	STRING: "STRING",

	ASSIGN:   "=",
	PLUS:     "+",
	MINUS:    "-",
	BANG:     "!",
	ASTERISK: "*",
	SLASH:    "/",
	MODULO:   "%",

	LT:        "<",
	LTE:       "<=",
	GT:        ">",
	GTE:       ">=",
	EQ:        "==",
	NOTEQ:     "!=",
	SPACESHIP: "<=>",

	NEWLINE:   "NEWLINE",
	COMMA:     ",",
	SEMICOLON: ";",
	HASH:      "#",

	DOT:      ".",
	COLON:    ":",
	LPAREN:   "(",
	RPAREN:   ")",
	LBRACE:   "{",
	RBRACE:   "}",
	LBRACKET: "[",
	RBRACKET: "]",
	PIPE:     "|",

	SCOPE:      "::",
	HASHROCKET: "=>",
	AT:         "@",

	DEF:             "def",
	SELF:            "self",
	END:             "end",
	UNLESS:          "unless",
	IF:              "if",
	THEN:            "then",
	ELSE:            "else",
	TRUE:            "true",
	FALSE:           "false",
	RETURN:          "return",
	NIL:             "nil",
	MODULE:          "module",
	CLASS:           "class",
	DO:              "do",
	YIELD:           "yield",
	BEGIN:           "begin",
	RESCUE:          "rescue",
	KEYWORD__FILE__: "__FILE__",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//
func (tok Type) String() string {
	s := ""
	if 0 <= tok && tok < Type(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Type

func init() {
	keywords = make(map[string]Type)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
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

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
//
func (tok Type) IsLiteral() bool { return literal_beg < tok && tok < literal_end }

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
//
func (tok Type) IsOperator() bool { return operator_beg < tok && tok < operator_end }

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
//
func (tok Type) IsKeyword() bool { return keyword_beg < tok && tok < keyword_end }
