package parser

import (
	"testing"

	"github.com/goruby/goruby/lexer"
)

func TestIsEOF(t *testing.T) {
	input := "def foo"
	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	errors := p.Errors()

	if len(errors) != 1 {
		t.Logf("Expected 1 error, got %d\n", len(errors))
		t.FailNow()
	}

	err := errors[0]

	isEOFErr := IsEOFError(err)

	if !isEOFErr {
		t.Logf("Expected an EOF error, got %T:%q\n", err, err)
		t.Fail()
	}
}

func TestIsEOFInsteadOfNewline(t *testing.T) {
	input := "def foo"
	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	errors := p.Errors()

	if len(errors) != 1 {
		t.Logf("Expected 1 error, got %d\n", len(errors))
		t.FailNow()
	}

	err := errors[0]

	needsMore := IsEOFInsteadOfNewlineError(err)

	if !needsMore {
		t.Logf("Expected an EOF error with an expected NEWLINE,\n\tgot %T:%q", err, err)
		t.Fail()
	}
}
