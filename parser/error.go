package parser

import (
	"bytes"
	"fmt"

	"github.com/goruby/goruby/token"
	"github.com/pkg/errors"
)

// Make sure unexpectedTokenError implements error interface
var _ error = &unexpectedTokenError{}

// Make sure Errors implements error interface
var _ error = &Errors{}

func NewErrors(context string, errors ...error) *Errors {
	return &Errors{context, errors}
}

type Errors struct {
	context string
	errors  []error
}

func (e *Errors) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s:\n", e.context)
	for _, err := range e.errors {
		fmt.Fprintf(&buf, "\t%s\n", err.Error())
	}
	return buf.String()
}

func IsEOFError(err error) bool {
	if errors, ok := err.(*Errors); ok {
		for _, e := range errors.errors {
			if IsEOFError(e) {
				return true
			}
		}
	}

	cause := errors.Cause(err)
	tokenErr, ok := cause.(*unexpectedTokenError)
	if !ok {
		return false
	}
	if tokenErr.actualToken != token.EOF {
		return false
	}

	return true
}

func IsEOFInsteadOfNewlineError(err error) bool {
	if !IsEOFError(err) {
		return false
	}

	if errors, ok := err.(*Errors); ok {
		for _, e := range errors.errors {
			if IsEOFInsteadOfNewlineError(e) {
				return true
			}
		}
	}

	tokenErr := errors.Cause(err).(*unexpectedTokenError)

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
