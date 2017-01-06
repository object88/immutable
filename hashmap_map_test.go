package immutable

import (
	"errors"
	"testing"
)

func Test_Hashmap_Map_WithUnassigned(t *testing.T) {
	var original *HashMap
	invokeCount := 0
	modified, err := original.Map(func(k Key, v Value) (Value, error) {
		invokeCount++
		return v.(int) * 2, nil
	})
	if err != nil {
		t.Error(err)
	}
	if modified != nil {
		t.Fatal("Did not return nil")
	}
	if invokeCount != 0 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Map_WithEmpty(t *testing.T) {
	contents := map[Key]Value{}
	original := NewHashMap(contents, nil)
	invokeCount := 0
	modified, err := original.Map(func(k Key, v Value) (Value, error) {
		invokeCount++
		return v.(int) * 2, nil
	})
	if err != nil {
		t.Error(err)
	}
	if modified == nil {
		t.Fatal("Did not return new hashmap")
	}
	if invokeCount != 0 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Map_WithContents(t *testing.T) {
	contents := map[Key]Value{
		IntKey(1): 1,
		IntKey(2): 2,
		IntKey(3): 3,
	}
	original := NewHashMap(contents, nil)
	invokeCount := 0
	modified, err := original.Map(func(k Key, v Value) (Value, error) {
		invokeCount++
		return v.(int) * 2, nil
	})
	if err != nil {
		t.Error(err)
	}
	for k, v := range contents {
		result := modified.Get(k)
		expected := v.(int) * 2
		if result != expected {
			t.Fatalf("At %s, got incorrect result, expected %d, got %d\n", k, expected, result)
		}
	}
	if invokeCount != 3 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Map_WithCancel(t *testing.T) {
	contents := map[Key]Value{
		IntKey(1): 1,
		IntKey(2): 2,
		IntKey(3): 3,
	}
	original := NewHashMap(contents, nil)
	modified, err := original.Map(func(k Key, v Value) (Value, error) {
		if k.(IntKey)%2 == 0 {
			return nil, errors.New("Found an even key")
		}
		return v.(int) * 2, nil
	})
	if err == nil {
		t.Fatalf("Failed to return error")
	}
	if modified != nil {
		t.Fatal("Did not return nil collection")
	}
}
