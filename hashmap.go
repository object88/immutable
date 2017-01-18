package immutable

import (
	"errors"
	"math"
	"math/rand"
	"time"
	"unsafe"

	"github.com/object88/immutable/memory"
)

type bucket struct {
	entryCount byte
	hobs       memory.Memories
	entries    memory.Memories
	overflow   *bucket
}

// HashMap is a read-only key-to-value collection
type HashMap struct {
	meta    KeyMetadata
	options *HashMapOptions
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

// NewHashMap creates a new instance of a HashMap
func NewHashMap(meta KeyMetadata, contents map[Key]Value, options ...HashMapOption) *HashMap {
	opts := defaultHashMapOptions()
	for _, fn := range options {
		fn(opts)
	}

	hash := createHashMap(len(contents), meta, opts)

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
			if b.hobs.Read(index) != maskedHash {
				continue
			}

			i := 2 * index
			if b.entries.Read(i) == key {
				v := *unsafe.Pointer(uintptr(b.entries.Read(i + 1)))
				return v.(interface{}), true, nil
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
		result = createHashMap(size, h.meta, h.options)
		for kvp := range h.iterate(abort) {
			insertValue := kvp.value
			if kvp.key == key {
				insertValue = value
			}
			result.internalSet(kvp.key, insertValue)
		}
	} else {
		size++
		result = createHashMap(size, h.meta, h.options)
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
		return createHashMap(0, h.meta, h.options), nil
	}

	result := createHashMap(size, h.meta, h.options)
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
	hash := createHashMap(size, h.meta, h.options)

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
	for b.entryCount == bucketCapacity {
		if b.overflow == nil {
			b.overflow = createEmptyBucket(h.options.BucketStrategy, hobSize)
		}
		b = b.overflow
	}
	i := uint64(2 * b.entryCount)
	b.entries.Assign(i, key)
	b.entries.Assign(i+1, uint64(uintptr(unsafe.Pointer(&value))))
	// b.entries[b.entryCount] = keyValuePair{key, value}
	b.hobs.Assign(uint64(b.entryCount), hashkey>>h.lobSize)
	b.entryCount++
}

func createHashMap(size int, meta KeyMetadata, options *HashMapOptions) *HashMap {
	initialCount := size
	initialSize := memory.NextPowerOfTwo(int(math.Ceil(float64(initialCount) / loadFactor)))
	lobSize := memory.PowerOf(initialSize)
	lobMask := uint32(^(0xffffffff << lobSize))
	buckets := make([]*bucket, initialSize)

	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)
	seed := uint32(random.Int31())

	return &HashMap{meta, options, initialCount, buckets, lobMask, lobSize, seed}
}

func createEmptyBucket(blockSize memory.BlockSize, hobSize uint32) *bucket {
	return &bucket{
		entryCount: 0,
		hobs:       memory.AllocateMemories(blockSize, hobSize, bucketCapacity),
		entries:    memory.AllocateMemories(memory.NoPacking, 0, 2*bucketCapacity),
		overflow:   nil,
	}
}
