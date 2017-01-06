package memory

import "testing"

func Test_Large_AllocateMemories(t *testing.T) {
	evaluateLargeAllocate(t, 32, 1, 1)
	evaluateLargeAllocate(t, 32, 2, 2)
	evaluateLargeAllocate(t, 25, 1, 1)
	evaluateLargeAllocate(t, 25, 2, 2)
}

func Test_NoPadding_AllocateMemories(t *testing.T) {
	evaluateNoPaddingAllocate(t, 1, 1)
	evaluateNoPaddingAllocate(t, 2, 2)
	evaluateNoPaddingAllocate(t, 4, 4)
}

func evaluateLargeAllocate(t *testing.T, bits, count uint32, expected int) {
	m := AllocateMemories(LargeBlock, bits, count)
	mem := m.(*Memories32).m
	evaluate(t, len(mem), expected)
}

func evaluateNoPaddingAllocate(t *testing.T, count uint32, expected int) {
	m := AllocateMemories(NoPadding, 0, count)
	mem := m.(*MemoriesNoPadding).m
	evaluate(t, len(mem), expected)
}

func evaluate(t *testing.T, result, expected int) {
	if result != expected {
		t.Fatalf("Incorrect byte count; expected %d, got %d", expected, result)
	}
}
