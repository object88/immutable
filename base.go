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

func (b *BaseStruct) mapping(predicate MapPredicate) (*BaseStruct, error) {
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

func (b *BaseStruct) reduce(predicate ReducePredicate, accumulator Value) (Value, error) {
	acc := accumulator
	var err error
	abort := make(chan struct{})
	for kvp := range b.Iterate(abort) {
		acc, err = predicate(acc, kvp.key, kvp.value)
		if err != nil {
			close(abort)
			return nil, err
		}
	}
	return acc, nil
}
