package immutable

import (
	"fmt"
	"strconv"
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

func Test_Hashmap_CustomKey_BadHash_Iterate(t *testing.T) {
	original, config, data := createHashmapAndData()

	size := original.Size()
	if size != len(data) {
		t.Fatalf("Incorrect size; expected %d, got %d\n", len(data), size)
	}
	original.ForEach(config, func(kp, _ unsafe.Pointer) {
		k := *(*int)(kp)
		if data[k] == "visited" {
			t.Fatalf("At %d, already visited\n", k)
		}
		data[k] = "visited"
	})

	for k, v := range data {
		if v != "visited" {
			t.Fatalf("At %d, not visited\n", k)
		}
	}
}

func Test_Hashmap_CustomKey_BadHash_Get(t *testing.T) {
	original, config, data := createHashmapAndData()

	for k, v := range data {
		rp, _, _ := original.Get(config, unsafe.Pointer(&k))
		r := *(*string)(rp)
		if r != v {
			t.Fatalf("At %d, expected %s, got %s\n", k, v, r)
		}
	}
}

func createHashmapAndData() (*core.InternalHashmap, *core.HashmapConfig, map[int]string) {
	max := 100
	data := make(map[int]string, max)
	for i := 0; i < max; i++ {
		data[i] = strconv.Itoa(i)
	}

	c := &core.HashmapConfig{
		KeyConfig:   MyBadHandler{},
		Options:     core.DefaultHashmapOptions(),
		ValueConfig: strings.GetHandler(),
	}

	original := core.CreateEmptyInternalHashmap(false, max)

	for k, v := range data {
		key, value := k, v
		kp, vp := unsafe.Pointer(&key), unsafe.Pointer(&value)
		original.InternalSet(c, kp, vp)
	}

	return original, c, data
}
