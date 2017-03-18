package object

import (
	"fmt"
	"strings"
	"testing"
)

func TestModuleAncestors(t *testing.T) {
	t.Run("class extending from BasicObject", func(t *testing.T) {
		context := &Class{name: "BasicObjectAsParent", superClass: BASIC_OBJECT_CLASS}

		result := moduleAncestors(context)

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
	t.Run("core class hierarchies", func(t *testing.T) {
		tests := []struct {
			class             RubyClassObject
			expectedAncestors []string
		}{
			{
				BASIC_OBJECT_CLASS,
				[]string{"BasicObject"},
			},
			{
				OBJECT_CLASS,
				[]string{"Object", "Kernel", "BasicObject"},
			},
		}

		for _, testCase := range tests {
			t.Run(testCase.class.Inspect(), func(t *testing.T) {
				result := moduleAncestors(testCase.class)

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
	context := &Class{
		superClass: mixin(BASIC_OBJECT_CLASS, KERNEL_MODULE),
	}

	result := moduleIncludedModules(context)

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
