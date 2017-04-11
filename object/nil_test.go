package object

import "testing"

func TestNilIsNil(t *testing.T) {
	context := &callContext{receiver: NIL}

	result, err := nilIsNil(context)

	checkError(t, err, nil)

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
