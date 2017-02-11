package strings

import (
	"sync"
	"unsafe"

	"github.com/object88/immutable/core"
	"github.com/object88/immutable/hasher"
)

var config *core.HandlerConfig
var once sync.Once

// WithStringKeyMetadata establishes the hydrator and dehydrator methods
// for working with integer keys.
func WithStringKeyMetadata(o *core.HashMapOptions) {
	var hc StringHandlerConfig
	o.KeyConfig = hc
}

func WithStringValueMetadata(o *core.HashMapOptions) {
	var hc StringHandlerConfig
	o.ValueConfig = hc
}

type StringHandlerConfig struct{}

func (StringHandlerConfig) Compare(a, b unsafe.Pointer) (match bool) {
	return *(*string)(a) == *(*string)(b)
}

func (shc StringHandlerConfig) CompareTo(memory unsafe.Pointer, index int, other unsafe.Pointer) (match bool) {
	// return (*(*[]string)(memory))[index] == *(*string)(other)
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
	// s := (*(*[]string)(memory))[index]
	// return unsafe.Pointer(&s)
	m := *(*[]unsafe.Pointer)(memory)
	return m[index]
}

func (StringHandlerConfig) Format(memory unsafe.Pointer) string {
	return *(*string)(memory)
}

func (StringHandlerConfig) Write(memory unsafe.Pointer, index int, value unsafe.Pointer) {
	// (*(*[]string)(memory))[index] = *(*string)(value)
	m := *(*[]unsafe.Pointer)(memory)
	m[index] = value
}

// func (StringHandlerConfig) CompareTo(memory unsafe.Pointer, index int, other unsafe.Pointer) (match bool) {
// 	return (*(*[]string)(memory))[index] == *(*string)(other)
// }
//
// func (StringHandlerConfig) CreateBucket(count int) unsafe.Pointer {
// 	m := make([]string, count)
// 	return unsafe.Pointer(&m)
// }
//
// func (StringHandlerConfig) Hash(element unsafe.Pointer, seed uint32) uint64 {
// 	s := *(*string)(element)
// 	return hasher.HashString(s, seed)
// }
//
// func (StringHandlerConfig) Read(memory unsafe.Pointer, index int) (result unsafe.Pointer) {
// 	s := (*(*[]string)(memory))[index]
// 	return unsafe.Pointer(&s)
// }
//
// func (StringHandlerConfig) Format(memory unsafe.Pointer) string {
// 	return *(*string)(memory)
// }
//
// func (StringHandlerConfig) Write(memory unsafe.Pointer, index int, value unsafe.Pointer) {
// 	(*(*[]string)(memory))[index] = *(*string)(value)
// }
