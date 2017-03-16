package object

import "fmt"

var (
	CLASS_EIGENCLASS RubyClass = newEigenClass(MODULE_CLASS, classClassMethods)
	CLASS_CLASS      RubyClass = &ClassClass{}
)

type ClassClass struct{}

func (c *ClassClass) Inspect() string            { return fmt.Sprintf("#<Class:%p>", c) }
func (c *ClassClass) Type() ObjectType           { return CLASS_OBJ }
func (c *ClassClass) Class() RubyClass           { return CLASS_CLASS }
func (c *ClassClass) Methods() map[string]method { return classMethods }
func (c *ClassClass) SuperClass() RubyClass      { return CLASS_EIGENCLASS }

type Class struct{}

func (c *Class) Inspect() string  { return fmt.Sprintf("#<Class:%p>", c) }
func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Class() RubyClass { return CLASS_CLASS }

var classClassMethods = map[string]method{}

var classMethods = map[string]method{}
