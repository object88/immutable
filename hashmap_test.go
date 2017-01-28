package immutable

import (
	"fmt"
	"testing"
)

func Test_Hashmap(t *testing.T) {
	data := map[Element]Element{
		IntElement(1):  StringElement("aa"),
		IntElement(2):  StringElement("bb"),
		IntElement(3):  StringElement("cc"),
		IntElement(4):  StringElement("dd"),
		IntElement(5):  StringElement("ee"),
		IntElement(6):  StringElement("ff"),
		IntElement(7):  StringElement("gg"),
		IntElement(8):  StringElement("hh"),
		IntElement(26): StringElement("zz"),
	}
	original := NewHashMap(data)
	if original.Size() != len(data) {
		t.Fatalf("Incorrect size")
	}
	for k, v := range data {
		result, ok, err := original.Get(k)
		if err != nil {
			t.Fatalf("Got err %s for key '%s'\n", err.Error(), k)
		}
		if !ok {
			t.Fatalf("Got !ok for key '%s'\n", k)
		}
		if result != v {
			t.Fatalf("Element '%s' -> '%s', expected '%s'\n", k, result, v)
		}
	}
}

func Test_Hashmap_Simple(t *testing.T) {
	contents := map[Element]Element{
		IntElement(1): StringElement("a"),
		IntElement(2): StringElement("b"),
	}
	original := NewHashMap(contents, WithIntegerKeyMetadata, WithStringValueMetadata)
	fmt.Printf("%#v", original)
	for k, v := range contents {
		result, _, _ := original.Get(k)
		if result != v {
			t.Fatalf("Element '%s' -> '%s', expected '%s'\n", k, result, v)
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
	original := NewHashMap(map[Element]Element{})
	if original == nil {
		t.Fatalf("Creating HashMap with empty map returned nil\n")
	}
}

func Test_HashMap_Empty_Size(t *testing.T) {
	original := NewHashMap(map[Element]Element{})
	size := original.Size()
	if size != 0 {
		t.Fatalf("Expected 0 size, got %d\n", size)
	}
}

func Test_Hashmap_Create_WithNilContents(t *testing.T) {
	original := NewHashMap(nil)
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
	modified, err := original.Insert(IntElement(4), StringElement("d"))
	if err == nil {
		t.Fatal("Insert to nil hashmap did not return error\n")
	}
	if modified != nil {
		t.Fatal("Insert to nil hashmap did not return nil\n")
	}
}

func Test_Hashmap_Insert_WithEmpty(t *testing.T) {
	original := NewHashMap(map[Element]Element{})
	key := IntElement(4)
	value := StringElement("d")
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
	contents := map[Element]Element{
		IntElement(0): StringElement("a"),
		IntElement(1): StringElement("b"),
		IntElement(2): StringElement("c"),
	}
	original := NewHashMap(contents)
	key := IntElement(4)
	value := StringElement("d")
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
			t.Fatalf("At %s, incorrect value; expected %s, got %s\n", k, v, result)
		}
	}
	result, _, _ := modified.Get(key)
	if result != value {
		t.Fatalf("At %s, incorrect value; expected %s, got %s\n", key, value, result)
	}
}

func Test_Hashmap_Insert_NilKey(t *testing.T) {
	contents := map[Element]Element{IntElement(1): StringElement("a")}
	original := NewHashMap(contents)
	modified, err := original.Insert(nil, StringElement("aa"))
	if err == nil {
		t.Fatal("Insert with nil key did not return error")
	}
	if modified != nil {
		t.Fatalf("Insert with nil key returned hashmap %s\n", modified)
	}
}

func Test_Hashmap_Insert_NilValue(t *testing.T) {
	contents := map[Element]Element{IntElement(1): StringElement("a")}
	original := NewHashMap(contents)
	modified, err := original.Insert(IntElement(2), nil)
	if err != nil {
		t.Fatalf("Received error during insert: %s\n", err)
	}
	size := modified.Size()
	if size != len(contents)+1 {
		t.Fatalf("Modified hash map contains wrong count; got %d, expected %d\n", size, len(contents)+1)
	}
	result, ok, _ := modified.Get(IntElement(2))
	if !ok {
		t.Fatal("Get returned not-ok")
	}
	if result != nil {
		t.Fatalf("Get returned wrong value; expected nil, got %s\n", result)
	}
}

func Test_Hashmap_Insert_WithSameKey(t *testing.T) {
	contents := map[Element]Element{IntElement(1): StringElement("a"), IntElement(2): StringElement("b")}
	original := NewHashMap(contents)
	modified, err := original.Insert(IntElement(1), StringElement("aa"))
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
	result, _, _ := modified.Get(IntElement(1))
	if result != StringElement("aa") {
		t.Fatalf("Modified hash map has wrong value for key; got %s, expected 'aa'\n", result)
	}
}

func Test_Hashmap_Insert_NilValue_WithSameKey(t *testing.T) {
	contents := map[Element]Element{IntElement(1): StringElement("a"), IntElement(2): StringElement("b")}
	original := NewHashMap(contents)
	modified, err := original.Insert(IntElement(1), nil)
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
	result, _, _ := modified.Get(IntElement(1))
	if result != nil {
		t.Fatalf("Modified hash map has wrong value for key; got %s, expected 'aa'\n", result)
	}
}

func Test_Hashmap_Insert_WithSameKeyAndSameValue(t *testing.T) {
	contents := map[Element]Element{IntElement(1): StringElement("a")}
	original := NewHashMap(contents)
	modified, err := original.Insert(IntElement(1), StringElement("a"))
	if err != nil {
		t.Fatalf("Received error during insert: %s\n", err)
	}
	if original != modified {
		t.Fatalf("Hash map returned from insert is not same hash map\n")
	}
}

