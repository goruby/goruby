package ast

import "fmt"

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
}

// Helper functions for common node lists. They may be empty.

func walkIdentList(v Visitor, list []*Identifier) {
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
	case *Identifier, *IntegerLiteral, *StringLiteral, *SymbolLiteral, *Boolean:
		// nothing to do

	case *FunctionLiteral:
		Walk(v, n.Name)
		walkIdentList(v, n.Parameters)
		Walk(v, n.Body)

	case *IndexExpression:
		Walk(v, n.Left)
		Walk(v, n.Index)

	case *ContextCallExpression:
		Walk(v, n.Context)
		Walk(v, n.Function)
		walkExprList(v, n.Arguments)

	case *PrefixExpression:
		Walk(v, n.Right)

	case *InfixExpression:
		Walk(v, n.Left)
		Walk(v, n.Right)

	// Types
	case *ArrayLiteral:
		walkExprList(v, n.Elements)

	case *ExpressionStatement:
		Walk(v, n.Expression)

	case *VariableAssignment:
		Walk(v, n.Name)
		Walk(v, n.Value)

	case *ReturnStatement:
		Walk(v, n.ReturnValue)

	case *BlockStatement:
		walkStmtList(v, n.Statements)

	case *IfExpression:
		Walk(v, n.Condition)
		Walk(v, n.Consequence)
		if n.Alternative != nil {
			Walk(v, n.Alternative)
		}

	// Program
	case *Program:
		walkStmtList(v, n.Statements)

	case nil:

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
