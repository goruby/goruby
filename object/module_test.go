package object

import (
	"fmt"
	"strings"
	"testing"
)

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
			receiver: mixin(
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
			superClass: mixin(basicObjectClass, kernelModule),
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
