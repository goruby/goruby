package ast

import "fmt"

// A Visitor's Visit method is invoked for each node encountered by Walk.
// If the result visitor w is not nil, Walk visits each of the children
// of node with the visitor w, followed by a call of w.Visit(nil).
type Visitor interface {
	Visit(node Node) (w Visitor)
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
		*Comment:
		// nothing to do

	case *BlockExpression:
		walkParameterList(v, n.Parameters)
		Walk(v, n.Body)

	case *FunctionLiteral:
		Walk(v, n.Name)
		walkParameterList(v, n.Parameters)
		Walk(v, n.Body)

	case *FunctionParameter:
		Walk(v, n.Name)
		Walk(v, n.Default)

	case *IndexExpression:
		Walk(v, n.Left)
		Walk(v, n.Index)

	case *ContextCallExpression:
		Walk(v, n.Context)
		Walk(v, n.Function)
		walkExprList(v, n.Arguments)

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

	// Types
	case *ArrayLiteral:
		walkExprList(v, n.Elements)

	case *ExpressionStatement:
		Walk(v, n.Expression)

	case *InstanceVariable:
		Walk(v, n.Name)

	case *Assignment:
		Walk(v, n.Left)
		Walk(v, n.Right)

	case *VariableAssignment:
		Walk(v, n.Name)
		Walk(v, n.Value)

	case *GlobalAssignment:
		Walk(v, n.Name)
		Walk(v, n.Value)

	case *ReturnStatement:
		Walk(v, n.ReturnValue)

	case *BlockStatement:
		walkStmtList(v, n.Statements)

	case *ScopedIdentifier:
		Walk(v, n.Outer)
		Walk(v, n.Inner)

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
