package object

var stringClass RubyClassObject = newClass("String", objectClass, stringMethods, stringClassMethods)

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
