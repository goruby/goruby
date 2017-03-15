package object

func Send(context RubyObject, method string, args ...RubyObject) RubyObject {
	class := context.Class()
	// search for the method in the ancestry tree
	for class != nil {
		fn, ok := class.Methods()[method]
		if ok {
			return fn(context, args...)
		}
		class = class.SuperClass()
	}

	methodMissingArgs := append(
		[]RubyObject{&Symbol{method}},
		args...,
	)

	return Send(context, "method_missing", methodMissingArgs...)
}
