package object

import (
	"os"
	"testing"
)

func TestFileExpandPath(t *testing.T) {
	t.Run("one arg flavour", func(t *testing.T) {
		env := NewEnvironment()
		context := &callContext{
			receiver: &Self{RubyObject: fileClass, Name: "File"},
			env:      env,
		}
		filename := &String{Value: "./fixtures/testfile.rb"}

		result, err := fileExpandPath(context, filename)

		checkError(t, err, nil)

		cwd, err := os.Getwd()
		if err != nil {
			t.Skip("Cannot determine working directory")
		}
		expected := &String{Value: cwd + "/fixtures/testfile.rb"}

		checkResult(t, result, expected)
	})
	t.Run("two arg flavour", func(t *testing.T) {
		env := NewEnvironment()
		context := &callContext{
			receiver: &Self{RubyObject: fileClass, Name: "File"},
			env:      env,
		}
		filename := &String{Value: "../../main.go"}
		dirname := &String{Value: "object/fixtures/"}

		result, err := fileExpandPath(context, filename, dirname)

		checkError(t, err, nil)

		cwd, err := os.Getwd()
		if err != nil {
			t.Skip("Cannot determine working directory")
		}
		expected := &String{Value: cwd + "/main.go"}

		checkResult(t, result, expected)
	})
}

func TestFileDirname(t *testing.T) {
	context := &callContext{
		receiver: &Self{RubyObject: fileClass, Name: "File"},
		env:      NewEnvironment(),
	}
	filename := &String{Value: "/var/log/foo.log"}

	result, err := fileDirname(context, filename)

	checkError(t, err, nil)

	expected := &String{Value: "/var/log"}

	checkResult(t, result, expected)
}
