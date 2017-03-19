package object

import "fmt"

var kernelFunctions = &Environment{
	store: map[string]Object{
		"puts": &Builtin{
			Fn: func(args ...Object) Object {
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
