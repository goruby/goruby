package object

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/goruby/goruby/ast"
)

func TestNewClass(t *testing.T) {
	t.Run("full constructor use", func(t *testing.T) {
		superClass := arrayClass
		env := NewEnvironment()

		classObject := NewClass("Abc", superClass, env)

		class, ok := classObject.(*class)
		if !ok {
			t.Logf("Expected returned object to be a class, got %T", classObject)
			t.Fail()
		}

		expectedInstanceMethodSet := NewMethodSet(map[string]RubyMethod{})
		if !reflect.DeepEqual(expectedInstanceMethodSet, class.instanceMethods) {
			t.Logf(
				"Expected class.instanceMethods to equal\n%+#v\n\tgot\n%+#v\n",
				expectedInstanceMethodSet,
				class.instanceMethods,
			)
			t.Fail()
		}

		if !reflect.DeepEqual(superClass, class.superClass) {
			t.Logf(
				"Expected class.superClass to equal\n%+#v\n\tgot\n%+#v\n",
				superClass,
				class.superClass,
			)
			t.Fail()
		}

		expectedClassMethodSet := NewMethodSet(map[string]RubyMethod{})
		actualClassMethods := class.class.Methods()
		if !reflect.DeepEqual(expectedClassMethodSet, actualClassMethods) {
			t.Logf(
				"Expected class.class.Methods to equal\n%+#v\n\tgot\n%+#v\n",
				expectedClassMethodSet,
				actualClassMethods,
			)
			t.Fail()
		}
	})
	t.Run("missing env", func(t *testing.T) {
		superClass := arrayClass

		classObject := NewClass("Abc", superClass, nil)

		class, ok := classObject.(*class)
		if !ok {
			t.Logf("Expected returned object to be a class, got %T", classObject)
			t.Fail()
		}

		expectedInstanceMethodSet := NewMethodSet(map[string]RubyMethod{})
		if !reflect.DeepEqual(expectedInstanceMethodSet, class.instanceMethods) {
			t.Logf(
				"Expected class.instanceMethods to equal\n%+#v\n\tgot\n%+#v\n",
				expectedInstanceMethodSet,
				class.instanceMethods,
			)
			t.Fail()
		}
	})
}

func TestClassInspect(t *testing.T) {
	context := &class{name: "Foo"}

	actual := context.Inspect()

	expected := "Foo"

	if expected != actual {
		t.Logf("Expected Inspect to equal %q, got %q", expected, actual)
		t.Fail()
	}
}

func TestClass_hashKey(t *testing.T) {
	hello1 := &class{name: "Hello World"}
	hello2 := &class{name: "Hello World"}
	diff1 := &class{name: "My name is johnny"}
	diff2 := &class{name: "My name is johnny"}

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

func TestClassInstanceInspect(t *testing.T) {
	context := &classInstance{class: &class{name: "Foo"}}

	actual := context.Inspect()

	expected := fmt.Sprintf("#<Foo:%p>", context)

	if expected != actual {
		t.Logf("Expected Inspect to equal %q, got %q", expected, actual)
		t.Fail()
	}
}

func TestClassClass(t *testing.T) {
	clazz := &class{}

	context := &class{class: clazz}

	actual := context.Class()

	expected := clazz

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Class to equal %v, got %v", expected, actual)
		t.Fail()
	}
}

func TestClassSuperClass(t *testing.T) {
	context := &class{superClass: basicObjectClass}

	actual := context.SuperClass()

	expected := basicObjectClass

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected SuperClass to equal %+#v, got %+#v", expected, actual)
		t.Fail()
	}
}

func TestClassMethods(t *testing.T) {
	contextMethods := map[string]RubyMethod{
		"a_method": nil,
	}

	context := &class{instanceMethods: NewMethodSet(contextMethods)}

	actual := context.Methods()

	expected := NewMethodSet(contextMethods)

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected Methods to equal %v, got %v", expected, actual)
		t.Fail()
	}
}

func TestClassType(t *testing.T) {
	context := &class{}

	actual := context.Type()

	expected := CLASS_OBJ

	if expected != actual {
		t.Logf("Expected Type to equal %v, got %v", expected, actual)
		t.Fail()
	}
}

