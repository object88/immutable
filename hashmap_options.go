package immutable

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

// WithIntegerKeyMetadata establishes the hydrator and dehydrator methods
// for working with integer keys.
func WithIntegerKeyMetadata(o *HashMapOptions) {
	o.KeyHandler = NewIntHandler()
}

func WithIntegerValueMetadata(o *HashMapOptions) {
	o.ValueHandler = NewIntHandler()
}

func WithStringValueMetadata(o *HashMapOptions) {
	o.ValueHandler = NewStringHandler()
}

// HashMapOptions contains the options which select the strategies used by
// a hash map for memory allocation.  Do not instantiate this directly; use
// the NewHashMapOptions function instead.
type HashMapOptions struct {
	BucketStrategy memory.BlockSize
	KeyHandler     BucketGenerator
	ValueHandler   BucketGenerator
}

func defaultHashMapOptions() *HashMapOptions {
	return &HashMapOptions{
		BucketStrategy: memory.LargeBlock,
	}
}
