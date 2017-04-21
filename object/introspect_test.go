package object

import (
	"reflect"
	"sort"
	"testing"
)

func TestGetMethods(t *testing.T) {
	superClassMethods := map[string]RubyMethod{
		"super_foo":           publicMethod(nil),
		"super_bar":           publicMethod(nil),
		"protected_super_foo": protectedMethod(nil),
		"private_super_foo":   privateMethod(nil),
	}
	contextMethods := map[string]RubyMethod{
		"foo":           publicMethod(nil),
		"bar":           publicMethod(nil),
		"protected_foo": protectedMethod(nil),
		"private_foo":   privateMethod(nil),
	}
	classWithoutSuperclass := &class{
		instanceMethods: NewMethodSet(contextMethods),
		superClass:      nil,
	}
	classWithSuperclass := &class{
		instanceMethods: NewMethodSet(contextMethods),
		superClass: &class{
			instanceMethods: NewMethodSet(superClassMethods),
			superClass:      nil,
		},
	}

	tests := []struct {
		name                 string
		class                RubyClass
		visibility           MethodVisibility
		addSuperclassMethods bool
		expectedMethods      []string
	}{
		{
			"no superclass public methods add super methods",
			classWithoutSuperclass,
			PUBLIC_METHOD,
			true,
			[]string{":foo", ":bar"},
		},
		{
			"no superclass public methods add no super methods",
			classWithoutSuperclass,
			PUBLIC_METHOD,
			false,
			[]string{":foo", ":bar"},
		},
		{
			"no superclass protected methods add super methods",
			classWithoutSuperclass,
			PROTECTED_METHOD,
			true,
			[]string{":protected_foo"},
		},
		{
			"no superclass protected methods add no super methods",
			classWithoutSuperclass,
			PROTECTED_METHOD,
			false,
			[]string{":protected_foo"},
		},
		{
			"no superclass private methods add super methods",
			classWithoutSuperclass,
			PRIVATE_METHOD,
			true,
			[]string{":private_foo"},
		},
		{
			"no superclass private methods add no super methods",
			classWithoutSuperclass,
			PRIVATE_METHOD,
			false,
			[]string{":private_foo"},
		},
		{
			"with superclass public methods add super methods",
			classWithSuperclass,
			PUBLIC_METHOD,
			true,
			[]string{":foo", ":bar", ":super_foo", ":super_bar"},
		},
		{
			"with superclass public methods add with super methods",
			classWithSuperclass,
			PUBLIC_METHOD,
			false,
			[]string{":foo", ":bar"},
		},
		{
			"with superclass protected methods add super methods",
			classWithSuperclass,
			PROTECTED_METHOD,
			true,
			[]string{":protected_foo", ":protected_super_foo"},
		},
		{
			"with superclass protected methods add with super methods",
			classWithSuperclass,
			PROTECTED_METHOD,
			false,
			[]string{":protected_foo"},
		},
		{
			"with superclass private methods add super methods",
			classWithSuperclass,
			PRIVATE_METHOD,
			true,
			[]string{":private_foo", ":private_super_foo"},
		},
		{
			"with superclass private methods add with super methods",
			classWithSuperclass,
			PRIVATE_METHOD,
			false,
			[]string{":private_foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMethods(tt.class, tt.visibility, tt.addSuperclassMethods)

			var methods []string
			for i, elem := range result.Elements {
				sym, ok := elem.(*Symbol)
				if !ok {
					t.Logf("Expected all elements to be symbols, got %T at index %d", elem, i)
					t.Fail()
				} else {
					methods = append(methods, sym.Inspect())
				}
			}

			expectedLen := len(tt.expectedMethods)

			if len(result.Elements) != expectedLen {
				t.Logf("Expected %d items, got %d", expectedLen, len(result.Elements))
				t.Fail()
			}

			sort.Strings(tt.expectedMethods)
			sort.Strings(methods)

			if !reflect.DeepEqual(tt.expectedMethods, methods) {
				t.Logf("Expected methods to equal\n%s\n\tgot\n%s\n", tt.expectedMethods, methods)
				t.Fail()
			}
		})
	}
}
