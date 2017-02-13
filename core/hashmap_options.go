package core

import "github.com/object88/immutable/memory"

// HashmapOption can be implemented to set a hash map option using the
// NewHashMap function
type HashmapOption func(*HashmapOptions)

// WithBucketStrategy selects a bucket strategy for the hash map
func WithBucketStrategy(blocksize memory.BlockSize) HashmapOption {
	return func(o *HashmapOptions) {
		o.BucketStrategy = blocksize
	}
}

// HashmapOptions contains the options which select the strategies used by
// a hash map for memory allocation.  Do not instantiate this directly; use
// the NewHashmapOptions function instead.
type HashmapOptions struct {
	BucketStrategy memory.BlockSize
}

func DefaultHashmapOptions() *HashmapOptions {
	return &HashmapOptions{
		BucketStrategy: memory.LargeBlock,
	}
}
