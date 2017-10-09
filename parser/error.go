package parser

import (
	"bytes"
	"fmt"
	gotoken "go/token"

	"github.com/goruby/goruby/token"
	"github.com/pkg/errors"
)

// Make sure unexpectedTokenError implements error interface
var _ error = &unexpectedTokenError{}

// Make sure Errors implements error interface
var _ error = &Errors{}

// NewErrors returns a composite Error object wrapping multiple errors into
// one.
func NewErrors(context string, errors ...error) *Errors {
	return &Errors{context, errors}
}

// Errors represents a group of errors and its context
//
// Errors implements the error interface to be used as an error in the code.
type Errors struct {
	context string
	errors  []error
}

// Error returns all error messages divided by newlines and prepended with the
// error context.
func (e *Errors) Error() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%s:\n", e.context)
	for _, err := range e.errors {
		fmt.Fprintf(&buf, "\t%s\n", err.Error())
	}
	return buf.String()
}

// IsEOFError returns true if err represents an unexpectedTokenError with its
// actual token type set to token.EOF.
//
// It returns false for any other error.
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

// IsEOFInsteadOfNewlineError returns true if err represents an unexpectedTokenError with its
// actual token type set to token.EOF and if its expected token types includes
// token.NEWLINE.
//
// It returns false for any other error.
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

type tokens []token.Type

func (t tokens) String() string {
	if len(t) == 1 {
		return t[0].String()
	}
	var s []token.Type = t
	return fmt.Sprintf("%s", s)
}

type unexpectedTokenError struct {
	Pos            gotoken.Position
	expectedTokens []token.Type
	actualToken    token.Type
}

func (e *unexpectedTokenError) Error() string {
	msg := fmt.Sprintf(
		"unexpected %s, expecting %s",
		e.actualToken,
		tokens(e.expectedTokens),
	)
	if e.Pos.Filename != "" || e.Pos.IsValid() {
		return e.Pos.String() + ": " + msg
	}
	return msg
}
