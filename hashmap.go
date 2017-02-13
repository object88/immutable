package immutable

import (
	"errors"
	"math"
	"math/rand"
	"time"
	"unsafe"

	"github.com/object88/immutable/core"
	"github.com/object88/immutable/memory"
)

type bucket struct {
	entryCount byte
	hobs       memory.Memories
	keys       unsafe.Pointer
	values     unsafe.Pointer
	overflow   *bucket
}

// HashMap is a read-only key-to-value collection
type HashMap struct {
	size    int
	buckets []*bucket
	lobMask uint32
	lobSize uint8
	seed    uint32
}

const (
	bucketCapacity = 8
	loadFactor     = 6.0
)

var emptyHashmap = &HashMap{0, nil, 0, 0, 0}

func CreateEmptyHashmap(size int) *HashMap {
	if size == 0 {
		return emptyHashmap
	}

	initialCount := size
	initialSize := memory.NextPowerOfTwo(int(math.Ceil(float64(initialCount) / loadFactor)))
	lobSize := memory.PowerOf(initialSize)
	lobMask := uint32(^(0xffffffff << lobSize))
	buckets := make([]*bucket, initialSize)

	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)
	seed := uint32(random.Int31())

	return &HashMap{initialCount, buckets, lobMask, lobSize, seed}
}

// Get returns the value for the given key
func (h *HashMap) Get(config *core.HashmapConfig, key unsafe.Pointer) (result unsafe.Pointer, ok bool, err error) {
	if key == nil {
		return nil, false, errors.New("Element is nil")
	}
	if h.size == 0 {
		return nil, false, nil
	}

	hashkey := config.KeyConfig.Hash(key, h.seed)

	selectedBucket := hashkey & uint64(h.lobMask)
	b := h.buckets[selectedBucket]
	maskedHash := hashkey >> h.lobSize

	totalEntries := uint64(b.entryCount)

	// fmt.Printf("\nlobSize: %d; h.lobMask: 0x%016x\n", h.lobSize, h.lobMask)
	// fmt.Printf(
	// 	"hashKey: 0x%016x / selectedBucket: %d / mashedHash: 0x%016x\n",
	// 	hashkey,
	// 	selectedBucket,
	// 	maskedHash)

	for b != nil {
		for index := 0; index < int(totalEntries); index++ {
			if b.hobs.Read(uint64(index)) != maskedHash {
				continue
			}

			// fmt.Printf(
			// 	"0x%016x <-> 0x%016x :: %#v <-> %s\n",
			// 	b.hobs.Read(uint64(index)),
			// 	maskedHash,
			// 	k, key)
			if config.KeyConfig.CompareTo(b.keys, index, key) {
				v := config.ValueConfig.Read(b.values, index)
				return v, true, nil
			}
		}
		b = b.overflow
	}
	return nil, false, nil
}

// GetKeys returns an array of keys in the hashmap.  If there are no entries,
// then an empty array is returned.  If the pointer reciever is nil, then
// nil is returned.  The array of keys is not ordered.
func (h *HashMap) GetKeys(config *core.HashmapConfig) (results []unsafe.Pointer, err error) {
	if h.size == 0 {
		return []unsafe.Pointer{}, nil
	}

	results = make([]unsafe.Pointer, h.size)
	count := 0
	for i := 0; i < len(h.buckets); i++ {
		b := h.buckets[i]
		if b == nil {
			continue
		}
		for j := 0; j < int(b.entryCount); j++ {
			v := config.KeyConfig.Read(b.keys, j)
			results[count] = v
			count++
		}
	}

	return results, nil
}

