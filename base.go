package immutable

// Base describes the low-level set of functions
type Base interface {
	Get(key Key) Value
	Size() int
	Instantiate(initialSize int, contents []*KeyValuePair) *BaseStruct
	InternalSet(key Key, value Value)
	Iterate(abort <-chan struct{}) <-chan KeyValuePair
}

// BaseStruct is what it say on the tin
type BaseStruct struct {
	Base
}

// Filter returns a subset of the collection, based on the predicate supplied
func (b *BaseStruct) Filter(predicate FilterPredicate) (*BaseStruct, error) {
	resultSet := make([]*KeyValuePair, b.Size())
	resultSetCount := 0
	abort := make(chan struct{})
	ch := b.Iterate(abort)
	for kvp := range ch {
		keep, err := predicate(kvp.Key, kvp.Value)
		if err != nil {
			close(abort)
			return nil, err
		}
		if keep {
			resultSet[resultSetCount] = &KeyValuePair{kvp.Key, kvp.Value}
			resultSetCount++
		}
	}

	mutated := b.Instantiate(resultSetCount, resultSet)

	return mutated, nil
}

func (b *BaseStruct) ForEach(predicate ForEachPredicate) {
	abort := make(chan struct{})
	ch := b.Iterate(abort)
	for kvp := range ch {
		predicate(kvp.Key, kvp.Value)
	}
}

func (b *BaseStruct) Mapping(predicate MapPredicate) (*BaseStruct, error) {
	mutated := b.Instantiate(b.Size(), nil)
	abort := make(chan struct{})
	ch := b.Iterate(abort)
	for kvp := range ch {
		newV, err := predicate(kvp.Key, kvp.Value)
		if err != nil {
			close(abort)
			return nil, err
		}
		mutated.InternalSet(kvp.Key, newV)
	}
	return mutated, nil
}

func (b *BaseStruct) Reduce(predicate ReducePredicate, accumulator Value) (Value, error) {
	acc := accumulator
	var err error
	abort := make(chan struct{})
	for kvp := range b.Iterate(abort) {
		acc, err = predicate(acc, kvp.Key, kvp.Value)
		if err != nil {
			close(abort)
			return nil, err
		}
	}
	return acc, nil
}
