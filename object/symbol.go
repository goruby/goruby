package object

var symbolClass RubyClassObject = NewClass("Symbol", objectClass, symbolMethods, symbolClassMethods)

func init() {
	classes.Set("Symbol", symbolClass)
}

// A Symbol represents a symbol in Ruby
type Symbol struct {
	Value string
}

// Inspect returns the value of the symbol
func (s *Symbol) Inspect() string { return ":" + s.Value }

// Type returns SYMBOL_OBJ
func (s *Symbol) Type() Type { return SYMBOL_OBJ }

// Class returns symbolClass
func (s *Symbol) Class() RubyClass { return symbolClass }

var symbolClassMethods = map[string]RubyMethod{}

var symbolMethods = map[string]RubyMethod{}
