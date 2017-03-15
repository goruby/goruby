package object

var STRING_EIGENCLASS RubyClass = &StringEigenclass{}
var STRING_CLASS RubyClass = &StringClass{}

type StringEigenclass struct{}

func (s *StringEigenclass) Inspect() string            { return "String" }
func (s *StringEigenclass) Type() ObjectType           { return STRING_CLASS_OBJ }
func (s *StringEigenclass) Methods() map[string]method { return stringClassMethods }
func (s *StringEigenclass) Class() RubyClass           { return OBJECT_CLASS }
func (s *StringEigenclass) SuperClass() RubyClass      { return BASIC_OBJECT_CLASS }

type StringClass struct{}

func (s *StringClass) Inspect() string            { return "String" }
func (s *StringClass) Type() ObjectType           { return STRING_CLASS_OBJ }
func (s *StringClass) Methods() map[string]method { return stringMethods }
func (s *StringClass) Class() RubyClass           { return STRING_EIGENCLASS }
func (s *StringClass) SuperClass() RubyClass      { return OBJECT_CLASS }

type String struct {
	Value string
}

func (s *String) Inspect() string            { return s.Value }
func (s *String) Type() ObjectType           { return STRING_OBJ }
func (s *String) Methods() map[string]method { return nil }
func (s *String) Class() RubyClass           { return STRING_CLASS }

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
