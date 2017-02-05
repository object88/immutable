package immutable

import "github.com/object88/immutable/core"

// Base describes the low-level set of functions
type Base interface {
	// Get(key Element) (result Element, ok bool)
	Size() int
	instantiate(initialSize int, contents []*core.KeyValuePair) *BaseStruct
	internalSet(key core.Element, value core.Element)
	iterate(abort <-chan struct{}) <-chan core.KeyValuePair
}

// BaseStruct is what it say on the tin
type BaseStruct struct {
	Base
}

// Filter returns a subset of the collection, based on the predicate supplied
func (b *BaseStruct) filter(predicate FilterPredicate) (*BaseStruct, error) {
	resultSet := make([]*core.KeyValuePair, b.Size())
	resultSetCount := 0
	abort := make(chan struct{})
	ch := b.iterate(abort)
	for kvp := range ch {
		keep, err := predicate(kvp.Key, kvp.Value)
		if err != nil {
			close(abort)
			return nil, err
		}
		if keep {
			resultSet[resultSetCount] = &core.KeyValuePair{Key: kvp.Key, Value: kvp.Value}
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
		predicate(kvp.Key, kvp.Value)
	}
}

func (b *BaseStruct) mapping(predicate MapPredicate) (*BaseStruct, error) {
	mutated := b.instantiate(b.Size(), nil)
	abort := make(chan struct{})
	ch := b.iterate(abort)
	for kvp := range ch {
		newV, err := predicate(kvp.Key, kvp.Value)
		if err != nil {
			close(abort)
			return nil, err
		}
		mutated.internalSet(kvp.Key, newV)
	}
	return mutated, nil
}

func (b *BaseStruct) reduce(predicate ReducePredicate, accumulator core.Element) (core.Element, error) {
	acc := accumulator
	var err error
	abort := make(chan struct{})
	for kvp := range b.iterate(abort) {
		acc, err = predicate(acc, kvp.Key, kvp.Value)
		if err != nil {
			close(abort)
			return nil, err
		}
	}
	return acc, nil
}
