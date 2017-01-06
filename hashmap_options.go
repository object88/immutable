package immutable

import "github.com/object88/immutable/memory"

// HashMapOptions contains the options which select the strategies used by
// a hash map for memory allocation.  Do not instantiate this directly; use
// the NewHashMapOptions function instead.
type HashMapOptions struct {
	BucketStrategy memory.BlockSize
}

// NewHashMapOptions creates a new options object with defaults set.
func NewHashMapOptions() *HashMapOptions {
	return &HashMapOptions{
		BucketStrategy: memory.LargeBlock,
	}
}

func (o *HashMapOptions) cloneHashMapOptions() *HashMapOptions {
	if o == nil {
		return NewHashMapOptions()
	}

	return &HashMapOptions{
		BucketStrategy: o.BucketStrategy,
	}
}
