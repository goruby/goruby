package ast

import (
	"container/list"
	"reflect"
	"testing"
)

func Test_Parent(t *testing.T) {
	t.Run("parent found", func(t *testing.T) {
		child := &Assignment{
			Left:  &Identifier{Value: "x"},
			Right: &IntegerLiteral{Value: 2},
		}
		parent := &ExpressionStatement{Expression: child}
		root := &Program{
			Statements: []Statement{
				&ExpressionStatement{
					Expression: &IntegerLiteral{Value: 3},
				},
				parent,
			},
		}

		p, ok := Parent(root, child)

		if !ok {
			t.Logf("Expected child not to be contained within root")
			t.Fail()
		}

		if !reflect.DeepEqual(parent, p) {
			t.Logf("Expected parent to equal\n%+#v\n\tgot\n%+#v\n", parent, p)
			t.Fail()
		}
	})
	t.Run("parent is root", func(t *testing.T) {
		root := &Program{
			Statements: []Statement{
				&ExpressionStatement{
					Expression: &IntegerLiteral{Value: 3},
				},
			},
		}

		_, ok := Parent(root, root)

		if ok {
			t.Logf("Expected bool to return false")
			t.Fail()
		}
	})
	t.Run("child not found", func(t *testing.T) {
		root := &Program{
			Statements: []Statement{
				&ExpressionStatement{
					Expression: &IntegerLiteral{Value: 3},
				},
				&ExpressionStatement{
					Expression: &Assignment{
						Left:  &Identifier{Value: "x"},
						Right: &IntegerLiteral{Value: 2},
					},
				},
			},
		}

		_, ok := Parent(root, &IntegerLiteral{Value: 3})

		if ok {
			t.Logf("Expected child to be contained within root, was not")
			t.Fail()
		}
	})
}

func Test_Path(t *testing.T) {
	t.Run("child found", func(t *testing.T) {
		child := &Assignment{
			Left:  &Identifier{Value: "x"},
			Right: &IntegerLiteral{Value: 2},
		}
		root := &Program{
			Statements: []Statement{
				&ExpressionStatement{
					Expression: &IntegerLiteral{Value: 3},
				},
				&ExpressionStatement{
					Expression: child,
				},
			},
		}

		path, ok := Path(root, child)

		if !ok {
			t.Logf("Expected child to be contained within root, was not")
			t.Fail()
		}

		expected := list.New()
		expected.PushBack(root)
		expected.PushBack(root.Statements[1])
		expected.PushBack(child)

		if !reflect.DeepEqual(expected, path) {
			t.Logf("Expected AST path to equal\n%+#v\n\tgot\n%+#v\n", expected, path)
			t.Fail()
		}
	})
	t.Run("child not found", func(t *testing.T) {
		root := &Program{
			Statements: []Statement{
				&ExpressionStatement{
					Expression: &IntegerLiteral{Value: 3},
				},
				&ExpressionStatement{
					Expression: &Assignment{
						Left:  &Identifier{Value: "x"},
						Right: &IntegerLiteral{Value: 2},
					},
				},
			},
		}

		_, ok := Path(root, &IntegerLiteral{Value: 3})

		if ok {
			t.Logf("Expected child not to be contained within root")
			t.Fail()
		}
	})
}

func Test_treeToList(t *testing.T) {
	root := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Expression: &IntegerLiteral{Value: 3},
			},
			&ExpressionStatement{
				Expression: &Assignment{
					Left:  &Identifier{Value: "x"},
					Right: &IntegerLiteral{Value: 2},
				},
			},
		},
	}

	actual := treeToLinkedList(root)

	expected := list.New()
	expected.PushBack(root)
	expected.PushBack(&ExpressionStatement{
		Expression: &IntegerLiteral{Value: 3},
	})
	expected.PushBack(&IntegerLiteral{Value: 3})
	expected.PushBack(&ExpressionStatement{
		Expression: &Assignment{
			Left:  &Identifier{Value: "x"},
			Right: &IntegerLiteral{Value: 2},
		},
	})
	expected.PushBack(&Assignment{
		Left:  &Identifier{Value: "x"},
		Right: &IntegerLiteral{Value: 2},
	})
	expected.PushBack(&Identifier{Value: "x"})
	expected.PushBack(&IntegerLiteral{Value: 2})

	if !reflect.DeepEqual(expected, actual) {
		t.Logf("Expected list to equal\n%+#v\n\tgot\n%+#v\n", expected, actual)
		t.Fail()
	}
}

func Test_Contains(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		needle := &IntegerLiteral{Value: 1}
		statement := &ExpressionStatement{
			Expression: needle,
		}

		ok := Contains(statement, needle)

		if !ok {
			t.Logf("Expected node to be within statement, was not")
			t.Fail()
		}
	})
	t.Run("not the same pointer", func(t *testing.T) {
		needle := &IntegerLiteral{Value: 1}
		statement := &ExpressionStatement{
			Expression: needle,
		}

		ok := Contains(statement, &IntegerLiteral{Value: 1})

		if ok {
			t.Logf("Expected node not to be within statement")
			t.Fail()
		}
	})
	t.Run("not found", func(t *testing.T) {
		needle := &IntegerLiteral{Value: 3}
		statement := &ExpressionStatement{
			Expression: &Assignment{
				Left:  &Identifier{Value: "foo"},
				Right: &StringLiteral{Value: "bar"},
			},
		}

		ok := Contains(statement, needle)

		if ok {
			t.Logf("Expected node not to be within statement")
			t.Fail()
		}
	})
}
