package ast

import (
	"go/ast"
	"io"
	"reflect"
)

// A FieldFilter may be provided to Fprint to control the output.
type FieldFilter func(name string, value reflect.Value) bool

// NotNilFilter returns true for field values that are not nil;
// it returns false otherwise.
func NotNilFilter(_ string, v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return !v.IsNil()
	}
	return true
}

// Fprint prints the (sub-)tree starting at AST node x to w.
//
// A non-nil FieldFilter f may be provided to control the output: struct fields
// for which f(fieldname, fieldvalue) is true are printed; all others are
// filtered from the output. Unexported struct fields are never printed.
func Fprint(w io.Writer, x interface{}, f FieldFilter) error {
	return ast.Fprint(w, nil, x, ast.FieldFilter(f))
}

// Print prints x to standard output, skipping nil fields.
// Print(x) is the same as Fprint(os.Stdout, x, NotNilFilter).
func Print(x interface{}) error {
	return ast.Print(nil, x)
}
