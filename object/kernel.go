package object

import "fmt"

var KERNEL_MODULE *Module = newModule("Kernel", kernelMethodSet)

func init() {
	classes.Set("Kernel", KERNEL_MODULE)
}

var kernelFunctions = &Environment{
	store: map[string]RubyObject{
		"puts": &Builtin{
			Fn: func(args ...RubyObject) RubyObject {
				out := ""
				for _, arg := range args {
					out += arg.Inspect()
				}
				fmt.Println(out)
				return NIL
			},
		},
	},
}

var kernelMethodSet = map[string]RubyMethod{
	"nil?":    withArity(0, publicMethod(kernelIsNil)),
	"methods": withArity(0, publicMethod(kernelMethods)),
	"class":   withArity(0, publicMethod(kernelClass)),
	"puts":    privateMethod(kernelPuts),
}

func kernelPuts(context RubyObject, args ...RubyObject) RubyObject {
	out := ""
	for _, arg := range args {
		out += arg.Inspect()
	}
	fmt.Println(out)
	return NIL
}

func kernelMethods(context RubyObject, args ...RubyObject) RubyObject {
	var methodSymbols []RubyObject
	class := context.Class()
	for class != nil {
		methods := class.Methods()
		for meth, fn := range methods {
			if fn.Visibility() == PUBLIC_METHOD {
				methodSymbols = append(methodSymbols, &Symbol{meth})
			}
		}
		class = class.SuperClass()
	}

	return &Array{Elements: methodSymbols}
}

func kernelIsNil(context RubyObject, args ...RubyObject) RubyObject {
	return FALSE
}

func kernelClass(context RubyObject, args ...RubyObject) RubyObject {
	class := context.Class()
	if eigenClass, ok := class.(*eigenclass); ok {
		class = eigenClass.Class()
	}
	classObj := class.(RubyClassObject)
	return classObj
}
