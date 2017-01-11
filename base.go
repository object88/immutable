package immutable

// Base describes the low-level set of functions
type Base interface {
	Get(key Key) (result Value, ok bool)
	Size() int
	instantiate(initialSize int, contents []*keyValuePair) *BaseStruct
	internalSet(key Key, value Value)
	iterate(abort <-chan struct{}) <-chan keyValuePair
}

// BaseStruct is what it say on the tin
type BaseStruct struct {
	Base
}

// Filter returns a subset of the collection, based on the predicate supplied
func (b *BaseStruct) filter(predicate FilterPredicate) (*BaseStruct, error) {
	resultSet := make([]*keyValuePair, b.Size())
	resultSetCount := 0
	abort := make(chan struct{})
	ch := b.iterate(abort)
	for kvp := range ch {
		keep, err := predicate(kvp.key, kvp.value)
		if err != nil {
			close(abort)
			return nil, err
		}
		if keep {
			resultSet[resultSetCount] = &keyValuePair{kvp.key, kvp.value}
			resultSetCount++
		}
	}

	mutated := b.instantiate(resultSetCount, resultSet)

	return mutated, nil
}

func (b *BaseStruct) forEach(predicate ForEachPredicate) {
	abort := make(chan struct{})
	ch := b.iterate(abort)
	for kvp := range ch {
		predicate(kvp.key, kvp.value)
	}
}

func (b *BaseStruct) mapping(predicate MapPredicate) (*BaseStruct, error) {
	mutated := b.instantiate(b.Size(), nil)
	abort := make(chan struct{})
	ch := b.iterate(abort)
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
	for kvp := range b.iterate(abort) {
		acc, err = predicate(acc, kvp.key, kvp.value)
		if err != nil {
			close(abort)
			return nil, err
		}
	}
	return acc, nil
}
