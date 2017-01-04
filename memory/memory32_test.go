package memory

import (
	"math/rand"
	"testing"
	"time"
)

func Test_AllocateMemories(t *testing.T) {
	m := AllocateMemories(LargeBlock, 32, 1)
	mem := m.(*Memories32).m
	blockCount := len(mem)
	if blockCount != 1 {
		t.Fatalf("Incorrect byte count; expected 1, got %d", blockCount)
	}
}

func Test_AllocateMemories_GreaterThanZeroCount(t *testing.T) {
	m := AllocateMemories(LargeBlock, 32, 2)
	mem := m.(*Memories32).m
	blockCount := len(mem)
	if blockCount != 2 {
		t.Fatal("Incorrect byte count")
	}
}

func Test_AllocateMemories_NonByteAlignedSize(t *testing.T) {
	m := AllocateMemories(LargeBlock, 25, 1)
	mem := m.(*Memories32).m
	blockCount := len(mem)
	if blockCount != 1 {
		t.Fatal("Incorrect byte count")
	}
}

func Test_AllocateMemories_NonByteAlignedSize_GreaterThanZeroCount(t *testing.T) {
	m := AllocateMemories(LargeBlock, 25, 2)
	mem := m.(*Memories32).m
	blockCount := len(mem)
	if blockCount != 2 {
		t.Fatal("Incorrect byte count")
	}
}

type memoryMap map[uint32]uint32

func Test_Read(t *testing.T) {
	evaluateRead(t, 4, 1, 0, 0x0f, memoryMap{0: 0x0f}, nil)
	evaluateRead(t, 4, 8, 1, 0x0f, memoryMap{0: 0xf0}, nil)
	evaluateRead(t, 4, 16, 8, 0x0f, memoryMap{1: 0x0f}, nil)
	evaluateRead(t, 4, 16, 9, 0x0f, memoryMap{1: 0xf0}, nil)
	evaluateRead(t, 6, 2, 1, 0x3f, memoryMap{0: 0xfc0}, nil)
	evaluateRead(t, 11, 2, 1, 0x7ff, memoryMap{0: 0x3ff800}, nil)
	evaluateRead(t, 24, 2, 1, 0xffffff, memoryMap{0: 0xff000000, 1: 0xffff}, nil)
	evaluateRead(t, 32, 1, 0, 0xffffffff, memoryMap{0: 0xffffffff}, nil)
	evaluateRead(t, 32, 1, 0, 0xaaaaaaaa, memoryMap{0: 0xaaaaaaaa}, nil)
	evaluateRead(t, 32, 1, 0, 0x55555555, memoryMap{0: 0x55555555}, nil)

	o := &readOptions{true}
	evaluateRead(t, 4, 1, 0, 0x0, memoryMap{0: 0xfffffff0}, o)
	evaluateRead(t, 4, 8, 1, 0x0, memoryMap{0: 0xffffff0f}, o)
}

type readOptions struct {
	initInvert bool
}

func evaluateRead(t *testing.T, bitCount, count, readIndex, expected uint32, assignment memoryMap, options *readOptions) {
	// fmt.Printf("\nReviewing %d/%d/0x%08x\n", bitCount, count, expected)
	m := AllocateMemories(LargeBlock, bitCount, count)
	mem := m.(*Memories32).m
	if options != nil && options.initInvert {
		for k := range mem {
			mem[k] = 0xffffffff
		}
	}
	for k, v := range assignment {
		// fmt.Printf("Writing 0x%08x to %d\n", v, k)
		mem[k] = largeBlock(v)
	}
	result := m.Read(readIndex)
	if result != expected {
		t.Fatalf("Incorrect result from read; expected 0x%08x, got 0x%08x\n", expected, result)
	}
}

