package object

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/goruby/goruby/lexer"
	"github.com/goruby/goruby/parser"
)

var kernelModule = NewModule("Kernel", kernelMethodSet)

func init() {
	classes.Set("Kernel", kernelModule)
}

var kernelMethodSet = map[string]RubyMethod{
	"nil?":    withArity(0, publicMethod(kernelIsNil)),
	"methods": withArity(0, publicMethod(kernelMethods)),
	"class":   withArity(0, publicMethod(kernelClass)),
	"puts":    privateMethod(kernelPuts),
	"require": withArity(1, privateMethod(kernelRequire)),
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
	var methodSymbols []RubyObject
	class := context.Receiver().Class()
	for class != nil {
		methods := class.Methods()
		for meth, fn := range methods {
			if fn.Visibility() == PUBLIC_METHOD {
				methodSymbols = append(methodSymbols, &Symbol{meth})
			}
		}
		class = class.SuperClass()
	}

	return &Array{Elements: methodSymbols}, nil
}

func kernelIsNil(context CallContext, args ...RubyObject) (RubyObject, error) {
	return FALSE, nil
}

func kernelClass(context CallContext, args ...RubyObject) (RubyObject, error) {
	class := context.Receiver().Class()
	if eigenClass, ok := class.(*eigenclass); ok {
		class = eigenClass.Class()
	}
	classObj := class.(RubyClassObject)
	return classObj, nil
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
		return nil, NewLoadError(name.Value)
	}
	l := lexer.New(string(file))
	p := parser.New(l)
	prog, err := p.ParseProgram()
	if err != nil {
		return nil, NewSyntaxError(err)
	}
	_, err = context.Eval(prog, NewEnclosedEnvironment(context.Env()))
	if err != nil {
		return nil, err
	}
	arr.Elements = append(arr.Elements, &String{Value: absolutePath})
	return TRUE, nil
}
