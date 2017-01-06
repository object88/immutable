package memory

import "fmt"

const fullExtraLargeBlock = ^uint64(0)

// Memories64 is all your memories.
type Memories64 struct {
	bitsPerEntry uint32
	m            []extraLargeBlock
}

// Assign sets a value to the internal memory at the given index
func (m *Memories64) Assign(index uint64, value uint64) {
	bitsRemaining := uint64(m.bitsPerEntry)
	offset := bitsRemaining * index
	byteOffset := offset / bitsInExtraLargeBlock
	bitOffset := offset % bitsInExtraLargeBlock

	// fmt.Printf("\nAssigning %064b to index %d\n", value, index)
	// fmt.Printf("byteOffset: %d, bitOffset: %d, bitsRemaining: %d\n", byteOffset, bitOffset, bitsRemaining)

	writeBitCount := bitsInExtraLargeBlock - bitOffset
	if writeBitCount > bitsRemaining {
		writeBitCount = bitsRemaining
	}
	initial := uint64(m.m[byteOffset])
	mask := uint64(fullExtraLargeBlock << writeBitCount)
	result := (initial & ^(^mask << bitOffset)) | ((value & ^mask) << bitOffset)
	m.m[byteOffset] = extraLargeBlock(result)

	bitsRemaining -= writeBitCount

	// fmt.Printf("result at %d: %032b\n", byteOffset, m.m[byteOffset])

	if bitsRemaining > 0 {
		initial := uint64(m.m[byteOffset+1])
		mask := uint64(fullExtraLargeBlock << bitsRemaining)
		result := (initial & mask) | ((value & (^mask << writeBitCount)) >> writeBitCount)
		m.m[byteOffset+1] = extraLargeBlock(result)

		// fmt.Printf("result at %d: %032b\n", byteOffset+1, m.m[byteOffset+1])
	}
}

// Reads the value at a particular offset
func (m *Memories64) Read(index uint64) (result uint64) {
	bitsRemaining := uint64(m.bitsPerEntry)
	offset := bitsRemaining * index
	bitOffset := offset % bitsInExtraLargeBlock
	byteOffset := offset / bitsInExtraLargeBlock
	fmt.Printf("\nbitOffset: %d, byteOffset: %d\n", bitOffset, byteOffset)
	fmt.Printf("m.m: %x\n", m.m)

	readBitCount := bitsInExtraLargeBlock - bitOffset
	if readBitCount > bitsRemaining {
		readBitCount = bitsRemaining
	}
	initial := uint64(m.m[byteOffset])
	mask := uint64(^(fullExtraLargeBlock << readBitCount)) << bitOffset
	result = uint64((initial & mask) >> bitOffset)

	bitsRemaining -= readBitCount

	fmt.Printf("--> %064b; %d\n", result, bitsRemaining)
	if bitsRemaining > 0 {
		initial := uint64(m.m[byteOffset+1])
		result |= ((initial & uint64(^(fullExtraLargeBlock << bitsRemaining))) << (uint64(m.bitsPerEntry) - bitsRemaining))
		fmt.Printf("--> %032b\n", result)
	}

	return result
}
