package memory

const fullBlock = ^uint32(0)

// Memories32 is all your memories.
type Memories32 struct {
	bitsPerEntry uint32
	m            []largeBlock
}

// Assign sets a value to the byte array at the given index
func (m *Memories32) Assign(index uint32, value uint32) {
	bitsRemaining := m.bitsPerEntry
	offset := m.bitsPerEntry * index
	byteOffset := offset / bitsInLargeBlock
	bitOffset := offset % bitsInLargeBlock

	// fmt.Printf("\nAssigning %032b to index %d\n", value, index)
	// fmt.Printf("byteOffset: %d, bitOffset: %d, bitsRemaining: %d\n", byteOffset, bitOffset, bitsRemaining)

	writeBitCount := bitsInLargeBlock - bitOffset
	if writeBitCount > bitsRemaining {
		writeBitCount = bitsRemaining
	}
	initial := uint32(m.m[byteOffset])
	mask := fullBlock << writeBitCount
	result := (initial & ^(^mask << bitOffset)) | ((value & ^mask) << bitOffset)
	m.m[byteOffset] = largeBlock(result)

	bitsRemaining -= writeBitCount

	// fmt.Printf("result at %d: %032b\n", byteOffset, m.m[byteOffset])

	if bitsRemaining > 0 {
		initial := uint32(m.m[byteOffset+1])
		mask := fullBlock << bitsRemaining
		result := (initial & mask) | ((value & (^mask << writeBitCount)) >> writeBitCount)
		m.m[byteOffset+1] = largeBlock(result)

		// fmt.Printf("result at %d: %032b\n", byteOffset+1, m.m[byteOffset+1])
	}
}

// Reads the value at a particular offset
func (m *Memories32) Read(index uint32) (result uint32) {
	bitsRemaining := m.bitsPerEntry
	offset := m.bitsPerEntry * index
	bitOffset := offset % bitsInLargeBlock
	byteOffset := offset / bitsInLargeBlock
	// fmt.Printf("\nbitOffset: %d, byteOffset: %d\n", bitOffset, byteOffset)
	// fmt.Printf("m.m: %x\n", m.m)

	readBitCount := bitsInLargeBlock - bitOffset
	if readBitCount > bitsRemaining {
		readBitCount = bitsRemaining
	}
	initial := uint32(m.m[byteOffset])
	mask := uint32(^(fullBlock << readBitCount)) << bitOffset
	result = (initial & mask) >> bitOffset

	bitsRemaining -= readBitCount

	// fmt.Printf("--> %032b; %d\n", result, bitsRemaining)
	if bitsRemaining > 0 {
		initial := uint32(m.m[byteOffset+1])
		result |= ((initial & ^(fullBlock << bitsRemaining)) << (m.bitsPerEntry - bitsRemaining))
		// fmt.Printf("--> %032b\n", result)
	}

	return result
}
