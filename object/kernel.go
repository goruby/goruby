package object

import "fmt"

var KERNEL_MODULE RubyClass = &KernelClass{}

type KernelClass struct{}

func (k *KernelClass) Inspect() string            { return "Kernel" }
func (k *KernelClass) Type() ObjectType           { return MODULE_OBJ }
func (k *KernelClass) Methods() map[string]method { return kernelMethods }
func (k *KernelClass) Class() RubyClass           { return newEigenClass(MODULE_CLASS, kernelMethods) }
func (k *KernelClass) SuperClass() RubyClass      { return MODULE_CLASS }

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

var kernelMethods = map[string]method{
	"puts": puts,
}

func puts(context RubyObject, args ...RubyObject) RubyObject {
	out := ""
	for _, arg := range args {
		out += arg.Inspect()
	}
	fmt.Println(out)
	return NIL
}
