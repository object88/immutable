package memory

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

type memoryMap32 map[uint32]uint32

func Test_Large_Read(t *testing.T) {
	o := &readOptions{true}

	testCases := []struct {
		bitCount   uint32
		count      uint32
		readIndex  uint64
		expected   uint64
		assignment memoryMap32
		options    *readOptions
	}{
		{4, 1, 0, 0x000000000000000f, memoryMap32{0: 0x0000000f}, nil},
		{4, 8, 1, 0x000000000000000f, memoryMap32{0: 0x000000f0}, nil},
		{4, 16, 8, 0x000000000000000f, memoryMap32{1: 0x0000000f}, nil},
		{4, 16, 9, 0x000000000000000f, memoryMap32{1: 0x000000f0}, nil},
		{6, 2, 1, 0x000000000000003f, memoryMap32{0: 0x00000fc0}, nil},
		{11, 2, 1, 0x00000000000007ff, memoryMap32{0: 0x003ff800}, nil},
		{24, 2, 1, 0x0000000000ffffff, memoryMap32{0: 0xff000000, 1: 0x0000ffff}, nil},
		{31, 4, 1, 0x000000007fffffff, memoryMap32{0: 0x80000000, 1: 0x3fffffff}, nil},
		{31, 4, 1, 0x000000002aaaaaaa, memoryMap32{1: 0x15555555}, nil},
		{31, 4, 1, 0x0000000055555555, memoryMap32{0: 0x80000000, 1: 0x2aaaaaaa}, nil},
		{32, 4, 1, 0x00000000ffffffff, memoryMap32{1: 0xffffffff}, nil},
		{32, 4, 1, 0x00000000aaaaaaaa, memoryMap32{1: 0xaaaaaaaa}, nil},
		{32, 4, 1, 0x0000000055555555, memoryMap32{1: 0x55555555}, nil},
		{63, 4, 1, 0x2aaaaaaaaaaaaaaa, memoryMap32{2: 0x55555555, 3: 0x55555555}, nil},
		{63, 4, 1, 0x5555555555555555, memoryMap32{1: 0x80000000, 2: 0xaaaaaaaa, 3: 0xaaaaaaaa}, nil},
		{63, 4, 1, 0x7fffffffffffffff, memoryMap32{1: 0xc0000000, 2: 0xffffffff, 3: 0x3fffffff}, nil},
		{64, 4, 1, 0xaaaaaaaaaaaaaaaa, memoryMap32{2: 0xaaaaaaaa, 3: 0xaaaaaaaa}, nil},
		{64, 4, 1, 0x5555555555555555, memoryMap32{2: 0x55555555, 3: 0x55555555}, nil},
		{64, 4, 1, 0xffffffffffffffff, memoryMap32{2: 0xffffffff, 3: 0xffffffff}, nil},

		{4, 1, 0, 0x0, memoryMap32{0: 0xfffffff0}, o},
		{4, 8, 1, 0x0, memoryMap32{0: 0xffffff0f}, o},
	}

	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			// fmt.Printf("\nReviewing %d/%d/0x%08x\n", bitCount, count, expected)
			m := AllocateMemories(LargeBlock, tc.bitCount, tc.count)
			mem := m.(*Memories32).m
			if tc.options != nil && tc.options.initInvert {
				for k := range mem {
					mem[k] = 0xffffffff
				}
			}
			for k, v := range tc.assignment {
				// fmt.Printf("Writing 0x%08x to %d\n", v, k)
				mem[k] = uint32(v)
			}
			result := m.Read(tc.readIndex)
			if result != tc.expected {
				t.Fatalf("Incorrect result from read; expected 0x%016x, got 0x%016x\n", tc.expected, result)
			}
		})
	}
}

