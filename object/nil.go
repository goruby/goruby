package object

var (
	nilClass RubyClassObject = newClass("NilClass", objectClass, nilMethods, nilClassMethods)
	// NIL represents the singleton object nil
	NIL RubyObject = &nilObject{}
)

func init() {
	classes.Set("NilClass", nilClass)
}

type nilObject struct{}

func (n *nilObject) Inspect() string  { return "nil" }
func (n *nilObject) Type() Type       { return NIL_OBJ }
func (n *nilObject) Class() RubyClass { return nilClass }

var nilClassMethods = map[string]RubyMethod{}

var nilMethods = map[string]RubyMethod{
	"nil?": withArity(0, publicMethod(nilIsNil)),
}

func nilIsNil(context CallContext, args ...RubyObject) (RubyObject, error) {
	return TRUE, nil
}
