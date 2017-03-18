package object

var (
	NIL_CLASS RubyClassObject = NewClass("NilClass", OBJECT_CLASS, nilMethods, nilClassMethods)
	NIL       RubyObject      = &Nil{}
)

type Nil struct{}

func (n *Nil) Inspect() string  { return "nil" }
func (n *Nil) Type() ObjectType { return NIL_OBJ }
func (n *Nil) Class() RubyClass { return NIL_CLASS }

var nilClassMethods = map[string]RubyMethod{}

var nilMethods = map[string]RubyMethod{
	"nil?": withArity(0, publicMethod(nilIsNil)),
}

func nilIsNil(context RubyObject, args ...RubyObject) RubyObject {
	return TRUE
}
