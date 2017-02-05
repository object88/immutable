package immutable_test

import (
	"errors"
	"testing"

	"github.com/object88/immutable"
	"github.com/object88/immutable/core"
	"github.com/object88/immutable/handlers/integers"
)

func Test_Hashmap_Filter_WithUnassigned(t *testing.T) {
	var original *immutable.HashMap
	invokeCount := 0
	modified, err := original.Filter(func(k core.Element, v core.Element) (bool, error) {
		invokeCount++
		return int(v.(integers.IntElement))%2 == 0, nil
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
	contents := map[core.Element]core.Element{}
	original := immutable.NewHashMap(contents, integers.WithIntKeyMetadata, integers.WithIntValueMetadata)
	invokeCount := 0
	modified, err := original.Filter(func(k core.Element, v core.Element) (bool, error) {
		invokeCount++
		return int(v.(integers.IntElement))%2 == 0, nil
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
	contents := map[core.Element]core.Element{
		integers.IntElement(1): integers.IntElement(1),
		integers.IntElement(2): integers.IntElement(2),
		integers.IntElement(3): integers.IntElement(3),
	}
	original := immutable.NewHashMap(contents, integers.WithIntKeyMetadata, integers.WithIntValueMetadata)
	invokeCount := 0
	modified, err := original.Filter(func(k core.Element, v core.Element) (bool, error) {
		invokeCount++
		return int(v.(integers.IntElement))%2 == 0, nil
	})
	if err != nil {
		t.Error(err)
	}
	size := modified.Size()
	if size != 1 {
		t.Fatalf("Incorrect number of elements in new collection; expected 1, got %d\n", size)
	}
	value, _, _ := modified.Get(integers.IntElement(2))
	if value == nil || int(value.(integers.IntElement)) != 2 {
		t.Fatalf("Incorrect contents of new collection:\n%s\n", modified)
	}
	if invokeCount != 3 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Filter_WithCancel(t *testing.T) {
	contents := map[core.Element]core.Element{
		integers.IntElement(1): integers.IntElement(1),
		integers.IntElement(2): integers.IntElement(2),
		integers.IntElement(3): integers.IntElement(3),
	}
	original := immutable.NewHashMap(contents, integers.WithIntKeyMetadata, integers.WithIntValueMetadata)
	modified, err := original.Filter(func(k core.Element, v core.Element) (bool, error) {
		if k.(integers.IntElement)%2 == 0 {
			return false, errors.New("Found an even key")
		}
		return int(v.(integers.IntElement))%2 == 0, nil
	})
	if err == nil {
		t.Fatalf("Failed to return error")
	}
	if modified != nil {
		t.Fatal("Did not return nil collection")
	}
}
