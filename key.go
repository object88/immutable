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

// // DehydrateIntKey takes an IntValue and returns a uint64
// func DehydrateIntKey(value Value) (result unsafe.Pointer, err error) {
// 	fmt.Printf("D: %#v\n", value.(IntKey))
// 	fmt.Printf("D: %d\n", int(value.(IntKey)))
// 	i := int(value.(IntKey))
// 	return unsafe.Pointer(&i), nil
// }
//
// // HydrateIntKey takes a uint64 and returns an IntValue
// func HydrateIntKey(value unsafe.Pointer) (result Value, err error) {
// 	fmt.Printf("H: %p\n", value)
// 	fmt.Printf("H: %p\n", (*int)(value))
// 	fmt.Printf("H: %d\n", *(*int)(value))
// 	return IntKey(*(*int)(value)), nil
// }

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
