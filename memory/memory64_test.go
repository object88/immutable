package memory

import "testing"

func Test_ExtraLarge_Read(t *testing.T) {
	evaluateExtraLargeRead(t, 4, 1, 0, 0xf, memoryMap64{0: 0xf}, nil)
	evaluateExtraLargeRead(t, 4, 16, 0, 0xf, memoryMap64{0: 0xf}, nil)
	evaluateExtraLargeRead(t, 4, 32, 0, 0xf, memoryMap64{0: 0xf}, nil)
	evaluateExtraLargeRead(t, 4, 32, 1, 0xf, memoryMap64{0: 0xf0}, nil)
	evaluateExtraLargeRead(t, 4, 32, 15, 0x000000000000000f, memoryMap64{0: 0xf000000000000000}, nil)
	evaluateExtraLargeRead(t, 4, 32, 16, 0x000000000000000f, memoryMap64{1: 0x000000000000000f}, nil)
	evaluateExtraLargeRead(t, 4, 32, 31, 0x000000000000000f, memoryMap64{1: 0xf000000000000000}, nil)
	evaluateExtraLargeRead(t, 8, 32, 0, 0x00000000000000ff, memoryMap64{0: 0x00000000000000ff}, nil)
	evaluateExtraLargeRead(t, 8, 32, 1, 0x00000000000000ff, memoryMap64{0: 0x000000000000ff00}, nil)
	evaluateExtraLargeRead(t, 31, 1, 0, 0x000000007fffffff, memoryMap64{0: 0x000000007fffffff}, nil)
	evaluateExtraLargeRead(t, 31, 4, 1, 0x000000007fffffff, memoryMap64{0: 0x3fffffff80000000}, nil)
	evaluateExtraLargeRead(t, 31, 4, 2, 0x000000007fffffff, memoryMap64{0: 0xc000000000000000, 1: 0x000000001fffffff}, nil)
	evaluateExtraLargeRead(t, 63, 4, 0, 0x7fffffffffffffff, memoryMap64{0: 0x7fffffffffffffff}, nil)
	evaluateExtraLargeRead(t, 63, 4, 1, 0x7fffffffffffffff, memoryMap64{0: 0x8000000000000000, 1: 0x3fffffffffffffff}, nil)
	evaluateExtraLargeRead(t, 63, 4, 2, 0x7fffffffffffffff, memoryMap64{1: 0xc000000000000000, 2: 0x1fffffffffffffff}, nil)
	// Need more here, with 555..., and aaa...
	// Also need some with bits flipped; writing 0s.
}

func evaluateExtraLargeRead(t *testing.T, bitCount, count uint32, readIndex, expected uint64, assignment memoryMap64, options *readOptions) {
	// fmt.Printf("\nReviewing %d/%d/0x%08x\n", bitCount, count, expected)
	m := AllocateMemories(ExtraLargeBlock, bitCount, count)
	mem := m.(*Memories64).m
	if options != nil && options.initInvert {
		for k := range mem {
			mem[k] = 0xffffffffffffffff
		}
	}
	for k, v := range assignment {
		// fmt.Printf("Writing 0x%08x to %d\n", v, k)
		mem[k] = extraLargeBlock(v)
	}
	result := m.Read(readIndex)
	if result != expected {
		t.Fatalf("Incorrect result from read; expected 0x%016x, got 0x%016x\n", expected, result)
	}
}

func Test_ExtraLarge_Assign(t *testing.T) {
	evaluateExtraLargeAssign(t, 4, 1, 0, 0xf, memoryMap64{0: 0xf}, nil)
	evaluateExtraLargeAssign(t, 4, 16, 0, 0xf, memoryMap64{0: 0xf}, nil)
	evaluateExtraLargeAssign(t, 4, 32, 0, 0xf, memoryMap64{0: 0xf}, nil)
	evaluateExtraLargeAssign(t, 4, 32, 1, 0xf, memoryMap64{0: 0xf0}, nil)
	evaluateExtraLargeAssign(t, 4, 32, 15, 0x000000000000000f, memoryMap64{0: 0xf000000000000000}, nil)
	evaluateExtraLargeAssign(t, 4, 32, 16, 0x000000000000000f, memoryMap64{1: 0x000000000000000f}, nil)
	evaluateExtraLargeAssign(t, 4, 32, 31, 0x000000000000000f, memoryMap64{1: 0xf000000000000000}, nil)
	evaluateExtraLargeAssign(t, 8, 32, 0, 0x00000000000000ff, memoryMap64{0: 0x00000000000000ff}, nil)
	evaluateExtraLargeAssign(t, 8, 32, 1, 0x00000000000000ff, memoryMap64{0: 0x000000000000ff00}, nil)
	evaluateExtraLargeAssign(t, 31, 1, 0, 0x000000007fffffff, memoryMap64{0: 0x000000007fffffff}, nil)
	evaluateExtraLargeAssign(t, 31, 4, 1, 0x000000007fffffff, memoryMap64{0: 0x3fffffff80000000}, nil)
	evaluateExtraLargeAssign(t, 31, 4, 2, 0x000000007fffffff, memoryMap64{0: 0xc000000000000000, 1: 0x000000001fffffff}, nil)
	evaluateExtraLargeAssign(t, 63, 4, 0, 0x7fffffffffffffff, memoryMap64{0: 0x7fffffffffffffff}, nil)
	evaluateExtraLargeAssign(t, 63, 4, 1, 0x7fffffffffffffff, memoryMap64{0: 0x8000000000000000, 1: 0x3fffffffffffffff}, nil)
	evaluateExtraLargeAssign(t, 63, 4, 2, 0x7fffffffffffffff, memoryMap64{1: 0xc000000000000000, 2: 0x1fffffffffffffff}, nil)
}

func evaluateExtraLargeAssign(t *testing.T, bitCount, count uint32, writeIndex, value uint64, assessment memoryMap64, options *assignOptions) {
	m := AllocateMemories(ExtraLargeBlock, bitCount, count)
	mem := m.(*Memories64).m
	if options != nil && options.initInvert {
		for k := range mem {
			mem[k] = 0xffffffffffffffff
		}
	}
	m.Assign(writeIndex, value)
	for k, v := range assessment {
		result := uint64(mem[k])
		if v != result {
			t.Fatalf("At %d, incorrect result from write; expected 0x%016x, got 0x%016x\n", k, v, result)
		}
	}
}
