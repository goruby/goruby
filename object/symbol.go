package object

import "hash/fnv"

var symbolClass RubyClassObject = newClass(
	"Symbol", objectClass, symbolMethods, symbolClassMethods, func(RubyClassObject) RubyObject { return &Symbol{} },
)

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

func (s *Symbol) hashKey() hashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return hashKey{Type: s.Type(), Value: h.Sum64()}
}

var symbolClassMethods = map[string]RubyMethod{}

var symbolMethods = map[string]RubyMethod{}
