package memory

import "testing"

type memoryMap64 map[uint64]uint64

func Test_NoPadding_Read(t *testing.T) {
	evaluateNoPaddingRead(t, 1, 0, 0x0f, memoryMap64{0: 0x0f}, nil)
	evaluateNoPaddingRead(t, 8, 1, 0x0f, memoryMap64{1: 0x0f}, nil)
	evaluateNoPaddingRead(t, 16, 8, 0x0f, memoryMap64{8: 0x0f}, nil)
	evaluateNoPaddingRead(t, 16, 9, 0x0f, memoryMap64{9: 0x0f}, nil)
	evaluateNoPaddingRead(t, 2, 1, 0x3f, memoryMap64{1: 0x3f}, nil)
	evaluateNoPaddingRead(t, 2, 1, 0x7ff, memoryMap64{1: 0x7ff}, nil)
	evaluateNoPaddingRead(t, 2, 1, 0xffffff, memoryMap64{1: 0xffffff}, nil)
	evaluateNoPaddingRead(t, 1, 0, 0xffffffffffffffff, memoryMap64{0: 0xffffffffffffffff}, nil)
	evaluateNoPaddingRead(t, 1, 0, 0xaaaaaaaaaaaaaaaa, memoryMap64{0: 0xaaaaaaaaaaaaaaaa}, nil)
	evaluateNoPaddingRead(t, 1, 0, 0x5555555555555555, memoryMap64{0: 0x5555555555555555}, nil)

	o := &readOptions{true}
	evaluateNoPaddingRead(t, 1, 0, 0x0, memoryMap64{0: 0x00}, o)
	evaluateNoPaddingRead(t, 8, 1, 0x0, memoryMap64{1: 0x00}, o)
}

func evaluateNoPaddingRead(t *testing.T, count uint32, readIndex, expected uint64, assignment memoryMap64, options *readOptions) {
	m := AllocateMemories(NoPadding, 0, count)
	mem := m.(*MemoriesNoPadding).m
	if options != nil && options.initInvert {
		for k := range mem {
			mem[k] = 0xffffffffffffffff
		}
	}
	for k, v := range assignment {
		mem[k] = extraLargeBlock(v)
	}
	result := m.Read(readIndex)
	if result != expected {
		t.Fatalf("Incorrect result from read; expected 0x%016x, got 0x%016x\n", expected, result)
	}
}

func Test_NoPadding_Assign(t *testing.T) {
	evaluateNoPaddingAssign(t, 1, 0, 0x000000000000000f, memoryMap64{0: 0x000000000000000f}, nil)
	evaluateNoPaddingAssign(t, 2, 1, 0x000000000000000f, memoryMap64{1: 0x000000000000000f}, nil)
	evaluateNoPaddingAssign(t, 3, 2, 0x000000000000000f, memoryMap64{2: 0x000000000000000f}, nil)
	evaluateNoPaddingAssign(t, 2, 1, 0x000000000000003f, memoryMap64{1: 0x000000000000003f}, nil)
	evaluateNoPaddingAssign(t, 2, 1, 0x00000000000007ff, memoryMap64{1: 0x00000000000007ff}, nil)
	evaluateNoPaddingAssign(t, 2, 1, 0x0000000000ffffff, memoryMap64{1: 0x0000000000ffffff}, nil)
	evaluateNoPaddingAssign(t, 1, 0, 0x000000007fffffff, memoryMap64{0: 0x000000007fffffff}, nil)
	evaluateNoPaddingAssign(t, 1, 0, 0x000000002aaaaaaa, memoryMap64{0: 0x000000002aaaaaaa}, nil)
	evaluateNoPaddingAssign(t, 1, 0, 0x0000000055555555, memoryMap64{0: 0x0000000055555555}, nil)
	evaluateNoPaddingAssign(t, 1, 0, 0xaaaaaaaaaaaaaaaa, memoryMap64{0: 0xaaaaaaaaaaaaaaaa}, nil)
	evaluateNoPaddingAssign(t, 1, 0, 0x5555555555555555, memoryMap64{0: 0x5555555555555555}, nil)
	evaluateNoPaddingAssign(t, 1, 0, 0xffffffffffffffff, memoryMap64{0: 0xffffffffffffffff}, nil)
	evaluateNoPaddingAssign(t, 2, 1, 0x7fffffffffffffff, memoryMap64{1: 0x7fffffffffffffff}, nil)
	evaluateNoPaddingAssign(t, 2, 1, 0x5555555555555555, memoryMap64{1: 0x5555555555555555}, nil)
	evaluateNoPaddingAssign(t, 2, 1, 0x2aaaaaaaaaaaaaaa, memoryMap64{1: 0x2aaaaaaaaaaaaaaa}, nil)
	evaluateNoPaddingAssign(t, 2, 1, 0x5555555555555555, memoryMap64{1: 0x5555555555555555}, nil)

	o := &assignOptions{true}
	evaluateNoPaddingAssign(t, 1, 0, 0x0, memoryMap64{0: 0x0}, o)
	evaluateNoPaddingAssign(t, 2, 1, 0x0, memoryMap64{1: 0x0}, o)
	evaluateNoPaddingAssign(t, 3, 2, 0x0, memoryMap64{2: 0x0}, o)
	evaluateNoPaddingAssign(t, 2, 1, 0x0, memoryMap64{1: 0x0}, o)
	evaluateNoPaddingAssign(t, 2, 1, 0x0, memoryMap64{1: 0x0}, o)
	evaluateNoPaddingAssign(t, 2, 1, 0x0, memoryMap64{1: 0x0}, o)
	evaluateNoPaddingAssign(t, 1, 0, 0x0, memoryMap64{0: 0x0}, o)
	evaluateNoPaddingAssign(t, 2, 1, 0x0, memoryMap64{1: 0x0}, o)
}

func evaluateNoPaddingAssign(t *testing.T, count uint32, writeIndex, value uint64, assessment memoryMap64, options *assignOptions) {
	m := AllocateMemories(NoPadding, 0, count)
	mem := m.(*MemoriesNoPadding).m
	if options != nil && options.initInvert {
		for k := range mem {
			mem[k] = 0xffffffffffffffff
		}
	}
	m.Assign(writeIndex, value)
	for k, v := range assessment {
		result := uint64(mem[k])
		if v != result {
			t.Fatalf("At %d, incorrect result from write; expected 0x%064x, got 0x%064x\n", k, v, result)
		}
	}
}
