package immutable_test

import (
	"errors"
	"testing"

	"github.com/object88/immutable"
	"github.com/object88/immutable/core"
	"github.com/object88/immutable/handlers/integers"
)

func Test_Hashmap_Reduce_WithUnassigned(t *testing.T) {
	var original *immutable.HashMap
	invokeCount := 0
	sum, err := original.Reduce(func(acc core.Element, k core.Element, v core.Element) (core.Element, error) {
		invokeCount++
		return integers.IntElement(int(acc.(integers.IntElement)) + int(v.(integers.IntElement))), nil
	}, integers.IntElement(0))
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
	contents := map[core.Element]core.Element{}
	original := immutable.NewHashMap(contents, integers.WithIntKeyMetadata, integers.WithIntValueMetadata)
	invokeCount := 0
	sum, err := original.Reduce(func(acc core.Element, k core.Element, v core.Element) (core.Element, error) {
		invokeCount++
		return integers.IntElement(int(acc.(integers.IntElement)) + int(v.(integers.IntElement))), nil
	}, integers.IntElement(0))
	if err != nil {
		t.Error(err)
	}
	if int(sum.(integers.IntElement)) != 0 {
		t.Fatal("Did not return initial accumulator")
	}
	if invokeCount != 0 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Reduce_WithContents(t *testing.T) {
	contents := map[core.Element]core.Element{
		integers.IntElement(1): integers.IntElement(1),
		integers.IntElement(2): integers.IntElement(2),
		integers.IntElement(3): integers.IntElement(3),
	}
	original := immutable.NewHashMap(contents, integers.WithIntKeyMetadata, integers.WithIntValueMetadata)
	invokeCount := 0
	sum, err := original.Reduce(func(acc core.Element, k core.Element, v core.Element) (core.Element, error) {
		invokeCount++
		return integers.IntElement(int(acc.(integers.IntElement)) + int(v.(integers.IntElement))), nil
	}, integers.IntElement(0))
	if err != nil {
		t.Error(err)
	}
	if int(sum.(integers.IntElement)) != 6 {
		t.Fatal("Did not return expected accumulator")
	}
	if invokeCount != 3 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Reduce_WithCancel(t *testing.T) {
	contents := map[core.Element]core.Element{
		integers.IntElement(1): integers.IntElement(1),
		integers.IntElement(2): integers.IntElement(2),
		integers.IntElement(3): integers.IntElement(3),
	}
	original := immutable.NewHashMap(contents, integers.WithIntKeyMetadata, integers.WithIntValueMetadata)
	sum, err := original.Reduce(func(acc core.Element, k core.Element, v core.Element) (core.Element, error) {
		if k.(integers.IntElement)%2 == 0 {
			return nil, errors.New("Found an even key")
		}
		return integers.IntElement(int(acc.(integers.IntElement)) + int(v.(integers.IntElement))), nil
	}, integers.IntElement(0))
	if err == nil {
		t.Fatalf("Failed to return error")
	}
	if sum != nil {
		t.Fatal("Did not return nil accumulator")
	}
}
