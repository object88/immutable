package immutable

import "unsafe"

// FilterPredicate describes the predicate function used by the Filter method
type FilterPredicate func(key unsafe.Pointer, value unsafe.Pointer) (bool, error)

// ForEachPredicate describes the predicate function used by the ForEach method
type ForEachPredicate func(key unsafe.Pointer, value unsafe.Pointer)

// MapPredicate describes the predicate function used by the Map method
type MapPredicate func(key unsafe.Pointer, value unsafe.Pointer) (unsafe.Pointer, error)

// ReducePredicate describes the predicate function used by the Reduce method
type ReducePredicate func(accumulator unsafe.Pointer, key unsafe.Pointer, value unsafe.Pointer) (unsafe.Pointer, error)
