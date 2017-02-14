package core

// HashmapOption can be implemented to set a hash map option using the
// NewInternalHashmap function
type HashmapOption func(*HashmapOptions)

// WithBucketStrategy selects a bucket strategy for the hash map
func WithBucketStrategy(packBuckets bool) HashmapOption {
	return func(o *HashmapOptions) {
		o.PackedBucket = packBuckets
	}
}

// HashmapOptions contains the options which select the strategies used by
// a hash map for memory allocation.  Do not instantiate this directly; use
// the NewHashmapOptions function instead.
type HashmapOptions struct {
	PackedBucket bool
}

func DefaultHashmapOptions() *HashmapOptions {
	return &HashmapOptions{
		PackedBucket: false,
	}
}
