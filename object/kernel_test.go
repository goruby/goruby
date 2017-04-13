package object

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/goruby/goruby/ast"
)

func TestKernelMethods(t *testing.T) {
	t.Run("without superclass", func(t *testing.T) {
		contextMethods := map[string]RubyMethod{
			"foo": publicMethod(nil),
			"bar": publicMethod(nil),
		}
		context := &callContext{
			receiver: &testRubyObject{
				class: &class{
					instanceMethods: contextMethods,
					superClass:      nil,
				},
			},
		}

		result, err := kernelMethods(context)

		checkError(t, err, nil)

		array, ok := result.(*Array)
		if !ok {
			t.Logf("Expected array, got %T", result)
			t.FailNow()
		}

		var methods []string
		for i, elem := range array.Elements {
			sym, ok := elem.(*Symbol)
			if !ok {
				t.Logf("Expected all elements to be symbols, got %T at index %d", elem, i)
				t.Fail()
			} else {
				methods = append(methods, sym.Inspect())
			}
		}

		var expectedMethods = []string{
			":foo", ":bar",
		}

		expectedLen := len(expectedMethods)

		if len(array.Elements) != expectedLen {
			t.Logf("Expected %d items, got %d", expectedLen, len(array.Elements))
			t.Fail()
		}

		sort.Strings(expectedMethods)
		sort.Strings(methods)

		if !reflect.DeepEqual(expectedMethods, methods) {
			t.Logf("Expected methods to equal\n%s\n\tgot\n%s\n", expectedMethods, methods)
			t.Fail()
		}
	})
	t.Run("with superclass", func(t *testing.T) {
		superClassMethods := map[string]RubyMethod{
			"super_foo": publicMethod(nil),
			"super_bar": publicMethod(nil),
		}
		contextMethods := map[string]RubyMethod{
			"foo": publicMethod(nil),
			"bar": publicMethod(nil),
		}
		context := &callContext{
			receiver: &testRubyObject{
				class: &class{
					instanceMethods: contextMethods,
					superClass: &class{
						instanceMethods: superClassMethods,
						superClass:      nil,
					},
				},
			},
		}

		result, err := kernelMethods(context)

		checkError(t, err, nil)

		array, ok := result.(*Array)
		if !ok {
			t.Logf("Expected array, got %T", result)
			t.FailNow()
		}

		var methods []string
		for i, elem := range array.Elements {
			sym, ok := elem.(*Symbol)
			if !ok {
				t.Logf("Expected all elements to be symbols, got %T at index %d", elem, i)
				t.Fail()
			} else {
				methods = append(methods, sym.Inspect())
			}
		}

		var expectedMethods = []string{
			":foo", ":bar", ":super_foo", ":super_bar",
		}

		expectedLen := len(expectedMethods)

		if len(array.Elements) != expectedLen {
			t.Logf("Expected %d items, got %d", expectedLen, len(array.Elements))
			t.Fail()
		}

		sort.Strings(expectedMethods)
		sort.Strings(methods)

		if !reflect.DeepEqual(expectedMethods, methods) {
			t.Logf("Expected methods to equal\n%s\n\tgot\n%s\n", expectedMethods, methods)
			t.Fail()
		}
	})
	t.Run("with private methods", func(t *testing.T) {
		contextMethods := map[string]RubyMethod{
			"foo":         publicMethod(nil),
			"bar":         publicMethod(nil),
			"private_foo": privateMethod(nil),
		}
		context := &callContext{
			receiver: &testRubyObject{
				class: &class{
					instanceMethods: contextMethods,
					superClass:      nil,
				},
			},
		}

		result, err := kernelMethods(context)

		checkError(t, err, nil)

		array, ok := result.(*Array)
		if !ok {
			t.Logf("Expected array, got %T", result)
			t.FailNow()
		}

		var methods []string
		for i, elem := range array.Elements {
			sym, ok := elem.(*Symbol)
			if !ok {
				t.Logf("Expected all elements to be symbols, got %T at index %d", elem, i)
				t.Fail()
			} else {
				methods = append(methods, sym.Inspect())
			}
		}

		var expectedMethods = []string{
			":foo", ":bar",
		}

		expectedLen := len(expectedMethods)

		if len(array.Elements) != expectedLen {
			t.Logf("Expected %d items, got %d", expectedLen, len(array.Elements))
			t.Fail()
		}

		sort.Strings(expectedMethods)
		sort.Strings(methods)

		if !reflect.DeepEqual(expectedMethods, methods) {
			t.Logf("Expected methods to equal\n%s\n\tgot\n%s\n", expectedMethods, methods)
			t.Fail()
		}
	})
}

