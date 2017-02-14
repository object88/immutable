package core

import "unsafe"

type HashmapConfig struct {
	KeyConfig   HandlerConfig
	Options     *HashmapOptions
	ValueConfig HandlerConfig
}

type HandlerConfig interface {
	Compare(a, b unsafe.Pointer) (match bool)
	CompareTo(memory unsafe.Pointer, index int, other unsafe.Pointer) (match bool)
	CreateBucket(count int) unsafe.Pointer
	Hash(element unsafe.Pointer, seed uint32) uint64
	Read(memory unsafe.Pointer, index int) (result unsafe.Pointer)
	Format(value unsafe.Pointer) string
	Write(memory unsafe.Pointer, index int, value unsafe.Pointer)
}

type KeyValuePair struct {
	Key   unsafe.Pointer
	Value unsafe.Pointer
}

// FilterPredicate describes the predicate function used by the Filter method
type FilterPredicate func(key unsafe.Pointer, value unsafe.Pointer) (bool, error)

// ForEachPredicate describes the predicate function used by the ForEach method
type ForEachPredicate func(key unsafe.Pointer, value unsafe.Pointer)

// MapPredicate describes the predicate function used by the Map method
type MapPredicate func(key unsafe.Pointer, value unsafe.Pointer) (unsafe.Pointer, error)

// ReducePredicate describes the predicate function used by the Reduce method
type ReducePredicate func(key unsafe.Pointer, value unsafe.Pointer) error
