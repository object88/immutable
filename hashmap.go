package immutable

import (
	"math"

	"github.com/object88/memory"
)

type bucket struct {
	entryCount byte
	hobs       memory.Memories
	entries    []entry
	overflow   *bucket
}

type entry struct {
	key   Key
	value Value
}

// HashMap is a read-only key-to-value collection
type HashMap struct {
	count   uint32
	size    uint32
	buckets []*bucket
	lobSize uint32
}

const (
	// lobSize        = 3
	bucketCapacity = 1 << 3 //lobSize
	loadFactor     = 6.0
)

// NewHashMap creates a new instance of a HashMap
func NewHashMap(contents map[Key]Value) *HashMap {
	initialCount := uint32(len(contents))
	initialSize := nextPowerOfTwo(uint32(math.Ceil(float64(initialCount) / loadFactor)))
	lobSize := fffff(initialSize)

	buckets := make([]*bucket, initialSize)

	hash := &HashMap{initialCount, initialSize, buckets, lobSize}

	for k, v := range contents {
		hash.internalSet(k, v)
	}

	return hash
}

// Get returns the value for the given key
func (h *HashMap) Get(key Key) Value {
	hashkey := key.Hash()

	lobSize := fffff(uint32(len(h.buckets)))
	lobMask := uint32(^(0xffffffff << lobSize))

	selectedBucket := hashkey & lobMask
	b := h.buckets[selectedBucket]
	index := uint32(0)
	maskedHash := hashkey >> lobSize

	// fmt.Printf("hashkey: 0x%08x, lobSize: %d, lobMask: 0x%d, selectedBucket: 0x%08x, maskedHash: 0x%08x\n", hashkey, lobSize, lobMask, selectedBucket, maskedHash)

	for ; index < bucketCapacity && b.hobs.Read(index) != maskedHash; index++ {
	}
	if index == bucketCapacity {
		// fmt.Printf("Returning nil...\n")
		return nil
	}
	return b.entries[index].value
}

// Iterate loops through all contents
func (h *HashMap) Iterate() Iterator {
	i, j := uint32(0), byte(0)
	// i_last := int(h.count)
	var iterator Iterator
	iterator = func() (key Key, value Value, next Iterator) {
		for ; i < uint32(len(h.buckets)); i++ {
			b := h.buckets[i]
			if b == nil {
				// fmt.Printf("At %d, empty bucket\n", i)
				continue
			}
			if j == b.entryCount {
				// fmt.Printf("At [%d:%d], at last entry\n", i, j)
				j = 0
				continue
			}
			e := b.entries[j]
			k := e.key
			v := e.value
			// fmt.Printf("At [%d:%d], got %s->%s\n", i, j, k, v)
			j++
			return k, v, iterator
		}

		return nil, nil, nil
	}
	return iterator
}

// Size returns the number of items in this collection
func (h *HashMap) Size() uint32 {
	if h == nil {
		return 0
	}
	return h.count
}

// Map iterates over the contents of a collection and calls the supplied predicate.
// The return value is a new map with the results of the predicate function.
func (h *HashMap) Map(predicate MapPredicate) (*HashMap, error) {
	b := &BaseStruct{h, h}
	n, e := b.Map(predicate)
	return n.Base.(*HashMap), e
}

func (h *HashMap) instantiate(size uint32) *BaseStruct {
	initialCount := size
	initialSize := nextPowerOfTwo(uint32(math.Ceil(float64(initialCount) / loadFactor)))
	lobSize := fffff(initialSize)
	buckets := make([]*bucket, initialSize)

	hash := &HashMap{initialCount, initialSize, buckets, lobSize}
	return &BaseStruct{Base: hash, internalFunctions: hash}
}

func (h *HashMap) internalSet(key Key, value Value) {
	lobSize := fffff(h.size)
	hobSize := uint32(32 - lobSize)
	lobMask := uint32(^(0xffffffff << lobSize))

	hashkey := key.Hash()
	selectedBucket := hashkey & lobMask
	// fmt.Printf("At [%s,%s]; h:0x%08x, sb: %d, lob: 0x%08x\n", key, value, hashkey, selectedBucket, hashkey>>h.lobSize)
	b := h.buckets[selectedBucket]
	if b == nil {
		// Create the bucket.
		b = &bucket{
			entryCount: 0,
			hobs:       memory.AllocateMemories(memory.LargeBlock, hobSize, 8),
			entries:    make([]entry, bucketCapacity),
			overflow:   nil,
		}
		h.buckets[selectedBucket] = b
	}
	b.entries[b.entryCount] = entry{key, value}
	b.hobs.Assign(uint32(b.entryCount), hashkey>>h.lobSize)
	b.entryCount++
}

// http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
func nextPowerOfTwo(size uint32) uint32 {
	size--
	size |= size >> 1
	size |= size >> 2
	size |= size >> 4
	size |= size >> 8
	size |= size >> 16
	size++
	return size
}

// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetKernighan
func fffff(value uint32) uint32 {
	v := value - 1

	c := uint32(0) // c accumulates the total bits set in v
	for ; v != 0; c++ {
		v &= v - 1 // clear the least significant bit set
	}
	return c
}
