package immutable

import (
	"fmt"
	"testing"
)

func Test_Hashmap(t *testing.T) {
	data := map[uint32]Value{1: "aa", 2: "bb", 3: "cc", 4: "dd", 5: "ee", 6: "ff", 7: "gg", 8: "hh", 26: "zz"}
	original := NewHashMap(data)
	if original.Size() != uint32(len(data)) {
		t.Fatalf("Incorrect size")
	}
	fmt.Println(original.String())
	for k, v := range data {
		result := original.Get(k)
		if result != v {
			t.Fatalf("Key '%d' -> '%s', expected '%s'\n", k, result, v)
		}
	}
}

func Test_HashMap_Nil_Size(t *testing.T) {
	var original *HashMap
	size := original.Size()
	if 0 != size {
		t.Fatalf("Expected 0 size, got %d\n", size)
	}
}

func Test_HashMap_Empty_Create(t *testing.T) {
	original := NewHashMap(map[uint32]Value{})
	if original == nil {
		t.Fatalf("Creating HashMap with empty map returned nil\n")
	}
}

func Test_HashMap_Empty_Size(t *testing.T) {
	original := NewHashMap(map[uint32]Value{})
	size := original.Size()
	if 0 != size {
		t.Fatalf("Expected 0 size, got %d\n", size)
	}
}

func Test_Hashmap_Iterate(t *testing.T) {
	data := map[uint32]Value{1: false, 2: false, 3: false, 4: false, 5: false, 6: false}
	original := NewHashMap(data)

	for k, v, i := original.Iterate()(); i != nil; k, v, i = i() {
		if v.(bool) {
			t.Fatalf("At %d, already visited\n", k)
		}
		data[k] = true
	}

	for k, v := range data {
		if !v.(bool) {
			t.Fatalf("At %d, not visited\n", k)
		}
	}
}

func Test_Hashmap_Map(t *testing.T) {
	// original := NewHashMap(map[Key]Value{StringKey("a"): "aa", StringKey("b"): "bb"})
	// modified, _ := original.ToBaseStruct().Map(func(k Key, v Value) (Value, error) {
	original := NewHashMap(map[uint32]Value{1: "aa", 2: "bb"})
	modified, _ := original.Map(func(k uint32, v Value) (Value, error) {
		return fmt.Sprintf("[%d -> %s]", k, v), nil
	})
	for k, v, i := modified.Iterate()(); i != nil; k, v, i = i() {
		fmt.Printf("%d -> %s\n", k, v)
	}
}
