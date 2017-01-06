package memory

import "testing"

func Test_Large_AllocateMemories(t *testing.T) {
	evaluateLargeAllocate(t, 32, 1, 1)
	evaluateLargeAllocate(t, 32, 2, 2)
	evaluateLargeAllocate(t, 25, 1, 1)
	evaluateLargeAllocate(t, 25, 2, 2)
}

func Test_ExtraLarge_AllocateMemories(t *testing.T) {
	evaluateExtraLargeAllocate(t, 3, 3, 1)
	evaluateExtraLargeAllocate(t, 3, 22, 2)
	evaluateExtraLargeAllocate(t, 31, 1, 1)
	evaluateExtraLargeAllocate(t, 31, 2, 1)
	evaluateExtraLargeAllocate(t, 31, 3, 2)
	evaluateExtraLargeAllocate(t, 32, 1, 1)
	evaluateExtraLargeAllocate(t, 32, 2, 1)
	evaluateExtraLargeAllocate(t, 32, 3, 2)
	evaluateExtraLargeAllocate(t, 33, 1, 1)
	evaluateExtraLargeAllocate(t, 33, 2, 2)
	evaluateExtraLargeAllocate(t, 33, 3, 2)
	evaluateExtraLargeAllocate(t, 63, 1, 1)
	evaluateExtraLargeAllocate(t, 63, 2, 2)
}

func Test_NoPacking_AllocateMemories(t *testing.T) {
	evaluateNoPackingAllocate(t, 1, 1)
	evaluateNoPackingAllocate(t, 2, 2)
	evaluateNoPackingAllocate(t, 4, 4)
}

func evaluateLargeAllocate(t *testing.T, bits, count uint32, expected int) {
	m := AllocateMemories(LargeBlock, bits, count)
	mem := m.(*Memories32).m
	evaluate(t, len(mem), expected)
}

func evaluateExtraLargeAllocate(t *testing.T, bits, count uint32, expected int) {
	m := AllocateMemories(ExtraLargeBlock, bits, count)
	mem := m.(*Memories64).m
	evaluate(t, len(mem), expected)
}

func evaluateNoPackingAllocate(t *testing.T, count uint32, expected int) {
	m := AllocateMemories(NoPacking, 0, count)
	mem := m.(*MemoriesNoPacking).m
	evaluate(t, len(mem), expected)
}

func evaluate(t *testing.T, result, expected int) {
	if result != expected {
		t.Fatalf("Incorrect byte count; expected %d, got %d", expected, result)
	}
}
