package core

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"
	"unsafe"

	"github.com/object88/immutable/memory"
)

type bucket struct {
	entryCount byte
	hobs       memory.Memories
	keys       unsafe.Pointer
	values     unsafe.Pointer
	overflow   *bucket
}

// InternalHashmap is a read-only key-to-value collection
type InternalHashmap struct {
	size      int
	buckets   []*bucket
	lobMask   uint64
	lobSize   uint8
	seed      uint32
	blockSize memory.BlockSize
}

const (
	bucketCapacity = 8
	loadFactor     = 6.0
)

var emptyHashmap = &InternalHashmap{0, nil, 0, 0, 0, memory.SmallBlockNoPacking}

func CreateEmptyInternalHashmap(packed bool, size int) *InternalHashmap {
	if size == 0 {
		return emptyHashmap
	}

	initialCount := size
	initialSize := memory.NextPowerOfTwo(int(math.Ceil(float64(initialCount) / loadFactor)))
	lobSize := memory.PowerOf(initialSize)
	lobMask := uint64(^(0xffffffffffffffff << lobSize))
	buckets := make([]*bucket, initialSize)

	hobSize := 64 - lobSize
	blockSize := memory.SelectBlockSize(packed, hobSize)

	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)
	seed := uint32(random.Int31())

	return &InternalHashmap{initialCount, buckets, lobMask, lobSize, seed, blockSize}
}

