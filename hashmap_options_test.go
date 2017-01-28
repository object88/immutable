package immutable

import (
	"testing"

	"github.com/object88/immutable/memory"
)

func Test_HashMapOptions_Defaults(t *testing.T) {
	options := defaultHashMapOptions()
	if options == nil {
		t.Fatalf("Default method returned nil\n")
	}
	if options.BucketStrategy != memory.LargeBlock {
		t.Fatalf("Default options has wrong bucket strategy; got %s, expected %s\n", options.BucketStrategy, memory.LargeBlock)
	}
}

func Test_HashMapOptions_WithBucketStrategy(t *testing.T) {
	options := defaultHashMapOptions()
	f := WithBucketStrategy(memory.ExtraLargeBlock)
	if nil == f {
		t.Fatalf("Fn returned from WithBucketStrategy is nil\n")
	}
	f(options)
	if options.BucketStrategy != memory.ExtraLargeBlock {
		t.Fatalf("Fn set wrong bucket strategy; got %s, expected %s\n", options.BucketStrategy, memory.ExtraLargeBlock)
	}
}

func Test_HashmapOptions_CreateHashmapWithOptions(t *testing.T) {
	original := NewHashMap(map[Element]Element{}, WithBucketStrategy(memory.ExtraLargeBlock))
	if original.options.BucketStrategy != memory.ExtraLargeBlock {
		t.Fatalf("Option function did not set bucket strategy; got %s, expected %s\n", original.options.BucketStrategy, memory.ExtraLargeBlock)
	}
}

func Test_HashmapOptions_SharedInternally(t *testing.T) {
	original := NewHashMap(map[Element]Element{IntElement(1): StringElement("a")}, WithIntegerKeyMetadata, WithStringValueMetadata)
	modified, _ := original.Insert(IntElement(2), StringElement("b"))

	if modified.options != original.options {
		t.Fatalf("Created new options object from mutating function\n")
	}
}
