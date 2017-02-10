package core

import "github.com/object88/immutable/memory"

// HashMapOption can be implemented to set a hash map option using the
// NewHashMap function
type HashMapOption func(*HashMapOptions)

// WithBucketStrategy selects a bucket strategy for the hash map
func WithBucketStrategy(blocksize memory.BlockSize) HashMapOption {
	return func(o *HashMapOptions) {
		o.BucketStrategy = blocksize
	}
}

// HashMapOptions contains the options which select the strategies used by
// a hash map for memory allocation.  Do not instantiate this directly; use
// the NewHashMapOptions function instead.
type HashMapOptions struct {
	BucketStrategy memory.BlockSize
	KeyConfig      HandlerConfig
	ValueConfig    HandlerConfig
}

func DefaultHashMapOptions() *HashMapOptions {
	return &HashMapOptions{
		BucketStrategy: memory.LargeBlock,
	}
}
