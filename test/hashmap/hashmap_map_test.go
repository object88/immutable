package immutable

import (
	"errors"
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap_Map_WithUnassigned(t *testing.T) {
	var original *immutable.HashMap
	invokeCount := 0
	modified, err := original.Map(func(k immutable.Key, v immutable.Value) (immutable.Value, error) {
		invokeCount++
		return v.(int) * 2, nil
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
	contents := map[immutable.Key]immutable.Value{}
	original := immutable.NewHashMap(contents)
	invokeCount := 0
	modified, err := original.Map(func(k immutable.Key, v immutable.Value) (immutable.Value, error) {
		invokeCount++
		return v.(int) * 2, nil
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
	contents := map[immutable.Key]immutable.Value{
		immutable.IntKey(1): 1,
		immutable.IntKey(2): 2,
		immutable.IntKey(3): 3,
	}
	original := immutable.NewHashMap(contents)
	invokeCount := 0
	modified, err := original.Map(func(k immutable.Key, v immutable.Value) (immutable.Value, error) {
		invokeCount++
		return v.(int) * 2, nil
	})
	if err != nil {
		t.Error(err)
	}
	for k, v := range contents {
		result, _, _ := modified.Get(k)
		expected := v.(int) * 2
		if result != expected {
			t.Fatalf("At %s, got incorrect result, expected %d, got %d\n", k, expected, result)
		}
	}
	if invokeCount != 3 {
		t.Fatalf("Function invoked %d times", invokeCount)
	}
}

func Test_Hashmap_Map_WithCancel(t *testing.T) {
	contents := map[immutable.Key]immutable.Value{
		immutable.IntKey(1): 1,
		immutable.IntKey(2): 2,
		immutable.IntKey(3): 3,
	}
	original := immutable.NewHashMap(contents)
	modified, err := original.Map(func(k immutable.Key, v immutable.Value) (immutable.Value, error) {
		if k.(immutable.IntKey)%2 == 0 {
			return nil, errors.New("Found an even key")
		}
		return v.(int) * 2, nil
	})
	if err == nil {
		t.Fatalf("Failed to return error")
	}
	if modified != nil {
		t.Fatal("Did not return nil collection")
	}
}
