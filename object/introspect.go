package object

func getMethods(class RubyClass, visibility MethodVisibility, addSuperMethods bool) *Array {
	var methodSymbols []RubyObject
	for class != nil {
		methods := class.Methods().GetAll()
		for meth, fn := range methods {
			if fn.Visibility() == visibility {
				methodSymbols = append(methodSymbols, &Symbol{meth})
			}
		}
		if !addSuperMethods {
			break
		}
		class = class.SuperClass()
	}

	return &Array{Elements: methodSymbols}
}
