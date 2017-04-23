package object

var stringClass RubyClassObject = newClass(
	"String", objectClass, stringMethods, stringClassMethods, func(RubyClassObject) RubyObject { return &String{} },
)

func init() {
	classes.Set("String", stringClass)
}

// String represents a string in Ruby
type String struct {
	Value string
}

// Inspect returns the Value
func (s *String) Inspect() string { return s.Value }

// Type returns STRING_OBJ
func (s *String) Type() Type { return STRING_OBJ }

// Class returns stringClass
func (s *String) Class() RubyClass { return stringClass }

var stringClassMethods = map[string]RubyMethod{}

var stringMethods = map[string]RubyMethod{
	"initialize": privateMethod(stringInitialize),
	"to_s":       withArity(0, publicMethod(stringToS)),
}

func stringInitialize(context CallContext, args ...RubyObject) (RubyObject, error) {
	self, _ := context.Receiver().(*Self)
	switch len(args) {
	case 0:
		self.RubyObject = &String{}
		return self, nil
	case 1:
		str, ok := args[0].(*String)
		if !ok {
			return nil, NewImplicitConversionTypeError(str, args[0])
		}
		self.RubyObject = &String{Value: str.Value}
		return self, nil
	default:
		return nil, NewWrongNumberOfArgumentsError(len(args), 1)
	}
}

func stringToS(context CallContext, args ...RubyObject) (RubyObject, error) {
	str := context.Receiver().(*String)
	return &String{str.Value}, nil
}
