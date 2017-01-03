package immutable

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
)

// IntKey is an integer-based Key
type IntKey int

// StringKey is a string-based Key
type StringKey string

// Hash calculates the 32-bit hash for the given IntKey
func (k IntKey) Hash() uint32 {
	hasher := fnv.New32a()

	binary.Write(hasher, binary.LittleEndian, uint32(k))
	hash := hasher.Sum32()
	return hash
}

func (k IntKey) String() string {
	return fmt.Sprintf("%d", int(k))
}

func (k StringKey) Hash() uint32 {
	return 0
}

func (k StringKey) String() string {
	return string(k)
}
