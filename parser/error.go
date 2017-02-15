package parser

import (
	"fmt"

	"github.com/goruby/goruby/token"
)

// Make sure unexpectedTokenError implements error interface
var _ error = &unexpectedTokenError{}

func IsEOFError(err error) bool {
	tokenErr, ok := err.(*unexpectedTokenError)
	if !ok {
		return false
	}
	if tokenErr.actualToken != token.EOF {
		return false
	}

	return true
}

func IsEOFInsteadOfNewlineError(err error) bool {
	tokenErr, ok := err.(*unexpectedTokenError)
	if !ok {
		return false
	}
	if tokenErr.actualToken != token.EOF {
		return false
	}

	for _, expectedToken := range tokenErr.expectedTokens {
		if expectedToken == token.NEWLINE {
			return true
		}
	}

	return false
}

type unexpectedTokenError struct {
	expectedTokens []token.TokenType
	actualToken    token.TokenType
}

func (u *unexpectedTokenError) Error() string {
	return fmt.Sprintf(
		"expected next token to be of type %s, got %s instead",
		u.expectedTokens,
		u.actualToken,
	)
}
