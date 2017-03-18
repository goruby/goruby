package object

var SYMBOL_CLASS RubyClassObject = NewClass("Symbol", OBJECT_CLASS, symbolMethods, symbolClassMethods)

type Symbol struct {
	Value string
}

func (s *Symbol) Inspect() string  { return ":" + s.Value }
func (s *Symbol) Type() ObjectType { return SYMBOL_OBJ }
func (s *Symbol) Class() RubyClass { return SYMBOL_CLASS }

var symbolClassMethods = map[string]RubyMethod{}

var symbolMethods = map[string]RubyMethod{}
