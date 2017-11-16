package object

import (
	"fmt"
	"hash/fnv"
)

var stringClass RubyClassObject = newClass(
	"String",
	objectClass,
	stringMethods,
	stringClassMethods,
	func(RubyClassObject, ...RubyObject) (RubyObject, error) {
		return &String{}, nil
	},
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

// hashKey returns a hash key to be used by Hashes
func (s *String) hashKey() hashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return hashKey{Type: s.Type(), Value: h.Sum64()}
}

func stringify(obj RubyObject) (string, error) {
	stringObj, err := Send(NewCallContext(nil, obj), "to_s")
	if err != nil {
		return "", NewTypeError(
			fmt.Sprintf(
				"can't convert %s into String",
				obj.Class().Name(),
			),
		)
	}
	str, ok := stringObj.(*String)
	if !ok {
		return "", NewTypeError(
			fmt.Sprintf(
				"can't convert %s to String (%s#to_s gives %s)",
				obj.Class().Name(),
				obj.Class().Name(),
				stringObj.Class().Name(),
			),
		)
	}
	return str.Value, nil
}

var stringClassMethods = map[string]RubyMethod{}

var stringMethods = map[string]RubyMethod{
	"initialize": privateMethod(stringInitialize),
	"to_s":       withArity(0, publicMethod(stringToS)),
	"+":          withArity(1, publicMethod(stringAdd)),
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

func stringAdd(context CallContext, args ...RubyObject) (RubyObject, error) {
	s := context.Receiver().(*String)
	add, ok := args[0].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(add, args[0])
	}
	return &String{s.Value + add.Value}, nil
}
