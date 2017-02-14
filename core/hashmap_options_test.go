package core_test

import (
	"testing"

	"github.com/object88/immutable/core"
)

func Test_HashmapOptions_Defaults(t *testing.T) {
	options := core.DefaultHashmapOptions()
	if options == nil {
		t.Fatalf("Default method returned nil\n")
	}
	if options.PackedBucket != false {
		t.Fatalf("Default options has wrong bucket strategy; got %s, expected %b\n", options.PackedBucket, false)
	}
}

func Test_HashmapOptions_WithBucketStrategy(t *testing.T) {
	options := core.DefaultHashmapOptions()
	f := core.WithBucketStrategy(true)
	if nil == f {
		t.Fatalf("Fn returned from WithBucketStrategy is nil\n")
	}
	f(options)
	if options.PackedBucket != true {
		t.Fatalf("Fn set wrong bucket strategy; got %s, expected %b\n", options.PackedBucket, true)
	}
}

func Test_HashmapOptions_CreateHashmapWithOptions(t *testing.T) {
	options := core.DefaultHashmapOptions()
	fn := core.WithBucketStrategy(true)
	fn(options)
	if options.PackedBucket != true {
		t.Fatalf("Option function did not set bucket strategy; got %s, expected %b\n", options.PackedBucket, true)
	}
}
