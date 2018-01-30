package evaluator

import (
	"fmt"
	"go/token"

	"github.com/goruby/goruby/ast"
)

func wrapStack(err error, node, root ast.Node, file *token.File) error {
	return &stackError{
		underlyingError: err,
		node:            node,
		root:            root,
		file:            file,
	}
}

type stackError struct {
	underlyingError error
	node            ast.Node
	root            ast.Node
	file            *token.File
}

func (s *stackError) Cause() error {
	return s.underlyingError
}

func (s *stackError) Error() string {
	return s.underlyingError.Error()
}

func (s *stackError) buildStack() {
	path, ok := ast.Path(s.root, s.node)
	if !ok {
		return
	}
	for elem := path.Front(); elem != nil; elem = elem.Next() {
		n := elem.Value.(ast.Node)
		pos := s.file.Position(token.Pos(n.Pos()))
		fmt.Println(pos.String())
		fmt.Printf("Node: %#v\n", n)
		// ast.Print(nil, n)
	}
}
