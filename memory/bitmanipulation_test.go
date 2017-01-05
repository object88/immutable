package memory

import "testing"

func Test_NextPowerOfTwo(t *testing.T) {
	evaluateNextPowerOfTwo(t, 0, 0)
	evaluateNextPowerOfTwo(t, 1, 1)
	evaluateNextPowerOfTwo(t, 2, 2)
	evaluateNextPowerOfTwo(t, 3, 4)
	evaluateNextPowerOfTwo(t, 4, 4)
	evaluateNextPowerOfTwo(t, 5, 8)
	evaluateNextPowerOfTwo(t, 8, 8)
	evaluateNextPowerOfTwo(t, 9, 16)
	evaluateNextPowerOfTwo(t, 16, 16)
	evaluateNextPowerOfTwo(t, 17, 32)
	evaluateNextPowerOfTwo(t, 32, 32)
	evaluateNextPowerOfTwo(t, 33, 64)
	evaluateNextPowerOfTwo(t, 64, 64)
	evaluateNextPowerOfTwo(t, 65, 128)
	evaluateNextPowerOfTwo(t, 128, 128)
}

func evaluateNextPowerOfTwo(t *testing.T, value, expected int) {
	result := NextPowerOfTwo(value)
	if result != expected {
		t.Fatalf("Testing %d; got %d, expected %d\n", value, result, expected)
	}
}

func Test_PowerOf(t *testing.T) {
	evaluatePowerOf(t, 1, 0)
	evaluatePowerOf(t, 2, 1)
	evaluatePowerOf(t, 4, 2)
	evaluatePowerOf(t, 8, 3)
	evaluatePowerOf(t, 16, 4)
	evaluatePowerOf(t, 32, 5)
	evaluatePowerOf(t, 64, 6)
	evaluatePowerOf(t, 128, 7)
}

func evaluatePowerOf(t *testing.T, value int, expected uint8) {
	result := PowerOf(value)
	if result != expected {
		t.Fatalf("Testing %d; got %d, expected %d\n", value, result, expected)
	}
}
