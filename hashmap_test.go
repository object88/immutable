package immutable

import (
	"fmt"
	"testing"
)

func Test_Hashmap(t *testing.T) {
	data := map[int]string{
		1:  "aa",
		2:  "bb",
		3:  "cc",
		4:  "dd",
		5:  "ee",
		6:  "ff",
		7:  "gg",
		8:  "hh",
		26: "zz",
	}
	original := NewIntToStringHashmap(data)
	if original.Size() != len(data) {
		t.Fatalf("Incorrect size")
	}
	for k, v := range data {
		result, ok, err := original.Get(k)
		if err != nil {
			t.Fatalf("Got err %s for key '%d'\n", err.Error(), k)
		}
		if !ok {
			t.Fatalf("Got !ok for key '%d'\n", k)
		}
		if result != v {
			t.Fatalf("Element '%d' -> '%s', expected '%s'\n", k, result, v)
		}
	}
}

func Test_HashMap_Nil_Size(t *testing.T) {
	var original *IntToStringHashmap
	size := original.Size()
	if 0 != size {
		t.Fatalf("Expected 0 size, got %d\n", size)
	}
}

func Test_HashMap_Empty_Create(t *testing.T) {
	original := NewIntToStringHashmap(map[int]string{})
	if original == nil {
		t.Fatalf("Creating HashMap with empty map returned nil\n")
	}
}

func Test_HashMap_Empty_Size(t *testing.T) {
	original := NewIntToStringHashmap(map[int]string{})
	size := original.Size()
	if size != 0 {
		t.Fatalf("Expected 0 size, got %d\n", size)
	}
}

func Test_Hashmap_Create_WithNilContents(t *testing.T) {
	original := NewIntToStringHashmap(nil)
	if nil == original {
		t.Fatal("NewHashMap with nil argument returned nil")
	}
	size := original.Size()
	if size != 0 {
		t.Fatalf("Expected 0 size, got %d\n", size)
	}
}

func Test_Hashmap_Insert_WithUnassigned(t *testing.T) {
	var original *IntToStringHashmap
	modified, err := original.Insert(4, "d")
	if err == nil {
		t.Fatal("Insert to nil hashmap did not return error\n")
	}
	if modified != nil {
		t.Fatal("Insert to nil hashmap did not return nil\n")
	}
}

func Test_Hashmap_Insert_WithEmpty(t *testing.T) {
	original := NewIntToStringHashmap(map[int]string{})
	key := 4
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
	returnedValue, _, _ := modified.Get(key)
	if returnedValue != value {
		t.Fatalf("Incorrect value stored in new hashmap; expected %s, got %s\n", value, returnedValue)
	}
}

func Test_Hashmap_Insert_WithContents(t *testing.T) {
	contents := map[int]string{
		0: "a",
		1: "b",
		2: "c",
	}
	original := NewIntToStringHashmap(contents)
	key := 4
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
		result, _, _ := modified.Get(k)
		if result != v {
			t.Fatalf("At %d, incorrect value; expected %s, got %s\n", k, v, result)
		}
	}
	result, _, _ := modified.Get(key)
	if result != value {
		t.Fatalf("At %d, incorrect value; expected %s, got %s\n", key, value, result)
	}
}

/*
func Test_Hashmap_Insert_NilKey(t *testing.T) {
	contents := map[int]string{1: "a"}
	original := NewIntToStringHashmap(contents)
	modified, err := original.Insert(nil, "aa")
	if err == nil {
		t.Fatal("Insert with nil key did not return error")
	}
	if modified != nil {
		t.Fatalf("Insert with nil key returned hashmap %s\n", modified)
	}
}

func Test_Hashmap_Insert_NilValue(t *testing.T) {
	contents := map[core.Element]core.Element{integers.IntElement(1): strings.StringElement("a")}
	original := NewHashMap(contents, integers.WithIntKeyMetadata, strings.WithStringPointerValueMetadata)
	modified, err := original.Insert(integers.IntElement(2), nil)
	if err != nil {
		t.Fatalf("Received error during insert: %s\n", err)
	}
	size := modified.Size()
	if size != len(contents)+1 {
		t.Fatalf("Modified hash map contains wrong count; got %d, expected %d\n", size, len(contents)+1)
	}
	result, ok, _ := modified.Get(integers.IntElement(2))
	if !ok {
		t.Fatal("Get returned not-ok")
	}
	if result != nil {
		t.Fatalf("Get returned wrong value; expected nil, got %s\n", result)
	}
}
*/

