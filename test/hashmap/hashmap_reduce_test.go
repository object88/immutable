package immutable

import (
	"errors"
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap_Reduce_WithUnassigned(t *testing.T) {
	var original *immutable.HashMap
	invokeCount := 0
	sum, err := original.Reduce(func(acc immutable.Value, k immutable.Key, v immutable.Value) (immutable.Value, error) {
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
	contents := map[immutable.Key]immutable.Value{}
	original := immutable.NewHashMap(contents)
	invokeCount := 0
	sum, err := original.Reduce(func(acc immutable.Value, k immutable.Key, v immutable.Value) (immutable.Value, error) {
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
	contents := map[immutable.Key]immutable.Value{
		immutable.IntKey(1): 1,
		immutable.IntKey(2): 2,
		immutable.IntKey(3): 3,
	}
	original := immutable.NewHashMap(contents)
	invokeCount := 0
	sum, err := original.Reduce(func(acc immutable.Value, k immutable.Key, v immutable.Value) (immutable.Value, error) {
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
	contents := map[immutable.Key]immutable.Value{
		immutable.IntKey(1): 1,
		immutable.IntKey(2): 2,
		immutable.IntKey(3): 3,
	}
	original := immutable.NewHashMap(contents)
	sum, err := original.Reduce(func(acc immutable.Value, k immutable.Key, v immutable.Value) (immutable.Value, error) {
		if k.(immutable.IntKey)%2 == 0 {
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
