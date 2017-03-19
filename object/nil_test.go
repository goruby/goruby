package object

import "testing"

func TestNilIsNil(t *testing.T) {
	result := nilIsNil(NIL)

	boolean, ok := result.(*Boolean)
	if !ok {
		t.Logf("Expected Boolean, got %T", result)
		t.FailNow()
	}

	if boolean.Value != true {
		t.Logf("Expected true, got false")
		t.Fail()
	}
}