func Test_Assign(t *testing.T) {
	evaluateAssign(t, 4, 1, 0, 0xf, memoryMap{0: 0xf}, nil)
	evaluateAssign(t, 4, 2, 1, 0xf, memoryMap{0: 0xf0}, nil)
	evaluateAssign(t, 4, 3, 2, 0xf, memoryMap{0: 0xf00}, nil)
	evaluateAssign(t, 6, 2, 1, 0x3f, memoryMap{0: 0xfc0}, nil)
	evaluateAssign(t, 11, 2, 1, 0x7ff, memoryMap{0: 0x3ff800}, nil)
	evaluateAssign(t, 24, 2, 1, 0xffffff, memoryMap{0: 0xff000000, 1: 0xffff}, nil)
	evaluateAssign(t, 31, 1, 0, 0x7fffffff, memoryMap{0: 0x7fffffff}, nil)
	evaluateAssign(t, 31, 1, 0, 0x2aaaaaaa, memoryMap{0: 0x2aaaaaaa}, nil)
	evaluateAssign(t, 31, 1, 0, 0x55555555, memoryMap{0: 0x55555555}, nil)
	evaluateAssign(t, 31, 2, 1, 0x7fffffff, memoryMap{0: 0x80000000, 1: 0x3fffffff}, nil)
	evaluateAssign(t, 31, 2, 1, 0x2aaaaaaa, memoryMap{1: 0x15555555}, nil)
	evaluateAssign(t, 31, 2, 1, 0x55555555, memoryMap{0: 0x80000000, 1: 0x2aaaaaaa}, nil)
	evaluateAssign(t, 32, 1, 0, 0xffffffff, memoryMap{0: 0xffffffff}, nil)
	evaluateAssign(t, 32, 1, 0, 0xaaaaaaaa, memoryMap{0: 0xaaaaaaaa}, nil)
	evaluateAssign(t, 32, 1, 0, 0x55555555, memoryMap{0: 0x55555555}, nil)

	o := &assignOptions{true}
	evaluateAssign(t, 4, 1, 0, 0x0, memoryMap{0: 0xfffffff0}, o)
	evaluateAssign(t, 4, 2, 1, 0x0, memoryMap{0: 0xffffff0f}, o)
	evaluateAssign(t, 4, 3, 2, 0x0, memoryMap{0: 0xfffff0ff}, o)
	evaluateAssign(t, 6, 2, 1, 0x0, memoryMap{0: 0xfffff03f}, o)
	evaluateAssign(t, 11, 2, 1, 0x0, memoryMap{0: 0xffc007ff}, o)
	evaluateAssign(t, 24, 2, 1, 0x0, memoryMap{0: 0x00ffffff, 1: 0xffff0000}, o)
	evaluateAssign(t, 31, 1, 0, 0x0, memoryMap{0: 0x80000000}, o)
	evaluateAssign(t, 31, 2, 1, 0x0, memoryMap{0: 0x7fffffff, 1: 0xc0000000}, o)
}

type assignOptions struct {
	initInvert bool
}

func evaluateAssign(t *testing.T, bitCount, count, writeIndex, value uint32, assessment memoryMap, options *assignOptions) {
	m := AllocateMemories(LargeBlock, bitCount, count)
	mem := m.(*Memories32).m
	if options != nil && options.initInvert {
		for k := range mem {
			mem[k] = 0xffffffff
		}
	}
	m.Assign(writeIndex, value)
	for k, v := range assessment {
		result := uint32(mem[k])
		if v != result {
			t.Fatalf("At %d, incorrect result from write; expected 0x%032x, got 0x%032x\n", k, v, result)
		}
	}
}

func Test_WriteAndRead2(t *testing.T) {
	count := 4
	set := make([]uint32, count)
	for i := 0; i < count; i++ {
		set[i] = 0x55555555
	}

	m := AllocateMemories(LargeBlock, 31, 4)

	for k, v := range set {
		m.Assign(uint32(k), v)
	}

	for k, v := range set {
		result := m.Read(uint32(k))
		if result != v {
			t.Fatalf("At %d\nexpected %032b\nreceived %032b\n", k, v, result)
		}
	}
}

func Test_WriteAndRead(t *testing.T) {
	count := 4
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	set := make([]uint32, count)
	for i := 0; i < count; i++ {
		set[i] = uint32(r.Int31()) & 0x7fffffff
	}

	m := AllocateMemories(LargeBlock, 31, 4)

	for k, v := range set {
		m.Assign(uint32(k), v)
	}

	for k, v := range set {
		result := m.Read(uint32(k))
		if result != v {
			t.Fatalf("At %d\nexpected %032b\nreceived %032b\n", k, v, result)
		}
	}
}
