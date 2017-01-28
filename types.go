package immutable

import "fmt"

// FilterPredicate describes the predicate function used by the Filter method
type FilterPredicate func(key Element, value Element) (bool, error)

// ForEachPredicate describes the predicate function used by the ForEach method
type ForEachPredicate func(key Element, value Element)

// Element may be a key or a value
type Element interface {
	fmt.Stringer
	Hash(seed uint32) uint64
}

// MapPredicate describes the predicate function used by the Map method
type MapPredicate func(key Element, value Element) (Element, error)

// ReducePredicate describes the predicate function used by the Reduce method
type ReducePredicate func(accumulator Element, key Element, value Element) (Element, error)

type keyValuePair struct {
	key   Element
	value Element
}

type BucketGenerator func(count int) SubBucket

type SubBucket interface {
	Hydrate(index int) (e Element)
	Dehydrate(index int, e Element)
}
