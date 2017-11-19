package object

import (
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/goruby/goruby/parser"
	"github.com/pkg/errors"
)

var kernelModule = newModule("Kernel", kernelMethodSet, nil)

func init() {
	classes.Set("Kernel", kernelModule)
}

var kernelMethodSet = map[string]RubyMethod{
	"to_s":              withArity(0, publicMethod(kernelToS)),
	"nil?":              withArity(0, publicMethod(kernelIsNil)),
	"methods":           publicMethod(kernelMethods),
	"public_methods":    publicMethod(kernelPublicMethods),
	"protected_methods": publicMethod(kernelProtectedMethods),
	"private_methods":   publicMethod(kernelPrivateMethods),
	"class":             withArity(0, publicMethod(kernelClass)),
	"puts":              privateMethod(kernelPuts),
	"require":           withArity(1, privateMethod(kernelRequire)),
	"extend":            publicMethod(kernelExtend),
	"block_given?":      withArity(0, privateMethod(kernelBlockGiven)),
	"tap":               publicMethod(kernelTap),
	"raise":             privateMethod(kernelRaise),
}

func kernelToS(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	if self, ok := receiver.(*Self); ok {
		receiver = self.RubyObject
	}
	val := fmt.Sprintf("#<%s:%p>", receiver.Class().Name(), receiver)
	return &String{Value: val}, nil
}

func kernelPuts(context CallContext, args ...RubyObject) (RubyObject, error) {
	out := ""
	for _, arg := range args {
		out += arg.Inspect()
	}
	fmt.Println(out)
	return NIL, nil
}

func kernelMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showInstanceMethods := true
	if len(args) == 1 {
		boolean, ok := args[0].(*Boolean)
		if !ok {
			boolean = TRUE.(*Boolean)
		}
		showInstanceMethods = boolean.Value
	}

	receiver := context.Receiver()
	class := context.Receiver().Class()

	extended, ok := receiver.(*extendedObject)

	if !showInstanceMethods && !ok {
		return &Array{}, nil
	}

	if !showInstanceMethods && ok {
		class = extended.class
	}

	publicMethods := getMethods(class, PUBLIC_METHOD, showInstanceMethods)
	protectedMethods := getMethods(class, PROTECTED_METHOD, showInstanceMethods)
	return &Array{Elements: append(publicMethods.Elements, protectedMethods.Elements...)}, nil
}

func kernelPublicMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showSuperClassMethods := true
	if len(args) == 1 {
		boolean, ok := args[0].(*Boolean)
		if !ok {
			boolean = TRUE.(*Boolean)
		}
		showSuperClassMethods = boolean.Value
	}
	class := context.Receiver().Class()
	return getMethods(class, PUBLIC_METHOD, showSuperClassMethods), nil
}

func kernelProtectedMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showSuperClassMethods := true
	if len(args) == 1 {
		boolean, ok := args[0].(*Boolean)
		if !ok {
			boolean = TRUE.(*Boolean)
		}
		showSuperClassMethods = boolean.Value
	}
	class := context.Receiver().Class()
	return getMethods(class, PROTECTED_METHOD, showSuperClassMethods), nil
}

func kernelPrivateMethods(context CallContext, args ...RubyObject) (RubyObject, error) {
	showSuperClassMethods := true
	if len(args) == 1 {
		boolean, ok := args[0].(*Boolean)
		if !ok {
			boolean = TRUE.(*Boolean)
		}
		showSuperClassMethods = boolean.Value
	}
	class := context.Receiver().Class()
	return getMethods(class, PRIVATE_METHOD, showSuperClassMethods), nil
}

func kernelIsNil(context CallContext, args ...RubyObject) (RubyObject, error) {
	return FALSE, nil
}

func kernelClass(context CallContext, args ...RubyObject) (RubyObject, error) {
	receiver := context.Receiver()
	if _, ok := receiver.(RubyClassObject); ok {
		return classClass, nil
	}
	return receiver.Class().(RubyClassObject), nil
}

