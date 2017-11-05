package object

import (
	"reflect"
	"testing"
)

func TestEnvStat(t *testing.T) {
	t.Run("top level env", func(t *testing.T) {
		obj := &Integer{42}

		env := &environment{store: map[string]RubyObject{"foo": obj}}

		info, ok := EnvStat(env, obj)

		if !ok {
			t.Logf("Expected obj to be found in env")
			t.FailNow()
		}

		expectedName := "foo"
		expectedEnv := env

		if expectedName != info.Name() {
			t.Logf("Expected info.Name to equal %q, got %q", expectedName, info.Name())
			t.Fail()
		}

		if expectedEnv != info.Env() {
			t.Logf("Expected env to equal\n%+#v\ngot\n\t%+#v\n", expectedEnv, info.Env())
			t.Fail()
		}
	})
	t.Run("two level nested", func(t *testing.T) {
		obj := &Integer{42}

		root := &environment{store: map[string]RubyObject{"foo": obj}}
		outer := &environment{store: make(map[string]RubyObject), outer: root}
		env := &environment{store: make(map[string]RubyObject), outer: outer}

		info, ok := EnvStat(env, obj)

		if !ok {
			t.Logf("Expected obj to be found in env")
			t.FailNow()
		}

		expectedName := "foo"
		expectedEnv := root

		if expectedName != info.Name() {
			t.Logf("Expected info.Name to equal %q, got %q", expectedName, info.Name())
			t.Fail()
		}

		if expectedEnv != info.Env() {
			t.Logf("Expected env to equal\n%+#v\ngot\n\t%+#v\n", expectedEnv, info.Env())
			t.Fail()
		}
	})
	t.Run("two level nested same value with different keys", func(t *testing.T) {
		obj := &Integer{42}

		root := &environment{store: map[string]RubyObject{"foo": obj}}
		outer := &environment{store: map[string]RubyObject{"bar": obj}, outer: root}
		env := &environment{store: make(map[string]RubyObject), outer: outer}

		info, ok := EnvStat(env, obj)

		if !ok {
			t.Logf("Expected obj to be found in env")
			t.FailNow()
		}

		expectedName := "bar"
		expectedEnv := outer

		if expectedName != info.Name() {
			t.Logf("Expected info.Name to equal %q, got %q", expectedName, info.Name())
			t.Fail()
		}

		if expectedEnv != info.Env() {
			t.Logf("Expected env to equal\n%+#v\ngot\n\t%+#v\n", expectedEnv, info.Env())
			t.Fail()
		}
	})
	t.Run("two level nested overshadowed key", func(t *testing.T) {
		obj := &Integer{42}

		root := &environment{store: map[string]RubyObject{"foo": obj}}
		outer := &environment{store: map[string]RubyObject{"foo": TRUE}, outer: root}
		env := &environment{store: make(map[string]RubyObject), outer: outer}

		info, ok := EnvStat(env, obj)

		if !ok {
			t.Logf("Expected obj to be found in env")
			t.FailNow()
		}

		expectedName := "foo"
		expectedEnv := root

		if expectedName != info.Name() {
			t.Logf("Expected info.Name to equal %q, got %q", expectedName, info.Name())
			t.Fail()
		}

		if expectedEnv != info.Env() {
			t.Logf("Expected env to equal\n%+#v\ngot\n\t%+#v\n", expectedEnv, info.Env())
			t.Fail()
		}
	})
	t.Run("two level nested not found", func(t *testing.T) {
		obj := &Integer{42}

		root := &environment{store: map[string]RubyObject{"foo": FALSE}}
		outer := &environment{store: map[string]RubyObject{"bar": TRUE}, outer: root}
		env := &environment{store: make(map[string]RubyObject), outer: outer}

		_, ok := EnvStat(env, obj)

		if ok {
			t.Logf("Expected obj not to be found in env")
			t.FailNow()
		}
	})
}

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

