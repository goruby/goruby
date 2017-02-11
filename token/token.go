package token

//go:generate stringer -type=TokenType

const (
	ILLEGAL TokenType = iota
	EOF

	NEWLINE // \n
	// Identifier + literals
	IDENT
	INT

	// Operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	BANG     // !
	ASTERISK // *
	SLASH    // /

	LT     // <
	GT     // >
	EQ     // ==
	NOT_EQ // !=

	// Delimiters
	COMMA
	SEMICOLON

	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }

	// Keywords
	DEF
	END
	IF
	ELSE
	TRUE
	FALSE
	RETURN
)

var keywords = map[string]TokenType{
	"def":    DEF,
	"end":    END,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type TokenType int

func NewToken(typ TokenType, literal string, pos int) Token {
	return Token{typ, literal, pos}
}

type Token struct {
	Type    TokenType
	Literal string
	Pos     int
}