func kernelRequire(context CallContext, args ...RubyObject) (RubyObject, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	name, ok := args[0].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(name, args[0])
	}
	filename := name.Value
	if !strings.HasSuffix(filename, "rb") {
		filename += ".rb"
	}
	absolutePath, _ := filepath.Abs(filename)
	loadedFeatures, ok := context.Env().Get("$LOADED_FEATURES")
	if !ok {
		loadedFeatures = NewArray()
		context.Env().SetGlobal("$LOADED_FEATURES", loadedFeatures)
	}
	arr, ok := loadedFeatures.(*Array)
	if !ok {
		arr = NewArray()
	}
	loaded := false
	for _, feat := range arr.Elements {
		if feat.Inspect() == absolutePath {
			loaded = true
			break
		}
	}
	if loaded {
		return FALSE, nil
	}

	file, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		found := false
		loadPath, _ := context.Env().Get("$:")
		for _, p := range loadPath.(*Array).Elements {
			newPath := path.Join(p.(*String).Value, filename)
			file, err = ioutil.ReadFile(newPath)
			if !os.IsNotExist(err) {
				absolutePath = newPath
				found = true
				break
			}
		}
		if !found {
			return nil, NewNoSuchFileLoadError(name.Value)
		}
	}

	prog, err := parser.ParseFile(token.NewFileSet(), absolutePath, file, 0)
	if err != nil {
		return nil, NewSyntaxError(err)
	}
	_, err = context.Eval(prog, WithScopedLocalVariables(context.Env()))
	if err != nil {
		return nil, errors.WithMessage(err, "require")
	}
	arr.Elements = append(arr.Elements, &String{Value: absolutePath})
	return TRUE, nil
}

func kernelExtend(context CallContext, args ...RubyObject) (RubyObject, error) {
	if len(args) == 0 {
		return nil, NewWrongNumberOfArgumentsError(1, 0)
	}
	modules := make([]*Module, len(args))
	for i, arg := range args {
		module, ok := arg.(*Module)
		if !ok {
			return nil, NewWrongArgumentTypeError(module, arg)
		}
		modules[i] = module
	}
	extended := &extendedObject{
		RubyObject: context.Receiver(),
		class: newEigenclass(
			newMixin(context.Receiver().Class().(RubyClassObject), modules...),
			map[string]RubyMethod{},
		),
	}
	info, _ := EnvStat(context.Env(), context.Receiver())
	info.Env().Set(info.Name(), extended)
	return extended, nil
}

func kernelBlockGiven(context CallContext, args ...RubyObject) (RubyObject, error) {
	self, _ := context.Receiver().(*Self)
	if self.Block == nil {
		return FALSE, nil
	}
	return TRUE, nil
}

func kernelTap(context CallContext, args ...RubyObject) (RubyObject, error) {
	block, remainingArgs, ok := extractBlockFromArgs(args)
	if !ok {
		return nil, NewNoBlockGivenLocalJumpError()
	}
	if len(remainingArgs) != 0 {
		return nil, NewWrongNumberOfArgumentsError(0, 1)
	}
	_, err := block.Call(context, context.Receiver())
	if err != nil {
		return nil, err
	}
	return context.Receiver(), nil
}

func kernelRaise(context CallContext, args ...RubyObject) (RubyObject, error) {
	switch len(args) {
	case 1:
		switch arg := args[0].(type) {
		case *String:
			return nil, NewRuntimeError(arg.Value)
		default:
			exc, err := Send(NewCallContext(context.Env(), arg), "exception")
			if err != nil {
				return nil, NewTypeError("exception class/object expected")
			}
			if excAsErr, ok := exc.(error); ok {
				return nil, excAsErr
			}
			return nil, nil
		}
	default:
		return nil, NewRuntimeError("")
	}
}
