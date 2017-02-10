package immutable

import (
	"math/rand"
	"testing"
	"time"

	"github.com/object88/immutable"
	"github.com/object88/immutable/memory"
)

const (
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIndexBits = 6                      // 6 bits to represent a letter index
	letterIndexMask = 1<<letterIndexBits - 1 // All 1-bits, as many as letterIndexBits
	letterIndexMax  = 63 / letterIndexBits   // # of letter indices fitting in 63 bits
	max             = 500000
)

var keys []immutable.IntKey
var contents map[immutable.Key]immutable.Value
var result string
var src = rand.NewSource(time.Now().UnixNano())

var hashmapSmallBlock, hashmapLargeBlock, hashmapExtraLargeBlock, hashmapNoPacked *immutable.HashMap

func init() {
	stringLength := 100
	contents = make(map[immutable.Key]immutable.Value, max)
	keys = make([]immutable.IntKey, max)
	for i := 0; i < max; i++ {
		keys[i] = immutable.IntKey(i)
		contents[keys[i]] = generateString(stringLength)
	}
	hashmapSmallBlock = createWithStragety(memory.SmallBlock)
	hashmapLargeBlock = createWithStragety(memory.LargeBlock)
	hashmapExtraLargeBlock = createWithStragety(memory.ExtraLargeBlock)
	hashmapNoPacked = createWithStragety(memory.NoPacking)
}

func Benchmark_Hashmap_Get_SmallBlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testStrategy(hashmapSmallBlock)
	}
}

func Benchmark_Hashmap_Get_LargeBlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testStrategy(hashmapLargeBlock)
	}
}

func Benchmark_Hashmap_Get_ExtraLargeBlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testStrategy(hashmapExtraLargeBlock)
	}
}

func Benchmark_Hashmap_Get_NoPacking(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testStrategy(hashmapNoPacked)
	}
}

func Benchmark_Hashmap_Get_NativeMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var r string
		for _, key := range keys {
			r = getStringFromMap(key)
		}
		result = r
	}
}

func getStringFromMap(key immutable.Key) string {
	r := contents[key].(string)
	return r
}

func createWithStragety(blocksize memory.BlockSize) *immutable.HashMap {
	return immutable.NewHashMap(contents, immutable.WithBucketStrategy(blocksize))
}

func testStrategy(original *immutable.HashMap) {
	var r string
	for _, key := range keys {
		r = getStringFromImmutable(original, key)
	}
	result = r
}

func getStringFromImmutable(original *immutable.HashMap, key immutable.Key) string {
	r, _, _ := original.Get(key).(string)
	return r
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
