package object

import (
	"fmt"
	"strings"
)

var hashClass RubyClassObject = newClass(
	"Hash",
	objectClass,
	hashMethods,
	hashClassMethods,
	func(RubyClassObject) RubyObject { return &Hash{make(map[RubyObject]RubyObject)} },
)

func init() {
	classes.Set("Hash", hashClass)
}

// A Hash represents a Ruby Hash
type Hash struct {
	Map map[RubyObject]RubyObject
}

// Type returns the ObjectType of the array
func (h *Hash) Type() Type { return HASH_OBJ }

// Inspect returns all elements within the array, divided by comma and
// surrounded by brackets
func (h *Hash) Inspect() string {
	elems := []string{}
	for key, val := range h.Map {
		elems = append(elems, fmt.Sprintf("%q => %q", key.Inspect(), val.Inspect()))
	}
	return "{" + strings.Join(elems, ", ") + "}"
}

// Class returns the class of the Array
func (h *Hash) Class() RubyClass { return hashClass }

var hashClassMethods = map[string]RubyMethod{}

var hashMethods = map[string]RubyMethod{}