func TestWithScopedLocalVariables(t *testing.T) {
	env := &environment{store: map[string]RubyObject{
		"foo": TRUE,
		"Bar": TRUE,
	}}

	scopedEnv := WithScopedLocalVariables(env)

	localVarsGuard, ok := scopedEnv.(*localVariableGuard)
	if !ok {
		t.Logf("Expected to get a %T environent, got %T\n", localVarsGuard, scopedEnv)
		t.FailNow()
	}

	if !reflect.DeepEqual(env, localVarsGuard.Environment) {
		t.Logf("Expected embedded env to equal\n%+#v\n\tgot\n%+#v\n", env, localVarsGuard.Environment)
		t.Fail()
	}

	expected := &environment{store: map[string]RubyObject{}}

	if !reflect.DeepEqual(expected, localVarsGuard.localVariables) {
		t.Logf("Expected local variable env to equal\n%+#v\n\tgot\n%+#v\n", expected, localVarsGuard.localVariables)
		t.Fail()
	}
}

func Test_localVariableGuardSet(t *testing.T) {
	env := &localVariableGuard{
		Environment:    &environment{store: map[string]RubyObject{}},
		localVariables: &environment{store: map[string]RubyObject{}},
	}

	t.Run("local variables", func(t *testing.T) {
		env.Set("foo", TRUE)

		_, ok := env.Environment.Get("foo")
		if ok {
			t.Logf("Expected embedded env not to contain 'foo'")
			t.Fail()
		}

		_, ok = env.localVariables.Get("foo")
		if !ok {
			t.Logf("Expected localVariables env to contain 'foo'")
			t.Fail()
		}
	})
	t.Run("instance variables", func(t *testing.T) {
		env.Set("@foo", TRUE)

		_, ok := env.Environment.Get("@foo")
		if !ok {
			t.Logf("Expected embedded env to contain '@foo'")
			t.Fail()
		}

		_, ok = env.localVariables.Get("@foo")
		if ok {
			t.Logf("Expected localVariables env not to contain '@foo'")
			t.Fail()
		}
	})
	t.Run("constants", func(t *testing.T) {
		env.Set("Foo", TRUE)

		_, ok := env.Environment.Get("Foo")
		if !ok {
			t.Logf("Expected embedded env to contain 'Foo'")
			t.Fail()
		}

		_, ok = env.localVariables.Get("Foo")
		if ok {
			t.Logf("Expected localVariables env not to contain 'Foo'")
			t.Fail()
		}
	})
}

func Test_localVariableGuardGet(t *testing.T) {
	embeddedEnv := &environment{store: map[string]RubyObject{
		"self": TRUE,
		"foo":  TRUE,
		"qux":  TRUE,
		"Foo":  TRUE,
		"Qux":  TRUE,
		"@foo": TRUE,
	}}
	env := &localVariableGuard{
		Environment: embeddedEnv,
		localVariables: &environment{store: map[string]RubyObject{
			"foo": FALSE,
			"bar": TRUE,
		}},
	}

	t.Run("constants", func(t *testing.T) {
		val, ok := env.Get("Foo")
		if !ok {
			t.Logf("Expected env to contain 'Foo'")
			t.Fail()
		}

		checkResult(t, TRUE, val)
	})

	t.Run("instance variables", func(t *testing.T) {
		val, ok := env.Get("@foo")
		if !ok {
			t.Logf("Expected env to contain '@foo'")
			t.Fail()
		}

		checkResult(t, TRUE, val)
	})

	t.Run("self", func(t *testing.T) {
		self, ok := env.Get("self")
		if !ok {
			t.Logf("Expected env to contain 'self'")
			t.Fail()
		}

		checkResult(t, TRUE, self)
	})

	t.Run("local variables", func(t *testing.T) {
		val, ok := env.Get("foo")
		if !ok {
			t.Logf("Expected env to contain 'foo'")
			t.Fail()
		}

		checkResult(t, FALSE, val)

		_, ok = env.Get("bar")
		if !ok {
			t.Logf("Expected env to contain 'bar'")
			t.Fail()
		}

		_, ok = env.Get("qux")
		if ok {
			t.Logf("Expected env not to contain 'qux'")
			t.Fail()
		}
	})
}

