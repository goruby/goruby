package object

var (
	SYMBOL_EIGENCLASS RubyClass       = newEigenClass(CLASS_CLASS, symbolClassMethods)
	SYMBOL_CLASS      RubyClassObject = &SymbolClass{}
)

type SymbolClass struct{}

func (s *SymbolClass) Inspect() string            { return "Symbol" }
func (s *SymbolClass) Type() ObjectType           { return SYMBOL_OBJ }
func (s *SymbolClass) Class() RubyClass           { return SYMBOL_EIGENCLASS }
func (s *SymbolClass) Methods() map[string]method { return symbolMethods }
func (s *SymbolClass) SuperClass() RubyClass      { return OBJECT_CLASS }

type Symbol struct {
	Value string
}

func (s *Symbol) Inspect() string  { return ":" + s.Value }
func (s *Symbol) Type() ObjectType { return SYMBOL_OBJ }
func (s *Symbol) Class() RubyClass { return SYMBOL_CLASS }

var symbolClassMethods = map[string]method{}

var symbolMethods = map[string]method{}
