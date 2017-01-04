package immutable

import (
	"errors"
	"testing"
)

func Test_Hashmap_Reduce_WithUnassigned(t *testing.T) {
	var original *HashMap
	invokeCount := 0
	sum, err := original.Reduce(func(acc Value, k Key, v Value) (Value, error) {
		invokeCount++
		return acc.(int) + v.(int), nil
	}, 0)
	if err != nil {
		t.Error(err)
	}
	if sum != nil {
		t.Fatal("Did not return nil")
	}
	if invokeCount != 0 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Reduce_WithEmpty(t *testing.T) {
	contents := map[Key]Value{}
	original := NewHashMap(contents)
	invokeCount := 0
	sum, err := original.Reduce(func(acc Value, k Key, v Value) (Value, error) {
		invokeCount++
		return acc.(int) + v.(int), nil
	}, 0)
	if err != nil {
		t.Error(err)
	}
	if sum != 0 {
		t.Fatal("Did not return initial accumulator")
	}
	if invokeCount != 0 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Reduce_WithContents(t *testing.T) {
	contents := map[Key]Value{
		IntKey(1): 1,
		IntKey(2): 2,
		IntKey(3): 3,
	}
	original := NewHashMap(contents)
	invokeCount := 0
	sum, err := original.Reduce(func(acc Value, k Key, v Value) (Value, error) {
		invokeCount++
		return acc.(int) + v.(int), nil
	}, 0)
	if err != nil {
		t.Error(err)
	}
	if sum != 6 {
		t.Fatal("Did not return expected accumulator")
	}
	if invokeCount != 3 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Reduce_WithCancel(t *testing.T) {
	contents := map[Key]Value{
		IntKey(1): 1,
		IntKey(2): 2,
		IntKey(3): 3,
	}
	original := NewHashMap(contents)
	sum, err := original.Reduce(func(acc Value, k Key, v Value) (Value, error) {
		if k.(IntKey)%2 == 0 {
			return nil, errors.New("Found an even key")
		}
		return acc.(int) + v.(int), nil
	}, 0)
	if err == nil {
		t.Fatalf("Failed to return error")
	}
	if sum != nil {
		t.Fatal("Did not return nil accumulator")
	}
}