func Test_localVariableGuardGetAll(t *testing.T) {
	embeddedEnv := &environment{store: map[string]RubyObject{
		"foo":  TRUE,
		"qux":  TRUE,
		"Foo":  TRUE,
		"Qux":  TRUE,
		"@foo": TRUE,
	}}
	env := &localVariableGuard{
		Environment: embeddedEnv,
		localVariables: &environment{store: map[string]RubyObject{
			"foo": FALSE,
			"bar": TRUE,
		}},
	}

	result := env.GetAll()

	expected := map[string]RubyObject{
		"foo":  FALSE,
		"bar":  TRUE,
		"qux":  TRUE,
		"Foo":  TRUE,
		"Qux":  TRUE,
		"@foo": TRUE,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Logf("Expected result to equal\n%+#v\n\tgot\n%+#v\n", expected, result)
		t.Fail()
	}
}

func Test_localVariableGuardUnset(t *testing.T) {
	embeddedEnv := &environment{store: map[string]RubyObject{
		"foo":  TRUE,
		"qux":  TRUE,
		"Foo":  TRUE,
		"Qux":  TRUE,
		"@foo": TRUE,
	}}
	env := &localVariableGuard{
		Environment: embeddedEnv,
		localVariables: &environment{store: map[string]RubyObject{
			"foo": FALSE,
			"bar": TRUE,
		}},
	}

	t.Run("constants", func(t *testing.T) {
		env.Unset("Foo")

		_, ok := env.Get("Foo")
		if ok {
			t.Logf("Expected env not to contain 'Foo'")
			t.Fail()
		}
	})

	t.Run("instance variables", func(t *testing.T) {
		env.Unset("@foo")

		_, ok := env.Get("@foo")
		if ok {
			t.Logf("Expected env not to contain '@foo'")
			t.Fail()
		}
	})

	t.Run("local variables", func(t *testing.T) {
		env.Unset("foo")

		_, ok := env.Get("foo")
		if ok {
			t.Logf("Expected env not to contain 'foo'")
			t.Fail()
		}

		_, ok = env.Environment.Get("foo")
		if !ok {
			t.Logf("Expected embedded env to contain 'foo'")
			t.Fail()
		}

		env.Unset("bar")

		_, ok = env.Get("bar")
		if ok {
			t.Logf("Expected env not to contain 'bar'")
			t.Fail()
		}

		env.Unset("qux")

		_, ok = env.Get("qux")
		if ok {
			t.Logf("Expected env not to contain 'qux'")
			t.Fail()
		}

		_, ok = env.Environment.Get("qux")
		if !ok {
			t.Logf("Expected embedded env to contain 'qux'")
			t.Fail()
		}
	})
}

func Test_localVariableGuardClone(t *testing.T) {
	embeddedEnv := &environment{store: map[string]RubyObject{
		"foo":  TRUE,
		"qux":  TRUE,
		"Foo":  TRUE,
		"Qux":  TRUE,
		"@foo": TRUE,
	}}
	env := &localVariableGuard{
		Environment: embeddedEnv,
		localVariables: &environment{store: map[string]RubyObject{
			"foo": FALSE,
			"bar": TRUE,
		}},
	}

	result := env.Clone()

	actual := result.GetAll()

	expected := map[string]RubyObject{
		"foo":  FALSE,
		"bar":  TRUE,
		"qux":  TRUE,
		"Foo":  TRUE,
		"Qux":  TRUE,
		"@foo": TRUE,
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected result to equal\n%+#v\n\tgot\n%+#v\n", expected, actual)
		t.Fail()
	}
}
