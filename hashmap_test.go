package immutable

import (
	"fmt"
	"testing"
)

func Test_Hashmap(t *testing.T) {
	data := map[Key]Value{IntKey(1): "aa", IntKey(2): "bb", IntKey(3): "cc", IntKey(4): "dd", IntKey(5): "ee", IntKey(6): "ff", IntKey(7): "gg", IntKey(8): "hh", IntKey(26): "zz"}
	original := NewHashMap(data)
	if original.Size() != uint32(len(data)) {
		t.Fatalf("Incorrect size")
	}
	fmt.Println(original.String())
	for k, v := range data {
		result := original.Get(k)
		if result != v {
			t.Fatalf("Key '%s' -> '%s', expected '%s'\n", k, result, v)
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
	original := NewHashMap(map[Key]Value{})
	if original == nil {
		t.Fatalf("Creating HashMap with empty map returned nil\n")
	}
}

func Test_HashMap_Empty_Size(t *testing.T) {
	original := NewHashMap(map[Key]Value{})
	size := original.Size()
	if 0 != size {
		t.Fatalf("Expected 0 size, got %d\n", size)
	}
}

func Test_Hashmap_ReadAndWriteLargeDataSet(t *testing.T) {
	max := 10000
	contents := make(map[Key]Value, max)
	for i := 0; i < max; i++ {
		contents[IntKey(i)] = i
	}

	original := NewHashMap(contents)
	for k, v := range contents {
		result := original.Get(k)
		if result != v {
			t.Fatalf("At %s; expected %d, got %d\n", k, v, result)
		}
	}
}
