package immutable

import (
	"fmt"
	"unsafe"
)

// FilterPredicate describes the predicate function used by the Filter method
type FilterPredicate func(key Key, value Value) (bool, error)

// ForEachPredicate describes the predicate function used by the ForEach method
type ForEachPredicate func(key Key, value Value)

// Key is a key
type Key interface {
	Value

	// Hash calculates the 64-bit hash value for a Key
	Hash(seed uint32) uint64
}

// Hydrater converts a raw uin64 to a Value
type Hydrater func(raw unsafe.Pointer) (result Value, err error)

// Dehydrater converts a Value into a raw uint64
type Dehydrater func(value Value) (result unsafe.Pointer, err error)

// MapPredicate describes the predicate function used by the Map method
type MapPredicate func(key Key, value Value) (Value, error)

// ReducePredicate describes the predicate function used by the Reduce method
type ReducePredicate func(accumulator Value, key Key, value Value) (Value, error)

// Value is a value
type Value interface {
	fmt.Stringer
}

type keyValuePair struct {
	key   Key
	value Value
}

type BucketGenerator func(count int) SubBucket

type SubBucket interface {
	Hydrate(index int) (value Value)
	Dehydrate(index int, value Value)
}
