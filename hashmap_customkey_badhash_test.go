package immutable

// import (
// 	"fmt"
// 	"testing"
//
// 	"github.com/object88/immutable/core"
// 	"github.com/object88/immutable/handlers/integers"
// )
//
// // This suite of tests is designed to test the bucket overflow behavior.
// // By using a `Hash` function which only ever resolves into one of two
// // hash keys, all entries will end up in one or two buckets.  With a
// // sufficiently large dataset, this will cause bucket overflow to be used.
//
// type MyBadElement int
//
// func (k MyBadElement) Hash(seed uint32) uint64 {
// 	if k%2 == 0 {
// 		return 0x0
// 	}
// 	return 0xffffffffffffffff
// }
//
// func (k MyBadElement) String() string {
// 	return fmt.Sprintf("%d", k)
// }
//
// func WithMyBadElementMetadata(o *core.HashMapOptions) {
// 	o.KeyConfig = &core.HandlerConfig{
// 		Compare: func(a, b core.Element) (match bool) {
// 			return int(a.(MyBadElement)) == int(b.(MyBadElement))
// 		},
// 		CreateBucket: func(count int) core.SubBucket {
// 			m := make([]int, count)
// 			return &MyBadElementSubBucket{m}
// 		},
// 	}
// 	integers.WithIntValueMetadata(o)
// }
//
// type MyBadElementSubBucket struct {
// 	m []int
// }
//
// func NewMyBadElementHashHandler() core.BucketGenerator {
// 	return NewMyBadElementSubBucket
// }
//
// func NewMyBadElementSubBucket(count int) core.SubBucket {
// 	m := make([]int, count)
// 	return &MyBadElementSubBucket{m}
// }
//
// func (b *MyBadElementSubBucket) Hydrate(index int) core.Element {
// 	u := b.m[index]
// 	v := MyBadElement(u)
// 	return v
// }
//
// func (b *MyBadElementSubBucket) Dehydrate(index int, v core.Element) {
// 	u := v.(MyBadElement)
// 	b.m[index] = int(u)
// }
//
// func Test_Hashmap_CustomKey_BadHash_Iterate(t *testing.T) {
// 	max := 100
// 	data := make(map[core.Element]core.Element, max)
// 	for i := 0; i < max; i++ {
// 		data[MyBadElement(i)] = integers.IntElement(0)
// 	}
// 	original := NewHashMap(data, WithMyBadElementMetadata)
// 	if original == nil {
// 		t.Fatal("NewHashMap returned nil\n")
// 	}
// 	size := original.Size()
// 	if size != len(data) {
// 		t.Fatalf("Incorrect size; expected %d, got %d\n", len(data), size)
// 	}
// 	original.ForEach(func(k core.Element, v core.Element) {
// 		if v.(integers.IntElement) == 1 {
// 			t.Fatalf("At %s, already visited\n", k)
// 		}
// 		data[k] = integers.IntElement(1)
// 	})
//
// 	for k, v := range data {
// 		if v.(integers.IntElement) != 1 {
// 			t.Fatalf("At %s, not visited\n", k)
// 		}
// 	}
// }
//
// func Test_Hashmap_CustomKey_BadHash_Get(t *testing.T) {
// 	max := 100
// 	contents := make(map[core.Element]core.Element, max)
// 	for i := 0; i < max; i++ {
// 		contents[MyBadElement(i)] = integers.IntElement(i)
// 	}
//
// 	original := NewHashMap(contents, WithMyBadElementMetadata)
//
// 	for k, v := range contents {
// 		result, _, _ := original.Get(k)
// 		if result != v {
// 			t.Fatalf("At %s, expected %s, got %s\n%#v", k, v, result, original)
// 		}
// 	}
// }
