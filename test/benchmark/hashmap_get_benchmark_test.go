package benchmark

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/object88/immutable"
	"github.com/object88/immutable/core"
)

const (
	letterBytes     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIndexBits = 6                      // 6 bits to represent a letter index
	letterIndexMask = 1<<letterIndexBits - 1 // All 1-bits, as many as letterIndexBits
	letterIndexMax  = 63 / letterIndexBits   // # of letter indices fitting in 63 bits
)

var max = []int{5, 50, 500, 5000, 50000, 500000, 5000000}
var result string
var src = rand.NewSource(time.Now().UnixNano())

func Benchmark_Hashmap_Get_One(b *testing.B) {
	hashmapSizes := []int{5, 50, 500, 5000, 50000, 500000, 5000000}
	stringLength := 100
	for _, hashmapSize := range hashmapSizes {
		contents := make(map[int]string, hashmapSize)
		key := int(src.Int63())
		contents[key] = generateString(stringLength)
		for len(contents) < hashmapSize {
			k := int(src.Int63())
			contents[k] = generateString(stringLength)
		}

		b.Run(fmt.Sprintf("n=%d/Packed", hashmapSize), func(b *testing.B) {
			hashmap := immutable.NewIntToStringHashmap(contents, core.WithBucketStrategy(true))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				testStrategy(key, hashmap)
			}
		})

		b.Run(fmt.Sprintf("n=%d/Unpacked", hashmapSize), func(b *testing.B) {
			hashmap := immutable.NewIntToStringHashmap(contents, core.WithBucketStrategy(false))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				testStrategy(key, hashmap)
			}
		})

		b.Run(fmt.Sprintf("n=%d/Native", hashmapSize), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testNative(key, contents)
			}
		})
	}
}

func testNative(key int, data map[int]string) {
	s := data[key]
	result = s
}

func testStrategy(key int, original *immutable.IntToStringHashmap) {
	s, _, _ := original.Get(key)
	result = s
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