func TestKernelIsNil(t *testing.T) {
	context := &callContext{receiver: TRUE}
	result, err := kernelIsNil(context)

	checkError(t, err, nil)

	boolean, ok := result.(*Boolean)
	if !ok {
		t.Logf("Expected Boolean, got %T", result)
		t.FailNow()
	}

	if boolean.Value != false {
		t.Logf("Expected false, got true")
		t.Fail()
	}
}

func TestKernelClass(t *testing.T) {
	t.Run("regular object", func(t *testing.T) {
		context := &callContext{receiver: &Integer{1}}

		result, err := kernelClass(context)

		checkError(t, err, nil)

		cl, ok := result.(*class)
		if !ok {
			t.Logf("Expected Class, got %T", result)
			t.Fail()
		}

		expected := integerClass

		if !reflect.DeepEqual(expected, cl) {
			t.Logf("Expected class to equal %+#v, got %+#v", expected, cl)
			t.Fail()
		}
	})
	t.Run("class object", func(t *testing.T) {
		context := &callContext{receiver: stringClass}

		result, err := kernelClass(context)

		checkError(t, err, nil)

		cl, ok := result.(*class)
		if !ok {
			t.Logf("Expected Class, got %T", result)
			t.Fail()
		}

		expected := classClass

		if !reflect.DeepEqual(expected, cl) {
			t.Logf("Expected class to equal %+#v, got %+#v", expected, cl)
			t.Fail()
		}
	})
	t.Run("class class", func(t *testing.T) {
		context := &callContext{receiver: classClass}

		result, err := kernelClass(context)

		checkError(t, err, nil)

		cl, ok := result.(*class)
		if !ok {
			t.Logf("Expected Class, got %T", result)
			t.Fail()
		}

		expected := classClass

		if !reflect.DeepEqual(expected, cl) {
			t.Logf("Expected class to equal %+#v, got %+#v", expected, cl)
			t.Fail()
		}
	})
}

