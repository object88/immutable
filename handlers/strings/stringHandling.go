package strings

import (
	"sync"

	"github.com/object88/immutable/core"
	"github.com/object88/immutable/hasher"
)

var config *core.HandlerConfig
var once sync.Once

// StringElement is a string.
type StringElement string

func (e StringElement) Hash(seed uint32) (hash uint64) {
	return hasher.HashString(string(e), seed)
}

func (e StringElement) String() string {
	return string(e)
}

// WithStringKeyMetadata establishes the hydrator and dehydrator methods
// for working with integer keys.
func WithStringKeyMetadata(o *core.HashMapOptions) {
	o.KeyConfig = createOneStringHandler()
}

func WithStringValueMetadata(o *core.HashMapOptions) {
	o.ValueConfig = createOneStringHandler()
}

func createOneStringHandler() *core.HandlerConfig {
	once.Do(func() {
		config = &core.HandlerConfig{
			Compare: func(a, b core.Element) (match bool) {
				return string(a.(StringElement)) == string(b.(StringElement))
			},
			CreateBucket: func(count int) core.SubBucket {
				m := make([]string, count)
				return &StringSubBucket{m}
			},
		}
	})
	return config
}

type StringSubBucket struct {
	m []string
}

func (h *StringSubBucket) Dehydrate(index int, value core.Element) {
	s := value.(StringElement)
	h.m[index] = string(s)
}

func (h *StringSubBucket) Hydrate(index int) core.Element {
	return StringElement(h.m[index])
}
