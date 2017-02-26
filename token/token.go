package token

//go:generate stringer -type=TokenType

const (
	ILLEGAL TokenType = iota
	EOF

	NEWLINE // \n
	// Identifier + literals
	IDENT
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

	LT     // <
	GT     // >
	EQ     // ==
	NOT_EQ // !=

	// Delimiters
	COMMA
	SEMICOLON

	COLON  // :
	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }

	// Keywords
	DEF
	END
	IF
	THEN
	ELSE
	TRUE
	FALSE
	RETURN
)

var keywords = map[string]TokenType{
	"def":    DEF,
	"end":    END,
	"if":     IF,
	"then":   THEN,
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
