package hasher

import (
	"math/rand"
	"time"

	"github.com/OneOfOne/xxhash"
)

const m1 = 194865226553
const m2 = 24574600569641

var hashkey [4]uintptr

func init() {
	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)
	hashkey[0] = uintptr(random.Int63() | 1) // make sure these numbers are odd
	hashkey[1] = uintptr(random.Int63() | 1)
	hashkey[2] = uintptr(random.Int63() | 1)
	hashkey[3] = uintptr(random.Int63() | 1)
}

// Hash8 calculates the hash for an 8-byte (64-bit) value
func Hash8(p uint64, seed uint32) uint64 {
	// k1 := uint64(*(*byte)(unsafe.Pointer(p)))
	k1 := p
	h := uint64(uintptr(seed) + 8*hashkey[0])
	h ^= (k1 & 0xffffffff)
	h ^= (k1 & 0xffffffff00000000) << 32
	h = rotl31(h*m1) * m2
	return h
}

func HashString(s string, seed uint32) uint64 {
	return xxhash.ChecksumString64S(s, uint64(seed))
}

func rotl31(x uint64) uint64 {
	return (x << 31) | (x >> (64 - 31))
}
