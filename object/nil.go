package object

var (
	NIL_EIGENCLASS RubyClass  = newEigenClass(OBJECT_CLASS, nilClassMethods)
	NIL_CLASS      RubyClass  = &NilClass{}
	NIL            RubyObject = &Nil{}
)

type NilClass struct{}

func (n *NilClass) Inspect() string            { return "NilClass" }
func (n *NilClass) Type() ObjectType           { return NIL_CLASS_OBJ }
func (n *NilClass) Methods() map[string]method { return nilMethods }
func (n *NilClass) Class() RubyClass           { return NIL_EIGENCLASS }
func (n *NilClass) SuperClass() RubyClass      { return OBJECT_CLASS }

type Nil struct{}

func (n *Nil) Inspect() string  { return "nil" }
func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Class() RubyClass { return NIL_CLASS }

var nilClassMethods = map[string]method{}

var nilMethods = map[string]method{
	"nil?": withArity(0, nilIsNil),
}

func nilIsNil(context RubyObject, args ...RubyObject) RubyObject {
	return TRUE
}
