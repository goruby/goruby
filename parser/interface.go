package parser

import (
	"bytes"
	"fmt"
	gotoken "go/token"
	"io"
	"io/ioutil"

	"github.com/goruby/goruby/ast"
	"github.com/pkg/errors"
)

// If src != nil, readSource converts src to a []byte if possible;
// otherwise it returns an error. If src == nil, readSource returns
// the result of reading the file specified by filename.
//
func readSource(filename string, src interface{}) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			// is io.Reader, but src is already available in []byte form
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, s); err != nil {
				return nil, err
			}
			return buf.Bytes(), nil
		}
		return nil, errors.New("invalid source")
	}
	return ioutil.ReadFile(filename)
}

// A Mode value is a set of flags (or 0).
// They control the amount of source code parsed and other optional
// parser functionality.
//
type Mode uint

// parser modes
const (
	ParseComments Mode = 1 << iota // parse comments and add them to AST
	Trace                          // print a trace of parsed productions
	AllErrors                      // report all errors (not just the first 10 on different lines)
)

var parseModes = map[string]Mode{
	"ParseComments": ParseComments,
	"Trace":         Trace,
	"AllErrors":     AllErrors,
}

// ParseFile parses the source code of a single Ruby source file and returns
// the corresponding ast.Program node. The source code may be provided via
// the filename of the source file, or via the src parameter.
//
// If src != nil, ParseFile parses the source from src and the filename is
// only used when recording position information. The type of the argument
// for the src parameter must be string, []byte, or io.Reader.
// If src == nil, ParseFile parses the file specified by filename.
//
// The mode parameter controls the amount of source text parsed and other
// optional parser functionality. Position information is recorded in the
// file set fset, which must not be nil.
//
// If the source couldn't be read or the source was read but syntax
// errors were found, the returned AST is nil and the error
// indicates the specific failure.
//
func ParseFile(fset *gotoken.FileSet, filename string, src interface{}, mode Mode) (*ast.Program, error) {
	if fset == nil {
		panic("parser.ParseFile: no token.FileSet provided (fset == nil)")
	}

	// get source
	text, err := readSource(filename, src)
	if err != nil {
		return nil, err
	}

	var p parser
	p.init(fset, filename, text, mode)

	return p.ParseProgram()
}

// ParseExprFrom is a convenience function for parsing an expression.
// The arguments have the same meaning as for ParseFile, but the source must
// be a valid Go (type or value) expression. Specifically, fset must not
// be nil.
//
func ParseExprFrom(fset *gotoken.FileSet, filename string, src interface{}, mode Mode) (ast.Expression, error) {
	if fset == nil {
		panic("parser.ParseExprFrom: no token.FileSet provided (fset == nil)")
	}

	// get source
	text, err := readSource(filename, src)
	if err != nil {
		return nil, err
	}

	var p parser
	p.init(fset, filename, text, mode)

	program, err := p.ParseProgram()
	if err != nil {
		return nil, err
	}

	if len(program.Statements) == 0 {
		return nil, fmt.Errorf("source did not contain any expressions to parse")
	}

	for _, stmt := range program.Statements {
		if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
			return exprStmt.Expression, nil
		}
	}

	return nil, fmt.Errorf("source only contains statements")
}

// ParseExpr is a convenience function for obtaining the AST of an expression x.
// The position information recorded in the AST is undefined. The filename used
// in error messages is the empty string.
//
func ParseExpr(x string) (ast.Expression, error) {
	return ParseExprFrom(gotoken.NewFileSet(), "", []byte(x), 0)
}
