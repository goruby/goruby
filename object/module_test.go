package object

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestModule_hashKey(t *testing.T) {
	hello1 := &Module{name: "Hello World"}
	hello2 := &Module{name: "Hello World"}
	diff1 := &Module{name: "My name is johnny"}
	diff2 := &Module{name: "My name is johnny"}

	if hello1.hashKey() != hello2.hashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if diff1.hashKey() != diff2.hashKey() {
		t.Errorf("strings with same content have different hash keys")
	}

	if hello1.hashKey() == diff1.hashKey() {
		t.Errorf("strings with different content have same hash keys")
	}
}

func TestModuleAncestors(t *testing.T) {
	t.Run("class extending from BasicObject", func(t *testing.T) {
		context := &callContext{
			receiver: &class{name: "BasicObjectAsParent", superClass: basicObjectClass},
		}

		result, err := moduleAncestors(context)

		checkError(t, err, nil)

		array, ok := result.(*Array)
		if !ok {
			t.Logf("Expected result to be an Array, got %T", result)
			t.FailNow()
		}

		if len(array.Elements) != 2 {
			t.Logf("Expected one ancestor, got %d", len(array.Elements))
			t.Fail()
		}

		expected := fmt.Sprintf("[%s]", strings.Join([]string{"BasicObjectAsParent", "BasicObject"}, ", "))
		actual := fmt.Sprintf("%s", array.Inspect())

		if expected != actual {
			t.Logf("Expected ancestors to equal %s, got %s", expected, actual)
			t.Fail()
		}
	})
	t.Run("class with mixed in modules", func(t *testing.T) {
		context := &callContext{
			receiver: newMixin(
				&class{name: "BasicObjectAsParent", superClass: basicObjectClass},
				kernelModule,
			),
		}

		result, err := moduleAncestors(context)

		checkError(t, err, nil)

		array, ok := result.(*Array)
		if !ok {
			t.Logf("Expected result to be an Array, got %T", result)
			t.FailNow()
		}

		expectedAncestors := 3

		if len(array.Elements) != expectedAncestors {
			t.Logf("Expected %d ancestors, got %d", expectedAncestors, len(array.Elements))
			t.Fail()
		}

		expected := fmt.Sprintf("[%s]", strings.Join([]string{"BasicObjectAsParent", "Kernel", "BasicObject"}, ", "))
		actual := fmt.Sprintf("%s", array.Inspect())

		if expected != actual {
			t.Logf("Expected ancestors to equal %s, got %s", expected, actual)
			t.Fail()
		}
	})
	t.Run("core class hierarchies", func(t *testing.T) {
		tests := []struct {
			class             RubyClassObject
			expectedAncestors []string
		}{
			{
				basicObjectClass,
				[]string{"BasicObject"},
			},
			{
				objectClass,
				[]string{"Object", "Kernel", "BasicObject"},
			},
		}

		for _, testCase := range tests {
			t.Run(testCase.class.Inspect(), func(t *testing.T) {
				context := &callContext{receiver: testCase.class}
				result, err := moduleAncestors(context)

				checkError(t, err, nil)

				array, ok := result.(*Array)
				if !ok {
					t.Logf("Expected result to be an Array, got %T", result)
					t.FailNow()
				}

				if len(array.Elements) != len(testCase.expectedAncestors) {
					t.Logf(
						"Expected ancestor count to equal %d, got %d",
						len(testCase.expectedAncestors),
						len(array.Elements),
					)
					t.Fail()
				}

				expected := fmt.Sprintf("[%s]", strings.Join(testCase.expectedAncestors, ", "))
				actual := fmt.Sprintf("%s", array.Inspect())

				if expected != actual {
					t.Logf("Expected ancestors to equal %s, got %s", expected, actual)
					t.Fail()
				}
			})
		}
	})
}

func TestModuleIncludedModules(t *testing.T) {
	context := &callContext{
		receiver: &class{
			superClass: newMixin(basicObjectClass, kernelModule),
		},
	}

	result, err := moduleIncludedModules(context)

	checkError(t, err, nil)

	array, ok := result.(*Array)
	if !ok {
		t.Logf("Expected result to be an Array, got %T", result)
		t.FailNow()
	}

	expectedModules := []string{"Kernel"}

	if len(array.Elements) != len(expectedModules) {
		t.Logf(
			"Expected %d module(s), got %d",
			len(expectedModules),
			len(array.Elements),
		)
		t.Fail()
	}

	expected := fmt.Sprintf("[%s]", strings.Join(expectedModules, ", "))
	actual := fmt.Sprintf("%s", array.Inspect())

	if expected != actual {
		t.Logf("Expected modules to equal %s, got %s", expected, actual)
		t.Fail()
	}
}

