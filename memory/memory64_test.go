package memory

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func Test_ExtraLarge_Read(t *testing.T) {
	// o := &readOptions{true}

	testCases := []struct {
		bitCount   uint32
		count      uint32
		readIndex  uint64
		expected   uint64
		assignment memoryMap64
		options    *readOptions
	}{
		{4, 1, 0, 0xf, memoryMap64{0: 0xf}, nil},
		{4, 16, 0, 0xf, memoryMap64{0: 0xf}, nil},
		{4, 32, 0, 0xf, memoryMap64{0: 0xf}, nil},
		{4, 32, 1, 0xf, memoryMap64{0: 0xf0}, nil},
		{4, 32, 15, 0x000000000000000f, memoryMap64{0: 0xf000000000000000}, nil},
		{4, 32, 16, 0x000000000000000f, memoryMap64{1: 0x000000000000000f}, nil},
		{4, 32, 31, 0x000000000000000f, memoryMap64{1: 0xf000000000000000}, nil},
		{8, 32, 0, 0x00000000000000ff, memoryMap64{0: 0x00000000000000ff}, nil},
		{8, 32, 1, 0x00000000000000ff, memoryMap64{0: 0x000000000000ff00}, nil},
		{31, 1, 0, 0x000000007fffffff, memoryMap64{0: 0x000000007fffffff}, nil},
		{31, 4, 1, 0x000000007fffffff, memoryMap64{0: 0x3fffffff80000000}, nil},
		{31, 4, 2, 0x000000007fffffff, memoryMap64{0: 0xc000000000000000, 1: 0x000000001fffffff}, nil},
		{63, 4, 0, 0x7fffffffffffffff, memoryMap64{0: 0x7fffffffffffffff}, nil},
		{63, 4, 1, 0x7fffffffffffffff, memoryMap64{0: 0x8000000000000000, 1: 0x3fffffffffffffff}, nil},
		{63, 4, 2, 0x7fffffffffffffff, memoryMap64{1: 0xc000000000000000, 2: 0x1fffffffffffffff}, nil},
		// Need more here, with 555..., and aaa...
		// Also need some with bits flipped; writing 0s.
	}

	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			// fmt.Printf("\nReviewing %d/%d/0x%08x\n", bitCount, count, expected)
			m := AllocateMemories(ExtraLargeBlock, tc.bitCount, tc.count)
			mem := m.(*Memories64).m
			if tc.options != nil && tc.options.initInvert {
				for k := range mem {
					mem[k] = 0xffffffffffffffff
				}
			}
			for k, v := range tc.assignment {
				// fmt.Printf("Writing 0x%08x to %d\n", v, k)
				mem[k] = uint64(v)
			}
			result := m.Read(tc.readIndex)
			if result != tc.expected {
				t.Fatalf("Incorrect result from read; expected 0x%016x, got 0x%016x\n", tc.expected, result)
			}
		})
	}
}

func Test_ExtraLarge_Assign(t *testing.T) {
	testCases := []struct {
		bitCount   uint32
		count      uint32
		writeIndex uint64
		expected   uint64
		assignment memoryMap64
		options    *readOptions
	}{
		{4, 1, 0, 0xf, memoryMap64{0: 0xf}, nil},
		{4, 16, 0, 0xf, memoryMap64{0: 0xf}, nil},
		{4, 32, 0, 0xf, memoryMap64{0: 0xf}, nil},
		{4, 32, 1, 0xf, memoryMap64{0: 0xf0}, nil},
		{4, 32, 15, 0x000000000000000f, memoryMap64{0: 0xf000000000000000}, nil},
		{4, 32, 16, 0x000000000000000f, memoryMap64{1: 0x000000000000000f}, nil},
		{4, 32, 31, 0x000000000000000f, memoryMap64{1: 0xf000000000000000}, nil},
		{8, 32, 0, 0x00000000000000ff, memoryMap64{0: 0x00000000000000ff}, nil},
		{8, 32, 1, 0x00000000000000ff, memoryMap64{0: 0x000000000000ff00}, nil},
		{31, 1, 0, 0x000000007fffffff, memoryMap64{0: 0x000000007fffffff}, nil},
		{31, 4, 1, 0x000000007fffffff, memoryMap64{0: 0x3fffffff80000000}, nil},
		{31, 4, 2, 0x000000007fffffff, memoryMap64{0: 0xc000000000000000, 1: 0x000000001fffffff}, nil},
		{63, 4, 0, 0x7fffffffffffffff, memoryMap64{0: 0x7fffffffffffffff}, nil},
		{63, 4, 1, 0x7fffffffffffffff, memoryMap64{0: 0x8000000000000000, 1: 0x3fffffffffffffff}, nil},
		{63, 4, 2, 0x7fffffffffffffff, memoryMap64{1: 0xc000000000000000, 2: 0x1fffffffffffffff}, nil},
	}

	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			m := AllocateMemories(ExtraLargeBlock, tc.bitCount, tc.count)
			mem := m.(*Memories64).m
			if tc.options != nil && tc.options.initInvert {
				for k := range mem {
					mem[k] = 0xffffffffffffffff
				}
			}
			m.Assign(tc.writeIndex, tc.expected)
			for k, v := range tc.assignment {
				result := uint64(mem[k])
				if v != result {
					t.Fatalf("At %d, incorrect result from write; expected 0x%016x, got 0x%016x\n", k, v, result)
				}
			}
		})
	}
}

func Test_ExtraLarge_Random(t *testing.T) {
	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)

	width := uint32(64 - 11)
	max := int64(math.Pow(2.0, float64(width)))

	count := uint64(8)
	contents := make([]uint64, count)
	for i := uint64(0); i < count; i++ {
		contents[i] = uint64(random.Int63n(max))
	}

	m := AllocateMemories(ExtraLargeBlock, width, 8)
	for k, v := range contents {
		m.Assign(uint64(k), v)
	}

	for k, v := range contents {
		result := m.Read(uint64(k))
		if result != v {
			t.Fatalf("At %d\nexpected %064b\nreceived %064b\n", k, v, result)
		}
	}
}
