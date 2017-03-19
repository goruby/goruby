package object

var symbolClass RubyClassObject = newClass("Symbol", objectClass, symbolMethods, symbolClassMethods)

func init() {
	classes.Set("Symbol", symbolClass)
}

type Symbol struct {
	Value string
}

func (s *Symbol) Inspect() string  { return ":" + s.Value }
func (s *Symbol) Type() Type       { return SYMBOL_OBJ }
func (s *Symbol) Class() RubyClass { return symbolClass }

var symbolClassMethods = map[string]RubyMethod{}

var symbolMethods = map[string]RubyMethod{}
