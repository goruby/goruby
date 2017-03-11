package object

import "fmt"

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
