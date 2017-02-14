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

var keys []int
var contents map[int]string
var result string
var src = rand.NewSource(time.Now().UnixNano())

func init() {
	stringLength := 100
	contents = make(map[int]string, max[len(max)-1])
	keys = make([]int, max[len(max)-1])
	for j := 0; j < max[len(max)-1]; j++ {
		k := int(src.Int63())
		keys[j] = k
		contents[k] = generateString(stringLength)
	}
}

func Benchmark_Hashmap_Get_PackedBlock(b *testing.B) {
	for _, tc := range max {
		data := prepareData(tc)
		hashmap := createWithStragety(data, true)
		b.Run(fmt.Sprintf("n=%d", tc), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testStrategy(hashmap)
			}
		})
	}
}

func Benchmark_Hashmap_Get_NotPackedBlock(b *testing.B) {
	for _, tc := range max {
		data := prepareData(tc)
		hashmap := createWithStragety(data, false)
		b.Run(fmt.Sprintf("n=%d", tc), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				testStrategy(hashmap)
			}
		})
	}
}

func Benchmark_Hashmap_Get_NativeMap(b *testing.B) {
	for _, tc := range max {
		data := prepareData(tc)
		b.Run(fmt.Sprintf("n=%d", tc), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				var r string
				for _, key := range keys {
					r = getStringFromMap(data, key)
				}
				result = r
			}
		})
	}
}

func getStringFromMap(data map[int]string, key int) string {
	s := data[key]
	return s
}

func createWithStragety(contents map[int]string, packed bool) *immutable.IntToStringHashmap {
	return immutable.NewIntToStringHashmap(contents, core.WithBucketStrategy(packed))
}

func testStrategy(original *immutable.IntToStringHashmap) {
	var r string
	for _, key := range keys {
		r = getStringFromImmutable(original, key)
	}
	result = r
}

func getStringFromImmutable(original *immutable.IntToStringHashmap, key int) string {
	s, _, _ := original.Get(key)
	return s
}

func prepareData(tc int) map[int]string {
	var data map[int]string
	if tc == max[len(max)-1] {
		data = contents
	} else {
		data = make(map[int]string, tc)
		for i := 0; i < tc; i++ {
			data[i] = contents[i]
		}
	}
	return data
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
