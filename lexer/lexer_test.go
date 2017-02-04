package lexer

import (
	"testing"

	"github.com/goruby/goruby/token"
)

func TestLexerNextToken(t *testing.T) {
	input := `five = 5
ten = 10

def add(x, y)
	x + y
end

result = add(five, ten)
!-/*5;
5 < 10 > 5
return
if 5 < 10
	true
else
	false
end

10 == 10
10 != 9
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.DEF, "def"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.END, "end"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.RETURN, "return"},
		{token.IF, "if"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.TRUE, "true"},
		{token.ELSE, "else"},
		{token.FALSE, "false"},
		{token.END, "end"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for pos, testCase := range tests {
		token := lexer.NextToken()

		if token.Type != testCase.expectedType {
			t.Logf("Expected token with type %q at position %d, got type %q\n", testCase.expectedType, pos, token.Type)
			t.Fail()
		}

		if token.Literal != testCase.expectedLiteral {
			t.Logf("Expected token with literal %q at position %d, got literal %q\n", testCase.expectedLiteral, pos, token.Literal)
			t.Fail()
		}
	}
}
