package immutable

// Base describes the low-level set of functions
type Base interface {
	Iterate(abort <-chan struct{}) <-chan keyValuePair
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
	abort := make(chan struct{})
	ch := b.Iterate(abort)
	for kvp := range ch {
		keep, err := predicate(kvp.key, kvp.value)
		if err != nil {
			close(abort)
			return nil, err
		}
		if keep {
			mutated.internalSet(kvp.key, kvp.value)
		}
	}

	return mutated, nil
}

// ForEach iterates over all collection contents
func (b *BaseStruct) ForEach(predicate ForEachPredicate) {
	if b == nil {
		return
	}

	abort := make(chan struct{})
	ch := b.Iterate(abort)
	for kvp := range ch {
		predicate(kvp.key, kvp.value)
	}
}

// Map iterates over the contents of a collection and calls the supplied predicate.
// The return value is a new map with the results of the predicate function.
func (b *BaseStruct) Map(predicate MapPredicate) (*BaseStruct, error) {
	if b == nil {
		return nil, nil
	}

	mutated := b.instantiate(b.Size())
	abort := make(chan struct{})
	ch := b.Iterate(abort)
	for kvp := range ch {
		newV, err := predicate(kvp.key, kvp.value)
		if err != nil {
			close(abort)
			return nil, err
		}
		mutated.internalSet(kvp.key, newV)
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
	abort := make(chan struct{})
	ch := b.Iterate(abort)
	for kvp := range ch {
		acc, err = predicate(acc, kvp.key, kvp.value)
		if err != nil {
			close(abort)
			return nil, err
		}
	}
	return acc, nil
}
