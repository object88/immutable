package immutable

import (
	"fmt"

	"github.com/object88/immutable/hasher"
)

// IntElement is an integer-based Element
type IntElement int

// // Hash calculates the 32-bit hash
// func (k IntElement) Hash() uint32 {
// 	hasher := fnv.New32a()
//
// 	binary.Write(hasher, binary.LittleEndian, uint32(k))
// 	hash := hasher.Sum32()
// 	return hash
// }

// Hash does what Hash cannot.
func (k IntElement) Hash(seed uint32) uint64 {
	return hasher.Hash8(uintptr(k), seed)
}

func (k IntElement) String() string {
	return fmt.Sprintf("%d", int(k))
}

// Hash calculates the 32-bit hash
func (k StringElement) Hash(seed uint32) uint64 {
	return 0
}

func (k StringElement) String() string {
	return string(k)
}
