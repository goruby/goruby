package object

var STRING_CLASS RubyClassObject = NewClass("String", OBJECT_CLASS, stringMethods, stringClassMethods)

func init() {
	classes.Set("String", STRING_CLASS)
}

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Class() RubyClass { return STRING_CLASS }

var stringClassMethods = map[string]RubyMethod{
	"new": publicMethod(func(context RubyObject, args ...RubyObject) RubyObject {
		switch len(args) {
		case 0:
			return &String{}
		case 1:
			str, ok := args[0].(*String)
			if !ok {
				return NewImplicitConversionTypeError(args[0], context)
			}
			return &String{Value: str.Value}
		default:
			return NewWrongNumberOfArgumentsError(len(args), 1)
		}
	}),
}

var stringMethods = map[string]RubyMethod{
	"to_s": withArity(0, publicMethod(stringToS)),
}

func stringToS(context RubyObject, args ...RubyObject) RubyObject {
	str := context.(*String)
	return &String{str.Value}
}
