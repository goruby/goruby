package ast

import (
	"container/list"
	"fmt"
)

// A Visitor's Visit method is invoked for each node encountered by Walk.  If
// the result visitor w is not nil, Walk visits each of the children of node
// with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}

// The VisitorFunc type is an adapter to allow the use of ordinary functions as
// AST Visitors. If f is a function with the appropriate signature,
// VisitorFunc(f) is a Visitor that calls f.
type VisitorFunc func(Node) Visitor

// Visit calls f(n)
func (f VisitorFunc) Visit(n Node) Visitor {
	return f(n)
}

// Parent returns the parent node of child. If child is not found within root,
// or child does not have a parent, i.e. equals root, the bool will be false
func Parent(root, child Node) (Node, bool) {
	if root == child {
		return nil, false
	}
	if !Contains(root, child) {
		return nil, false
	}
	path, ok := Path(root, child)
	if !ok {
		return nil, false
	}
	parentElement := path.Back().Prev()
	if parentElement == nil {
		return nil, false
	}
	parent, ok := parentElement.Value.(Node)
	if !ok {
		return nil, false
	}
	return parent, true
}

// Path returns the path from the root of the AST till the child as a doubly
// linked list. If child is not found within root, the bool will be false and
// the list nil.
func Path(root, child Node) (*list.List, bool) {
	if !Contains(root, child) {
		return nil, false
	}
	childTree := list.New()
	l := treeToLinkedList(root)
	for e := l.Front(); e != nil; e = e.Next() {
		n, ok := e.Value.(Node)
		if !ok {
			continue
		}
		if Contains(n, child) {
			childTree.PushBack(n)
		}
	}
	return childTree, true
}

func treeToLinkedList(node Node) *list.List {
	list := list.New()
	for n := range WalkEmit(node) {
		if n != nil {
			list.PushBack(n)
		}
	}
	return list
}

// Contains reports whether root contains child or not. It matches child via
// pointer equality.
func Contains(root Node, child Node) bool {
	var contains bool
	filter := func(n Node) bool {
		if n == child {
			contains = true
		}
		return !contains
	}
	Inspect(root, filter)
	return contains
}

// WalkEmit traverses node in depth-first order and emits each visited node
// into the channel
func WalkEmit(root Node) <-chan Node {
	out := make(chan Node)
	var visitor Visitor
	visitor = VisitorFunc(func(n Node) Visitor {
		out <- n
		return visitor
	})
	go func() {
		defer close(out)
		Walk(visitor, root)
	}()
	return out
}

// Helper functions for common node lists. They may be empty.

func walkParameterList(v Visitor, list []*FunctionParameter) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkIdentifierList(v Visitor, list []*Identifier) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkExprList(v Visitor, list []Expression) {
	for _, x := range list {
		Walk(v, x)
	}
}

func walkStmtList(v Visitor, list []Statement) {
	for _, x := range list {
		Walk(v, x)
	}
}

// Walk traverses an AST in depth-first order: It starts by calling
// v.Visit(node); node must not be nil. If the visitor w returned by
// v.Visit(node) is not nil, Walk is invoked recursively with visitor
// w for each of the non-nil children of node, followed by a call of
// w.Visit(nil).
//
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := node.(type) {
	// Expressions
	case *Identifier,
		*Global,
		*IntegerLiteral,
		*StringLiteral,
		*SymbolLiteral,
		*Boolean,
		*Nil,
		*Self,
		*BlockCapture,
		*Keyword__FILE__,
		*Comment:
		// nothing to do

	case *BlockExpression:
		walkParameterList(v, n.Parameters)
		Walk(v, n.Body)

	case *ExceptionHandlingBlock:
		Walk(v, n.TryBody)
		for _, r := range n.Rescues {
			Walk(v, r)
		}

	case *RescueBlock:
		if len(n.ExceptionClasses) != 0 {
			for _, e := range n.ExceptionClasses {
				Walk(v, e)
			}
		}
		if n.Exception != nil {
			Walk(v, n.Exception)
		}
		Walk(v, n.Body)

	case *FunctionLiteral:
		if n.Receiver != nil {
			Walk(v, n.Receiver)
		}
		Walk(v, n.Name)
		walkParameterList(v, n.Parameters)
		if n.CapturedBlock != nil {
			Walk(v, n.CapturedBlock)
		}
		Walk(v, n.Body)
		for _, r := range n.Rescues {
			Walk(v, r)
		}

	case *FunctionParameter:
		Walk(v, n.Name)
		Walk(v, n.Default)

	case *IndexExpression:
		Walk(v, n.Left)
		Walk(v, n.Index)
		if n.Length != nil {
			Walk(v, n.Length)
		}

	case *ContextCallExpression:
		Walk(v, n.Context)
		Walk(v, n.Function)
		walkExprList(v, n.Arguments)
		// TODO: examine why it is not working
		// Walk(v, n.Block)

	case *ModuleExpression:
		Walk(v, n.Name)
		Walk(v, n.Body)

	case *ClassExpression:
		Walk(v, n.Name)
		if n.SuperClass != nil {
			Walk(v, n.SuperClass)
		}
		Walk(v, n.Body)

	case *YieldExpression:
		walkExprList(v, n.Arguments)

	case *PrefixExpression:
		Walk(v, n.Right)

	case *InfixExpression:
		Walk(v, n.Left)
		Walk(v, n.Right)

	case *MultiAssignment:
		walkIdentifierList(v, n.Variables)
		walkExprList(v, n.Values)

	case ExpressionList:
		walkExprList(v, n)

	// Types
	case *ArrayLiteral:
		walkExprList(v, n.Elements)

	case *HashLiteral:
		for k, val := range n.Map {
			Walk(v, k)
			Walk(v, val)
		}

	case *ExpressionStatement:
		Walk(v, n.Expression)

	case *InstanceVariable:
		Walk(v, n.Name)

	case *Assignment:
		Walk(v, n.Left)
		Walk(v, n.Right)

	case *ReturnStatement:
		Walk(v, n.ReturnValue)

	case *BlockStatement:
		walkStmtList(v, n.Statements)

	case *ScopedIdentifier:
		Walk(v, n.Outer)
		Walk(v, n.Inner)

	case *ConditionalExpression:
		Walk(v, n.Condition)
		Walk(v, n.Consequence)
		if n.Alternative != nil {
			Walk(v, n.Alternative)
		}

	case *LoopExpression:
		Walk(v, n.Condition)
		Walk(v, n.Block)

	// Program
	case *Program:
		walkStmtList(v, n.Statements)

	case nil:
		// nothing to do

	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", n))
	}

	v.Visit(nil)
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Inspect traverses an AST in depth-first order: It starts by calling
// f(node); node must not be nil. If f returns true, Inspect invokes f
// recursively for each of the non-nil children of node, followed by a
// call of f(nil).
//
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}