func Test_Large_Assign(t *testing.T) {
	o := &assignOptions{true}

	testCases := []struct {
		bitCount   uint32
		count      uint32
		writeIndex uint64
		value      uint64
		assessment memoryMap32
		options    *assignOptions
	}{
		{4, 1, 0, 0x000000000000000f, memoryMap32{0: 0x0000000f}, nil},
		{4, 2, 1, 0x000000000000000f, memoryMap32{0: 0x000000f0}, nil},
		{4, 3, 2, 0x000000000000000f, memoryMap32{0: 0x00000f00}, nil},
		{6, 2, 1, 0x000000000000003f, memoryMap32{0: 0x00000fc0}, nil},
		{11, 2, 1, 0x00000000000007ff, memoryMap32{0: 0x003ff800}, nil},
		{24, 2, 1, 0x0000000000ffffff, memoryMap32{0: 0xff000000, 1: 0xffff}, nil},
		{31, 1, 0, 0x000000002aaaaaaa, memoryMap32{0: 0x2aaaaaaa}, nil},
		{31, 1, 0, 0x0000000055555555, memoryMap32{0: 0x55555555}, nil},
		{31, 1, 0, 0x000000007fffffff, memoryMap32{0: 0x7fffffff}, nil},
		{31, 2, 1, 0x000000002aaaaaaa, memoryMap32{1: 0x15555555}, nil},
		{31, 2, 1, 0x0000000055555555, memoryMap32{0: 0x80000000, 1: 0x2aaaaaaa}, nil},
		{31, 2, 1, 0x000000007fffffff, memoryMap32{0: 0x80000000, 1: 0x3fffffff}, nil},
		{32, 1, 0, 0x0000000055555555, memoryMap32{0: 0x55555555}, nil},
		{32, 1, 0, 0x00000000aaaaaaaa, memoryMap32{0: 0xaaaaaaaa}, nil},
		{32, 1, 0, 0x00000000ffffffff, memoryMap32{0: 0xffffffff}, nil},
		{32, 2, 1, 0x0000000055555555, memoryMap32{1: 0x55555555}, nil},
		{32, 2, 1, 0x00000000aaaaaaaa, memoryMap32{1: 0xaaaaaaaa}, nil},
		{32, 2, 1, 0x00000000ffffffff, memoryMap32{1: 0xffffffff}, nil},
		{63, 1, 0, 0x5555555555555555, memoryMap32{0: 0x55555555, 1: 0x55555555}, nil},
		{63, 1, 0, 0x2aaaaaaaaaaaaaaa, memoryMap32{0: 0xaaaaaaaa, 1: 0x2aaaaaaa}, nil},
		{63, 1, 0, 0x7fffffffffffffff, memoryMap32{0: 0xffffffff, 1: 0x7fffffff}, nil},
		{63, 2, 1, 0x5555555555555555, memoryMap32{1: 0x80000000, 2: 0xaaaaaaaa, 3: 0x2aaaaaaa}, nil},
		{63, 2, 1, 0x2aaaaaaaaaaaaaaa, memoryMap32{2: 0x55555555, 3: 0x15555555}, nil},
		{63, 2, 1, 0x7fffffffffffffff, memoryMap32{1: 0x80000000, 2: 0xffffffff, 3: 0x3fffffff}, nil},
		{64, 2, 1, 0x5555555555555555, memoryMap32{2: 0x55555555, 3: 0x55555555}, nil},
		{64, 2, 1, 0xaaaaaaaaaaaaaaaa, memoryMap32{2: 0xaaaaaaaa, 3: 0xaaaaaaaa}, nil},
		{64, 2, 1, 0xffffffffffffffff, memoryMap32{2: 0xffffffff, 3: 0xffffffff}, nil},

		{4, 1, 0, 0x0, memoryMap32{0: 0xfffffff0}, o},
		{4, 2, 1, 0x0, memoryMap32{0: 0xffffff0f}, o},
		{4, 3, 2, 0x0, memoryMap32{0: 0xfffff0ff}, o},
		{6, 2, 1, 0x0, memoryMap32{0: 0xfffff03f}, o},
		{11, 2, 1, 0x0, memoryMap32{0: 0xffc007ff}, o},
		{24, 2, 1, 0x0, memoryMap32{0: 0x00ffffff, 1: 0xffff0000}, o},
		{31, 1, 0, 0x0, memoryMap32{0: 0x80000000}, o},
		{31, 2, 1, 0x0, memoryMap32{0: 0x7fffffff, 1: 0xc0000000}, o},
	}

	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			m := AllocateMemories(LargeBlock, tc.bitCount, tc.count)
			mem := m.(*Memories32).m
			if tc.options != nil && tc.options.initInvert {
				for k := range mem {
					mem[k] = 0xffffffff
				}
			}
			m.Assign(tc.writeIndex, tc.value)
			for k, v := range tc.assessment {
				result := uint32(mem[k])
				if v != result {
					t.Fatalf("At %d, incorrect result from write; expected 0x%08x, got 0x%08x\n", k, v, result)
				}
			}
		})
	}
}

func Test_Large_Random(t *testing.T) {
	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)

	width := uint32(64 - 11)
	max := int64(math.Pow(2.0, float64(width)))

	count := 8
	contents := make([]uint64, count)
	for i := 0; i < count; i++ {
		contents[i] = uint64(random.Int63n(max))
	}

	m := AllocateMemories(LargeBlock, width, uint32(count))
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
