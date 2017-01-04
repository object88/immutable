package immutable

import (
	"fmt"
	"testing"
)

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

func Test_Hashmap_CustomKey_BadHash(t *testing.T) {
	max := 100
	data := make(map[Key]Value, max)
	for i := 0; i < max; i++ {
		data[MyBadKey{i}] = false
	}
	original := NewHashMap(data)
	if original == nil {
		t.Fatal("NewHashMap returned nil\n")
	}
	size := original.Size()
	if size != uint32(len(data)) {
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
