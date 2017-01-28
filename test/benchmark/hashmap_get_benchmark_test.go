package benchmark

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

var keys []int
var contents map[int]string
var result string
var src = rand.NewSource(time.Now().UnixNano())

var hashmapSmallBlock, hashmapLargeBlock, hashmapExtraLargeBlock, hashmapNoPacked *immutable.IntToStringHashmap

func init() {
	stringLength := 100
	contents = make(map[int]string, max)
	keys = make([]int, max)
	for i := 0; i < max; i++ {
		keys[i] = i
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
			r = contents[key]
		}
		result = r
	}
}

func createWithStragety(blocksize memory.BlockSize) *immutable.IntToStringHashmap {
	return immutable.NewIntToStringHashmap(contents)
}

func testStrategy(original *immutable.IntToStringHashmap) {
	var r string
	for _, key := range keys {
		r = original.Get(key)
	}
	result = r
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
