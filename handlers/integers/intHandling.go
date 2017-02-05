package integers

import (
	"fmt"
	"sync"

	"github.com/object88/immutable/core"
	"github.com/object88/immutable/hasher"
)

var config *core.HandlerConfig
var once sync.Once

type IntElement int

func (e IntElement) Hash(seed uint32) (hash uint64) {
	i := int(e)
	p := uintptr(i) // unsafe.Pointer(&i)
	return hasher.Hash8(p, seed)
}

func (e IntElement) String() string {
	return fmt.Sprintf("%d", int(e))
}

// WithIntKeyMetadata establishes the hydrator and dehydrator methods
// for working with integer keys.
func WithIntKeyMetadata(o *core.HashMapOptions) {
	o.KeyConfig = createOneIntHandler()
}

func WithIntValueMetadata(o *core.HashMapOptions) {
	o.ValueConfig = createOneIntHandler()
}

func createOneIntHandler() *core.HandlerConfig {
	once.Do(func() {
		config = &core.HandlerConfig{
			Compare: func(a, b core.Element) (match bool) {
				return int(a.(IntElement)) == int(b.(IntElement))
			},
			CreateBucket: func(count int) core.SubBucket {
				m := make([]int, count)
				return &IntSubBucket{m}
			},
		}
	})
	return config
}

type IntSubBucket struct {
	m []int
}

func (h *IntSubBucket) Dehydrate(index int, value core.Element) {
	h.m[index] = int(value.(IntElement))
}

func (h *IntSubBucket) Hydrate(index int) core.Element {
	return IntElement(h.m[index])
}
