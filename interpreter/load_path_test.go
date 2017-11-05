package interpreter_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/goruby/goruby/interpreter"
	"github.com/goruby/goruby/object"
	"github.com/pkg/errors"
)

func TestLoadPath(t *testing.T) {
	content := []byte("$foo = 12")
	tmpDir, err := ioutil.TempDir("", "example")
	if err != nil {
		panic(err)
	}
	tmpAbsPath := path.Join(tmpDir, "example.rb")
	err = ioutil.WriteFile(tmpAbsPath, content, 0600)
	if err != nil {
		panic(err)
	}
	tmpBase := path.Base(tmpAbsPath)

	defer os.RemoveAll(tmpDir) // clean up

	t.Run("load absolute file", func(t *testing.T) {
		input := fmt.Sprintf(`
		require '%s'
		$foo
		`, tmpAbsPath)

		i := interpreter.New()

		_, err = i.Interpret("", input)

		if err != nil {
			t.Logf("Expected no error, got %T:%s\n", err, err)
			t.Fail()
		}
	})
	t.Run("load outside the load path", func(t *testing.T) {
		input := fmt.Sprintf(`
		require '%s'
		$foo
		`, tmpBase)

		i := interpreter.New()

		_, err = i.Interpret("", input)

		expectedError := object.NewNoSuchFileLoadError(tmpBase)

		if !reflect.DeepEqual(expectedError, errors.Cause(err)) {
			t.Logf("Expected err to equal\n%s\n\tgot\n%s\n", expectedError, err)
			t.Fail()
		}
	})
	t.Run("load within the load path", func(t *testing.T) {
		input := fmt.Sprintf(`
		$:.push '%s'

		require '%s'
		$foo
		`, tmpAbsPath, tmpBase)

		i := interpreter.New()

		_, err = i.Interpret("", input)

		if err != nil {
			t.Logf("Expected no error, got %T:%s\n", err, err)
			t.Fail()
		}
	})
}
