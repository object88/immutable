package immutable_test

import (
	"errors"
	"testing"

	"github.com/object88/immutable"
	"github.com/object88/immutable/core"
	"github.com/object88/immutable/handlers/integers"
)

func Test_Hashmap_Map_WithUnassigned(t *testing.T) {
	var original *immutable.HashMap
	invokeCount := 0
	modified, err := original.Map(func(k core.Element, v core.Element) (core.Element, error) {
		invokeCount++
		return integers.IntElement(int(v.(integers.IntElement)) * 2), nil
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
	contents := map[core.Element]core.Element{}
	original := immutable.NewHashMap(contents, integers.WithIntKeyMetadata, integers.WithIntValueMetadata)
	invokeCount := 0
	modified, err := original.Map(func(k core.Element, v core.Element) (core.Element, error) {
		invokeCount++
		return integers.IntElement(int(v.(integers.IntElement)) * 2), nil
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
	contents := map[core.Element]core.Element{
		integers.IntElement(1): integers.IntElement(1),
		integers.IntElement(2): integers.IntElement(2),
		integers.IntElement(3): integers.IntElement(3),
	}
	original := immutable.NewHashMap(contents, integers.WithIntKeyMetadata, integers.WithIntValueMetadata)
	invokeCount := 0
	modified, err := original.Map(func(k core.Element, v core.Element) (core.Element, error) {
		invokeCount++
		return integers.IntElement(int(v.(integers.IntElement)) * 2), nil
	})
	if err != nil {
		t.Error(err)
	}
	for k, v := range contents {
		result, _, _ := modified.Get(k)
		expected := int(v.(integers.IntElement)) * 2
		if int(result.(integers.IntElement)) != expected {
			t.Fatalf("At %s, got incorrect result, expected %d, got %s\n", k, expected, result)
		}
	}
	if invokeCount != 3 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Map_WithCancel(t *testing.T) {
	contents := map[core.Element]core.Element{
		integers.IntElement(1): integers.IntElement(1),
		integers.IntElement(2): integers.IntElement(2),
		integers.IntElement(3): integers.IntElement(3),
	}
	original := immutable.NewHashMap(contents, integers.WithIntKeyMetadata, integers.WithIntValueMetadata)
	modified, err := original.Map(func(k core.Element, v core.Element) (core.Element, error) {
		if k.(integers.IntElement)%2 == 0 {
			return nil, errors.New("Found an even key")
		}
		return integers.IntElement(int(v.(integers.IntElement)) * 2), nil
	})
	if err == nil {
		t.Fatalf("Failed to return error")
	}
	if modified != nil {
		t.Fatal("Did not return nil collection")
	}
}
