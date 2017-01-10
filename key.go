package immutable

import (
	"fmt"
	"math/rand"
	"time"
)

// IntKey is an integer-based Key
type IntKey int

// StringKey is a string-based Key
type StringKey string

// // Hash calculates the 32-bit hash
// func (k IntKey) Hash() uint32 {
// 	hasher := fnv.New32a()
//
// 	binary.Write(hasher, binary.LittleEndian, uint32(k))
// 	hash := hasher.Sum32()
// 	return hash
// }

const m1 = 194865226553
const m2 = 24574600569641

var hashkey [4]uintptr

func init() {
	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)
	// getRandomData((*[len(hashkey) * sys.PtrSize]byte)(unsafe.Pointer(&hashkey))[:])
	hashkey[0] = uintptr(random.Int63() | 1) // make sure these numbers are odd
	hashkey[1] = uintptr(random.Int63() | 1)
	hashkey[2] = uintptr(random.Int63() | 1)
	hashkey[3] = uintptr(random.Int63() | 1)
}

// Hash does what Hash cannot.
func (k IntKey) Hash(seed uint32) uint64 {
	k1 := uint64(k)
	h := uint64(uintptr(seed) + 8*hashkey[0])
	h ^= (k1 & 0xffffffff)
	h ^= (k1 & 0xffffffff00000000) << 32
	h = rotl31(h*m1) * m2
	return h
}

func rotl31(x uint64) uint64 {
	return (x << 31) | (x >> (64 - 31))
}

func (k IntKey) String() string {
	return fmt.Sprintf("%d", int(k))
}

// Hash calculates the 32-bit hash
func (k StringKey) Hash() uint32 {
	return 0
}

func (k StringKey) String() string {
	return string(k)
}
