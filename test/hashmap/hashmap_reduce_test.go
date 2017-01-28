package immutable

import (
	"errors"
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap_Reduce_WithUnassigned(t *testing.T) {
	var original *immutable.HashMap
	invokeCount := 0
	sum, err := original.Reduce(func(acc immutable.Element, k immutable.Element, v immutable.Element) (immutable.Element, error) {
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
	contents := map[immutable.Element]immutable.Element{}
	original := immutable.NewHashMap(contents)
	invokeCount := 0
	sum, err := original.Reduce(func(acc immutable.Element, k immutable.Element, v immutable.Element) (immutable.Element, error) {
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
	contents := map[immutable.Element]immutable.Element{
		immutable.IntElement(1): 1,
		immutable.IntElement(2): 2,
		immutable.IntElement(3): 3,
	}
	original := immutable.NewHashMap(contents)
	invokeCount := 0
	sum, err := original.Reduce(func(acc immutable.Element, k immutable.Element, v immutable.Element) (immutable.Element, error) {
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
	contents := map[immutable.Element]immutable.Element{
		immutable.IntElement(1): 1,
		immutable.IntElement(2): 2,
		immutable.IntElement(3): 3,
	}
	original := immutable.NewHashMap(contents)
	sum, err := original.Reduce(func(acc immutable.Element, k immutable.Element, v immutable.Element) (immutable.Element, error) {
		if k.(immutable.IntElement)%2 == 0 {
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
