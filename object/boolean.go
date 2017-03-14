package object

import "fmt"

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string                                 { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType                                { return BOOLEAN_OBJ }
func (b *Boolean) Send(name string, args ...RubyObject) RubyObject { return NIL }
