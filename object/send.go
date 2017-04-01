package object

// Send sends message method with args to context and returns its result
func Send(context RubyObject, method string, args ...RubyObject) RubyObject {
	class := context.Class()

	// search for the method in the ancestry tree
	for class != nil {
		fn, ok := class.Methods()[method]
		if !ok {
			class = class.SuperClass()
			continue
		}

		if fn.Visibility() == PRIVATE_METHOD && context.Type() != SELF {
			return NewPrivateNoMethodError(context, method)
		}

		return fn.Call(context, args...)
	}

	methodMissingArgs := append(
		[]RubyObject{&Symbol{method}},
		args...,
	)

	return methodMissing(context, methodMissingArgs...)
}

// AddMethod adds a method to a given object. It returns the object with the modified method set
func AddMethod(context RubyObject, methodName string, method *Function) RubyObject {
	objectToExtend := context
	self, contextIsSelf := context.(*Self)
	if contextIsSelf {
		objectToExtend = self.RubyObject
	}
	extended, ok := objectToExtend.(*extendedObject)
	if !ok {
		extended = &extendedObject{
			RubyObject: context,
			class:      newEigenclass(context.Class(), map[string]RubyMethod{}),
		}
	}
	extended.addMethod(methodName, method)
	if contextIsSelf {
		self.RubyObject = extended
		return self
	}
	return extended
}

func methodMissing(context RubyObject, args ...RubyObject) RubyObject {
	class := context.Class()

	// search for method_missing in the ancestry tree
	for class != nil {
		fn, ok := class.Methods()["method_missing"]
		if !ok {
			class = class.SuperClass()
			continue
		}
		return fn.Call(context, args...)
	}
	return NewNoMethodError(context, args[0].(*Symbol).Value)
}
