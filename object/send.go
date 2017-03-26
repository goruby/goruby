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

		if fn.Visibility() == PRIVATE_METHOD {
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
