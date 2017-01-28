package immutable

import (
	"fmt"
	"testing"
)

// This suite of tests is designed to test the bucket overflow behavior.
// By using a `Hash` function which only ever resolves into one of two
// hash keys, all entries will end up in one or two buckets.  With a
// sufficiently large dataset, this will cause bucket overflow to be used.

type MyBadKey int

func (k MyBadKey) Hash(seed uint32) uint64 {
	if k%2 == 0 {
		return 0x0
	}
	return 0xffffffffffffffff
}

func (k MyBadKey) String() string {
	return fmt.Sprintf("%d", k)
}

func WithMyBadKeyMetadata(o *HashMapOptions) {
	o.KeyHandler = NewMyBadElementHashHandler()
	o.ValueHandler = NewIntHandler()
}

type MyBadElementSubBucket struct {
	m []int
}

func NewMyBadElementHashHandler() BucketGenerator {
	return NewMyBadElementSubBucket
}

func NewMyBadElementSubBucket(count int) SubBucket {
	m := make([]int, count)
	return &MyBadElementSubBucket{m}
}

func (b *MyBadElementSubBucket) Hydrate(index int) Element {
	u := b.m[index]
	v := MyBadKey(u)
	return v
}

func (b *MyBadElementSubBucket) Dehydrate(index int, v Element) {
	u := v.(MyBadKey)
	b.m[index] = int(u)
}

func Test_Hashmap_CustomKey_BadHash_Iterate(t *testing.T) {
	max := 100
	data := make(map[Element]Element, max)
	for i := 0; i < max; i++ {
		data[MyBadKey(i)] = IntElement(0)
	}
	original := NewHashMap(data, WithMyBadKeyMetadata)
	if original == nil {
		t.Fatal("NewHashMap returned nil\n")
	}
	size := original.Size()
	if size != len(data) {
		t.Fatalf("Incorrect size; expected %d, got %d\n", len(data), size)
	}
	original.ForEach(func(k Element, v Element) {
		if v.(IntElement) == 1 {
			t.Fatalf("At %s, already visited\n", k)
		}
		data[k] = IntElement(1)
	})

	for k, v := range data {
		if v.(IntElement) != 1 {
			t.Fatalf("At %s, not visited\n", k)
		}
	}
}

func Test_Hashmap_CustomKey_BadHash_Get(t *testing.T) {
	max := 100
	contents := make(map[Element]Element, max)
	for i := 0; i < max; i++ {
		contents[MyBadKey(i)] = IntElement(i)
	}

	original := NewHashMap(contents, WithMyBadKeyMetadata)

	for k, v := range contents {
		result, _, _ := original.Get(k)
		if result != v {
			t.Fatalf("At %s, expected %s, got %s\n%#v", k, v, result, original)
		}
	}
}
