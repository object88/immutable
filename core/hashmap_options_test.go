package core

import (
	"testing"

	"github.com/object88/immutable/handlers/integers"
	"github.com/object88/immutable/handlers/strings"
	"github.com/object88/immutable/memory"
)

func Test_HashMapOptions_Defaults(t *testing.T) {
	options := DefaultHashMapOptions()
	if options == nil {
		t.Fatalf("Default method returned nil\n")
	}
	if options.BucketStrategy != memory.LargeBlock {
		t.Fatalf("Default options has wrong bucket strategy; got %s, expected %s\n", options.BucketStrategy, memory.LargeBlock)
	}
}

func Test_HashMapOptions_WithBucketStrategy(t *testing.T) {
	options := DefaultHashMapOptions()
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
	original := NewHashMap(map[Element]Element{integers.IntElement(1): strings.StringElement("a")}, integers.WithIntKeyMetadata, strings.WithStringValueMetadata)
	modified, _ := original.Insert(integers.IntElement(2), strings.StringElement("b"))

	if modified.options != original.options {
		t.Fatalf("Created new options object from mutating function\n")
	}
}
