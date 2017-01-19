package immutable

import (
	"errors"
	"math"
	"math/rand"
	"time"

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
	options *HashMapOptions
	size    int
	buckets []*bucket
	lobMask uint32
	lobSize uint8
	seed    uint32
}

const (
	bucketCapacity = 1 << 3
	loadFactor     = 6.0
)

// NewHashMap creates a new instance of a HashMap
func NewHashMap(contents map[Key]Value, options ...HashMapOption) *HashMap {
	opts := defaultHashMapOptions()
	for _, fn := range options {
		fn(opts)
	}

	hash := createHashMap(len(contents), opts)

	for k, v := range contents {
		hash.internalSet(k, v)
	}

	return hash
}

// Get returns the value for the given key
func (h *HashMap) Get(key Key) (result Value, ok bool, err error) {
	if h == nil {
		return nil, false, errors.New("Pointer receiver is nil")
	}
	if key == nil {
		return nil, false, errors.New("Key is nil")
	}
	if h.size == 0 {
		return nil, false, nil
	}

	hashkey := key.Hash(h.seed)

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
		for index := uint64(0); index < totalEntries; index++ {
			// fmt.Printf(
			// 	"0x%016x <-> 0x%016x :: %s <-> %s\n",
			// 	b.hobs.Read(index),
			// 	maskedHash,
			// 	b.entries[index].key, key)
			if b.hobs.Read(index) == maskedHash && b.entries[index].key == key {
				return b.entries[index].value, true, nil
			}
		}
		b = b.overflow
	}
	return nil, false, nil
}

// GetKeys returns an array of keys in the hashmap.  If there are no entries,
// then an empty array is returned.  If the pointer reciever is nil, then
// nil is returned.  The array of keys is not ordered.
func (h *HashMap) GetKeys() (results []Key, err error) {
	if h == nil {
		return nil, errors.New("Pointer receiver is nil")
	}

	if h.size == 0 {
		return []Key{}, nil
	}

	results = make([]Key, h.size)
	count := 0
	for i := 0; i < len(h.buckets); i++ {
		b := h.buckets[i]
		if b == nil {
			continue
		}
		for j := byte(0); j < b.entryCount; j++ {
			results[count] = b.entries[j].key
			count++
		}
	}

	return results, nil
}

func (h *HashMap) iterate(abort <-chan struct{}) <-chan keyValuePair {
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

	b := &BaseStruct{h}
	result, err := b.filter(predicate)
	if err != nil {
		return nil, err
	}
	return result.Base.(*HashMap), nil
}

// ForEach iterates over each key-value pair in this collection
func (h *HashMap) ForEach(predicate ForEachPredicate) {
	b := &BaseStruct{h}
	b.forEach(predicate)
}

// Insert returns a new collection with the provided key-value pair added.
// The pointer reciever may be nil; it will be treated as a instance with
// no contents.
func (h *HashMap) Insert(key Key, value Value) (*HashMap, error) {
	if h == nil {
		return nil, errors.New("Pointer receiver is nil")
	}
	if key == nil {
		return nil, errors.New("Key is nil")
	}

	if h.size == 0 {
		result := createHashMap(1, h.options)
		result.internalSet(key, value)
		return result, nil
	}

	foundValue, _, _ := h.Get(key)
	matched := foundValue != nil
	if matched && foundValue == value {
		return h, nil
	}

	var result *HashMap
	abort := make(chan struct{})
	size := h.Size()
	if matched {
		result = createHashMap(size, h.options)
		for kvp := range h.iterate(abort) {
			insertValue := kvp.value
			if kvp.key == key {
				insertValue = value
			}
			result.internalSet(kvp.key, insertValue)
		}
	} else {
		size++
		result = createHashMap(size, h.options)
		for kvp := range h.iterate(abort) {
			result.internalSet(kvp.key, kvp.value)
		}

		result.internalSet(key, value)
	}

	return result, nil
}

// Map iterates over the contents of a collection and calls the supplied predicate.
// The return value is a new map with the results of the predicate function.
func (h *HashMap) Map(predicate MapPredicate) (*HashMap, error) {
	if h == nil {
		return nil, nil
	}

	b := &BaseStruct{h}
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

	b := &BaseStruct{h}
	return b.reduce(predicate, accumulator)
}

// Remove returns a copy of the provided HashMap with the specified element
// removed.
func (h *HashMap) Remove(key Key) (*HashMap, error) {
	if h == nil {
		return nil, errors.New("Pointer receiver is nil")
	}

	if h.size == 0 {
		return h, nil
	}

	if _, ok, _ := h.Get(key); !ok {
		return h, nil
	}

	size := h.Size() - 1
	if size == 0 {
		return createHashMap(0, h.options), nil
	}

	result := createHashMap(size, h.options)
	abort := make(chan struct{})
	for kvp := range h.iterate(abort) {
		if kvp.key != key {
			result.internalSet(kvp.key, kvp.value)
		}
	}

	return result, nil
}

// Size returns the number of items in this collection
func (h *HashMap) Size() int {
	if h == nil {
		return 0
	}
	return h.size
}

func (h *HashMap) instantiate(size int, contents []*keyValuePair) *BaseStruct {
	hash := createHashMap(size, h.options)

	for _, v := range contents {
		if v != nil {
			hash.internalSet(v.key, v.value)
		}
	}

	return &BaseStruct{hash}
}

func (h *HashMap) internalSet(key Key, value Value) {
	hobSize := uint32(64 - h.lobSize)

	hashkey := key.Hash(h.seed)
	selectedBucket := hashkey & uint64(h.lobMask)
	b := h.buckets[selectedBucket]
	if b == nil {
		// Create the bucket.
		b = createEmptyBucket(h.options.BucketStrategy, hobSize)
		h.buckets[selectedBucket] = b
	}
	for b.entryCount == 8 {
		if b.overflow == nil {
			b.overflow = createEmptyBucket(h.options.BucketStrategy, hobSize)
		}
		b = b.overflow
	}
	b.entries[b.entryCount] = entry{key, value}
	b.hobs.Assign(uint64(b.entryCount), hashkey>>h.lobSize)
	b.entryCount++
}

func createHashMap(size int, options *HashMapOptions) *HashMap {
	initialCount := size
	initialSize := memory.NextPowerOfTwo(int(math.Ceil(float64(initialCount) / loadFactor)))
	lobSize := memory.PowerOf(initialSize)
	lobMask := uint32(^(0xffffffff << lobSize))
	buckets := make([]*bucket, initialSize)

	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)
	seed := uint32(random.Int31())

	return &HashMap{options, initialCount, buckets, lobMask, lobSize, seed}
}

func createEmptyBucket(blockSize memory.BlockSize, hobSize uint32) *bucket {
	return &bucket{
		entryCount: 0,
		hobs:       memory.AllocateMemories(blockSize, hobSize, 8),
		entries:    make([]entry, bucketCapacity),
		overflow:   nil,
	}
}
