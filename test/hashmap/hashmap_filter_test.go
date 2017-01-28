package immutable_test

import (
	"errors"
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap_Filter_WithUnassigned(t *testing.T) {
	var original *immutable.HashMap
	invokeCount := 0
	modified, err := original.Filter(func(k immutable.Element, v immutable.Element) (bool, error) {
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
	contents := map[immutable.Element]immutable.Element{}
	original := immutable.NewHashMap(contents)
	invokeCount := 0
	modified, err := original.Filter(func(k immutable.Element, v immutable.Element) (bool, error) {
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
	contents := map[immutable.Element]immutable.Element{
		immutable.IntElement(1): 1,
		immutable.IntElement(2): 2,
		immutable.IntElement(3): 3,
	}
	original := immutable.NewHashMap(contents)
	invokeCount := 0
	modified, err := original.Filter(func(k immutable.Element, v immutable.Element) (bool, error) {
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
	value, _, _ := modified.Get(immutable.IntElement(2))
	if value == nil || value.(int) != 2 {
		t.Fatalf("Incorrect contents of new collection:\n%s\n", modified)
	}
	if invokeCount != 3 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Filter_WithCancel(t *testing.T) {
	contents := map[immutable.Element]immutable.Element{
		immutable.IntElement(1): 1,
		immutable.IntElement(2): 2,
		immutable.IntElement(3): 3,
	}
	original := immutable.NewHashMap(contents)
	modified, err := original.Filter(func(k immutable.Element, v immutable.Element) (bool, error) {
		if k.(immutable.IntElement)%2 == 0 {
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
