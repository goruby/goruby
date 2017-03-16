package object

var (
	STRING_EIGENCLASS RubyClass       = newEigenClass(CLASS_CLASS, stringClassMethods)
	STRING_CLASS      RubyClassObject = &StringClass{}
)

type StringClass struct{}

func (s *StringClass) Inspect() string            { return "String" }
func (s *StringClass) Type() ObjectType           { return STRING_CLASS_OBJ }
func (s *StringClass) Class() RubyClass           { return STRING_EIGENCLASS }
func (s *StringClass) Methods() map[string]method { return stringMethods }
func (s *StringClass) SuperClass() RubyClass      { return OBJECT_CLASS }

type String struct {
	Value string
}

func (s *String) Inspect() string  { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Class() RubyClass { return STRING_CLASS }

var stringClassMethods = map[string]method{
	"new": func(context RubyObject, args ...RubyObject) RubyObject {
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
	},
}

var stringMethods = map[string]method{
	"to_s": func(context RubyObject, args ...RubyObject) RubyObject {
		str := context.(*String)
		return &String{str.Value}
	},
}
