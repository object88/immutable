package immutable

// Base describes the low-level set of functions
type Base interface {
	Iterate() Iterator
	Get(key Key) Value
	Size() uint32
}

type internalFunctions interface {
	instantiate(initialSize uint32) *BaseStruct
	internalSet(key Key, value Value)
}

// BaseStruct is what it say on the tin
type BaseStruct struct {
	Base
	internalFunctions
}

// Filter returns a subset of the collection, based on the predicate supplied
func (b *BaseStruct) Filter(predicate FilterPredicate) (*BaseStruct, error) {
	if b == nil {
		return nil, nil
	}

	mutated := b.instantiate(0)
	for k, v, i := b.Iterate()(); i != nil; k, v, i = i() {
		keep, err := predicate(k, v)
		if err != nil {
			return nil, err
		}
		if keep {
			mutated.internalSet(k, v)
		}
	}

	return mutated, nil
}

// ForEach iterates over all collection contents
func (b *BaseStruct) ForEach(predicate ForEachPredicate) {
	if b == nil {
		return
	}

	for k, v, i := b.Iterate()(); i != nil; k, v, i = i() {
		predicate(k, v)
	}
}

// Map iterates over the contents of a collection and calls the supplied predicate.
// The return value is a new map with the results of the predicate function.
func (b *BaseStruct) Map(predicate MapPredicate) (*BaseStruct, error) {
	if b == nil {
		return nil, nil
	}

	mutated := b.instantiate(b.Size())
	for k, v, i := b.Iterate()(); i != nil; k, v, i = i() {
		newV, err := predicate(k, v)
		if err != nil {
			return nil, err
		}
		mutated.internalSet(k, newV)
	}
	return mutated, nil
}

// Reduce operates over the collection contents to produce a single value
func (b *BaseStruct) Reduce(predicate ReducePredicate, accumulator Value) (Value, error) {
	if b == nil {
		return nil, nil
	}

	acc := accumulator
	var err error
	for k, v, i := b.Iterate()(); i != nil; k, v, i = i() {
		acc, err = predicate(acc, k, v)
		if err != nil {
			return nil, err
		}
	}
	return acc, nil
}
