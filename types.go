package immutable

import "github.com/object88/immutable/core"

// FilterPredicate describes the predicate function used by the Filter method
type FilterPredicate func(key core.Element, value core.Element) (bool, error)

// ForEachPredicate describes the predicate function used by the ForEach method
type ForEachPredicate func(key core.Element, value core.Element)

// MapPredicate describes the predicate function used by the Map method
type MapPredicate func(key core.Element, value core.Element) (core.Element, error)

// ReducePredicate describes the predicate function used by the Reduce method
type ReducePredicate func(accumulator core.Element, key core.Element, value core.Element) (core.Element, error)
