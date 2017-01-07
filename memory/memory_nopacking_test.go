package memory

import "testing"

func Test_NoPacking_Read(t *testing.T) {
	evaluateNoPackingRead(t, 1, 0, 0x0f, memoryMap64{0: 0x0f}, nil)
	evaluateNoPackingRead(t, 8, 1, 0x0f, memoryMap64{1: 0x0f}, nil)
	evaluateNoPackingRead(t, 16, 8, 0x0f, memoryMap64{8: 0x0f}, nil)
	evaluateNoPackingRead(t, 16, 9, 0x0f, memoryMap64{9: 0x0f}, nil)
	evaluateNoPackingRead(t, 2, 1, 0x3f, memoryMap64{1: 0x3f}, nil)
	evaluateNoPackingRead(t, 2, 1, 0x7ff, memoryMap64{1: 0x7ff}, nil)
	evaluateNoPackingRead(t, 2, 1, 0xffffff, memoryMap64{1: 0xffffff}, nil)
	evaluateNoPackingRead(t, 1, 0, 0xffffffffffffffff, memoryMap64{0: 0xffffffffffffffff}, nil)
	evaluateNoPackingRead(t, 1, 0, 0xaaaaaaaaaaaaaaaa, memoryMap64{0: 0xaaaaaaaaaaaaaaaa}, nil)
	evaluateNoPackingRead(t, 1, 0, 0x5555555555555555, memoryMap64{0: 0x5555555555555555}, nil)

	o := &readOptions{true}
	evaluateNoPackingRead(t, 1, 0, 0x0, memoryMap64{0: 0x00}, o)
	evaluateNoPackingRead(t, 8, 1, 0x0, memoryMap64{1: 0x00}, o)
}

func evaluateNoPackingRead(t *testing.T, count uint32, readIndex, expected uint64, assignment memoryMap64, options *readOptions) {
	m := AllocateMemories(NoPacking, 0, count)
	mem := m.(*MemoriesNoPacking).m
	if options != nil && options.initInvert {
		for k := range mem {
			mem[k] = 0xffffffffffffffff
		}
	}
	for k, v := range assignment {
		mem[k] = uint64(v)
	}
	result := m.Read(readIndex)
	if result != expected {
		t.Fatalf("Incorrect result from read; expected 0x%016x, got 0x%016x\n", expected, result)
	}
}

func Test_NoPacking_Assign(t *testing.T) {
	evaluateNoPackingAssign(t, 1, 0, 0x000000000000000f, memoryMap64{0: 0x000000000000000f}, nil)
	evaluateNoPackingAssign(t, 2, 1, 0x000000000000000f, memoryMap64{1: 0x000000000000000f}, nil)
	evaluateNoPackingAssign(t, 3, 2, 0x000000000000000f, memoryMap64{2: 0x000000000000000f}, nil)
	evaluateNoPackingAssign(t, 2, 1, 0x000000000000003f, memoryMap64{1: 0x000000000000003f}, nil)
	evaluateNoPackingAssign(t, 2, 1, 0x00000000000007ff, memoryMap64{1: 0x00000000000007ff}, nil)
	evaluateNoPackingAssign(t, 2, 1, 0x0000000000ffffff, memoryMap64{1: 0x0000000000ffffff}, nil)
	evaluateNoPackingAssign(t, 1, 0, 0x000000007fffffff, memoryMap64{0: 0x000000007fffffff}, nil)
	evaluateNoPackingAssign(t, 1, 0, 0x000000002aaaaaaa, memoryMap64{0: 0x000000002aaaaaaa}, nil)
	evaluateNoPackingAssign(t, 1, 0, 0x0000000055555555, memoryMap64{0: 0x0000000055555555}, nil)
	evaluateNoPackingAssign(t, 1, 0, 0xaaaaaaaaaaaaaaaa, memoryMap64{0: 0xaaaaaaaaaaaaaaaa}, nil)
	evaluateNoPackingAssign(t, 1, 0, 0x5555555555555555, memoryMap64{0: 0x5555555555555555}, nil)
	evaluateNoPackingAssign(t, 1, 0, 0xffffffffffffffff, memoryMap64{0: 0xffffffffffffffff}, nil)
	evaluateNoPackingAssign(t, 2, 1, 0x7fffffffffffffff, memoryMap64{1: 0x7fffffffffffffff}, nil)
	evaluateNoPackingAssign(t, 2, 1, 0x5555555555555555, memoryMap64{1: 0x5555555555555555}, nil)
	evaluateNoPackingAssign(t, 2, 1, 0x2aaaaaaaaaaaaaaa, memoryMap64{1: 0x2aaaaaaaaaaaaaaa}, nil)
	evaluateNoPackingAssign(t, 2, 1, 0x5555555555555555, memoryMap64{1: 0x5555555555555555}, nil)

	o := &assignOptions{true}
	evaluateNoPackingAssign(t, 1, 0, 0x0, memoryMap64{0: 0x0}, o)
	evaluateNoPackingAssign(t, 2, 1, 0x0, memoryMap64{1: 0x0}, o)
	evaluateNoPackingAssign(t, 3, 2, 0x0, memoryMap64{2: 0x0}, o)
	evaluateNoPackingAssign(t, 2, 1, 0x0, memoryMap64{1: 0x0}, o)
	evaluateNoPackingAssign(t, 2, 1, 0x0, memoryMap64{1: 0x0}, o)
	evaluateNoPackingAssign(t, 2, 1, 0x0, memoryMap64{1: 0x0}, o)
	evaluateNoPackingAssign(t, 1, 0, 0x0, memoryMap64{0: 0x0}, o)
	evaluateNoPackingAssign(t, 2, 1, 0x0, memoryMap64{1: 0x0}, o)
}

func evaluateNoPackingAssign(t *testing.T, count uint32, writeIndex, value uint64, assessment memoryMap64, options *assignOptions) {
	m := AllocateMemories(NoPacking, 0, count)
	mem := m.(*MemoriesNoPacking).m
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
