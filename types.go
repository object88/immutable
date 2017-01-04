package immutable

import "fmt"

// FilterPredicate describes the predicate function used by the Filter method
type FilterPredicate func(key Key, value Value) (bool, error)

// ForEachPredicate describes the predicate function used by the ForEach method
type ForEachPredicate func(key Key, value Value)

// Iterator is a function returned from Iterate
type Iterator func() (key Key, value Value, next Iterator)

// Key is a key
type Key interface {
	fmt.Stringer

	// Hash calculates the 32-bit hash value for a Key
	Hash() uint32
}

// MapPredicate describes the predicate function used by the Map method
type MapPredicate func(key Key, value Value) (Value, error)

// ReducePredicate describes the predicate function used by the Reduce method
type ReducePredicate func(accumulator Value, key Key, value Value) (Value, error)

// Value is a value
type Value interface{}

type keyValuePair struct {
	key   Key
	value Value
}
