package immutable

import (
	"fmt"
	"testing"
)

type MyKey struct {
	value int
}

func (k MyKey) Hash(seed uint32) uint64 {
	return uint64(k.value)
}

func (k MyKey) String() string {
	return fmt.Sprintf("%d", k.value)
}

func Test_HashMap_CustomKey(t *testing.T) {
	data := map[Key]Value{MyKey{1}: "a", MyKey{2}: "b"}
	original := NewHashMap(data, nil)
	if original == nil {
		t.Fatal("NewHashMap returned nil\n")
	}

	size := original.Size()
	if size != len(data) {
		t.Fatalf("Incorrect size; expected %d, got %d\n", len(data), size)
	}
}
