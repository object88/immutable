package memory

// NextPowerOfTwo takes a positive integer and finds the next greater number
// that's a power of two.  Code directly copied from
// http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
func NextPowerOfTwo(size uint) uint {
	size--
	size |= size >> 1
	size |= size >> 2
	size |= size >> 4
	size |= size >> 8
	size |= size >> 16
	size |= size >> 32
	size++
	return size
}

// PowerOf takes a power-of-two, and returns the power.  For example, if
// 8 is passed in, 3 it returned, because 2^3 is 8.  Behavior is undefined if
// value is not a positive power-of-two.  Code largely inspired by
// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetKernighan
func PowerOf(value uint) uint {
	v := value - 1

	c := uint(0) // c accumulates the total bits set in v
	for ; v != 0; c++ {
		v &= v - 1 // clear the least significant bit set
	}
	return c
}
