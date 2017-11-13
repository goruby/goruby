package object

import (
	"path/filepath"
)

var fileClass RubyClassObject = newClass(
	"File",
	objectClass,
	fileMethods,
	fileClassMethods,
	func(RubyClassObject, ...RubyObject) (RubyObject, error) {
		return &File{make(map[RubyObject]RubyObject)}, nil
	},
)

func init() {
	classes.Set("File", fileClass)
}

// A File represents the Ruby class File
type File struct {
	Map map[RubyObject]RubyObject
}

// Type returns the ObjectType of the array
func (f *File) Type() Type { return OBJECT_OBJ }

// Inspect returns all elements within the array, divided by comma and
// surrounded by brackets
func (f *File) Inspect() string {
	return ""
}

// Class returns the class of the Array
func (f *File) Class() RubyClass { return fileClass }

var fileClassMethods = map[string]RubyMethod{
	"expand_path": publicMethod(fileExpandPath),
	"dirname":     publicMethod(fileDirname),
}

var fileMethods = map[string]RubyMethod{}

func fileExpandPath(context CallContext, args ...RubyObject) (RubyObject, error) {
	switch len(args) {
	case 1:
		str, ok := args[0].(*String)
		if !ok {
			return nil, NewImplicitConversionTypeError(str, args[0])
		}
		path, err := filepath.Abs(str.Value)

		if err == nil {
			return &String{Value: path}, nil
		}

		return nil, NewNotImplementedError("Cannot determine working directory")
	case 2:
		filename, ok := args[0].(*String)
		if !ok {
			return nil, NewImplicitConversionTypeError(filename, args[0])
		}
		dirname, ok := args[1].(*String)
		if !ok {
			return nil, NewImplicitConversionTypeError(filename, args[0])
		}
		// TODO: make sure this is really the wanted behaviour
		abs, err := filepath.Abs(filepath.Join(dirname.Value, filename.Value))
		if err != nil {
			return nil, NewNotImplementedError(err.Error())
		}

		return &String{Value: abs}, nil
	default:
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
}

func fileDirname(context CallContext, args ...RubyObject) (RubyObject, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError(1, len(args))
	}
	filename, ok := args[0].(*String)
	if !ok {
		return nil, NewImplicitConversionTypeError(filename, args[0])
	}

	dirname := filepath.Dir(filename.Value)

	return &String{Value: dirname}, nil
}