func Test_Hashmap_Insert_WithSameKey(t *testing.T) {
	contents := map[int]string{1: "a", 2: "b"}
	original := NewIntToStringHashmap(contents)
	modified, err := original.Insert(1, "aa")
	if err != nil {
		t.Fatalf("Received error during insert: %s\n", err)
	}
	if modified == nil {
		t.Fatal("Received nil from insert\n")
	}
	count := modified.Size()
	if count != len(contents) {
		t.Fatalf("Modified hash map contains wrong count; got %d, expected %d\n", count, len(contents))
	}
	result, _, _ := modified.Get(1)
	if result != "aa" {
		t.Fatalf("Modified hash map has wrong value for key; got %s, expected 'aa'\n", result)
	}
}

/*
func Test_Hashmap_Insert_NilValue_WithSameKey(t *testing.T) {
	contents := map[int]string{1: "a", 2: "b"}
	original := NewIntToStringHashmap(contents)
	modified, err := original.Insert(1, nil)
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
	result, _, _ := modified.Get(1)
	if result != nil {
		t.Fatalf("Modified hash map has wrong value for key; got %s, expected 'aa'\n", result)
	}
}
*/

func Test_Hashmap_Insert_WithSameKeyAndSameValue(t *testing.T) {
	contents := map[int]string{1: "a"}
	original := NewIntToStringHashmap(contents)
	modified, err := original.Insert(1, "a")
	if err != nil {
		t.Fatalf("Received error during insert: %s\n", err)
	}
	if original != modified {
		t.Fatalf("Hash map returned from insert is not same hash map\n")
	}
}

// func Test_Hashmap_Insert_SharesOptions(t *testing.T) {
// 	original := NewHashMap(map[core.Element]core.Element{integers.IntElement(1): strings.StringElement("a")}, integers.WithIntKeyMetadata, strings.WithStringValueMetadata)
// 	modified, _ := original.Insert(integers.IntElement(2), strings.StringElement("b"))
//
// 	if modified.options != original.options {
// 		t.Fatalf("Created new options object from mutating function\n")
// 	}
// }
//
func Test_Hashmap_Get_WithUnassigned(t *testing.T) {
	var original *IntToStringHashmap
	result, ok, err := original.Get(2)
	if err == nil {
		t.Fatal("Get from nil did not return error\n")
	}
	if ok {
		t.Fatal("Get from nil returned ok\n")
	}
	if result != "" {
		t.Fatalf("Request from nil hashmap returned %s", result)
	}
}

func Test_Hashmap_Get_WithEmpty(t *testing.T) {
	original := NewIntToStringHashmap(map[int]string{})
	result, ok, err := original.Get(2)
	if err != nil {
		t.Fatalf("Get from empty hashmap returned error '%s'\n", err)
	}
	if ok {
		t.Fatal("Get from empty hashmap returned ok")
	}
	if result != "" {
		t.Fatalf("Request from empty hashmap returned %s", result)
	}
}

func Test_Hashmap_Get_WithContents(t *testing.T) {
	value := "a"
	original := NewIntToStringHashmap(map[int]string{1: value})
	result, ok, err := original.Get(1)
	if err != nil {
		t.Fatalf("Get with contents returned error '%s'\n", err)
	}
	if !ok {
		t.Fatal("Get with contents returned !ok")
	}
	if result != value {
		t.Fatalf("Expected %s, got %s", value, result)
	}
}

func Test_Hashmap_Get_Miss(t *testing.T) {
	original := NewIntToStringHashmap(map[int]string{1: "a"})
	result, ok, err := original.Get(2)
	if err != nil {
		t.Fatalf("Get with contents returned error '%s'\n", err)
	}
	if ok {
		t.Fatal("Get from nil returned ok\n")
	}
	if result != "" {
		t.Fatalf("Request for miss key returned %s\n", result)
	}
}

func Test_Hashmap_Remove_WithUnassigned(t *testing.T) {
	var original *IntToStringHashmap
	key := 4
	modified, err := original.Remove(key)
	if nil == err {
		t.Fatalf("Remove from nil hashmap did not return error\n")
	}
	if nil != modified {
		t.Fatal("Remove from nil hashmap did not return nil\n")
	}
}

func Test_Hashmap_Remove_WithEmpty(t *testing.T) {
	original := NewIntToStringHashmap(map[int]string{})
	key := 4
	modified, error := original.Remove(key)
	if nil != error {
		t.Fatalf("Remove from empty hashmap returned error %s\n", error)
	}
	if nil == modified {
		t.Fatal("Remove from empty hashmap returned nil\n")
	}
	size := modified.Size()
	if size != 0 {
		t.Fatalf("Remove from empty hashmap returned a non-empty hashmap: %s\n", modified)
	}
	if modified != original {
		t.Fatal("Modified hashmap is not the same as original")
	}
}

