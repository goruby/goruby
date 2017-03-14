package object

type String struct {
	Value string
}

func (s *String) Inspect() string                                 { return s.Value }
func (s *String) Type() ObjectType                                { return STRING_OBJ }
func (s *String) Send(name string, args ...RubyObject) RubyObject { return NIL }