func TestModuleInstanceMethods(t *testing.T) {
	superClassMethods := map[string]RubyMethod{
		"super_foo":         publicMethod(nil),
		"super_bar":         publicMethod(nil),
		"private_super_foo": privateMethod(nil),
	}
	contextMethods := map[string]RubyMethod{
		"foo":         publicMethod(nil),
		"bar":         publicMethod(nil),
		"private_foo": privateMethod(nil),
	}
	t.Run("without superclass", func(t *testing.T) {
		context := &callContext{
			receiver: &class{
				instanceMethods: NewMethodSet(contextMethods),
				superClass:      nil,
			},
		}

		result, err := modulePublicInstanceMethods(context)

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
		context := &callContext{
			receiver: &class{
				instanceMethods: NewMethodSet(contextMethods),
				superClass: &class{
					instanceMethods: NewMethodSet(superClassMethods),
					superClass:      nil,
				},
			},
		}

		result, err := modulePublicInstanceMethods(context)

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
	t.Run("with superclass but boolean arg false", func(t *testing.T) {
		context := &callContext{
			receiver: &class{
				instanceMethods: NewMethodSet(contextMethods),
				superClass: &class{
					instanceMethods: NewMethodSet(superClassMethods),
					superClass:      nil,
				},
			},
		}

		result, err := modulePublicInstanceMethods(context, FALSE)

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
	t.Run("with superclass and boolean arg true", func(t *testing.T) {
		context := &callContext{
			receiver: &class{
				instanceMethods: NewMethodSet(contextMethods),
				superClass: &class{
					instanceMethods: NewMethodSet(superClassMethods),
					superClass:      nil,
				},
			},
		}

		result, err := modulePublicInstanceMethods(context, TRUE)

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
}

func TestModuleAppendFeatures(t *testing.T) {
	t.Run("argument validation", func(t *testing.T) {
		context := &callContext{
			receiver: &Self{
				RubyObject: &class{
					superClass:  objectClass,
					Environment: NewEnvironment(),
				},
				Name: "X",
			},
		}

		_, err := moduleAppendFeatures(context, &Integer{Value: 2})

		checkError(t, err, NewWrongArgumentTypeError(&Module{}, &Integer{}))
	})
	t.Run("return value", func(t *testing.T) {
		context := &callContext{
			receiver: &Self{
				RubyObject: &class{
					superClass:  objectClass,
					Environment: NewEnvironment(),
				},
				Name: "X",
			},
		}
		module := NewModule("foo", nil)

		result, err := moduleAppendFeatures(context, module)

		checkError(t, err, nil)

		checkResult(t, result, module)
	})
	t.Run("add constants", func(t *testing.T) {
		outer := NewEnvironment()
		outer.Set("A", &String{Value: "foo"})
		outer.Set("C", &Symbol{Value: "bar"})
		env := NewEnclosedEnvironment(outer)
		env.Set("A", &Integer{Value: 4})
		env.Set("B", &Integer{Value: 6})
		receiver := NewModule("X", env)
		context := &callContext{
			receiver: &Self{RubyObject: receiver, Name: "X"},
		}

		module := NewModule("foo", nil)

		moduleAppendFeatures(context, module)

		a, ok := module.Get("A")
		if !ok {
			t.Logf("Expected constant A to be within module env")
			t.Fail()
		}

		checkResult(t, &Integer{Value: 4}, a)

		b, ok := module.Get("B")
		if !ok {
			t.Logf("Expected constant B to be within module env")
			t.Fail()
		}

		checkResult(t, &Integer{Value: 6}, b)

		c, ok := module.Get("C")
		if !ok {
			t.Logf("Expected constant C to be within module env")
			t.Fail()
		}

		checkResult(t, &Symbol{Value: "bar"}, c)
	})
	t.Run("add instance variables", func(t *testing.T) {
		outer := NewEnvironment()
		outer.Set("@foo", &String{Value: "foo"})
		outer.Set("@qux", &Symbol{Value: "bar"})
		env := NewEnclosedEnvironment(outer)
		env.Set("@foo", &Integer{Value: 4})
		env.Set("@bar", &Integer{Value: 6})
		receiver := NewModule("X", env)
		context := &callContext{
			receiver: &Self{RubyObject: receiver, Name: "X"},
		}

		module := NewModule("foo", nil)

		moduleAppendFeatures(context, module)

		a, ok := module.Get("@foo")
		if !ok {
			t.Logf("Expected module variable foo to be within module env")
			t.Fail()
		}

		checkResult(t, &Integer{Value: 4}, a)

		b, ok := module.Get("@bar")
		if !ok {
			t.Logf("Expected module variable bar to be within module env")
			t.Fail()
		}

		checkResult(t, &Integer{Value: 6}, b)

		c, ok := module.Get("@qux")
		if !ok {
			t.Logf("Expected module variable qux to be within module env")
			t.Fail()
		}

		checkResult(t, &Symbol{Value: "bar"}, c)
	})
	t.Run("does not add local variables", func(t *testing.T) {
		outer := NewEnvironment()
		outer.Set("foo", &String{Value: "foo"})
		outer.Set("qux", &Symbol{Value: "bar"})
		env := NewEnclosedEnvironment(outer)
		env.Set("foo", &Integer{Value: 4})
		env.Set("bar", &Integer{Value: 6})
		receiver := NewModule("X", env)
		context := &callContext{
			receiver: &Self{RubyObject: receiver, Name: "X"},
		}

		module := NewModule("foo", nil)

		moduleAppendFeatures(context, module)

		_, ok := module.Get("foo")
		if ok {
			t.Logf("Expected local variable foo not to be within module env")
			t.Fail()
		}

		_, ok = module.Get("bar")
		if ok {
			t.Logf("Expected local variable bar not to be within module env")
			t.Fail()
		}

		_, ok = module.Get("qux")
		if ok {
			t.Logf("Expected local variable qux not to be within module env")
			t.Fail()
		}
	})
	t.Run("add methods", func(t *testing.T) {
		methods := map[string]RubyMethod{
			"a": nil,
			"b": nil,
			"c": nil,
		}
		receiver := newModule("X", methods, nil)
		context := &callContext{
			receiver: &Self{RubyObject: receiver, Name: "X"},
		}

		module := NewModule("foo", nil)

		moduleAppendFeatures(context, module)

		moduleMethods := module.class.Methods()

		_, ok := moduleMethods.Get("a")
		if !ok {
			t.Logf("Expected method a to be within module class")
			t.Fail()
		}

		_, ok = moduleMethods.Get("b")
		if !ok {
			t.Logf("Expected method b to be within module class")
			t.Fail()
		}

		_, ok = moduleMethods.Get("c")
		if !ok {
			t.Logf("Expected method c to be within module class")
			t.Fail()
		}
	})
}

func TestModuleInclude(t *testing.T) {
	t.Run("argument validation", func(t *testing.T) {
		context := &callContext{
			receiver: &Self{
				RubyObject: NewModule("X", nil),
				Name:       "X",
			},
		}

		_, err := moduleInclude(context)

		checkError(t, err, NewWrongNumberOfArgumentsError(1, 0))

		_, err = moduleInclude(context, &Integer{Value: 2}, NewModule("X", nil))

		checkError(t, err, NewWrongArgumentTypeError(&Module{}, &Integer{}))

		_, err = moduleInclude(context, NewModule("X", nil), &Symbol{Value: "d"})

		checkError(t, err, NewWrongArgumentTypeError(&Module{}, &Symbol{}))
	})
	t.Run("return value", func(t *testing.T) {
		context := &callContext{
			receiver: &Self{
				RubyObject: NewModule("X", nil),
				Name:       "X",
			},
		}
		module := NewModule("foo", nil)

		result, err := moduleInclude(context, module)

		checkError(t, err, nil)

		checkResult(t, result, context.receiver)
	})
	t.Run("calls #append_features on each parameter in reverse order", func(t *testing.T) {
		calls := objects{}
		mockMethod := withArity(1, privateMethod(func(ctx CallContext, args ...RubyObject) (RubyObject, error) {
			mod, ok := args[0].(*Module)
			if !ok {
				return nil, NewWrongArgumentTypeError(mod, args[0])
			}
			calls = append(calls, mod)
			return mod, nil
		}))
		context := &callContext{
			receiver: &Self{
				RubyObject: newModule("X", map[string]RubyMethod{"append_features": mockMethod}, nil),
				Name:       "X",
			},
		}
		module1 := NewModule("foo", nil)
		module2 := NewModule("bar", nil)
		module3 := NewModule("qux", nil)

		_, err := moduleInclude(context, module1, module2, module3)

		checkError(t, err, nil)

		if len(calls) != 3 {
			t.Logf("Expected #append_features to be called 3 times, got %d\n", len(calls))
			t.Fail()
		}

		expected := objects{module3, module2, module1}
		if !reflect.DeepEqual(expected, calls) {
			t.Logf("Expected #append_features to be called with\n%s\n\tgot\n%s\n", expected, calls)
			t.Fail()
		}
	})
}
