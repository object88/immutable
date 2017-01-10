package hashmap

import (
	"fmt"
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap(t *testing.T) {
	data := map[immutable.Key]immutable.Value{
		immutable.IntKey(1):  "aa",
		immutable.IntKey(2):  "bb",
		immutable.IntKey(3):  "cc",
		immutable.IntKey(4):  "dd",
		immutable.IntKey(5):  "ee",
		immutable.IntKey(6):  "ff",
		immutable.IntKey(7):  "gg",
		immutable.IntKey(8):  "hh",
		immutable.IntKey(26): "zz",
	}
	original := NewHashMap(data, nil)
	if original.Size() != len(data) {
		t.Fatalf("Incorrect size")
	}
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
	original := NewHashMap(map[immutable.Key]immutable.Value{}, nil)
	if original == nil {
		t.Fatalf("Creating HashMap with empty map returned nil\n")
	}
}

func Test_HashMap_Empty_Size(t *testing.T) {
	original := NewHashMap(map[immutable.Key]immutable.Value{}, nil)
	size := original.Size()
	if size != 0 {
		t.Fatalf("Expected 0 size, got %d\n", size)
	}
}

func Test_Hashmap_Create_WithNilContents(t *testing.T) {
	original := NewHashMap(nil, nil)
	if nil == original {
		t.Fatal("NewHashMap with nil argument returned nil")
	}
	size := original.Size()
	if size != 0 {
		t.Fatalf("Expected 0 size, got %d\n", size)
	}
}

func Test_Hashmap_Insert_WithUnassigned(t *testing.T) {
	var original *HashMap
	key := immutable.IntKey(4)
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
	original := NewHashMap(map[immutable.Key]immutable.Value{}, nil)
	key := immutable.IntKey(4)
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
	contents := map[immutable.Key]immutable.Value{
		immutable.IntKey(0): "a",
		immutable.IntKey(1): "b",
		immutable.IntKey(2): "c",
	}
	original := NewHashMap(contents, nil)
	key := immutable.IntKey(4)
	value := "d"
	modified, error := original.Insert(key, value)
	if nil != error {
		t.Fatalf("Insert to empty hashmap returned error %s\n", error)
	}
	if nil == modified {
		t.Fatal("Insert to empty hashmap did not create a new hashmap\n")
	}
	size := modified.Size()
	if size != len(contents)+1 {
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

func Test_Hashmap_Insert_WithSameKey(t *testing.T) {
	contents := map[immutable.Key]immutable.Value{immutable.IntKey(1): "a", immutable.IntKey(2): "b"}
	original := NewHashMap(contents, nil)
	modified, err := original.Insert(immutable.IntKey(1), "aa")
	if err != nil {
		t.Fatalf("Received error during insert: %s\n", err)
	}
	if modified == nil {
		t.Fatalf("Received nil from insert\n")
	}
	count := modified.Size()
	if count != len(contents) {
		t.Fatalf("Modified hash map contains wrong count; got %d, expected %d\n", count, len(contents))
	}
	result := modified.Get(immutable.IntKey(1))
	if result != "aa" {
		t.Fatalf("Modified hash map has wrong value for key; got %d, expected 'aa'\n", result)
	}
}

func Test_Hashmap_Insert_WithSameKeyAndSameValue(t *testing.T) {
	contents := map[immutable.Key]immutable.Value{immutable.IntKey(1): "a"}
	original := NewHashMap(contents, nil)
	modified, err := original.Insert(immutable.IntKey(1), "a")
	if err != nil {
		t.Fatalf("Received error during insert: %s\n", err)
	}
	if original != modified {
		t.Fatalf("Hash map returned from insert is not same hash map\n")
	}
}

func Test_Hashmap_Get_WithUnassigned(t *testing.T) {
	var original *HashMap
	result := original.Get(immutable.IntKey(2))
	if result != nil {
		t.Fatalf("Request from nil hashmap returned %s", result)
	}
}

func Test_Hashmap_Get_WithEmpty(t *testing.T) {
	original := NewHashMap(map[immutable.Key]immutable.Value{}, nil)
	result := original.Get(immutable.IntKey(2))
	if result != nil {
		t.Fatalf("Request from empty hashmap returned %s", result)
	}
}

func Test_Hashmap_Get_WithContents(t *testing.T) {
	value := "a"
	original := NewHashMap(map[immutable.Key]immutable.Value{immutable.IntKey(1): value}, nil)
	result := original.Get(immutable.IntKey(1))
	if result != value {
		t.Fatalf("Expected %s, got %s", value, result)
	}
}

func Test_Hashmap_Get_Miss(t *testing.T) {
	original := NewHashMap(map[immutable.Key]immutable.Value{immutable.IntKey(1): "a"}, nil)
	result := original.Get(immutable.IntKey(2))
	if result != nil {
		t.Fatalf("Request for miss key returned %s", result)
	}
}

func Test_Hashmap_Remove_WithUnassigned(t *testing.T) {
	var original *HashMap
	key := immutable.IntKey(4)
	modified, error := original.Remove(key)
	if nil != error {
		t.Fatalf("Remove from nil hashmap returned error %s\n", error)
	}
	if nil != modified {
		t.Fatal("Remove from nil hashmap did not return nil\n")
	}
}

func Test_Hashmap_Remove_WithEmpty(t *testing.T) {
	original := NewHashMap(map[immutable.Key]immutable.Value{}, nil)
	key := immutable.IntKey(4)
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
	key1, key2 := immutable.IntKey(0), immutable.IntKey(1)
	value1, value2 := "a", "b"
	contents := map[immutable.Key]immutable.Value{
		key1: value1,
		key2: value2,
	}
	original := NewHashMap(contents, nil)
	modified, error := original.Remove(key2)
	if error != nil {
		t.Fatalf("Error returned from remove: %s", error)
	}
	if modified == nil {
		t.Fatal("Nil returned from remove")
	}
	size := modified.Size()
	if size != len(contents)-1 {
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

func Test_Hashmap_Remove_WithContents_ToEmpty(t *testing.T) {
	key1 := immutable.IntKey(0)
	value1 := "a"
	contents := map[immutable.Key]immutable.Value{
		key1: value1,
	}
	original := NewHashMap(contents, nil)
	modified, error := original.Remove(key1)
	if error != nil {
		t.Fatalf("Error returned from remove: %s", error)
	}
	if modified == nil {
		t.Fatal("Nil returned from remove")
	}
	size := modified.Size()
	if size != 0 {
		t.Fatalf("Incorrect number of entries in returned collection; expected 0, got %d\n", size)
	}
}

func Test_Hashmap_Remove_Miss(t *testing.T) {
	contents := map[immutable.Key]immutable.Value{
		immutable.IntKey(0): "a",
		immutable.IntKey(1): "b",
	}
	original := NewHashMap(contents, nil)
	modified, error := original.Remove(immutable.IntKey(2))
	if error != nil {
		t.Fatalf("Error returned from remove: %s", error)
	}
	if modified != original {
		t.Fatalf("Return from remove did not return the same instance")
	}
}

func Test_Hashmap_ReadAndWriteLargeDataSet(t *testing.T) {
	max := 10000
	contents := make(map[immutable.Key]immutable.Value, max)
	for i := 0; i < max; i++ {
		contents[immutable.IntKey(i)] = i
	}

	original := NewHashMap(contents, nil)
	for i := 0; i < max; i++ {
		k := immutable.IntKey(i)
		result := original.Get(k)
		v := contents[k]
		if result != v {
			t.Fatalf("At %s; expected %d, got %d\n", k, v, result)
		}
	}
}
