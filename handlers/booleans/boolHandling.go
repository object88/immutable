package booleans

import (
	"fmt"
	"sync"

	"github.com/object88/immutable/core"
)

var config *core.HandlerConfig
var once sync.Once

type BoolElement bool

func (e BoolElement) Hash(seed uint32) (hash uint64) {
	b := bool(e)
	if b {
		return 0x0
	}
	return 0xffffffffffffffff
}

func (e BoolElement) String() string {
	return fmt.Sprintf("%t", bool(e))
}

// WithBoolKeyMetadata establishes the hydrator and dehydrator methods
// for working with integer keys.
func WithBoolKeyMetadata(o *core.HashMapOptions) {
	o.KeyConfig = createOneBoolHandler()
}

func WithBoolValueMetadata(o *core.HashMapOptions) {
	o.ValueConfig = createOneBoolHandler()
}

func createOneBoolHandler() *core.HandlerConfig {
	once.Do(func() {
		config = &core.HandlerConfig{
			Compare: func(a, b core.Element) (match bool) {
				return false
			},
			CreateBucket: func(count int) core.SubBucket {
				m := make([]bool, count)
				return &BoolSubBucket{m}
			},
		}
	})
	return config
}

type BoolSubBucket struct {
	m []bool
}

func (h *BoolSubBucket) Dehydrate(index int, value core.Element) {
	h.m[index] = bool(value.(BoolElement))
}

func (h *BoolSubBucket) Hydrate(index int) core.Element {
	return BoolElement(h.m[index])
}
