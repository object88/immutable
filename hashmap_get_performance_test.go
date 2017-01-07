package immutable

import (
	"math/rand"
	"testing"
	"time"

	"github.com/object88/immutable/memory"
)

const (
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIndexBits = 6                      // 6 bits to represent a letter index
	letterIndexMask = 1<<letterIndexBits - 1 // All 1-bits, as many as letterIndexBits
	letterIndexMax  = 63 / letterIndexBits   // # of letter indices fitting in 63 bits
	max             = 500000
)

var keys []IntKey
var contents map[Key]Value
var result string
var src = rand.NewSource(time.Now().UnixNano())

func init() {
	stringLength := 100
	contents = make(map[Key]Value, max)
	keys = make([]IntKey, max)
	for i := 0; i < max; i++ {
		keys[i] = IntKey(i)
		contents[keys[i]] = generateString(stringLength)
	}
}

func compareBucketStrategy(original *HashMap) {
	var r string
	for key := range contents {
		r = original.Get(key).(string)
	}
	result = r
}

func Benchmark_LargeBlock(b *testing.B) {
	options := NewHashMapOptions()
	options.BucketStrategy = memory.LargeBlock
	original := NewHashMap(contents, options)
	for i := 0; i < b.N; i++ {
		compareBucketStrategy(original)
	}
}

func Benchmark_ExtraLargeBlock(b *testing.B) {
	options := NewHashMapOptions()
	options.BucketStrategy = memory.ExtraLargeBlock
	original := NewHashMap(contents, options)
	for i := 0; i < b.N; i++ {
		compareBucketStrategy(original)
	}
}

func Benchmark_NoPackingBlock(b *testing.B) {
	options := NewHashMapOptions()
	options.BucketStrategy = memory.NoPacking
	original := NewHashMap(contents, options)
	for i := 0; i < b.N; i++ {
		compareBucketStrategy(original)
	}
}

func Benchmark_NativeMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r string
		for _, v := range keys {
			r = contents[v].(string)
		}
		result = r
	}
}

// This code copied directly from StackOverflow; see randStringBytesMaskImprSrc
// function:
// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func generateString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIndexMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIndexMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIndexMax
		}
		if idx := int(cache & letterIndexMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIndexBits
		remain--
	}

	return string(b)
}
