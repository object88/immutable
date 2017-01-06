package memory

// Memories is a interface over Memories32 and Memories8
type Memories interface {
	Assign(index uint64, value uint64)
	Read(index uint64) uint64
}

// AllocateMemories creates a block of memory using the given block size
func AllocateMemories(blockSize BlockSize, bits uint32, count uint32) Memories {
	totalBitCount := bits * count
	requiredBlocks := calculateBlocksRequired(totalBitCount, blockSize)
	var m Memories
	switch blockSize {
	case SmallBlock:
		b := make([]smallBlock, requiredBlocks)
		m = &Memories8{bits, b}
	case LargeBlock:
		b := make([]largeBlock, requiredBlocks)
		m = &Memories32{bits, b}
	case NoPacking:
		b := make([]extraLargeBlock, count)
		m = &MemoriesNoPacking{b}
	}
	return m
}

func calculateBlocksRequired(totalBits uint32, blockSize BlockSize) uint32 {
	var size uint32
	switch blockSize {
	case SmallBlock:
		size = bitsInSmallBlock
	case LargeBlock:
		size = bitsInLargeBlock
	case NoPacking:
		size = bitsInExtraLargeBlock
	}
	requiredBlocks := totalBits / size
	if totalBits%size != 0 {
		requiredBlocks++
	}
	return requiredBlocks
}
