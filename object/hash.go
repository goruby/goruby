package object

import (
	"fmt"
	"hash/fnv"
	"strings"
)

var hashClass RubyClassObject = newClass(
	"Hash",
	objectClass,
	hashMethods,
	hashClassMethods,
	func(RubyClassObject, ...RubyObject) (RubyObject, error) {
		return &Hash{hashMap: make(map[hashKey]hashPair)}, nil
	},
)

func init() {
	classes.Set("Hash", hashClass)
}

type hashKey struct {
	Type  Type
	Value uint64
}

func (h hashKey) bytes() []byte {
	return append([]byte(h.Type), byte(h.Value))
}

func hash(obj RubyObject) hashKey {
	if hashable, ok := obj.(hashable); ok {
		return hashable.hashKey()
	}
	pointer := fmt.Sprintf("%p", obj)
	h := fnv.New64a()
	h.Write([]byte(pointer))
	return hashKey{Type: obj.Type(), Value: h.Sum64()}
}

type hashPair struct {
	Key   RubyObject
	Value RubyObject
}

// A Hash represents a Ruby Hash
type Hash struct {
	hashMap map[hashKey]hashPair
}

func (h *Hash) init() {
	if h.hashMap == nil {
		h.hashMap = make(map[hashKey]hashPair)
	}
}

// Set puts the object obj into the Hash
func (h *Hash) Set(key, value RubyObject) RubyObject {
	h.init()
	h.hashMap[hash(key)] = hashPair{Key: key, Value: value}
	return value
}

// Get retrieves the object for key within the hash. If not found, the boolean will be false
func (h *Hash) Get(key RubyObject) (RubyObject, bool) {
	v, ok := h.hashMap[hash(key)]
	if !ok {
		return nil, false
	}
	return v.Value, true
}

// Map returns a map of RubyObject to RubyObject
func (h *Hash) Map() map[RubyObject]RubyObject {
	hashmap := make(map[RubyObject]RubyObject)
	for _, v := range h.hashMap {
		hashmap[v.Key] = v.Value
	}
	return hashmap
}

// Type returns the ObjectType of the array
func (h *Hash) Type() Type { return HASH_OBJ }

// Inspect returns all elements within the array, divided by comma and
// surrounded by brackets
func (h *Hash) Inspect() string {
	elems := []string{}
	for _, v := range h.hashMap {
		elems = append(elems, fmt.Sprintf("%q => %q", v.Key.Inspect(), v.Value.Inspect()))
	}
	return "{" + strings.Join(elems, ", ") + "}"
}

// Class returns the class of the Array
func (h *Hash) Class() RubyClass { return hashClass }

func (h *Hash) hashKey() hashKey {
	hash := fnv.New64a()
	for k := range h.hashMap {
		hash.Write(k.bytes())
	}
	return hashKey{Type: h.Type(), Value: hash.Sum64()}
}

var hashClassMethods = map[string]RubyMethod{}

var hashMethods = map[string]RubyMethod{}