func Test_Hashmap_Get_WithUnassigned(t *testing.T) {
	var original *HashMap
	result, ok, err := original.Get(IntElement(2))
	if err == nil {
		t.Fatal("Get from nil did not return error\n")
	}
	if ok {
		t.Fatal("Get from nil returned ok\n")
	}
	if result != nil {
		t.Fatalf("Request from nil hashmap returned %s", result)
	}
}

func Test_Hashmap_Get_WithEmpty(t *testing.T) {
	original := NewHashMap(map[Element]Element{})
	result, ok, err := original.Get(IntElement(2))
	if err != nil {
		t.Fatalf("Get from empty hashmap returned error '%s'\n", err)
	}
	if ok {
		t.Fatal("Get from empty hashmap returned ok")
	}
	if result != nil {
		t.Fatalf("Request from empty hashmap returned %s", result)
	}
}

func Test_Hashmap_Get_WithContents(t *testing.T) {
	value := StringElement("a")
	original := NewHashMap(map[Element]Element{IntElement(1): value})
	result, ok, err := original.Get(IntElement(1))
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
	original := NewHashMap(map[Element]Element{IntElement(1): StringElement("a")})
	result, ok, err := original.Get(IntElement(2))
	if err != nil {
		t.Fatalf("Get with contents returned error '%s'\n", err)
	}
	if ok {
		t.Fatal("Get from nil returned ok\n")
	}
	if result != nil {
		t.Fatalf("Request for miss key returned %s\n", result)
	}
}

func Test_Hashmap_Remove_WithUnassigned(t *testing.T) {
	var original *HashMap
	key := IntElement(4)
	modified, err := original.Remove(key)
	if nil == err {
		t.Fatalf("Remove from nil hashmap did not return error\n")
	}
	if nil != modified {
		t.Fatal("Remove from nil hashmap did not return nil\n")
	}
}

func Test_Hashmap_Remove_WithEmpty(t *testing.T) {
	original := NewHashMap(map[Element]Element{})
	key := IntElement(4)
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
	key1, key2 := IntElement(0), IntElement(1)
	value1, value2 := StringElement("a"), StringElement("b")
	contents := map[Element]Element{
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
	if size != len(contents)-1 {
		t.Fatalf("Incorrect number of entries in returned collection; expected %d, got %d\n", len(contents)-1, size)
	}
	result1, _, _ := modified.Get(key1)
	if result1 != value1 {
		t.Fatalf("Element returned from key %s is incorrect; expected %s, got %s\n", key1, value1, result1)
	}
	result2, _, _ := modified.Get(key2)
	if result2 != nil {
		t.Fatalf("Got value from key %s; got %s, expected nil\n", key2, result2)
	}
}

func Test_Hashmap_Remove_WithContents_ToEmpty(t *testing.T) {
	key1 := IntElement(0)
	value1 := StringElement("a")
	contents := map[Element]Element{
		key1: value1,
	}
	original := NewHashMap(contents)
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
	contents := map[Element]Element{
		IntElement(0): StringElement("a"),
		IntElement(1): StringElement("b"),
	}
	original := NewHashMap(contents)
	modified, error := original.Remove(IntElement(2))
	if error != nil {
		t.Fatalf("Error returned from remove: %s", error)
	}
	if modified != original {
		t.Fatalf("Return from remove did not return the same instance")
	}
}

func Test_Hashmap_ReadAndWriteLargeDataSet(t *testing.T) {
	max := 10000
	contents := make(map[Element]Element, max)
	for i := 0; i < max; i++ {
		contents[IntElement(i)] = IntElement(i)
	}

	original := NewHashMap(contents)
	for i := 0; i < max; i++ {
		k := IntElement(i)
		result, _, _ := original.Get(k)
		v := contents[k]
		if result != v {
			t.Fatalf("At %s; expected %s, got %s\n", k, v, result)
		}
	}
}

func Test_Hashmap_String_WithUnassigned(t *testing.T) {
	var original *HashMap

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
	contents := map[Element]Element{
		IntElement(0): StringElement("a"),
		IntElement(1): StringElement("b"),
	}
	original := NewHashMap(contents)

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
	var original *HashMap
	keys, err := original.GetKeys()
	if err == nil {
		t.Fatalf("Got nil err from unassigned hashmap\n")
	}
	if keys != nil {
		t.Fatalf("Got keys back from unassigned hashmap: %s\n", keys)
	}
}

func Test_Hashmap_GetKeys_WithEmpty(t *testing.T) {
	original := NewHashMap(map[Element]Element{})
	keys, err := original.GetKeys()
	if err != nil {
		t.Fatalf("Got error '%s'\n", err)
	}
	if keys == nil {
		t.Fatal("Got nil back from empty hashmap\n")
	}
	if len(keys) != 0 {
		t.Fatalf("Got keys back from empty hashmap: %s\n", keys)
	}
}

func Test_Hashmap_GetKeys_WithContents(t *testing.T) {
	contents := map[Element]Element{
		IntElement(1): StringElement("a"),
		IntElement(2): StringElement("b"),
	}
	original := NewHashMap(contents)
	keys, err := original.GetKeys()
	if err != nil {
		t.Fatalf("Got error '%s'\n", err)
	}
	if keys == nil {
		t.Fatal("Got nil back from empty hashmap\n")
	}
	if len(keys) != len(contents) {
		t.Fatalf("Got wrong number of keys back from empty hashmap: %s\n", keys)
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
			t.Fatalf("Element %s not found in keys\n", k)
		}
	}
}
