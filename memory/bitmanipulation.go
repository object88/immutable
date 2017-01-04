package memory

// http://graphics.stanford.edu/~seander/bithacks.html#RoundUpPowerOf2
func NextPowerOfTwo(size uint32) uint32 {
	size--
	size |= size >> 1
	size |= size >> 2
	size |= size >> 4
	size |= size >> 8
	size |= size >> 16
	size++
	return size
}

// http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetKernighan
func Fffff(value uint32) uint32 {
	v := value - 1

	c := uint32(0) // c accumulates the total bits set in v
	for ; v != 0; c++ {
		v &= v - 1 // clear the least significant bit set
	}
	return c
}
