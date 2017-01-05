package immutable

import (
	"math"

	"github.com/object88/immutable/memory"
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
	initialSize := memory.NextPowerOfTwo(uint32(math.Ceil(float64(initialCount) / loadFactor)))
	lobSize := memory.PowerOf(initialSize)

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

	lobSize := memory.PowerOf(uint32(len(h.buckets)))
	lobMask := uint32(^(0xffffffff << lobSize))

	selectedBucket := hashkey & lobMask
	b := h.buckets[selectedBucket]
	maskedHash := hashkey >> lobSize

	// fmt.Printf(
	// 	"hashkey: 0x%08x, lobSize: %d, lobMask: 0x%d, selectedBucket: 0x%08x, maskedHash: 0x%08x\n",
	// 	hashkey, lobSize, lobMask, selectedBucket, maskedHash)

	for b != nil {
		for index := byte(0); index < b.entryCount; index++ {
			if b.hobs.Read(uint32(index)) == maskedHash && b.entries[index].key == key {
				return b.entries[index].value
			}
		}
		b = b.overflow
	}
	return nil
}

// Iterate is a generator function
func (h *HashMap) Iterate(abort <-chan struct{}) <-chan keyValuePair {
	ch := make(chan keyValuePair)

	go func() {
		defer close(ch)
		for i := 0; i < len(h.buckets); i++ {
			b := h.buckets[i]
			for b != nil {
				for j := byte(0); j < b.entryCount; j++ {
					select {
					case ch <- keyValuePair{b.entries[j].key, b.entries[j].value}:
					case <-abort:
						return
					}
				}

				b = b.overflow
			}
		}
	}()
	return ch
}

// Filter returns a subset of the collection, based on the predicate supplied
func (h *HashMap) Filter(predicate FilterPredicate) (*HashMap, error) {
	if h == nil {
		return nil, nil
	}

	b := &BaseStruct{h, h}
	result, err := b.filter(predicate)
	if err != nil {
		return nil, err
	}
	return result.Base.(*HashMap), nil
}

// ForEach iterates over each key-value pair in this collection
func (h *HashMap) ForEach(predicate ForEachPredicate) {
	b := &BaseStruct{h, h}
	b.forEach(predicate)
}

// Map iterates over the contents of a collection and calls the supplied predicate.
// The return value is a new map with the results of the predicate function.
func (h *HashMap) Map(predicate MapPredicate) (*HashMap, error) {
	if h == nil {
		return nil, nil
	}

	b := &BaseStruct{h, h}
	result, err := b.mapping(predicate)
	if err != nil {
		return nil, err
	}
	return result.Base.(*HashMap), nil
}

// Reduce operates over the collection contents to produce a single value
func (h *HashMap) Reduce(predicate ReducePredicate, accumulator Value) (Value, error) {
	if h == nil {
		return nil, nil
	}

	b := &BaseStruct{h, h}
	return b.reduce(predicate, accumulator)
}

// Size returns the number of items in this collection
func (h *HashMap) Size() uint32 {
	if h == nil {
		return 0
	}
	return h.count
}

func (h *HashMap) instantiate(size uint32) *BaseStruct {
	initialCount := size
	initialSize := memory.NextPowerOfTwo(uint32(math.Ceil(float64(initialCount) / loadFactor)))
	lobSize := memory.PowerOf(initialSize)
	buckets := make([]*bucket, initialSize)

	hash := &HashMap{initialCount, initialSize, buckets, lobSize}
	return &BaseStruct{Base: hash, internalFunctions: hash}
}

func (h *HashMap) instantiateWithContents(size uint32, contents []*keyValuePair) *BaseStruct {
	newHashMap := h.instantiate(size)

	for _, v := range contents {
		if v != nil {
			newHashMap.internalSet(v.key, v.value)
		}
	}

	return newHashMap
}

func (h *HashMap) internalSet(key Key, value Value) {
	lobSize := memory.PowerOf(h.size)
	hobSize := uint32(32 - lobSize)
	lobMask := uint32(^(0xffffffff << lobSize))

	hashkey := key.Hash()
	selectedBucket := hashkey & lobMask
	// fmt.Printf("At [%s,%s]; h:0x%08x, sb: %d, lob: 0x%08x\n", key, value, hashkey, selectedBucket, hashkey>>h.lobSize)
	b := h.buckets[selectedBucket]
	if b == nil {
		// Create the bucket.
		b = createEmptyBucket(memory.LargeBlock, hobSize)
		h.buckets[selectedBucket] = b
	}
	for b.entryCount == 8 {
		if b.overflow == nil {
			b.overflow = createEmptyBucket(memory.LargeBlock, hobSize)
		}
		b = b.overflow
	}
	b.entries[b.entryCount] = entry{key, value}
	b.hobs.Assign(uint32(b.entryCount), hashkey>>h.lobSize)
	b.entryCount++
}

func createEmptyBucket(blockSize memory.BlockSize, hobSize uint32) *bucket {
	return &bucket{
		entryCount: 0,
		hobs:       memory.AllocateMemories(blockSize, hobSize, 8),
		entries:    make([]entry, bucketCapacity),
		overflow:   nil,
	}
}
