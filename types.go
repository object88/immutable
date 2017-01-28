package immutable

import (
	"fmt"
	"unsafe"
)

// FilterPredicate describes the predicate function used by the Filter method
type FilterPredicate func(key Element, value Element) (bool, error)

// ForEachPredicate describes the predicate function used by the ForEach method
type ForEachPredicate func(key Element, value Element)

type Element interface {
	fmt.Stringer
	Hash(seed uint32) uint64
}

// // Element is a key
// type Element interface {
// 	Element
//
// 	// Hash calculates the 64-bit hash value for a Element
// 	Hash(seed uint32) uint64
// }

// Hydrater converts a raw uin64 to a Element
type Hydrater func(raw unsafe.Pointer) (result Element, err error)

// Dehydrater converts a Element into a raw uint64
type Dehydrater func(value Element) (result unsafe.Pointer, err error)

// MapPredicate describes the predicate function used by the Map method
type MapPredicate func(key Element, value Element) (Element, error)

// ReducePredicate describes the predicate function used by the Reduce method
type ReducePredicate func(accumulator Element, key Element, value Element) (Element, error)

// // Element is a value
// type Element interface {
// 	fmt.Stringer
// }

type keyValuePair struct {
	key   Element
	value Element
}

type BucketGenerator func(count int) SubBucket

type SubBucket interface {
	Hydrate(index int) (e Element)
	Dehydrate(index int, e Element)
}
