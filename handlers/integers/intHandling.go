package integers

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/object88/immutable/core"
	"github.com/object88/immutable/hasher"
)

var config core.HandlerConfig
var once sync.Once

func GetHandler() core.HandlerConfig {
	once.Do(func() {
		config = &IntHandlerConfig{}
	})
	return config
}

// WithIntKeyMetadata establishes the hydrator and dehydrator methods
// for working with integer keys.
func WithIntKeyMetadata(o *core.HashMapOptions) {
	var ihc IntHandlerConfig
	o.KeyConfig = ihc
}

func WithIntValueMetadata(o *core.HashMapOptions) {
	var ihc IntHandlerConfig
	o.ValueConfig = ihc
}

type IntHandlerConfig struct{}

func (IntHandlerConfig) Compare(a, b unsafe.Pointer) (match bool) {
	return *(*int)(a) == *(*int)(b)
}

func (IntHandlerConfig) CompareTo(memory unsafe.Pointer, index int, other unsafe.Pointer) (match bool) {
	return (*(*[]int)(memory))[index] == *(*int)(other)
}

func (IntHandlerConfig) CreateBucket(count int) unsafe.Pointer {
	m := make([]int, count)
	return unsafe.Pointer(&m)
}

func (IntHandlerConfig) Hash(e unsafe.Pointer, seed uint32) uint64 {
	i := *(*int)(e)
	p := uint64(i)
	return hasher.Hash8(p, seed)
}

func (IntHandlerConfig) Read(memory unsafe.Pointer, index int) (result unsafe.Pointer) {
	i := (*(*[]int)(memory))[index]
	return unsafe.Pointer(&i)
}

func (IntHandlerConfig) Format(value unsafe.Pointer) string {
	return fmt.Sprintf("%d", *(*int)(value))
}

func (IntHandlerConfig) Write(memory unsafe.Pointer, index int, value unsafe.Pointer) {
	v := *(*int)(value)
	m := *(*[]int)(memory)
	m[index] = v
}
