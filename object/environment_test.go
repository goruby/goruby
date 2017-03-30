package object

import "testing"

func TestEnvironmentSet(t *testing.T) {
	env := &environment{store: make(map[string]RubyObject)}

	env.Set("foo", NIL)

	val, ok := env.store["foo"]
	if !ok {
		t.Logf("Expected store to contain 'foo'")
		t.Fail()
	}

	if val != NIL {
		t.Logf("Expected value to equal NIL, got %v", val)
		t.Fail()
	}
}

func TestEnvironmentSetGlobal(t *testing.T) {
	t.Run("toplevel env", func(t *testing.T) {
		env := &environment{store: make(map[string]RubyObject)}

		env.SetGlobal("$foo", NIL)

		val, ok := env.store["$foo"]
		if !ok {
			t.Logf("Expected store to contain '$foo'")
			t.Fail()
		}

		if val != NIL {
			t.Logf("Expected value to equal NIL, got %v", val)
			t.Fail()
		}
	})
	t.Run("inner env one level", func(t *testing.T) {
		outer := &environment{store: make(map[string]RubyObject)}
		env := &environment{store: make(map[string]RubyObject), outer: outer}

		env.SetGlobal("$foo", NIL)

		_, ok := env.store["$foo"]
		if ok {
			t.Logf("Expected env store to not contain '$foo'")
			t.Fail()
		}

		val, ok := outer.store["$foo"]
		if !ok {
			t.Logf("Expected outer store to contain '$foo'")
			t.Fail()
		}

		if val != NIL {
			t.Logf("Expected value to equal NIL, got %v", val)
			t.Fail()
		}
	})
	t.Run("inner env two level", func(t *testing.T) {
		root := &environment{store: make(map[string]RubyObject)}
		outer := &environment{store: make(map[string]RubyObject), outer: root}
		env := &environment{store: make(map[string]RubyObject), outer: outer}

		env.SetGlobal("$foo", NIL)

		_, ok := env.store["$foo"]
		if ok {
			t.Logf("Expected env store to not contain '$foo'")
			t.Fail()
		}

		_, ok = outer.store["$foo"]
		if ok {
			t.Logf("Expected outer store to not contain '$foo'")
			t.Fail()
		}

		val, ok := root.store["$foo"]
		if !ok {
			t.Logf("Expected root store to contain '$foo'")
			t.Fail()
		}

		if val != NIL {
			t.Logf("Expected value to equal NIL, got %v", val)
			t.Fail()
		}
	})
}

func TestEnvironmentGet(t *testing.T) {
	t.Run("toplevel env", func(t *testing.T) {
		env := &environment{store: map[string]RubyObject{"foo": TRUE}}

		val, ok := env.Get("foo")
		if !ok {
			t.Logf("Expected store to contain 'foo'")
			t.Fail()
		}

		if val != TRUE {
			t.Logf("Expected value to equal TRUE, got %v", val)
			t.Fail()
		}
	})
	t.Run("inner env one level", func(t *testing.T) {
		outer := &environment{store: map[string]RubyObject{"foo": TRUE}}
		env := &environment{store: make(map[string]RubyObject), outer: outer}

		val, ok := env.Get("foo")
		if !ok {
			t.Logf("Expected store to contain 'foo'")
			t.Fail()
		}

		if val != TRUE {
			t.Logf("Expected value to equal TRUE, got %v", val)
			t.Fail()
		}
	})
	t.Run("inner env two level", func(t *testing.T) {
		root := &environment{store: map[string]RubyObject{"foo": TRUE}}
		outer := &environment{store: make(map[string]RubyObject), outer: root}
		env := &environment{store: make(map[string]RubyObject), outer: outer}

		val, ok := env.Get("foo")
		if !ok {
			t.Logf("Expected store to contain 'foo'")
			t.Fail()
		}

		if val != TRUE {
			t.Logf("Expected value to equal TRUE, got %v", val)
			t.Fail()
		}
	})
	t.Run("inner env two level overridden key", func(t *testing.T) {
		root := &environment{store: map[string]RubyObject{"foo": FALSE}}
		outer := &environment{store: map[string]RubyObject{"foo": TRUE}, outer: root}
		env := &environment{store: make(map[string]RubyObject), outer: outer}

		val, ok := env.Get("foo")
		if !ok {
			t.Logf("Expected store to contain 'foo'")
			t.Fail()
		}

		if val != TRUE {
			t.Logf("Expected value to equal TRUE, got %v", val)
			t.Fail()
		}
	})
}
