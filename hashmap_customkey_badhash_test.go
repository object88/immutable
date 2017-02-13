package immutable

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/object88/immutable/core"
	"github.com/object88/immutable/handlers/strings"
)

// This suite of tests is designed to test the bucket overflow behavior.
// By using a `Hash` function which only ever resolves into one of two
// hash keys, all entries will end up in one or two buckets.  With a
// sufficiently large dataset, this will cause bucket overflow to be used.

type MyBadHandler struct{}

func (MyBadHandler) Compare(a, b unsafe.Pointer) (match bool) {
	return *(*int)(a) == *(*int)(b)
}

func (h MyBadHandler) CompareTo(memory unsafe.Pointer, index int, other unsafe.Pointer) (match bool) {
	p := h.Read(memory, index)
	return h.Compare(p, other)
}

func (MyBadHandler) CreateBucket(count int) unsafe.Pointer {
	m := make([]int, count)
	return unsafe.Pointer(&m)
}

func (MyBadHandler) Hash(ep unsafe.Pointer, seed uint32) uint64 {
	e := *(*int)(ep)
	if e%2 == 0 {
		return 0x0
	}
	return 0xffffffffffffffff
}

func (MyBadHandler) Read(memory unsafe.Pointer, index int) (result unsafe.Pointer) {
	m := *(*[]int)(memory)
	return unsafe.Pointer(&(m[index]))
}

func (MyBadHandler) Format(value unsafe.Pointer) string {
	return fmt.Sprintf("%d", *(*int)(value))
}

func (MyBadHandler) Write(memory unsafe.Pointer, index int, value unsafe.Pointer) {
	m := *(*[]int)(memory)
	m[index] = *(*int)(value)
}

func WithMyBadElementMetadata(o *core.HashMapOptions) {
	var handler MyBadHandler
	o.KeyConfig = handler
	strings.WithStringValueMetadata(o)
}

func Test_Hashmap_CustomKey_BadHash_Iterate(t *testing.T) {
	max := 100
	data := make(map[int]string, max)
	for i := 0; i < max; i++ {
		data[i] = "0"
	}
	contents := make(map[unsafe.Pointer]unsafe.Pointer, max)
	for k, v := range data {
		key, value := k, v
		contents[unsafe.Pointer(&key)] = unsafe.Pointer(&value)
	}
	original := NewHashMap(contents, WithMyBadElementMetadata)
	if original == nil {
		t.Fatal("NewHashMap returned nil\n")
	}
	size := original.Size()
	if size != len(data) {
		t.Fatalf("Incorrect size; expected %d, got %d\n", len(data), size)
	}
	original.ForEach(func(kp, _ unsafe.Pointer) {
		k := *(*int)(kp)
		if data[k] == "1" {
			t.Fatalf("At %d, already visited\n", k)
		}
		data[k] = "1"
	})

	for k, v := range data {
		if v != "1" {
			t.Fatalf("At %d, not visited\n", k)
		}
	}
}

func Test_Hashmap_CustomKey_BadHash_Get(t *testing.T) {
	max := 100
	contents := make(map[unsafe.Pointer]unsafe.Pointer, max)
	for i := 0; i < max; i++ {
		s := fmt.Sprintf("%d", i)
		contents[unsafe.Pointer(&i)] = unsafe.Pointer(&s)
	}

	original := NewHashMap(contents, WithMyBadElementMetadata)

	for k, v := range contents {
		result, _, _ := original.Get(k)
		if result != v {
			t.Fatalf("At %d, expected %s, got %s\n%#v", k, *(*string)(v), *(*string)(result), original)
		}
	}
}
