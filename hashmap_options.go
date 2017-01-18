package immutable

import "github.com/object88/immutable/memory"

type HashMapOption func(*HashMapOptions)

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
}

func defaultHashMapOptions() *HashMapOptions {
	return &HashMapOptions{
		BucketStrategy: memory.LargeBlock,
	}
}

// func (o *HashMapOptions) cloneHashMapOptions() *HashMapOptions {
// 	if o == nil {
// 		return NewHashMapOptions()
// 	}
//
// 	return &HashMapOptions{
// 		BucketStrategy: o.BucketStrategy,
// 	}
// }
