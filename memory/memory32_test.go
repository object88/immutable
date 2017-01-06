package memory

import (
	"math/rand"
	"testing"
	"time"
)

type memoryMap32 map[uint32]uint32

func Test_Large_Read(t *testing.T) {
	evaluateLargeRead(t, 4, 1, 0, 0x000000000000000f, memoryMap32{0: 0x0000000f}, nil)
	evaluateLargeRead(t, 4, 8, 1, 0x000000000000000f, memoryMap32{0: 0x000000f0}, nil)
	evaluateLargeRead(t, 4, 16, 8, 0x000000000000000f, memoryMap32{1: 0x0000000f}, nil)
	evaluateLargeRead(t, 4, 16, 9, 0x000000000000000f, memoryMap32{1: 0x000000f0}, nil)
	evaluateLargeRead(t, 6, 2, 1, 0x000000000000003f, memoryMap32{0: 0x00000fc0}, nil)
	evaluateLargeRead(t, 11, 2, 1, 0x00000000000007ff, memoryMap32{0: 0x003ff800}, nil)
	evaluateLargeRead(t, 24, 2, 1, 0x0000000000ffffff, memoryMap32{0: 0xff000000, 1: 0x0000ffff}, nil)
	evaluateLargeRead(t, 31, 4, 1, 0x000000007fffffff, memoryMap32{0: 0x80000000, 1: 0x3fffffff}, nil)
	evaluateLargeRead(t, 31, 4, 1, 0x000000002aaaaaaa, memoryMap32{1: 0x15555555}, nil)
	evaluateLargeRead(t, 31, 4, 1, 0x0000000055555555, memoryMap32{0: 0x80000000, 1: 0x2aaaaaaa}, nil)
	evaluateLargeRead(t, 32, 4, 1, 0x00000000ffffffff, memoryMap32{1: 0xffffffff}, nil)
	evaluateLargeRead(t, 32, 4, 1, 0x00000000aaaaaaaa, memoryMap32{1: 0xaaaaaaaa}, nil)
	evaluateLargeRead(t, 32, 4, 1, 0x0000000055555555, memoryMap32{1: 0x55555555}, nil)
	evaluateLargeRead(t, 63, 4, 1, 0x2aaaaaaaaaaaaaaa, memoryMap32{2: 0x55555555, 3: 0x55555555}, nil)
	evaluateLargeRead(t, 63, 4, 1, 0x5555555555555555, memoryMap32{1: 0x80000000, 2: 0xaaaaaaaa, 3: 0xaaaaaaaa}, nil)
	evaluateLargeRead(t, 63, 4, 1, 0x7fffffffffffffff, memoryMap32{1: 0xc0000000, 2: 0xffffffff, 3: 0x3fffffff}, nil)
	evaluateLargeRead(t, 64, 4, 1, 0xaaaaaaaaaaaaaaaa, memoryMap32{2: 0xaaaaaaaa, 3: 0xaaaaaaaa}, nil)
	evaluateLargeRead(t, 64, 4, 1, 0x5555555555555555, memoryMap32{2: 0x55555555, 3: 0x55555555}, nil)
	evaluateLargeRead(t, 64, 4, 1, 0xffffffffffffffff, memoryMap32{2: 0xffffffff, 3: 0xffffffff}, nil)

	o := &readOptions{true}
	evaluateLargeRead(t, 4, 1, 0, 0x0, memoryMap32{0: 0xfffffff0}, o)
	evaluateLargeRead(t, 4, 8, 1, 0x0, memoryMap32{0: 0xffffff0f}, o)
}

func evaluateLargeRead(t *testing.T, bitCount, count uint32, readIndex, expected uint64, assignment memoryMap32, options *readOptions) {
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
		t.Fatalf("Incorrect result from read; expected 0x%016x, got 0x%016x\n", expected, result)
	}
}

