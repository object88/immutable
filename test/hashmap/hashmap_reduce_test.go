package immutable_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap_Reduce_WithUnassigned(t *testing.T) {
	var original *immutable.IntToStringHashmap
	invokeCount := 0
	sum, err := original.Reduce(0, func(acc interface{}, k int, v string) (interface{}, error) {
		invokeCount++
		i, _ := strconv.Atoi(v)
		return acc.(int) + i, nil
	})
	if err == nil {
		t.Error("No error returned")
	}
	if sum != nil {
		t.Fatal("Did not return nil")
	}
	if invokeCount != 0 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Reduce_WithEmpty(t *testing.T) {
	contents := map[int]string{}
	original := immutable.NewIntToStringHashmap(contents)
	invokeCount := 0
	sum, err := original.Reduce(0, func(acc interface{}, k int, v string) (interface{}, error) {
		invokeCount++
		i, _ := strconv.Atoi(v)
		return acc.(int) + i, nil
	})
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
	contents := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	original := immutable.NewIntToStringHashmap(contents)
	invokeCount := 0
	sum, err := original.Reduce(0, func(acc interface{}, k int, v string) (interface{}, error) {
		invokeCount++
		i, _ := strconv.Atoi(v)
		return acc.(int) + i, nil
	})
	if err != nil {
		t.Error(err)
	}
	if sum != 6 {
		t.Fatalf("Did not return expected accumulator; got %s; expected 123", sum)
	}
	if invokeCount != 3 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Reduce_WithCancel(t *testing.T) {
	contents := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	original := immutable.NewIntToStringHashmap(contents)
	sum, err := original.Reduce(0, func(acc interface{}, k int, v string) (interface{}, error) {
		if k%2 == 0 {
			return nil, errors.New("Found an even key")
		}
		i, _ := strconv.Atoi(v)
		return acc.(int) + i, nil
	})
	if err == nil {
		t.Fatalf("Failed to return error")
	}
	if sum != nil {
		t.Fatal("Did not return nil accumulator")
	}
}
