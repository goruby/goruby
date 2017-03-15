package object

var (
	NIL_EIGENCLASS RubyClass = &NilEigenClass{}
	NIL_CLASS      RubyClass = &NilClass{}
)

type NilEigenClass struct{}

func (n *NilEigenClass) Inspect() string            { return "nil" }
func (n *NilEigenClass) Type() ObjectType           { return NIL_OBJ }
func (n *NilEigenClass) Methods() map[string]method { return nil }
func (n *NilEigenClass) Class() RubyClass           { return nil }
func (n *NilEigenClass) SuperClass() RubyClass      { return nil }

type NilClass struct{}

func (n *NilClass) Inspect() string            { return "NilClass" }
func (n *NilClass) Type() ObjectType           { return NIL_OBJ }
func (n *NilClass) Methods() map[string]method { return nil }
func (n *NilClass) Class() RubyClass           { return nil }
func (n *NilClass) SuperClass() RubyClass      { return nil }

type Nil struct{}

func (n *Nil) Inspect() string            { return "nil" }
func (n *Nil) Type() ObjectType           { return NIL_OBJ }
func (n *Nil) Methods() map[string]method { return nil }
func (n *Nil) Class() RubyClass           { return nil }
