package immutable

import "testing"

func Test_HashMapOptions_Create(t *testing.T) {
	options := NewHashMapOptions()
	if options == nil {
		t.Fatalf("Create method returned nil")
	}
}

func Test_HashMapOptions_Clone_FromNil(t *testing.T) {
	defaultOptions := NewHashMapOptions()
	var options *HashMapOptions
	clone := options.cloneHashMapOptions()
	if clone == nil {
		t.Fatalf("Clone method with nil returned nil")
	}
	if clone.BucketStrategy != defaultOptions.BucketStrategy {
		t.Fatalf("Incorrect default for BucketStrategy; expected %s, got %s\n", defaultOptions.BucketStrategy, clone.BucketStrategy)
	}
}

func Test_HashMapOptions_Clone_FromExisting(t *testing.T) {
	defaultOptions := NewHashMapOptions()
	clone := defaultOptions.cloneHashMapOptions()
	if clone == nil {
		t.Fatalf("Clone method with nil returned nil")
	}
	if clone.BucketStrategy != defaultOptions.BucketStrategy {
		t.Fatalf("Incorrect default for BucketStrategy; expected %s, got %s\n", defaultOptions.BucketStrategy, clone.BucketStrategy)
	}
}