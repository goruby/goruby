package lexer

import (
	"testing"

	"github.com/goruby/goruby/token"
)

func TestLexerNextToken(t *testing.T) {
	input := `five = 5
while x < y do
	x += x
end
seven,
# just comment
fifty = 5_0
ten = 10
?-
?\n
? foo : bar

def add(x, y)
	x + y
end
|
||

result = add(five, ten)
!-/*%5;
+= -= *= /= %=
5 < 10 > 5
return
if 5 < 10 then
	true
else
	false
end

begin
rescue
end

10 == 10
10 != 9
10 <= 9
10 >= 9
10 <=> 9
10 << 9
""
"foobar"
'foobar'
"foo bar"
'foo bar'
:sym
:"sym"
:'sym'
.
&foo
&
&&
:dotAfter.

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
add { |x| x }
add do |x|
end
yield
while
A::B
=>
__FILE__
@
$foo,
$foo;
$Foo
$dotAfter.
$@
$a`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.NEWLINE, "\n"},
		{token.WHILE, "while"},
		{token.IDENT, "x"},
		{token.LT, "<"},
		{token.IDENT, "y"},
		{token.DO, "do"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "x"},
		{token.ADDASSIGN, "+="},
		{token.IDENT, "x"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "seven"},
		{token.COMMA, ","},
		{token.NEWLINE, "\n"},
		{token.HASH, "#"},
		{token.STRING, " just comment"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "fifty"},
		{token.ASSIGN, "="},
		{token.INT, "5_0"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.NEWLINE, "\n"},
		{token.STRING, "-"},
		{token.NEWLINE, "\n"},
		{token.STRING, "\\n"},
		{token.NEWLINE, "\n"},
		{token.QMARK, "?"},
		{token.IDENT, "foo"},
		{token.COLON, ":"},
		{token.IDENT, "bar"},
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
		{token.PIPE, "|"},
		{token.NEWLINE, "\n"},
		{token.LOGICALOR, "||"},
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
		{token.MODULO, "%"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.ADDASSIGN, "+="},
		{token.SUBASSIGN, "-="},
		{token.MULASSIGN, "*="},
		{token.DIVASSIGN, "/="},
		{token.MODASSIGN, "%="},
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
		{token.BEGIN, "begin"},
		{token.NEWLINE, "\n"},
		{token.RESCUE, "rescue"},
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
		{token.INT, "10"},
		{token.LTE, "<="},
		{token.INT, "9"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.GTE, ">="},
		{token.INT, "9"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.SPACESHIP, "<=>"},
		{token.INT, "9"},
		{token.NEWLINE, "\n"},
		{token.INT, "10"},
		{token.LSHIFT, "<<"},
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
		{token.SYMBEG, ":"},
		{token.IDENT, "sym"},
		{token.NEWLINE, "\n"},
		{token.SYMBEG, ":"},
		{token.STRING, "sym"},
		{token.NEWLINE, "\n"},
		{token.SYMBEG, ":"},
		{token.STRING, "sym"},
		{token.NEWLINE, "\n"},
		{token.DOT, "."},
		{token.NEWLINE, "\n"},
		{token.CAPTURE, "&"},
		{token.IDENT, "foo"},
		{token.NEWLINE, "\n"},
		{token.AND, "&"},
		{token.NEWLINE, "\n"},
		{token.LOGICALAND, "&&"},
		{token.NEWLINE, "\n"},
		{token.SYMBEG, ":"},
		{token.IDENT, "dotAfter"},
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
		{token.IDENT, "add"},
		{token.LBRACE, "{"},
		{token.PIPE, "|"},
		{token.IDENT, "x"},
		{token.PIPE, "|"},
		{token.IDENT, "x"},
		{token.RBRACE, "}"},
		{token.NEWLINE, "\n"},
		{token.IDENT, "add"},
		{token.DO, "do"},
		{token.PIPE, "|"},
		{token.IDENT, "x"},
		{token.PIPE, "|"},
		{token.NEWLINE, "\n"},
		{token.END, "end"},
		{token.NEWLINE, "\n"},
		{token.YIELD, "yield"},
		{token.NEWLINE, "\n"},
		{token.WHILE, "while"},
		{token.NEWLINE, "\n"},
		{token.CONST, "A"},
		{token.SCOPE, "::"},
		{token.CONST, "B"},
		{token.NEWLINE, "\n"},
		{token.HASHROCKET, "=>"},
		{token.NEWLINE, "\n"},
		{token.KEYWORD__FILE__, "__FILE__"},
		{token.NEWLINE, "\n"},
		{token.AT, "@"},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$foo"},
		{token.COMMA, ","},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$foo"},
		{token.SEMICOLON, ";"},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$Foo"},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$dotAfter"},
		{token.DOT, "."},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$@"},
		{token.NEWLINE, "\n"},
		{token.GLOBAL, "$a"},
		{token.EOF, ""},
	}

	lexer := New(input)

	for pos, testCase := range tests {
		if !lexer.HasNext() {
			t.Logf("Unexpected EOF at %d\n", lexer.pos)
			t.FailNow()
		}
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
