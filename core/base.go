package core

import "unsafe"

// Base describes the low-level set of functions
type Base interface {
	InternalSet(config *HashmapConfig, key unsafe.Pointer, value unsafe.Pointer)
	Size() int

	instantiate(config *HashmapConfig, initialSize int, contents []*KeyValuePair) *BaseStruct
	iterate(config *HashmapConfig, abort <-chan struct{}) <-chan KeyValuePair
}

// BaseStruct is what it say on the tin
type BaseStruct struct {
	Base
}

// Filter returns a subset of the collection, based on the predicate supplied
func (b *BaseStruct) filter(config *HashmapConfig, predicate FilterPredicate) (*BaseStruct, error) {
	resultSet := make([]*KeyValuePair, b.Size())
	resultSetCount := 0
	abort := make(chan struct{})
	ch := b.iterate(config, abort)
	for kvp := range ch {
		keep, err := predicate(kvp.Key, kvp.Value)
		if err != nil {
			close(abort)
			return nil, err
		}
		if keep {
			resultSet[resultSetCount] = &KeyValuePair{Key: kvp.Key, Value: kvp.Value}
			resultSetCount++
		}
	}

	mutated := b.instantiate(config, resultSetCount, resultSet)

	return mutated, nil
}

func (b *BaseStruct) forEach(config *HashmapConfig, predicate ForEachPredicate) {
	abort := make(chan struct{})
	ch := b.iterate(config, abort)
	for kvp := range ch {
		predicate(kvp.Key, kvp.Value)
	}
}

func (b *BaseStruct) mapping(config *HashmapConfig, predicate MapPredicate) (*BaseStruct, error) {
	mutated := b.instantiate(config, b.Size(), nil)
	abort := make(chan struct{})
	ch := b.iterate(config, abort)
	for kvp := range ch {
		newV, err := predicate(kvp.Key, kvp.Value)
		if err != nil {
			close(abort)
			return nil, err
		}
		mutated.InternalSet(config, kvp.Key, newV)
	}
	return mutated, nil
}

func (b *BaseStruct) reduce(config *HashmapConfig, predicate ReducePredicate, accumulator unsafe.Pointer) (unsafe.Pointer, error) {
	acc := accumulator
	var err error
	abort := make(chan struct{})
	for kvp := range b.iterate(config, abort) {
		acc, err = predicate(acc, kvp.Key, kvp.Value)
		if err != nil {
			close(abort)
			return nil, err
		}
	}
	return acc, nil
}
