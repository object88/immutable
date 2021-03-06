package strings

import (
	"sync"
	"unsafe"

	"github.com/object88/immutable/core"
	"github.com/object88/immutable/hasher"
)

var config core.HandlerConfig
var once sync.Once

func GetHandler() core.HandlerConfig {
	once.Do(func() {
		config = &StringHandlerConfig{}
	})
	return config
}

type StringHandlerConfig struct{}

func (StringHandlerConfig) Compare(a, b unsafe.Pointer) (match bool) {
	return *(*string)(a) == *(*string)(b)
}

func (shc StringHandlerConfig) CompareTo(memory unsafe.Pointer, index int, other unsafe.Pointer) (match bool) {
	u := shc.Read(memory, index)
	return *(*string)(u) == *(*string)(other)
}

func (StringHandlerConfig) CreateBucket(count int) unsafe.Pointer {
	m := make([]unsafe.Pointer, count)
	return unsafe.Pointer(&m)
}

func (StringHandlerConfig) Hash(element unsafe.Pointer, seed uint32) uint64 {
	s := *(*string)(element)
	return hasher.HashString(s, seed)
}

func (StringHandlerConfig) Read(memory unsafe.Pointer, index int) (result unsafe.Pointer) {
	m := *(*[]unsafe.Pointer)(memory)
	return m[index]
}

func (StringHandlerConfig) Format(memory unsafe.Pointer) string {
	return *(*string)(memory)
}

func (StringHandlerConfig) Write(memory unsafe.Pointer, index int, value unsafe.Pointer) {
	m := *(*[]unsafe.Pointer)(memory)
	m[index] = value
}
