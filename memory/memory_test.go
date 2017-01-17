package memory

import (
	"fmt"
	"testing"
)

func Test_Small_AllocateMemories(t *testing.T) {
	testCases := []struct {
		bits     uint32
		count    uint32
		expected int
	}{
		{1, 8, 1},
		{8, 1, 1},
		{5, 2, 2},
		{7, 7, 7},
		{7, 8, 7},
		{7, 9, 8},
	}
	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			m := AllocateMemories(SmallBlock, tc.bits, tc.count)
			mem := m.(*Memories8).m
			evaluate(t, len(mem), tc.expected)
		})
	}
}

func Test_Large_AllocateMemories(t *testing.T) {
	testCases := []struct {
		bits     uint32
		count    uint32
		expected int
	}{
		{32, 1, 1},
		{32, 2, 2},
		{25, 1, 1},
		{25, 2, 2},
	}
	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			m := AllocateMemories(LargeBlock, tc.bits, tc.count)
			mem := m.(*Memories32).m
			evaluate(t, len(mem), tc.expected)
		})
	}
}

func Test_ExtraLarge_AllocateMemories(t *testing.T) {
	testCases := []struct {
		bits     uint32
		count    uint32
		expected int
	}{
		{3, 3, 1},
		{3, 22, 2},
		{31, 1, 1},
		{31, 2, 1},
		{31, 3, 2},
		{32, 1, 1},
		{32, 2, 1},
		{32, 3, 2},
		{33, 1, 1},
		{33, 2, 2},
		{33, 3, 2},
		{63, 1, 1},
		{63, 2, 2},
	}
	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			m := AllocateMemories(ExtraLargeBlock, tc.bits, tc.count)
			mem := m.(*Memories64).m
			evaluate(t, len(mem), tc.expected)
		})
	}
}

func Test_NoPacking_AllocateMemories(t *testing.T) {
	testCases := []struct {
		count    uint32
		expected int
	}{
		{1, 1},
		{2, 2},
		{4, 4},
	}
	for index, tc := range testCases {
		t.Run(fmt.Sprintf("%d", index), func(t *testing.T) {
			m := AllocateMemories(NoPacking, 0, tc.count)
			mem := m.(*MemoriesNoPacking).m
			evaluate(t, len(mem), tc.expected)
		})
	}
}

func evaluate(t *testing.T, result, expected int) {
	if result != expected {
		t.Fatalf("Incorrect byte count; expected %d, got %d", expected, result)
	}
}
