package immutable

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/object88/immutable/memory"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	max           = 500000
)

var keys []IntKey
var contents map[Key]Value
var src = rand.NewSource(time.Now().UnixNano())

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func setup() {
	stringLength := 100
	contents = make(map[Key]Value, max)
	keys = make([]IntKey, max)
	for i := 0; i < max; i++ {
		keys[i] = IntKey(i)
		contents[keys[i]] = randStringBytesMaskImprSrc(stringLength)
	}
}

func shutdown() {
}

var result string

func compareBucketStrategy(blockSize memory.BlockSize) {
	options := NewHashMapOptions()
	options.BucketStrategy = blockSize
	original := NewHashMap(contents, options)
	var r string
	for key := range contents {
		r = original.Get(key).(string)
	}
	result = r
}

func Benchmark_LargeBlock(b *testing.B) {
	compareBucketStrategy(memory.LargeBlock)
}

func Benchmark_ExtraLargeBlock(b *testing.B) {
	compareBucketStrategy(memory.ExtraLargeBlock)
}

func Benchmark_NoPackingBlock(b *testing.B) {
	compareBucketStrategy(memory.NoPacking)
}

func Benchmark_NativeMap(b *testing.B) {
	var r string
	for _, v := range keys {
		r = contents[v].(string)
	}
	result = r
}

// This code copied directly from StackOverflow:
// http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func randStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