func (h *HashMap) iterate(config *core.HashmapConfig, abort <-chan struct{}) <-chan core.KeyValuePair {
	ch := make(chan core.KeyValuePair)

	go func() {
		defer close(ch)
		for i := 0; i < len(h.buckets); i++ {
			b := h.buckets[i]
			for b != nil {
				for j := 0; j < int(b.entryCount); j++ {
					k := config.KeyConfig.Read(b.keys, j)
					v := config.ValueConfig.Read(b.values, j)

					select {
					case ch <- core.KeyValuePair{Key: k, Value: v}:
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
func (h *HashMap) Filter(config *core.HashmapConfig, predicate FilterPredicate) (*HashMap, error) {
	b := &BaseStruct{h}
	result, err := b.filter(config, predicate)
	if err != nil {
		return nil, err
	}
	return result.Base.(*HashMap), nil
}

// ForEach iterates over each key-value pair in this collection
func (h *HashMap) ForEach(config *core.HashmapConfig, predicate ForEachPredicate) {
	b := &BaseStruct{h}
	b.forEach(config, predicate)
}

// Insert returns a new collection with the provided key-value pair added.
func (h *HashMap) Insert(config *core.HashmapConfig, key unsafe.Pointer, value unsafe.Pointer) (*HashMap, error) {
	if h.size == 0 {
		result := CreateEmptyHashmap(1)
		result.internalSet(config, key, value)
		return result, nil
	}

	foundValue, ok, _ := h.Get(config, key)
	matched := ok
	if matched && config.ValueConfig.Compare(foundValue, value) {
		return h, nil
	}

	var result *HashMap
	abort := make(chan struct{})
	size := h.Size()
	if matched {
		result = CreateEmptyHashmap(size)
		for kvp := range h.iterate(config, abort) {
			insertValue := kvp.Value
			if config.KeyConfig.Compare(kvp.Key, key) {
				insertValue = value
			}
			result.internalSet(config, kvp.Key, insertValue)
		}
	} else {
		size++
		result = CreateEmptyHashmap(size)
		for kvp := range h.iterate(config, abort) {
			result.internalSet(config, kvp.Key, kvp.Value)
		}

		result.internalSet(config, key, value)
	}

	return result, nil
}

// Map iterates over the contents of a collection and calls the supplied predicate.
// The return value is a new map with the results of the predicate function.
func (h *HashMap) Map(config *core.HashmapConfig, predicate MapPredicate) (*HashMap, error) {
	b := &BaseStruct{h}
	result, err := b.mapping(config, predicate)
	if err != nil {
		return nil, err
	}
	return result.Base.(*HashMap), nil
}

// Reduce operates over the collection contents to produce a single value
func (h *HashMap) Reduce(config *core.HashmapConfig, predicate ReducePredicate, accumulator unsafe.Pointer) (unsafe.Pointer, error) {
	b := &BaseStruct{h}
	return b.reduce(config, predicate, accumulator)
}

// Remove returns a copy of the provided HashMap with the specified element
// removed.
func (h *HashMap) Remove(config *core.HashmapConfig, key unsafe.Pointer) (*HashMap, error) {
	if h.size == 0 {
		return h, nil
	}

	if _, ok, _ := h.Get(config, key); !ok {
		return h, nil
	}

	newSize := h.Size() - 1
	if newSize == 0 {
		return CreateEmptyHashmap(0), nil
	}

	result := CreateEmptyHashmap(newSize)
	abort := make(chan struct{})
	for kvp := range h.iterate(config, abort) {
		if !config.KeyConfig.Compare(kvp.Key, key) {
			result.internalSet(config, kvp.Key, kvp.Value)
		}
	}

	return result, nil
}

// Size returns the number of items in this collection
func (h *HashMap) Size() int {
	return h.size
}

func (h *HashMap) instantiate(config *core.HashmapConfig, size int, contents []*core.KeyValuePair) *BaseStruct {
	hash := CreateEmptyHashmap(size)

	for _, v := range contents {
		if v != nil {
			hash.internalSet(config, v.Key, v.Value)
		}
	}

	return &BaseStruct{hash}
}

func (h *HashMap) internalSet(config *core.HashmapConfig, key unsafe.Pointer, value unsafe.Pointer) {
	hobSize := uint32(64 - h.lobSize)

	hashkey := config.KeyConfig.Hash(key, h.seed)
	selectedBucket := hashkey & uint64(h.lobMask)
	b := h.buckets[selectedBucket]
	if b == nil {
		b = createEmptyBucket(config, hobSize)
		h.buckets[selectedBucket] = b
	}
	for b.entryCount == bucketCapacity {
		if b.overflow == nil {
			b.overflow = createEmptyBucket(config, hobSize)
		}
		b = b.overflow
	}
	index := int(b.entryCount)
	config.KeyConfig.Write(b.keys, index, key)
	config.ValueConfig.Write(b.values, index, value)
	b.hobs.Assign(uint64(b.entryCount), hashkey>>h.lobSize)
	b.entryCount++
}

func createEmptyBucket(config *core.HashmapConfig, hobSize uint32) *bucket {
	return &bucket{
		entryCount: 0,
		hobs:       memory.AllocateMemories(memory.LargeBlock, hobSize, bucketCapacity),
		keys:       config.KeyConfig.CreateBucket(bucketCapacity),
		values:     config.ValueConfig.CreateBucket(bucketCapacity),
		overflow:   nil,
	}
}
