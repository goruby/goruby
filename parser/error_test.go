package parser

import (
	"testing"

	"github.com/goruby/goruby/token"
	"github.com/pkg/errors"
)

func TestIsEOF(t *testing.T) {
	t.Run("Errors with unexpected token EOF", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.IDENT},
			actualToken:    token.EOF,
		}

		isEOFErr := IsEOFError(NewErrors("", err))

		if !isEOFErr {
			t.Logf("Expected an EOF error, got %T:%q\n", err, err)
			t.Fail()
		}
	})
	t.Run("Errors with unexpected token not EOF", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.IDENT},
			actualToken:    token.NEWLINE,
		}

		isEOFErr := IsEOFError(NewErrors("", err))

		if isEOFErr {
			t.Logf("Expected no EOF error, got %T:%q\n", err, err)
			t.Fail()
		}
	})
	t.Run("unexpected token EOF", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.IDENT},
			actualToken:    token.EOF,
		}

		isEOFErr := IsEOFError(err)

		if !isEOFErr {
			t.Logf("Expected an EOF error, got %T:%q\n", err, err)
			t.Fail()
		}
	})
	t.Run("unexpected token not EOF", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.IDENT},
			actualToken:    token.NEWLINE,
		}

		isEOFErr := IsEOFError(err)

		if isEOFErr {
			t.Logf("Expected no EOF error, got %T:%q\n", err, err)
			t.Fail()
		}
	})
	t.Run("unexpected token EOF wrapped error", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.IDENT},
			actualToken:    token.EOF,
		}
		wrapped := errors.Wrap(err, "some error")

		isEOFErr := IsEOFError(wrapped)

		if !isEOFErr {
			t.Logf("Expected an EOF error, got %T:%q\n", wrapped, wrapped)
			t.Fail()
		}
	})
	t.Run("unexpected token not EOF wrapped error", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.IDENT},
			actualToken:    token.NEWLINE,
		}
		wrapped := errors.Wrap(err, "some error")

		isEOFErr := IsEOFError(wrapped)

		if isEOFErr {
			t.Logf("Expected no EOF error, got %T:%q\n", wrapped, wrapped)
			t.Fail()
		}
	})
	t.Run("random error", func(t *testing.T) {
		err := errors.New("some error")

		isEOFErr := IsEOFError(err)

		if isEOFErr {
			t.Logf("Expected no EOF error, got %T:%q\n", err, err)
			t.Fail()
		}
	})
}

func TestIsEOFInsteadOfNewline(t *testing.T) {
	t.Run("Errors with unexpected token EOF, expected token NEWLINE", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.NEWLINE},
			actualToken:    token.EOF,
		}

		isEOFNewlineErr := IsEOFInsteadOfNewlineError(NewErrors("", err))

		if !isEOFNewlineErr {
			t.Logf("Expected an EOF NEWLINE error, got %T:%q\n", err, err)
			t.Fail()
		}
	})
	t.Run("unexpected token EOF, expected token NEWLINE", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.NEWLINE},
			actualToken:    token.EOF,
		}

		isEOFNewlineErr := IsEOFInsteadOfNewlineError(err)

		if !isEOFNewlineErr {
			t.Logf("Expected an EOF NEWLINE error, got %T:%q\n", err, err)
			t.Fail()
		}
	})
	t.Run("Errors with unexpected token not EOF", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.IDENT},
			actualToken:    token.NEWLINE,
		}

		isEOFNewlineErr := IsEOFInsteadOfNewlineError(NewErrors("", err))

		if isEOFNewlineErr {
			t.Logf("Expected no EOF NEWLINE error, got %T:%q\n", err, err)
			t.Fail()
		}
	})
	t.Run("unexpected token not EOF", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.IDENT},
			actualToken:    token.NEWLINE,
		}

		isEOFNewlineErr := IsEOFInsteadOfNewlineError(err)

		if isEOFNewlineErr {
			t.Logf("Expected no EOF NEWLINE error, got %T:%q\n", err, err)
			t.Fail()
		}
	})
	t.Run("unexpected token EOF expected NEWLINE wrapped error", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.NEWLINE},
			actualToken:    token.EOF,
		}
		wrapped := errors.Wrap(err, "some error")

		isEOFNewlineErr := IsEOFInsteadOfNewlineError(wrapped)

		if !isEOFNewlineErr {
			t.Logf("Expected an EOF NEWLINE error, got %T:%q\n", wrapped, wrapped)
			t.Fail()
		}
	})
	t.Run("unexpected token EOF expected not NEWLINE wrapped error", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.IDENT},
			actualToken:    token.EOF,
		}
		wrapped := errors.Wrap(err, "some error")

		isEOFNewlineErr := IsEOFInsteadOfNewlineError(wrapped)

		if isEOFNewlineErr {
			t.Logf("Expected no EOF NEWLINE error, got %T:%q\n", wrapped, wrapped)
			t.Fail()
		}
	})
	t.Run("unexpected token not EOF wrapped error", func(t *testing.T) {
		err := &unexpectedTokenError{
			expectedTokens: []token.TokenType{token.IDENT},
			actualToken:    token.NEWLINE,
		}
		wrapped := errors.Wrap(err, "some error")

		isEOFNewlineErr := IsEOFInsteadOfNewlineError(wrapped)

		if isEOFNewlineErr {
			t.Logf("Expected no EOF NEWLINE error, got %T:%q\n", wrapped, wrapped)
			t.Fail()
		}
	})
	t.Run("random error", func(t *testing.T) {
		err := errors.New("some error")

		isEOFErr := IsEOFError(err)

		if isEOFErr {
			t.Logf("Expected no EOF NEWLINE error, got %T:%q\n", err, err)
			t.Fail()
		}
	})
}
