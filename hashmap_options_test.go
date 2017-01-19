package immutable

import (
	"testing"

	"github.com/object88/immutable/memory"
)

func Test_HashMapOptions_Create(t *testing.T) {
	options := defaultHashMapOptions()
	if options == nil {
		t.Fatalf("Create method returned nil")
	}
}

func Test_HashMapOptions_WithBucketStrategy(t *testing.T) {
	options := defaultHashMapOptions()
	f := WithBucketStrategy(memory.ExtraLargeBlock)
	if nil == f {
		t.Fatalf("Fn returned from WithBucketStrategy is nil")
	}
	f(options)
	if options.BucketStrategy != memory.ExtraLargeBlock {
		t.Fatalf("Fn set wrong bucket strategy; got %s, expected %s", options.BucketStrategy, memory.ExtraLargeBlock)
	}
}

// func Test_HashMapOptions_Clone_FromNil(t *testing.T) {
// 	defaultOptions := NewHashMapOptions()
// 	var options *HashMapOptions
// 	clone := options.cloneHashMapOptions()
// 	if clone == nil {
// 		t.Fatalf("Clone method with nil returned nil")
// 	}
// 	if clone.BucketStrategy != defaultOptions.BucketStrategy {
// 		t.Fatalf("Incorrect default for BucketStrategy; expected %s, got %s\n", defaultOptions.BucketStrategy, clone.BucketStrategy)
// 	}
// }
//
// func Test_HashMapOptions_Clone_FromExisting(t *testing.T) {
// 	defaultOptions := NewHashMapOptions()
// 	clone := defaultOptions.cloneHashMapOptions()
// 	if clone == nil {
// 		t.Fatalf("Clone method with nil returned nil")
// 	}
// 	if clone.BucketStrategy != defaultOptions.BucketStrategy {
// 		t.Fatalf("Incorrect default for BucketStrategy; expected %s, got %s\n", defaultOptions.BucketStrategy, clone.BucketStrategy)
// 	}
// }
//
// func Test_HashmapOptions_AttemptChange(t *testing.T) {
// 	// options := NewHashMapOptions()
// 	// options.BucketStrategy = memory.SmallBlock
// 	fn := func(o *HashMapOptions) {
// 		o.BucketStrategy = memory.SmallBlock
// 	}
// 	original := NewHashMap(map[Key]Value{}, fn)
// 	if original.options.BucketStrategy != memory.SmallBlock {
// 		t.Fatalf("Passed in options were not honored for BucketStrategy; expected %s, got %s\n", memory.SmallBlock, original.options.BucketStrategy)
// 	}
//
// 	options.BucketStrategy = memory.LargeBlock
// 	if original.options.BucketStrategy != memory.SmallBlock {
// 		t.Fatalf("Changing bucket strategy on options altered hashmap; expected %s, got %s\n", memory.SmallBlock, original.options.BucketStrategy)
// 	}
// }

func Test_HashmapOptions_SharedInternally(t *testing.T) {
	original := NewHashMap(map[Key]Value{IntKey(1): "a"})
	modified, _ := original.Insert(IntKey(2), "b")

	if modified.options != original.options {
		t.Fatalf("Created new options object from mutating function")
	}
}