func Test_Hashmap_Remove_WithContents(t *testing.T) {
	key1, key2 := 0, 1
	value1, value2 := "a", "b"
	contents := map[int]string{
		key1: value1,
		key2: value2,
	}
	original := NewIntToStringHashmap(contents)
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
	result1, ok1, err1 := modified.Get(key1)
	if err1 != nil {
		t.Fatalf("Got err when getting key %d: %s\n", key1, err1)
	}
	if !ok1 {
		t.Fatalf("Got !ok when getting key %d\n", key1)
	}
	if result1 != value1 {
		t.Fatalf("Element returned from key %d is incorrect; expected %s, got %s\n", key1, value1, result1)
	}
	result2, ok2, err2 := modified.Get(key2)
	if err2 != nil {
		t.Fatalf("Got err when getting key %d: %s\n", key2, err2)
	}
	if ok2 {
		t.Fatalf("Got !ok when getting key %d\n", key2)
	}
	if result2 != "" {
		t.Fatalf("Got value from key %d; got %s, expected empty string\n", key2, result2)
	}
}

func Test_Hashmap_Remove_WithContents_ToEmpty(t *testing.T) {
	key1 := 0
	value1 := "a"
	contents := map[int]string{
		key1: value1,
	}
	original := NewIntToStringHashmap(contents)
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
	contents := map[int]string{
		0: "a",
		1: "b",
	}
	original := NewIntToStringHashmap(contents)
	modified, error := original.Remove(2)
	if error != nil {
		t.Fatalf("Error returned from remove: %s", error)
	}
	if modified != original {
		t.Fatalf("Return from remove did not return the same instance")
	}
}

func Test_Hashmap_ReadAndWriteLargeDataSet(t *testing.T) {
	max := 10000
	contents := make(map[int]string, max)
	for i := 0; i < max; i++ {
		contents[i] = fmt.Sprintf("#%d", i)
	}

	original := NewIntToStringHashmap(contents)
	for i := 0; i < max; i++ {
		k := i
		result, _, _ := original.Get(k)
		v := contents[k]
		if result != v {
			t.Fatalf("At %d; expected %s, got %s\n", k, v, result)
		}
	}
}

func Test_Hashmap_String_WithUnassigned(t *testing.T) {
	var original *IntToStringHashmap

	s := original.String()
	if s != "(nil)" {
		t.Fatal("Got wrong string\n")
	}

	s = original.GoString()
	if s != "(nil)" {
		t.Fatalf("Got wring go-string\n")
	}
}

func Test_Hashmap_String_WithContents(t *testing.T) {
	contents := map[int]string{
		0: "a",
		1: "b",
	}
	original := NewIntToStringHashmap(contents)

	s := original.String()
	if len(s) == 0 {
		t.Fatal("Failed to get string\n")
	}

	s = original.GoString()
	if len(s) == 0 {
		t.Fatal("Failed to get go-string\n")
	}
}

func Test_Hashmap_GetKeys_WithUnassigned(t *testing.T) {
	var original *IntToStringHashmap
	keys, err := original.GetKeys()
	if err == nil {
		t.Fatalf("Got nil err from unassigned hashmap\n")
	}
	if keys != nil {
		t.Fatalf("Got keys back from unassigned hashmap: %d\n", keys)
	}
}

func Test_Hashmap_GetKeys_WithEmpty(t *testing.T) {
	original := NewIntToStringHashmap(map[int]string{})
	keys, err := original.GetKeys()
	if err != nil {
		t.Fatalf("Got error '%s'\n", err)
	}
	if keys == nil {
		t.Fatal("Got nil back from empty hashmap\n")
	}
	if len(keys) != 0 {
		t.Fatalf("Got keys back from empty hashmap: %d\n", keys)
	}
}

func Test_Hashmap_GetKeys_WithContents(t *testing.T) {
	contents := map[int]string{
		1: "a",
		2: "b",
	}
	original := NewIntToStringHashmap(contents)
	keys, err := original.GetKeys()
	if err != nil {
		t.Fatalf("Got error '%s'\n", err)
	}
	if keys == nil {
		t.Fatal("Got nil back from empty hashmap\n")
	}
	if len(keys) != len(contents) {
		t.Fatalf("Got wrong number of keys back from empty hashmap: %d\n", keys)
	}
	for k := range contents {
		found := false
		for i := 0; i < len(keys); i++ {
			if keys[i] == k {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Element %d not found in keys\n", k)
		}
	}
}