func Test_Large_Assign(t *testing.T) {
	evaluateLargeAssign(t, 4, 1, 0, 0x000000000000000f, memoryMap32{0: 0x0000000f}, nil)
	evaluateLargeAssign(t, 4, 2, 1, 0x000000000000000f, memoryMap32{0: 0x000000f0}, nil)
	evaluateLargeAssign(t, 4, 3, 2, 0x000000000000000f, memoryMap32{0: 0x00000f00}, nil)
	evaluateLargeAssign(t, 6, 2, 1, 0x000000000000003f, memoryMap32{0: 0x00000fc0}, nil)
	evaluateLargeAssign(t, 11, 2, 1, 0x00000000000007ff, memoryMap32{0: 0x003ff800}, nil)
	evaluateLargeAssign(t, 24, 2, 1, 0x0000000000ffffff, memoryMap32{0: 0xff000000, 1: 0xffff}, nil)
	evaluateLargeAssign(t, 31, 1, 0, 0x000000002aaaaaaa, memoryMap32{0: 0x2aaaaaaa}, nil)
	evaluateLargeAssign(t, 31, 1, 0, 0x0000000055555555, memoryMap32{0: 0x55555555}, nil)
	evaluateLargeAssign(t, 31, 1, 0, 0x000000007fffffff, memoryMap32{0: 0x7fffffff}, nil)
	evaluateLargeAssign(t, 31, 2, 1, 0x000000002aaaaaaa, memoryMap32{1: 0x15555555}, nil)
	evaluateLargeAssign(t, 31, 2, 1, 0x0000000055555555, memoryMap32{0: 0x80000000, 1: 0x2aaaaaaa}, nil)
	evaluateLargeAssign(t, 31, 2, 1, 0x000000007fffffff, memoryMap32{0: 0x80000000, 1: 0x3fffffff}, nil)
	evaluateLargeAssign(t, 32, 1, 0, 0x0000000055555555, memoryMap32{0: 0x55555555}, nil)
	evaluateLargeAssign(t, 32, 1, 0, 0x00000000aaaaaaaa, memoryMap32{0: 0xaaaaaaaa}, nil)
	evaluateLargeAssign(t, 32, 1, 0, 0x00000000ffffffff, memoryMap32{0: 0xffffffff}, nil)
	evaluateLargeAssign(t, 32, 2, 1, 0x0000000055555555, memoryMap32{1: 0x55555555}, nil)
	evaluateLargeAssign(t, 32, 2, 1, 0x00000000aaaaaaaa, memoryMap32{1: 0xaaaaaaaa}, nil)
	evaluateLargeAssign(t, 32, 2, 1, 0x00000000ffffffff, memoryMap32{1: 0xffffffff}, nil)
	evaluateLargeAssign(t, 63, 1, 0, 0x5555555555555555, memoryMap32{0: 0x55555555, 1: 0x55555555}, nil)
	evaluateLargeAssign(t, 63, 1, 0, 0x2aaaaaaaaaaaaaaa, memoryMap32{0: 0xaaaaaaaa, 1: 0x2aaaaaaa}, nil)
	evaluateLargeAssign(t, 63, 1, 0, 0x7fffffffffffffff, memoryMap32{0: 0xffffffff, 1: 0x7fffffff}, nil)
	evaluateLargeAssign(t, 63, 2, 1, 0x5555555555555555, memoryMap32{1: 0x80000000, 2: 0xaaaaaaaa, 3: 0x2aaaaaaa}, nil)
	evaluateLargeAssign(t, 63, 2, 1, 0x2aaaaaaaaaaaaaaa, memoryMap32{2: 0x55555555, 3: 0x15555555}, nil)
	evaluateLargeAssign(t, 63, 2, 1, 0x7fffffffffffffff, memoryMap32{1: 0x80000000, 2: 0xffffffff, 3: 0x3fffffff}, nil)
	evaluateLargeAssign(t, 64, 2, 1, 0x5555555555555555, memoryMap32{2: 0x55555555, 3: 0x55555555}, nil)
	evaluateLargeAssign(t, 64, 2, 1, 0xaaaaaaaaaaaaaaaa, memoryMap32{2: 0xaaaaaaaa, 3: 0xaaaaaaaa}, nil)
	evaluateLargeAssign(t, 64, 2, 1, 0xffffffffffffffff, memoryMap32{2: 0xffffffff, 3: 0xffffffff}, nil)

	o := &assignOptions{true}
	evaluateLargeAssign(t, 4, 1, 0, 0x0, memoryMap32{0: 0xfffffff0}, o)
	evaluateLargeAssign(t, 4, 2, 1, 0x0, memoryMap32{0: 0xffffff0f}, o)
	evaluateLargeAssign(t, 4, 3, 2, 0x0, memoryMap32{0: 0xfffff0ff}, o)
	evaluateLargeAssign(t, 6, 2, 1, 0x0, memoryMap32{0: 0xfffff03f}, o)
	evaluateLargeAssign(t, 11, 2, 1, 0x0, memoryMap32{0: 0xffc007ff}, o)
	evaluateLargeAssign(t, 24, 2, 1, 0x0, memoryMap32{0: 0x00ffffff, 1: 0xffff0000}, o)
	evaluateLargeAssign(t, 31, 1, 0, 0x0, memoryMap32{0: 0x80000000}, o)
	evaluateLargeAssign(t, 31, 2, 1, 0x0, memoryMap32{0: 0x7fffffff, 1: 0xc0000000}, o)
}

func evaluateLargeAssign(t *testing.T, bitCount, count uint32, writeIndex, value uint64, assessment memoryMap32, options *assignOptions) {
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
			t.Fatalf("At %d, incorrect result from write; expected 0x%016x, got 0x%016x\n", k, v, result)
		}
	}
}

func Test_WriteAndRead2(t *testing.T) {
	count := 4
	set := make([]uint64, count)
	for i := 0; i < count; i++ {
		set[i] = 0x55555555
	}

	m := AllocateMemories(LargeBlock, 31, 4)

	for k, v := range set {
		m.Assign(uint64(k), v)
	}

	for k, v := range set {
		result := m.Read(uint64(k))
		if result != v {
			t.Fatalf("At %d\nexpected %032b\nreceived %032b\n", k, v, result)
		}
	}
}

func Test_WriteAndRead(t *testing.T) {
	count := 4
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	set := make([]uint64, count)
	for i := 0; i < count; i++ {
		set[i] = uint64(r.Int31()) & 0x7fffffff
	}

	m := AllocateMemories(LargeBlock, 31, 4)

	for k, v := range set {
		m.Assign(uint64(k), v)
	}

	for k, v := range set {
		result := m.Read(uint64(k))
		if result != v {
			t.Fatalf("At %d\nexpected %032b\nreceived %032b\n", k, v, result)
		}
	}
}
