package immutable

import (
	"fmt"
	"testing"
)

func Test_Hashmap(t *testing.T) {
	data := map[Key]Value{
		IntKey(1):  "aa",
		IntKey(2):  "bb",
		IntKey(3):  "cc",
		IntKey(4):  "dd",
		IntKey(5):  "ee",
		IntKey(6):  "ff",
		IntKey(7):  "gg",
		IntKey(8):  "hh",
		IntKey(26): "zz",
	}
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

func Test_Hashmap_Insert_WithUnassigned(t *testing.T) {
	var original *HashMap
	key := IntKey(4)
	value := "d"
	modified, error := original.Insert(key, value)
	if nil != error {
		t.Fatalf("Insert to nil hashmap returned error %s\n", error)
	}
	if nil == modified {
		t.Fatal("Insert to nil hashmap did not create a new hashmap\n")
	}
	size := modified.Size()
	if size != 1 {
		t.Fatalf("New hashmap has size %d; expected 1", size)
	}
	returnedValue := modified.Get(key)
	if returnedValue != value {
		t.Fatalf("Incorrect value stored in new hashmap; expected %s, got %s\n", value, returnedValue)
	}
}

func Test_Hashmap_Insert_WithEmpty(t *testing.T) {
	original := NewHashMap(map[Key]Value{})
	key := IntKey(4)
	value := "d"
	modified, error := original.Insert(key, value)
	if nil != error {
		t.Fatalf("Insert to empty hashmap returned error %s\n", error)
	}
	if nil == modified {
		t.Fatal("Insert to empty hashmap did not create a new hashmap\n")
	}
	size := modified.Size()
	if size != 1 {
		t.Fatalf("New hashmap has size %d; expected 1", size)
	}
	returnedValue := modified.Get(key)
	if returnedValue != value {
		t.Fatalf("Incorrect value stored in new hashmap; expected %s, got %s\n", value, returnedValue)
	}
}

func Test_Hashmap_Insert_WithContents(t *testing.T) {
	contents := map[Key]Value{
		IntKey(0): "a",
		IntKey(1): "b",
		IntKey(2): "c",
	}
	original := NewHashMap(contents)
	key := IntKey(4)
	value := "d"
	modified, error := original.Insert(key, value)
	if nil != error {
		t.Fatalf("Insert to empty hashmap returned error %s\n", error)
	}
	if nil == modified {
		t.Fatal("Insert to empty hashmap did not create a new hashmap\n")
	}
	size := modified.Size()
	if size != uint32(len(contents)+1) {
		t.Fatalf("New hashmap has size %d; expected %d", size, len(contents)+1)
	}
	for k, v := range contents {
		result := modified.Get(k)
		if result != v {
			t.Fatalf("At %s, incorrect value; expected %s, got %s\n", k, v, result)
		}
	}
	result := modified.Get(key)
	if result != value {
		t.Fatalf("At %s, incorrect value; expected %s, got %s\n", key, value, result)
	}
}

func Test_Hashmap_Get_WithUnassigned(t *testing.T) {
	var original *HashMap
	result := original.Get(IntKey(2))
	if result != nil {
		t.Fatalf("Request from nil hashmap returned %s", result)
	}
}

func Test_Hashmap_Get_WithEmpty(t *testing.T) {
	original := NewHashMap(map[Key]Value{})
	result := original.Get(IntKey(2))
	if result != nil {
		t.Fatalf("Request from empty hashmap returned %s", result)
	}
}

func Test_Hashmap_Get_WithContents(t *testing.T) {
	value := "a"
	original := NewHashMap(map[Key]Value{IntKey(1): value})
	result := original.Get(IntKey(1))
	if result != value {
		t.Fatalf("Expected %s, got %s", value, result)
	}
}

func Test_Hashmap_Get_Miss(t *testing.T) {
	original := NewHashMap(map[Key]Value{IntKey(1): "a"})
	result := original.Get(IntKey(2))
	if result != nil {
		t.Fatalf("Request for miss key returned %s", result)
	}
}

func Test_Hashmap_Remove_WithUnassigned(t *testing.T) {
	var original *HashMap
	key := IntKey(4)
	modified, error := original.Remove(key)
	if nil != error {
		t.Fatalf("Remove from nil hashmap returned error %s\n", error)
	}
	if nil != modified {
		t.Fatal("Remove from nil hashmap did not return nil\n")
	}
}

func Test_Hashmap_Remove_WithEmpty(t *testing.T) {
	original := NewHashMap(map[Key]Value{})
	key := IntKey(4)
	modified, error := original.Remove(key)
	if nil != error {
		t.Fatalf("Remove from empty hashmap returned error %s\n", error)
	}
	if nil == modified {
		t.Fatal("Remove from empty hashmap returned nil\n")
	}
	size := modified.Size()
	if size != 0 {
		fmt.Printf("Remove from empty hashmap returned a non-empty hashmap: %s\n", modified)
	}
}

func Test_Hashmap_Remove_WithContents(t *testing.T) {
	key1, key2 := IntKey(0), IntKey(1)
	value1, value2 := "a", "b"
	contents := map[Key]Value{
		key1: value1,
		key2: value2,
	}
	original := NewHashMap(contents)
	modified, error := original.Remove(key2)
	if error != nil {
		t.Fatalf("Error returned from remove: %s", error)
	}
	if modified == nil {
		t.Fatal("Nil returned from remove")
	}
	size := modified.Size()
	if size != uint32(len(contents)-1) {
		t.Fatalf("Incorrect number of entries in returned collection; expected %d, got %d\n", len(contents)-1, size)
	}
	result1 := modified.Get(key1)
	if result1 != value1 {
		t.Fatalf("Value returned from key %s is incorrect; expected %s, got %s\n", key1, value1, result1)
	}
	result2 := modified.Get(key2)
	if result2 != nil {
		t.Fatalf("Got value from key %s; got %s, expected nil\n", key2, result2)
	}
}

func Test_Hashmap_Remove_Miss(t *testing.T) {
	key1, key2, key3 := IntKey(0), IntKey(1), IntKey(2)
	value1, value2 := "a", "b"
	contents := map[Key]Value{
		key1: value1,
		key2: value2,
	}
	original := NewHashMap(contents)
	modified, error := original.Remove(key3)
	if error != nil {
		t.Fatalf("Error returned from remove: %s", error)
	}
	if modified == nil {
		t.Fatal("Nil returned from remove")
	}
	size := modified.Size()
	if size != uint32(len(contents)) {
		t.Fatalf("Incorrect number of entries in returned collection; expected %d, got %d\n", len(contents), size)
	}
	for k, v := range contents {
		result := modified.Get(k)
		if result != v {
			t.Fatalf("Value returned from key %s is incorrect; expected %s, got %s\n", k, v, result)
		}
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