// Get returns the value for the given key
func (h *InternalHashmap) Get(config *HashmapConfig, key unsafe.Pointer) (result unsafe.Pointer, ok bool, err error) {
	if key == nil {
		return nil, false, errors.New("Key is nil")
	}
	if h.size == 0 {
		return nil, false, nil
	}

	hashkey := config.KeyConfig.Hash(key, h.seed)

	selectedBucket := hashkey & h.lobMask
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
func (h *InternalHashmap) GetKeys(config *HashmapConfig) (results []unsafe.Pointer, err error) {
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

// GoString provides a programmatic view into a InternalHashmap.  This may be used,
// for example, with the '%#v' operand to fmt.Printf, fmt.Sprintf, etc.
func (h *InternalHashmap) GoString(config *HashmapConfig) string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Size: %d\n[\n", h.size))
	for k, v := range h.buckets {
		b := v
		if b == nil {
			buffer.WriteString(fmt.Sprintf("  bucket #%d: nil\n", k))
			continue
		}
		buffer.WriteString(fmt.Sprintf("  bucket #%d: {\n    entryCount: %d\n    entries: [\n", k, b.entryCount))
		for b != nil {
			for i := 0; i < int(b.entryCount); i++ {
				key := config.KeyConfig.Read(b.keys, i)
				value := config.ValueConfig.Read(b.values, i)

				ks := config.KeyConfig.Format(key)
				vs := config.ValueConfig.Format(value)
				buffer.WriteString(fmt.Sprintf("      [0x%016x,%s] -> %s\n", b.hobs.Read(uint64(i)), ks, vs))
			}

			b = b.overflow
		}
		buffer.WriteString("    ]\n  },\n")
	}
	buffer.WriteString("]\n")
	return buffer.String()
}

func (h *InternalHashmap) iterate(config *HashmapConfig, abort <-chan struct{}) <-chan KeyValuePair {
	ch := make(chan KeyValuePair)

	go func() {
		defer close(ch)
		for i := 0; i < len(h.buckets); i++ {
			b := h.buckets[i]
			for b != nil {
				for j := 0; j < int(b.entryCount); j++ {
					k := config.KeyConfig.Read(b.keys, j)
					v := config.ValueConfig.Read(b.values, j)

					select {
					case ch <- KeyValuePair{Key: k, Value: v}:
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
func (h *InternalHashmap) Filter(config *HashmapConfig, predicate FilterPredicate) (*InternalHashmap, error) {
	b := &BaseStruct{h}
	result, err := b.filter(config, predicate)
	if err != nil {
		return nil, err
	}
	return result.Base.(*InternalHashmap), nil
}

// ForEach iterates over each key-value pair in this collection
func (h *InternalHashmap) ForEach(config *HashmapConfig, predicate ForEachPredicate) {
	b := &BaseStruct{h}
	b.forEach(config, predicate)
}

// Insert returns a new collection with the provided key-value pair added.
func (h *InternalHashmap) Insert(config *HashmapConfig, key unsafe.Pointer, value unsafe.Pointer) (*InternalHashmap, error) {
	if h.size == 0 {
		result := CreateEmptyInternalHashmap(config.Options.PackedBucket, 1)
		result.InternalSet(config, key, value)
		return result, nil
	}

	foundValue, ok, _ := h.Get(config, key)
	matched := ok
	if matched && config.ValueConfig.Compare(foundValue, value) {
		return h, nil
	}

	var result *InternalHashmap
	abort := make(chan struct{})
	size := h.Size()
	if matched {
		result = CreateEmptyInternalHashmap(config.Options.PackedBucket, size)
		for kvp := range h.iterate(config, abort) {
			insertValue := kvp.Value
			if config.KeyConfig.Compare(kvp.Key, key) {
				insertValue = value
			}
			result.InternalSet(config, kvp.Key, insertValue)
		}
	} else {
		size++
		result = CreateEmptyInternalHashmap(config.Options.PackedBucket, size)
		for kvp := range h.iterate(config, abort) {
			result.InternalSet(config, kvp.Key, kvp.Value)
		}

		result.InternalSet(config, key, value)
	}

	return result, nil
}

// Map iterates over the contents of a collection and calls the supplied predicate.
// The return value is a new map with the results of the predicate function.
func (h *InternalHashmap) Map(config *HashmapConfig, predicate MapPredicate) (*InternalHashmap, error) {
	b := &BaseStruct{h}
	result, err := b.mapping(config, predicate)
	if err != nil {
		return nil, err
	}
	return result.Base.(*InternalHashmap), nil
}

// Reduce operates over the collection contents to produce a single value
func (h *InternalHashmap) Reduce(config *HashmapConfig, predicate ReducePredicate) error {
	b := &BaseStruct{h}
	return b.reduce(config, predicate)
}

// Remove returns a copy of the provided InternalHashmap with the specified element
// removed.
func (h *InternalHashmap) Remove(config *HashmapConfig, key unsafe.Pointer) (*InternalHashmap, error) {
	if h.size == 0 {
		return h, nil
	}

	if _, ok, _ := h.Get(config, key); !ok {
		return h, nil
	}

	newSize := h.Size() - 1
	if newSize == 0 {
		return CreateEmptyInternalHashmap(config.Options.PackedBucket, 0), nil
	}

	result := CreateEmptyInternalHashmap(config.Options.PackedBucket, newSize)
	abort := make(chan struct{})
	for kvp := range h.iterate(config, abort) {
		if !config.KeyConfig.Compare(kvp.Key, key) {
			result.InternalSet(config, kvp.Key, kvp.Value)
		}
	}

	return result, nil
}

// Size returns the number of items in this collection
func (h *InternalHashmap) Size() int {
	return h.size
}

func (h *InternalHashmap) String(config *HashmapConfig) string {
	var buffer bytes.Buffer
	buffer.WriteString("Size: ")
	buffer.WriteString(fmt.Sprintf("%d", h.size))
	buffer.WriteString("\n[\n")
	h.ForEach(config, func(k unsafe.Pointer, v unsafe.Pointer) {
		ks := config.KeyConfig.Format(k)
		vs := config.ValueConfig.Format(v)
		buffer.WriteString(fmt.Sprintf("  %s: %s\n", ks, vs))
	})
	buffer.WriteString("]\n")
	return buffer.String()
}

func (h *InternalHashmap) instantiate(config *HashmapConfig, size int, contents []*KeyValuePair) *BaseStruct {
	hash := CreateEmptyInternalHashmap(config.Options.PackedBucket, size)

	for _, v := range contents {
		if v != nil {
			hash.InternalSet(config, v.Key, v.Value)
		}
	}

	return &BaseStruct{hash}
}

func (h *InternalHashmap) InternalSet(config *HashmapConfig, key unsafe.Pointer, value unsafe.Pointer) {
	hobSize := uint32(64 - h.lobSize)

	hashkey := config.KeyConfig.Hash(key, h.seed)
	selectedBucket := hashkey & h.lobMask
	b := h.buckets[selectedBucket]
	if b == nil {
		b = h.createEmptyBucket(config, hobSize)
		h.buckets[selectedBucket] = b
	}
	for b.entryCount == bucketCapacity {
		if b.overflow == nil {
			b.overflow = h.createEmptyBucket(config, hobSize)
		}
		b = b.overflow
	}
	index := int(b.entryCount)
	config.KeyConfig.Write(b.keys, index, key)
	config.ValueConfig.Write(b.values, index, value)
	b.hobs.Assign(uint64(b.entryCount), hashkey>>h.lobSize)
	b.entryCount++
}

func (h *InternalHashmap) createEmptyBucket(config *HashmapConfig, hobSize uint32) *bucket {
	return &bucket{
		entryCount: 0,
		hobs:       memory.AllocateMemories(h.blockSize, hobSize, bucketCapacity),
		keys:       config.KeyConfig.CreateBucket(bucketCapacity),
		values:     config.ValueConfig.CreateBucket(bucketCapacity),
		overflow:   nil,
	}
}
