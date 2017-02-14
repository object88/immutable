package memory

import "fmt"

// Memories is a interface over Memories32 and Memories8
type Memories interface {
	Assign(index uint64, value uint64)
	Read(index uint64) uint64
}

// SelectBlockSize returns the BlockSize necesary for the number of bits to
// store per value
func SelectBlockSize(packed bool, bitsRequired uint8) (result BlockSize) {
	if packed {
		if bitsRequired <= 8 {
			return SmallBlock
		} else if bitsRequired <= 16 {
			return MediumBlock
		} else if bitsRequired <= 32 {
			return LargeBlock
		} else if bitsRequired <= 64 {
			return ExtraLargeBlock
		}
	} else {
		if bitsRequired <= 8 {
			return SmallBlockNoPacking
		} else if bitsRequired <= 16 {
			return MediumBlockNoPacking
		} else if bitsRequired <= 32 {
			return LargeBlockNoPacking
		} else if bitsRequired <= 64 {
			return ExtraLargeBlockNoPacking
		}
	}
	panic(fmt.Sprintf("Cannot create a memory block with %d bits", bitsRequired))
}

// AllocateMemories creates a block of memory using the given block size
func AllocateMemories(blockSize BlockSize, bits uint32, count uint32) Memories {
	totalBitCount := bits * count
	requiredBlocks := calculateBlocksRequired(totalBitCount, blockSize)
	var m Memories
	switch blockSize {
	case SmallBlock:
		b := make([]uint8, requiredBlocks)
		m = &Memories8{bits, b}
	case SmallBlockNoPacking:
		b := make([]uint8, count)
		m = &Memories8NoPacking{b}
	// case MediumBlock:
	// 	b := make([]uint16, requiredBlocks)
	// 	m = &Memories16{bits, b}
	case MediumBlockNoPacking:
		b := make([]uint16, count)
		m = &Memories16NoPacking{b}
	case LargeBlock:
		b := make([]uint32, requiredBlocks)
		m = &Memories32{bits, b}
	case LargeBlockNoPacking:
		b := make([]uint32, count)
		m = &Memories32NoPacking{b}
	case ExtraLargeBlock:
		b := make([]uint64, requiredBlocks)
		m = &Memories64{bits, b}
	case ExtraLargeBlockNoPacking:
		b := make([]uint64, count)
		m = &Memories64NoPacking{b}
	}
	return m
}

func calculateBlocksRequired(totalBits uint32, blockSize BlockSize) uint32 {
	var size uint32
	switch blockSize {
	case SmallBlock, SmallBlockNoPacking:
		size = bitsInSmallBlock
	case MediumBlock, MediumBlockNoPacking:
		size = bitsInMediumBlock
	case LargeBlock, LargeBlockNoPacking:
		size = bitsInLargeBlock
	case ExtraLargeBlock, ExtraLargeBlockNoPacking:
		size = bitsInExtraLargeBlock
	}
	requiredBlocks := totalBits / size
	if totalBits%size != 0 {
		requiredBlocks++
	}
	return requiredBlocks
}
