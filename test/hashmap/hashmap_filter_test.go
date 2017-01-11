package immutable_test

import (
	"errors"
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap_Filter_WithUnassigned(t *testing.T) {
	var original *immutable.HashMap
	invokeCount := 0
	modified, err := original.Filter(func(k immutable.Key, v immutable.Value) (bool, error) {
		invokeCount++
		return v.(int)%2 == 0, nil
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

func Test_Hashmap_Filter_WithEmpty(t *testing.T) {
	contents := map[immutable.Key]immutable.Value{}
	original := immutable.NewHashMap(contents, nil)
	invokeCount := 0
	modified, err := original.Filter(func(k immutable.Key, v immutable.Value) (bool, error) {
		invokeCount++
		return v.(int)%2 == 0, nil
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

func Test_Hashmap_Filter_WithContents(t *testing.T) {
	contents := map[immutable.Key]immutable.Value{
		immutable.IntKey(1): 1,
		immutable.IntKey(2): 2,
		immutable.IntKey(3): 3,
	}
	original := immutable.NewHashMap(contents, nil)
	invokeCount := 0
	modified, err := original.Filter(func(k immutable.Key, v immutable.Value) (bool, error) {
		invokeCount++
		return v.(int)%2 == 0, nil
	})
	if err != nil {
		t.Error(err)
	}
	size := modified.Size()
	if size != 1 {
		t.Fatalf("Incorrect number of elements in new collection; expected 1, got %d\n", size)
	}
	value := modified.Get(immutable.IntKey(2))
	if value == nil || value.(int) != 2 {
		t.Fatalf("Incorrect contents of new collection:\n%s\n", modified)
	}
	if invokeCount != 3 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Filter_WithCancel(t *testing.T) {
	contents := map[immutable.Key]immutable.Value{
		immutable.IntKey(1): 1,
		immutable.IntKey(2): 2,
		immutable.IntKey(3): 3,
	}
	original := immutable.NewHashMap(contents, nil)
	modified, err := original.Filter(func(k immutable.Key, v immutable.Value) (bool, error) {
		if k.(immutable.IntKey)%2 == 0 {
			return false, errors.New("Found an even key")
		}
		return v.(int)%2 == 0, nil
	})
	if err == nil {
		t.Fatalf("Failed to return error")
	}
	if modified != nil {
		t.Fatal("Did not return nil collection")
	}
}
