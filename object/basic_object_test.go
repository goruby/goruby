package object

import (
	"reflect"
	"testing"
)

func TestMethodSetDefine(t *testing.T) {
	set := methodSet{context: &Integer{3}, methods: make(map[string]function)}

	sym := set.Define("foo", nil)

	expectedSymbol := &Symbol{"foo"}

	if !reflect.DeepEqual(expectedSymbol, sym) {
		t.Logf("Expected symbol to equal %s, got %s\n", expectedSymbol, sym)
		t.Fail()
	}

	expectedMethodsLength := 1
	actualMethodsLength := len(set.methods)
	if actualMethodsLength != expectedMethodsLength {
		t.Logf("Expected methods to have %d items, got %d\n", expectedMethodsLength, actualMethodsLength)
		t.Fail()
	}
}
