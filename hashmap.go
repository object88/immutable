package immutable

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/object88/immutable/core"
	"github.com/object88/immutable/memory"
)

type bucket struct {
	entryCount byte
	hobs       memory.Memories
	keys       core.SubBucket
	values     core.SubBucket
	overflow   *bucket
}

// HashMap is a read-only key-to-value collection
type HashMap struct {
	options *core.HashMapOptions
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
func NewHashMap(contents map[core.Element]core.Element, options ...core.HashMapOption) *HashMap {
	opts := core.DefaultHashMapOptions()
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
func (h *HashMap) Get(key core.Element) (result core.Element, ok bool, err error) {
	if h == nil {
		return nil, false, errors.New("Pointer receiver is nil")
	}
	if key == nil {
		return nil, false, errors.New("Element is nil")
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
		for index := 0; index < int(totalEntries); index++ {
			if b.hobs.Read(uint64(index)) != maskedHash {
				continue
			}

			k, ok := b.keys.Hydrate(index).(core.Element)
			if !ok {
				panic("NOPE")
			}
			// fmt.Printf(
			// 	"0x%016x <-> 0x%016x :: %#v <-> %s\n",
			// 	b.hobs.Read(uint64(index)),
			// 	maskedHash,
			// 	k, key)
			if k == key {
				v := b.values.Hydrate(index)

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
func (h *HashMap) GetKeys() (results []core.Element, err error) {
	if h == nil {
		return nil, errors.New("Pointer receiver is nil")
	}

	if h.size == 0 {
		return []core.Element{}, nil
	}

	results = make([]core.Element, h.size)
	count := 0
	for i := 0; i < len(h.buckets); i++ {
		b := h.buckets[i]
		if b == nil {
			continue
		}
		for j := 0; j < int(b.entryCount); j++ {
			v, _ := b.keys.Hydrate(j).(core.Element)
			results[count] = v.(core.Element)
			count++
		}
	}

	return results, nil
}

func (h *HashMap) iterate(abort <-chan struct{}) <-chan core.KeyValuePair {
	ch := make(chan core.KeyValuePair)

	go func() {
		defer close(ch)
		for i := 0; i < len(h.buckets); i++ {
			b := h.buckets[i]
			for b != nil {
				for j := 0; j < int(b.entryCount); j++ {
					k, _ := b.keys.Hydrate(j).(core.Element)
					v := b.values.Hydrate(j)

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
func (h *HashMap) Insert(key core.Element, value core.Element) (*HashMap, error) {
	if h == nil {
		return nil, errors.New("Pointer receiver is nil")
	}
	if key == nil {
		return nil, errors.New("Element is nil")
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
			insertValue := kvp.Value
			if kvp.Key == key {
				insertValue = value
			}
			result.internalSet(kvp.Key, insertValue)
		}
	} else {
		size++
		result = createHashMap(size, h.options)
		for kvp := range h.iterate(abort) {
			result.internalSet(kvp.Key, kvp.Value)
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
func (h *HashMap) Reduce(predicate ReducePredicate, accumulator core.Element) (core.Element, error) {
	if h == nil {
		return nil, nil
	}

	b := &BaseStruct{h}
	return b.reduce(predicate, accumulator)
}

// Remove returns a copy of the provided HashMap with the specified element
// removed.
func (h *HashMap) Remove(key core.Element) (*HashMap, error) {
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
		if kvp.Key != key {
			result.internalSet(kvp.Key, kvp.Value)
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

func (h *HashMap) instantiate(size int, contents []*core.KeyValuePair) *BaseStruct {
	hash := createHashMap(size, h.options)

	for _, v := range contents {
		if v != nil {
			hash.internalSet(v.Key, v.Value)
		}
	}

	return &BaseStruct{hash}
}

func (h *HashMap) internalSet(key core.Element, value core.Element) {
	hobSize := uint32(64 - h.lobSize)

	hashkey := key.Hash(h.seed)
	selectedBucket := hashkey & uint64(h.lobMask)
	b := h.buckets[selectedBucket]
	if b == nil {
		// Create the bucket.
		b = createEmptyBucket(h.options, hobSize)
		h.buckets[selectedBucket] = b
	}
	for b.entryCount == bucketCapacity {
		if b.overflow == nil {
			b.overflow = createEmptyBucket(h.options, hobSize)
		}
		b = b.overflow
	}
	index := int(b.entryCount)
	b.keys.Dehydrate(index, key)
	b.values.Dehydrate(index, value)
	b.hobs.Assign(uint64(b.entryCount), hashkey>>h.lobSize)
	b.entryCount++
}

func createHashMap(size int, options *core.HashMapOptions) *HashMap {
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

func createEmptyBucket(options *core.HashMapOptions, hobSize uint32) *bucket {
	// var keys, values valueInterface
	// if options.KeyDirect {
	// 	keys = NewDirect(8, bucketCapacity)
	// } else {
	// 	keys = NewIndirect(bucketCapacity)
	// }
	// if options.ValueDirect {
	// 	values = NewDirect(8, bucketCapacity)
	// } else {
	// 	values = NewIndirect(bucketCapacity)
	// }
	return &bucket{
		entryCount: 0,
		hobs:       memory.AllocateMemories(options.BucketStrategy, hobSize, bucketCapacity),
		keys:       options.KeyConfig.CreateBucket(bucketCapacity),
		values:     options.ValueConfig.CreateBucket(bucketCapacity),
		overflow:   nil,
	}
}
