package parser

import (
	"fmt"
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

	fmt.Printf("Error: %+#v\n", err)

	if !isEOFErr {
		t.Logf("Expected an EOF error.")
		t.Fail()
	}
}
