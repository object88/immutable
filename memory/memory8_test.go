package memory

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

type memoryMap8 map[uint32]uint8

func Test_Small_Read(t *testing.T) {
	o := &readOptions{true}

	testCases := []struct {
		bitCount   uint32
		count      uint32
		readIndex  uint64
		expected   uint64
		assignment memoryMap8
		options    *readOptions
	}{
		{4, 1, 0, 0x05, memoryMap8{0: 0x05}, nil},
		{4, 1, 0, 0x0a, memoryMap8{0: 0x0a}, nil},
		{4, 1, 0, 0x0f, memoryMap8{0: 0x0f}, nil},
		{4, 2, 1, 0x0f, memoryMap8{0: 0xf0}, nil},
		{7, 4, 0, 0x7f, memoryMap8{0: 0x7f}, nil},
		{7, 4, 1, 0x7f, memoryMap8{0: 0x80, 1: 0x3f}, nil},
		{7, 4, 2, 0x7f, memoryMap8{1: 0xc0, 2: 0x1f}, nil},
		{7, 4, 3, 0x7f, memoryMap8{2: 0xe0, 3: 0x0f}, nil},
		{11, 2, 0, 0x00000000000007ff, memoryMap8{0: 0xff, 1: 0x07}, nil},
		{11, 2, 1, 0x00000000000007ff, memoryMap8{1: 0xf8, 2: 0x3f}, nil},

		{4, 1, 0, 0x00, memoryMap8{0: 0xf0}, o},
	}

	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			// fmt.Printf("\nReviewing %d/%d/0x%08x\n", bitCount, count, expected)
			m := AllocateMemories(SmallBlock, tc.bitCount, tc.count)
			mem := m.(*Memories8).m
			if tc.options != nil && tc.options.initInvert {
				for k := range mem {
					mem[k] = 0xff
				}
			}
			for k, v := range tc.assignment {
				// fmt.Printf("Writing 0x%08x to %d\n", v, k)
				mem[k] = uint8(v)
			}
			result := m.Read(tc.readIndex)
			if result != tc.expected {
				t.Fatalf("Incorrect result from read; expected 0x%016x, got 0x%016x\n", tc.expected, result)
			}
		})
	}
}

func Test_Small_Assign(t *testing.T) {
	o := &assignOptions{true}

	testCases := []struct {
		bitCount   uint32
		count      uint32
		writeIndex uint64
		value      uint64
		assessment memoryMap8
		options    *assignOptions
	}{
		{4, 1, 0, 0x05, memoryMap8{0: 0x05}, nil},
		{4, 1, 0, 0x0a, memoryMap8{0: 0x0a}, nil},
		{4, 1, 0, 0x0f, memoryMap8{0: 0x0f}, nil},
		{4, 2, 1, 0x0f, memoryMap8{0: 0xf0}, nil},
		{7, 4, 0, 0x7f, memoryMap8{0: 0x7f}, nil},
		{7, 4, 1, 0x7f, memoryMap8{0: 0x80, 1: 0x3f}, nil},
		{7, 4, 2, 0x7f, memoryMap8{1: 0xc0, 2: 0x1f}, nil},
		{7, 4, 3, 0x7f, memoryMap8{2: 0xe0, 3: 0x0f}, nil},
		{11, 2, 0, 0x00000000000007ff, memoryMap8{0: 0xff, 1: 0x07}, nil},
		{11, 2, 1, 0x00000000000007ff, memoryMap8{1: 0xf8, 2: 0x3f}, nil},

		{4, 1, 0, 0x00, memoryMap8{0: 0xf0}, o},
	}

	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			m := AllocateMemories(SmallBlock, tc.bitCount, tc.count)
			mem := m.(*Memories8).m
			if tc.options != nil && tc.options.initInvert {
				for k := range mem {
					mem[k] = 0xff
				}
			}
			m.Assign(tc.writeIndex, tc.value)
			for k, v := range tc.assessment {
				result := uint8(mem[k])
				if v != result {
					t.Fatalf("At %d, incorrect result from write; expected 0x%08x, got 0x%08x\n", k, v, result)
				}
			}
		})
	}
}

func Test_Small_Random(t *testing.T) {
	src := rand.NewSource(time.Now().UnixNano())
	random := rand.New(src)

	width := uint32(64 - 11)
	max := int64(math.Pow(2.0, float64(width)))

	count := 8
	contents := make([]uint64, count)
	for i := 0; i < count; i++ {
		contents[i] = uint64(random.Int63n(max))
	}

	m := AllocateMemories(SmallBlock, width, uint32(count))
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
