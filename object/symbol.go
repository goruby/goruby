package object

var SYMBOL_EIGENCLASS RubyClass = &SymbolEigenclass{}
var SYMBOL_CLASS RubyClass = &SymbolClass{}

type SymbolEigenclass struct{}

func (s *SymbolEigenclass) Inspect() string            { return "(Symbol)" }
func (s *SymbolEigenclass) Type() ObjectType           { return SYMBOL_OBJ }
func (s *SymbolEigenclass) Methods() map[string]method { return nil }
func (s *SymbolEigenclass) Class() RubyClass           { return OBJECT_CLASS }
func (s *SymbolEigenclass) SuperClass() RubyClass      { return BASIC_OBJECT_CLASS }

type SymbolClass struct{}

func (s *SymbolClass) Inspect() string            { return "Symbol" }
func (s *SymbolClass) Type() ObjectType           { return SYMBOL_OBJ }
func (s *SymbolClass) Methods() map[string]method { return nil }
func (s *SymbolClass) Class() RubyClass           { return SYMBOL_EIGENCLASS }
func (s *SymbolClass) SuperClass() RubyClass      { return OBJECT_CLASS }

type Symbol struct {
	Value string
}

func (s *Symbol) Inspect() string            { return ":" + s.Value }
func (s *Symbol) Type() ObjectType           { return SYMBOL_OBJ }
func (s *Symbol) Methods() map[string]method { return nil }
func (s *Symbol) Class() RubyClass           { return SYMBOL_CLASS }
