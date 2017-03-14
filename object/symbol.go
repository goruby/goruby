package object

type Symbol struct {
	Value string
}

func (s *Symbol) Inspect() string                                 { return ":" + s.Value }
func (s *Symbol) Type() ObjectType                                { return SYMBOL_OBJ }
func (s *Symbol) Send(name string, args ...RubyObject) RubyObject { return NIL }
