package object

import "testing"

func TestExceptionClassNew(t *testing.T) {

}

func TestExceptionInitialize(t *testing.T) {
	context := &callContext{
		receiver: &Self{RubyObject: &Exception{}},
		env:      NewMainEnvironment(),
	}
	t.Run("without args", func(t *testing.T) {
		result, err := exceptionInitialize(context)

		checkError(t, err, nil)

		checkResult(t, result, &Exception{message: "Exception"})
	})
	t.Run("with arg", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			result, err := exceptionInitialize(context, &String{Value: "err"})

			checkError(t, err, nil)

			checkResult(t, result, &Exception{message: "err"})
		})
		t.Run("other object", func(t *testing.T) {
			result, err := exceptionInitialize(context, &Symbol{Value: "symbol"})

			checkError(t, err, nil)

			checkResult(t, result, &Exception{message: "symbol"})
		})
	})
}

func TestExceptionClassException(t *testing.T) {
	context := &callContext{
		receiver: exceptionClass,
		env:      NewMainEnvironment(),
	}
	t.Run("without args", func(t *testing.T) {
		result, err := exceptionClassException(context)

		checkError(t, err, nil)

		checkResult(t, result, &Exception{message: "Exception"})
	})
	t.Run("with arg", func(t *testing.T) {
		result, err := exceptionClassException(context, &String{Value: "err"})

		checkError(t, err, nil)

		checkResult(t, result, &Exception{message: "err"})
	})
}

func TestExceptionException(t *testing.T) {
	contextObject := &Exception{message: "x"}
	context := &callContext{
		receiver: contextObject,
		env:      NewMainEnvironment(),
	}
	t.Run("without args", func(t *testing.T) {
		result, err := exceptionException(context)

		checkError(t, err, nil)

		if contextObject != result {
			t.Logf("Expected result to pointer equal context\n")
			t.Fail()
		}
		checkResult(t, result, &Exception{message: "x"})
	})
	t.Run("with arg", func(t *testing.T) {
		result, err := exceptionException(context, &String{Value: "x"})

		checkError(t, err, nil)

		if contextObject != result {
			t.Logf("Expected result to pointer equal context\n")
			t.Fail()
		}

		checkResult(t, result, &Exception{message: "x"})
	})
	t.Run("with arg but different message", func(t *testing.T) {
		result, err := exceptionException(context, &String{Value: "err"})

		checkError(t, err, nil)

		checkResult(t, result, &Exception{message: "err"})
	})
}
