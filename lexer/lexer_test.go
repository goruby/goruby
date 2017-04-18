package lexer

import (
	"testing"

	"github.com/goruby/goruby/token"
)

func TestLexerNextToken(t *testing.T) {
	input := `five = 5
# just comment
fifty = 5_0
ten = 10

def add(x, y)
	x + y
end

result = add(five, ten)
!-/*5;
5 < 10 > 5
return
if 5 < 10 then
	true
else
	false
end

10 == 10
10 != 9
""
"foobar"
'foobar'
"foo bar"
'foo bar'
:sym
.

def nil?
end

def run!
end
[1, 2]
nil
self
Ten = 10
module Abc
end
class Abc
end
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "fifty"},
		{token.ASSIGN, "="},
		{token.INT, "5_0"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.DEF, "def"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.NEWLINE, "\n"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.NEWLINE, "\n"},
		{token.RETURN, "return"},
		{token.NEWLINE, "\n"},
		{token.IF, "if"},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.THEN, "then"},
		{token.NEWLINE, "\n"},
		{token.TRUE, "true"},
		{token.NEWLINE, "\n"},
		{token.ELSE, "else"},
		{token.NEWLINE, "\n"},
		{token.FALSE, "false"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.NOTEQ, "!="},
		{token.INT, "9"},
		{token.NEWLINE, "\n"},
		{token.STRING, ""},
		{token.NEWLINE, "\n"},
		{token.STRING, "foobar"},
		{token.NEWLINE, "\n"},
		{token.STRING, "foobar"},
		{token.NEWLINE, "\n"},
		{token.STRING, "foo bar"},
		{token.NEWLINE, "\n"},
		{token.STRING, "foo bar"},
		{token.NEWLINE, "\n"},
		{token.SYMBOL, "sym"},
		{token.NEWLINE, "\n"},
		{token.DOT, "."},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.DEF, "def"},
		{token.IDENT, "nil?"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.NEWLINE, "\n"},
		{token.DEF, "def"},
		{token.IDENT, "run!"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.NEWLINE, "\n"},
		{token.NIL, "nil"},
		{token.NEWLINE, "\n"},
		{token.SELF, "self"},
		{token.NEWLINE, "\n"},
		{token.CONST, "Ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.NEWLINE, "\n"},
		{token.MODULE, "module"},
		{token.CONST, "Abc"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.CLASS, "class"},
		{token.CONST, "Abc"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
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
