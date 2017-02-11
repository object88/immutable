package immutable_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap_Map_WithUnassigned(t *testing.T) {
	var original *immutable.IntToStringHashmap
	invokeCount := 0
	modified, err := original.Map(func(k int, v string) (string, error) {
		invokeCount++
		return fmt.Sprintf("%s%s", v), nil
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

func Test_Hashmap_Map_WithEmpty(t *testing.T) {
	contents := map[int]string{}
	original := immutable.NewIntToStringHashmap(contents)
	invokeCount := 0
	modified, err := original.Map(func(k int, v string) (string, error) {
		invokeCount++
		return fmt.Sprintf("%s%s", v), nil
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
	contents := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	original := immutable.NewIntToStringHashmap(contents)
	invokeCount := 0
	modified, err := original.Map(func(k int, v string) (string, error) {
		invokeCount++
		return fmt.Sprintf("%s%s", v), nil
	})
	if err != nil {
		t.Error(err)
	}
	for k, v := range contents {
		result, _, _ := modified.Get(k)
		expected := fmt.Sprintf("%s%s", v)
		if result != expected {
			t.Fatalf("At %s, got incorrect result, expected %s, got %s\n", k, expected, result)
		}
	}
	if invokeCount != 3 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Map_WithCancel(t *testing.T) {
	contents := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	original := immutable.NewIntToStringHashmap(contents)
	modified, err := original.Map(func(k int, v string) (string, error) {
		if k%2 == 0 {
			return "", errors.New("Found an even key")
		}
		return fmt.Sprintf("%s%s", v), nil
	})
	if err == nil {
		t.Fatalf("Failed to return error")
	}
	if modified != nil {
		t.Fatal("Did not return nil collection")
	}
}
