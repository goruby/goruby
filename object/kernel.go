package object

import "fmt"

var KERNEL_MODULE *Module = newModule("Kernel", newEigenclass(MODULE_CLASS, kernelMethods))

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

var kernelMethods = map[string]RubyMethod{
	"puts": publicMethod(puts),
}

func puts(context RubyObject, args ...RubyObject) RubyObject {
	out := ""
	for _, arg := range args {
		out += arg.Inspect()
	}
	fmt.Println(out)
	return NIL
}
