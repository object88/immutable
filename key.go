package immutable

import (
	"fmt"

	"github.com/object88/immutable/hasher"
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

// Hash does what Hash cannot.
func (k IntKey) Hash(seed uint32) uint64 {
	return hasher.Hash8(uintptr(k), seed)
}

type IntKeyMetadata struct{}

func (IntKeyMetadata) Indirect() bool {
	return false
}

func (IntKeyMetadata) StorageSize() int {
	return 8
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
