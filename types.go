package immutable

// FilterPredicate describes the predicate function used by the Filter method
type FilterPredicate func(key uint32, value Value) (bool, error)

// ForEachPredicate describes the predicate function used by the ForEach method
type ForEachPredicate func(key uint32, value Value)

// Iterator is a function returned from Iterate
type Iterator func() (key uint32, value Value, next Iterator)

// // Key is a key
// type Key interface {
// 	// Hash() int
// }

// MapPredicate describes the predicate function used by the Map method
type MapPredicate func(key uint32, value Value) (Value, error)

// ReducePredicate describes the predicate function used by the Reduce method
type ReducePredicate func(accumulator Value, key uint32, value Value) (Value, error)

// Value is a value
type Value interface{}
