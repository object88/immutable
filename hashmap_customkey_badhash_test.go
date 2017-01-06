package immutable

import (
	"fmt"
	"testing"
)

// This suite of tests is designed to test the bucket overflow behavior.
// By using a `Hash` function which only ever resolves into one of two
// hash keys, all entries will end up in one or two buckets.  With a
// sufficiently large dataset, this will cause bucket overflow to be used.

type MyBadKey struct {
	value int
}

func (k MyBadKey) Hash() uint32 {
	if k.value%2 == 0 {
		return 0x0
	}
	return 0xffffffff
}

func (k MyBadKey) String() string {
	return fmt.Sprintf("%d", k.value)
}

func Test_Hashmap_CustomKey_BadHash_Iterate(t *testing.T) {
	max := 100
	data := make(map[Key]Value, max)
	for i := 0; i < max; i++ {
		data[MyBadKey{i}] = false
	}
	original := NewHashMap(data, nil)
	if original == nil {
		t.Fatal("NewHashMap returned nil\n")
	}
	size := original.Size()
	if size != len(data) {
		t.Fatalf("Incorrect size; expected %d, got %d\n", len(data), size)
	}
	original.ForEach(func(k Key, v Value) {
		if v.(bool) {
			t.Fatalf("At %s, already visited\n", k)
		}
		data[k] = true
	})

	for k, v := range data {
		if !v.(bool) {
			t.Fatalf("At %s, not visited\n", k)
		}
	}
}

type MyIntValue int

func (v MyIntValue) String() string {
	return fmt.Sprintf("%d", v)
}

func Test_Hashmap_CustomKey_BadHash_Get(t *testing.T) {
	max := 100
	contents := make(map[Key]Value, max)
	for i := 0; i < max; i++ {
		contents[MyBadKey{i}] = MyIntValue(i)
	}

	original := NewHashMap(contents, nil)

	for k, v := range contents {
		result := original.Get(k)
		if result != v {
			t.Fatalf("At %s, expected %d, got %d\n%#v", k, v, result, original)
		}
	}
}