func TestClassSuperclass(t *testing.T) {
	t.Run("anything else than BasicObject", func(t *testing.T) {
		context := &callContext{receiver: &class{superClass: objectClass}}

		result, err := classSuperclass(context)
		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		_, ok := result.(*class)
		if !ok {
			t.Logf("Expected Class object, got %T\n", result)
			t.Fail()
		}
	})
	t.Run("BasicObject", func(t *testing.T) {
		context := &callContext{receiver: basicObjectClass}

		result, err := classSuperclass(context)
		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		_, ok := result.(*nilObject)
		if !ok {
			t.Logf("Expected Nil object, got %T\n", result)
			t.Fail()
		}
	})
	t.Run("Eigenclass", func(t *testing.T) {
		context := &callContext{receiver: &class{superClass: newEigenclass(objectClass, nil)}}

		result, err := classSuperclass(context)
		if err != nil {
			t.Logf("Expected no error, got %T:%v\n", err, err)
			t.Fail()
		}

		_, ok := result.(*eigenclass)
		if !ok {
			t.Logf("Expected eigenClass object, got %T\n", result)
			t.Fail()
		}
	})
}

func TestClassNew(t *testing.T) {
	var initializeContext CallContext
	var initializeArgs []RubyObject
	callCount := 0
	initializeStub := func(context CallContext, args ...RubyObject) (RubyObject, error) {
		callCount++
		initializeContext = context
		initializeArgs = args
		return context.Receiver(), nil
	}
	fooClass := newClass(
		"Foo",
		objectClass,
		map[string]RubyMethod{"initialize": privateMethod(initializeStub)},
		nil,
		func(c RubyClassObject) RubyObject { return &classInstance{c} })
	env := NewEnvironment()
	env.Set("Class", classClass)
	env.Set("Foo", fooClass)
	context := &callContext{
		receiver: fooClass,
		env:      env,
		eval:     func(ast.Node, Environment) (RubyObject, error) { return nil, nil },
	}

	args := []RubyObject{&String{"foo"}, &Symbol{"bar"}, &Integer{7}}

	result, err := classNew(context, args...)
	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	instance, ok := result.(*classInstance)
	if !ok {
		t.Logf("Expected classInstance object, got %T\n", result)
		t.Fail()
	}

	if !reflect.DeepEqual(instance.class, context.receiver) {
		t.Logf(
			"Expected instance.class to equal\n%+#v\n\tgot\n%+#v\n",
			instance.class,
			context.receiver,
		)
		t.Fail()
	}
	if callCount != 1 {
		t.Logf("Expected `initialize` of Foo class to be called once, was %d", callCount)
		t.Fail()
	}

	if !reflect.DeepEqual(args, initializeArgs) {
		t.Logf(
			"Expected initialize args to equal\n%+#v\n\tgot\n%+#v\n",
			args,
			initializeArgs,
		)
		t.Fail()
	}

	expectedReceiver := &Self{RubyObject: instance, Name: "Foo"}
	if !reflect.DeepEqual(expectedReceiver, initializeContext.Receiver()) {
		t.Logf(
			"Expected initialize context receiver to equal\n%+#v\n\tgot\n%+#v\n",
			expectedReceiver,
			initializeContext.Receiver(),
		)
		t.Fail()
	}

	if !reflect.DeepEqual(context.Env(), initializeContext.Env()) {
		t.Logf(
			"Expected initialize context env to equal\n%+#v\n\tgot\n%+#v\n",
			context.Env(),
			initializeContext.Env(),
		)
		t.Fail()
	}
}

func TestClassInitialize(t *testing.T) {
	env := NewEnvironment()
	context := &callContext{
		receiver: &classInstance{class: &class{name: "Foo"}},
		env:      env,
	}

	result, err := classInitialize(context, &String{"foo"}, &Symbol{"bar"}, &Integer{7})
	if err != nil {
		t.Logf("Expected no error, got %T:%v\n", err, err)
		t.Fail()
	}

	instance, ok := result.(*classInstance)
	if !ok {
		t.Logf("Expected classInstance object, got %T\n", result)
		t.Fail()
	}

	if !reflect.DeepEqual(instance, context.Receiver()) {
		t.Logf(
			"Expected instance to equal\n%+#v\n\tgot\n%+#v\n",
			instance,
			context.Receiver(),
		)
		t.Fail()
	}
}
