package core_test

// import (
// 	"testing"
//
// 	"github.com/object88/immutable/core"
// 	"github.com/object88/immutable/memory"
// )
//
// func Test_HashMapOptions_Defaults(t *testing.T) {
// 	options := core.DefaultHashMapOptions()
// 	if options == nil {
// 		t.Fatalf("Default method returned nil\n")
// 	}
// 	if options.BucketStrategy != memory.LargeBlock {
// 		t.Fatalf("Default options has wrong bucket strategy; got %s, expected %s\n", options.BucketStrategy, memory.LargeBlock)
// 	}
// }
//
// func Test_HashMapOptions_WithBucketStrategy(t *testing.T) {
// 	options := core.DefaultHashMapOptions()
// 	f := core.WithBucketStrategy(memory.ExtraLargeBlock)
// 	if nil == f {
// 		t.Fatalf("Fn returned from WithBucketStrategy is nil\n")
// 	}
// 	f(options)
// 	if options.BucketStrategy != memory.ExtraLargeBlock {
// 		t.Fatalf("Fn set wrong bucket strategy; got %s, expected %s\n", options.BucketStrategy, memory.ExtraLargeBlock)
// 	}
// }
//
// func Test_HashmapOptions_CreateHashmapWithOptions(t *testing.T) {
// 	options := core.DefaultHashMapOptions()
// 	fn := core.WithBucketStrategy(memory.ExtraLargeBlock)
// 	fn(options)
// 	if options.BucketStrategy != memory.ExtraLargeBlock {
// 		t.Fatalf("Option function did not set bucket strategy; got %s, expected %s\n", options.BucketStrategy, memory.ExtraLargeBlock)
// 	}
// }
