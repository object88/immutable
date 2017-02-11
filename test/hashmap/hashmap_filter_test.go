package immutable_test

import (
	"errors"
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap_Filter_WithUnassigned(t *testing.T) {
	var original *immutable.IntToStringHashmap
	invokeCount := 0
	modified, err := original.Filter(func(k int, v string) (bool, error) {
		invokeCount++
		return k%2 == 0, nil
	})
	if err == nil {
		t.Error("No error returned")
	}
	if modified != nil {
		t.Fatal("Did not return nil")
	}
	if invokeCount != 0 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Filter_WithEmpty(t *testing.T) {
	contents := map[int]string{}
	original := immutable.NewIntToStringHashmap(contents)
	invokeCount := 0
	modified, err := original.Filter(func(k int, v string) (bool, error) {
		invokeCount++
		return k%2 == 0, nil
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
	contents := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	original := immutable.NewIntToStringHashmap(contents)
	invokeCount := 0
	modified, err := original.Filter(func(k int, v string) (bool, error) {
		invokeCount++
		return k%2 == 0, nil
	})
	if err != nil {
		t.Error(err)
	}
	size := modified.Size()
	if size != 1 {
		t.Fatalf("Incorrect number of elements in new collection; expected 1, got %d\n", size)
	}
	value, _, _ := modified.Get(2)
	if value != "2" {
		t.Fatalf("Incorrect contents of new collection:\n%s\n", modified)
	}
	if invokeCount != 3 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Filter_WithCancel(t *testing.T) {
	contents := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	original := immutable.NewIntToStringHashmap(contents)
	modified, err := original.Filter(func(k int, v string) (bool, error) {
		if k > 1 {
			return false, errors.New("Found a key greater than 1")
		}
		return k%2 == 0, nil
	})
	if err == nil {
		t.Fatalf("Failed to return error")
	}
	if modified != nil {
		t.Fatal("Did not return nil collection")
	}
}
