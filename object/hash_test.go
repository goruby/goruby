package object

import (
	"reflect"
	"testing"
)

func TestHashSet(t *testing.T) {
	t.Run("Set on initialized hash", func(t *testing.T) {
		hash := &Hash{hashMap: make(map[hashKey]hashPair)}

		key := &String{Value: "foo"}
		value := &Integer{Value: 42}

		result := hash.Set(key, value)

		if len(hash.hashMap) != 1 {
			t.Logf("Expected hashMap to contain 1 item, got %d\n", len(hash.hashMap))
			t.Fail()
		}

		var values []hashPair
		for _, v := range hash.hashMap {
			values = append(values, v)
		}

		if !reflect.DeepEqual(values[0].Key, key) {
			t.Logf("Expect hashPair Key to equal\n%v\n\tgot\n%v\n", key, values[0].Key)
			t.Fail()
		}
		if !reflect.DeepEqual(values[0].Value, value) {
			t.Logf("Expect hashPair Value to equal\n%v\n\tgot\n%v\n", value, values[0].Value)
			t.Fail()
		}

		if !reflect.DeepEqual(result, value) {
			t.Logf("Expect returned value to equal\n%v\n\tgot\n%v\n", value, result)
			t.Fail()
		}
	})
	t.Run("Set on uninitialized hash", func(t *testing.T) {
		var hash Hash

		key := &String{Value: "foo"}
		value := &Integer{Value: 42}

		result := hash.Set(key, value)

		if len(hash.hashMap) != 1 {
			t.Logf("Expected hashMap to contain 1 item, got %d\n", len(hash.hashMap))
			t.Fail()
		}

		var values []hashPair
		for _, v := range hash.hashMap {
			values = append(values, v)
		}

		if !reflect.DeepEqual(values[0].Key, key) {
			t.Logf("Expect hashPair Key to equal\n%v\n\tgot\n%v\n", key, values[0].Key)
			t.Fail()
		}
		if !reflect.DeepEqual(values[0].Value, value) {
			t.Logf("Expect hashPair Value to equal\n%v\n\tgot\n%v\n", value, values[0].Value)
			t.Fail()
		}

		if !reflect.DeepEqual(result, value) {
			t.Logf("Expect returned value to equal\n%v\n\tgot\n%v\n", value, result)
			t.Fail()
		}
	})
}

func TestHashGet(t *testing.T) {
	t.Run("value found", func(t *testing.T) {
		key := &String{Value: "foo"}
		value := &Integer{Value: 42}

		hash := &Hash{hashMap: map[hashKey]hashPair{
			key.hashKey(): hashPair{Key: key, Value: value},
		}}

		result, ok := hash.Get(key)

		if !ok {
			t.Logf("Expected returned bool to be true, got false")
			t.Fail()
		}

		if !reflect.DeepEqual(result, value) {
			t.Logf("Expect returned value to equal\n%v\n\tgot\n%v\n", value, result)
			t.Fail()
		}
	})
	t.Run("value not found", func(t *testing.T) {
		key := &String{Value: "foo"}

		hash := &Hash{hashMap: map[hashKey]hashPair{}}

		result, ok := hash.Get(key)

		if ok {
			t.Logf("Expected returned bool to be false, got true")
			t.Fail()
		}

		if result != nil {
			t.Logf("Expect returned value to be nil\n")
			t.Fail()
		}
	})
	t.Run("on uninitalized hash", func(t *testing.T) {
		key := &String{Value: "foo"}

		var hash Hash

		result, ok := hash.Get(key)

		if ok {
			t.Logf("Expected returned bool to be false, got true")
			t.Fail()
		}

		if result != nil {
			t.Logf("Expect returned value to be nil\n")
			t.Fail()
		}
	})
}

func TestHashMap(t *testing.T) {
	t.Run("on initialized hash", func(t *testing.T) {
		key := &String{Value: "foo"}
		value := &Integer{Value: 42}

		hash := &Hash{hashMap: map[hashKey]hashPair{
			key.hashKey(): hashPair{Key: key, Value: value},
		}}

		var result map[RubyObject]RubyObject = hash.Map()

		expected := map[string]RubyObject{
			"foo": value,
		}
		actual := make(map[string]RubyObject)
		for k, v := range result {
			actual[k.Inspect()] = v
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Logf("Expected hash to equal\n%s\n\tgot\n%s\n", expected, actual)
			t.Fail()
		}
	})
	t.Run("on uninitialized hash", func(t *testing.T) {
		var hash Hash

		var result map[RubyObject]RubyObject = hash.Map()

		expected := map[string]RubyObject{}
		actual := make(map[string]RubyObject)
		for k, v := range result {
			actual[k.Inspect()] = v
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Logf("Expected hash to equal\n%s\n\tgot\n%s\n", expected, actual)
			t.Fail()
		}
	})
}

func Test_hash(t *testing.T) {
	t.Run("hashable object", func(t *testing.T) {
		obj := &String{Value: "bar"}

		hashKey := hash(obj)

		if hashKey != obj.hashKey() {
			t.Logf("Expected to get the same hashKey as from the hashKey fn, got %v", hashKey)
			t.Fail()
		}
	})
	t.Run("object which is not hashable", func(t *testing.T) {
		obj := &basicObject{}

		key := hash(obj)

		expectedKeyType := obj.Type()

		if key.Type != expectedKeyType {
			t.Logf("Expected to get the same hashKey Type as the object Type, got %v", key.Type)
			t.Fail()
		}

		obj2 := &basicObject{}

		key2 := hash(obj2)
		t.Logf("pointer obj1: %p, pointer obj2: %p\n", obj, obj2)

		if key == key2 {
			t.Logf("Expected different keys for different object instances")
			t.Logf("obj == obj2: %t\n", obj == obj2)
			t.Fail()
		}
	})
}
