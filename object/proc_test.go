package object

import (
	"reflect"
	"testing"
)

func TestExtractBlockFromArgs(t *testing.T) {
	t.Run("args empty", func(t *testing.T) {
		args := []RubyObject{}

		block, remaining, ok := extractBlockFromArgs(args)

		if ok {
			t.Logf("Expected no block found")
			t.Fail()
		}

		if block != nil {
			t.Logf("Expected block to be nil, got %+#v\n", block)
			t.Fail()
		}

		if len(remaining) != 0 {
			t.Logf("Expected remaining args to have length %d, got %d", 0, len(remaining))
			t.Fail()
		}
	})
	t.Run("args with block only", func(t *testing.T) {
		args := []RubyObject{&Proc{}}

		block, remaining, ok := extractBlockFromArgs(args)

		if !ok {
			t.Logf("Expected block found")
			t.Fail()
		}

		if block == nil {
			t.Logf("Expected block not to be nil")
			t.Fail()
		}

		if len(remaining) != 0 {
			t.Logf("Expected remaining args to have length %d, got %d", 0, len(remaining))
			t.Fail()
		}
	})
	t.Run("args with nil and block", func(t *testing.T) {
		args := []RubyObject{NIL, &Proc{}}

		block, remaining, ok := extractBlockFromArgs(args)

		if !ok {
			t.Logf("Expected block found")
			t.Fail()
		}

		if block == nil {
			t.Logf("Expected block not to be nil")
			t.Fail()
		}

		if len(remaining) != 1 {
			t.Logf("Expected remaining args to have length %d, got %d", 1, len(remaining))
			t.Fail()
		}

		expected := []RubyObject{NIL}

		if !reflect.DeepEqual(expected, remaining) {
			t.Logf("Expected remaining args to equal\n%+#v\n\tgot\n%+#v\n", expected, remaining)
			t.Fail()
		}
	})
	t.Run("args with nil and block but block not at the end", func(t *testing.T) {
		args := []RubyObject{&Proc{}, NIL}

		block, remaining, ok := extractBlockFromArgs(args)

		if ok {
			t.Logf("Expected no block found")
			t.Fail()
		}

		if block != nil {
			t.Logf("Expected block to be nil, got %+#v\n", block)
			t.Fail()
		}

		if len(remaining) != 2 {
			t.Logf("Expected remaining args to have length %d, got %d", 2, len(remaining))
			t.Fail()
		}

		expected := []RubyObject{&Proc{}, NIL}

		if !reflect.DeepEqual(expected, remaining) {
			t.Logf("Expected remaining args to equal\n%+#v\n\tgot\n%+#v\n", expected, remaining)
			t.Fail()
		}
	})
}
