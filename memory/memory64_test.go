package memory

import "testing"

func Test_ExtraLarge_Read(t *testing.T) {
	evaluateExtraLargeRead(t, 4, 1, 0, 0xf, memoryMap64{0: 0xf}, nil)
	evaluateExtraLargeRead(t, 4, 16, 0, 0xf, memoryMap64{0: 0xf}, nil)
	evaluateExtraLargeRead(t, 4, 32, 0, 0xf, memoryMap64{0: 0xf}, nil)
	evaluateExtraLargeRead(t, 4, 32, 1, 0xf, memoryMap64{0: 0xf0}, nil)
	evaluateExtraLargeRead(t, 4, 32, 15, 0xf, memoryMap64{0: 0xf000000000000000}, nil)
	evaluateExtraLargeRead(t, 4, 32, 16, 0xf, memoryMap64{1: 0x000000000000000f}, nil)
	evaluateExtraLargeRead(t, 4, 32, 31, 0xf, memoryMap64{1: 0xf000000000000000}, nil)
	evaluateExtraLargeRead(t, 8, 32, 0, 0xff, memoryMap64{0: 0x00000000000000ff}, nil)
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