func TestKernelRequire(t *testing.T) {
	t.Run("wiring together", func(t *testing.T) {
		evalCallCount := 0
		var evalCallASTNode ast.Node
		var evalCallEnv Environment
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			evalCallCount++
			evalCallASTNode = node
			evalCallEnv = env
			return TRUE, nil
		}

		context := &callContext{
			env:      NewEnvironment(),
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile.rb"}

		result, err := kernelRequire(context, name)

		if err != nil {
			t.Logf("expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		boolean, ok := result.(*Boolean)
		if !ok {
			t.Logf("Expected Boolean, got %#v", result)
			t.FailNow()
		}

		if boolean != TRUE {
			t.Logf("Expected return to equal TRUE, got FALSE")
			t.Fail()
		}

		if evalCallCount != 1 {
			t.Logf("Expected context.Eval to be called once, was %d\n", evalCallCount)
			t.Fail()
		}

		expectedASTNodeString := "x = 5"
		actualASTNodeString := evalCallASTNode.String()
		if expectedASTNodeString != actualASTNodeString {
			t.Logf("Expected Eval AST param to equal %q, got %q\n", expectedASTNodeString, actualASTNodeString)
			t.Fail()
		}
	})
	t.Run("env side effects no $LOADED_FEATURES", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile.rb"}

		_, err := kernelRequire(context, name)
		if err != nil {
			panic(err)
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		expected := NewArray(&String{abs})

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
	t.Run("env side effects missing suffix", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile"}

		_, err := kernelRequire(context, name)
		if err != nil {
			panic(err)
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		expected := NewArray(&String{abs})

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
	t.Run("env side effects $LOADED_FEATURES exist", func(t *testing.T) {
		env := NewEnvironment()
		env.SetGlobal("$LOADED_FEATURES", NewArray(&String{"foo"}))
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile"}

		_, err := kernelRequire(context, name)
		if err != nil {
			panic(err)
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		expected := NewArray(&String{"foo"}, &String{abs})

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
	t.Run("env side effects local variables", func(t *testing.T) {
		env := NewEnvironment()
		var eval func(node ast.Node, env Environment) (RubyObject, error)
		eval = func(node ast.Node, env Environment) (RubyObject, error) {
			switch node := node.(type) {
			case *ast.Program:
				var result RubyObject
				var err error
				for _, statement := range node.Statements {
					result, err = eval(statement, env)

					if err != nil {
						return nil, err
					}
				}
				return result, nil
			case *ast.ExpressionStatement:
				return eval(node.Expression, env)
			case *ast.VariableAssignment:
				val, err := eval(node.Value, env)
				if err != nil {
					return nil, err
				}
				env.Set(node.Name.Value, val)
				return val, nil
			}
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile"}

		_, err := kernelRequire(context, name)
		if err != nil {
			panic(err)
		}

		_, ok := env.Get("x")

		if ok {
			t.Logf("Expected local variable not to leak over require")
			t.Fail()
		}
	})
	t.Run("file does not exist", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"file/not/exist"}

		_, err := kernelRequire(context, name)
		if err == nil {
			t.Logf("Expected error not to be nil")
			t.Fail()
		}

		expectedErr := NewLoadError("file/not/exist")
		if !reflect.DeepEqual(expectedErr, err) {
			t.Logf("Expected error to equal %v, got %v", expectedErr, err)
			t.Fail()
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		expected := NewArray()

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
	t.Run("syntax error", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile_syntax_error.rb"}

		_, err := kernelRequire(context, name)
		if err == nil {
			t.Logf("Expected error not to be nil")
			t.Fail()
		}

		expectedErr := "syntax error, Parsing errors:\n\texpected next token to be of type [END], got EOF instead\n"
		if !reflect.DeepEqual(expectedErr, err.Error()) {
			t.Logf("Expected error to equal\n%q\n\tgot\n%q", expectedErr, err.Error())
			t.Fail()
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		expected := NewArray()

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
	t.Run("thrown error", func(t *testing.T) {
		env := NewEnvironment()
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return nil, NewException("something went wrong")
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile_name_error.rb"}

		_, err := kernelRequire(context, name)
		if err == nil {
			t.Logf("Expected error not to be nil")
			t.Fail()
		}

		expectedErr := NewException("something went wrong")
		if !reflect.DeepEqual(expectedErr, err) {
			t.Logf("Expected error to equal\n%q\n\tgot\n%q", expectedErr, err)
			t.Fail()
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		expected := NewArray()

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
	t.Run("already loaded", func(t *testing.T) {
		abs, _ := filepath.Abs("./fixtures/testfile.rb")
		env := NewEnvironment()
		env.SetGlobal("$LOADED_FEATURES", NewArray(&String{abs}))
		eval := func(node ast.Node, env Environment) (RubyObject, error) {
			return TRUE, nil
		}

		context := &callContext{
			env:      env,
			eval:     eval,
			receiver: &Object{},
		}
		name := &String{"./fixtures/testfile.rb"}

		result, err := kernelRequire(context, name)
		if err != nil {
			t.Logf("Expected no error, got %T:%v", err, err)
			t.Fail()
		}

		boolean, ok := result.(*Boolean)
		if !ok {
			t.Logf("Expected Boolean, got %#v", result)
			t.FailNow()
		}

		if boolean != FALSE {
			t.Logf("Expected return to equal FALSE, got TRUE")
			t.Fail()
		}

		loadedFeatures, ok := env.Get("$LOADED_FEATURES")

		if !ok {
			t.Logf("Expected env to contain global $LOADED_FEATURES")
			t.Fail()
		}

		arr, ok := loadedFeatures.(*Array)

		if !ok {
			t.Logf("Expected $LOADED_FEATURES to be an Array, got %T", loadedFeatures)
			t.FailNow()
		}

		expected := NewArray(&String{abs})

		if !reflect.DeepEqual(expected, arr) {
			t.Logf("Expected $LOADED_FEATURES to equal\n%#v\n\tgot\n%#v\n", expected.Inspect(), arr.Inspect())
			t.Fail()
		}
	})
}

func TestKernelExtend(t *testing.T) {
	objectToExtend := &Object{}
	env := NewEnvironment()
	env.Set("foo", objectToExtend)
	context := &callContext{
		receiver: objectToExtend,
		env:      env,
	}

	module := newModule("Ext", map[string]RubyMethod{
		"foo": publicMethod(nil),
	})

	result, err := kernelExtend(context, module)

	checkError(t, err, nil)

	extended, ok := result.(*extendedObject)
	if !ok {
		t.Logf("Expected result to be an extendedObject, got %T", result)
		t.Fail()
	}

	if !reflect.DeepEqual(objectToExtend, extended.RubyObject) {
		t.Logf("Expected result to equal %+#v, got %+#v\n", objectToExtend, extended.RubyObject)
		t.Fail()
	}

	expectedClass := &eigenclass{
		map[string]RubyMethod{},
		&methodSet{
			objectToExtend.Class().(RubyClassObject),
			[]*Module{module},
		},
	}

	if !reflect.DeepEqual(expectedClass, extended.Class()) {
		t.Logf("Expected wrapped class to equal\n%+#v\n\tgot\n%+#v\n", expectedClass, extended.Class())
		t.Fail()
	}

	actual, ok := env.Get("foo")
	if !ok {
		panic("Not found in env")
	}

	if !reflect.DeepEqual(extended, actual) {
		t.Logf("Expected context receiver to equal\n%+#v\n\tgot\n%+#v\n", extended, actual)
		t.Fail()
	}
}
