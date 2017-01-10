package memory

const fullBlock = ^uint32(0)

// Memories32 is all your memories.
type Memories32 struct {
	bitsPerEntry uint32
	m            []uint32
}

// Assign sets a value to the internal memory at the given index
func (m *Memories32) Assign(index uint64, value uint64) {
	bitsRemaining := m.bitsPerEntry
	offset := bitsRemaining * uint32(index)
	byteOffset := offset / bitsInLargeBlock
	bitOffset := offset % bitsInLargeBlock

	// fmt.Printf("\nAssigning %064b to index %d\n", value, index)

	writeBitCount := bitsInLargeBlock - bitOffset
	if writeBitCount > bitsRemaining {
		writeBitCount = bitsRemaining
	}
	// fmt.Printf("byteOffset: %d, bitOffset: %d, bitsRemaining: %d, writeBitCount: %d\n", byteOffset, bitOffset, bitsRemaining, writeBitCount)
	initial := m.m[byteOffset]
	mask := ^(fullExtraLargeBlock << writeBitCount)
	result := uint32(value&mask)<<bitOffset | initial&^((^(fullBlock << writeBitCount))<<bitOffset)
	m.m[byteOffset] = result

	// fmt.Printf("result at %d: %032b ->  %032b\n", byteOffset, initial, result)

	bitsRemaining -= writeBitCount
	byteOffset++

	if bitsRemaining >= 32 {
		o := m.bitsPerEntry - bitsRemaining
		result := uint32((value & (fullExtraLargeBlock << o)) >> o)
		m.m[byteOffset] = result
		// fmt.Printf("result at %d: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx ->  %032b\n", byteOffset, result)

		bitsRemaining -= 32
		byteOffset++
	}

	if bitsRemaining > 0 {
		writeBitCount = m.bitsPerEntry - bitsRemaining
		initial := m.m[byteOffset]
		mask := fullExtraLargeBlock << bitsRemaining
		result := (initial & (fullBlock << bitsRemaining)) | uint32((value&((^mask)<<writeBitCount))>>writeBitCount)
		m.m[byteOffset] = result

		// fmt.Printf("result at %d: %032b ->  %032b\n", byteOffset, initial, result)
	}
}

// Reads the value at a particular offset
func (m *Memories32) Read(index uint64) (result uint64) {
	bitsRemaining := uint64(m.bitsPerEntry)
	offset := bitsRemaining * index
	bitOffset := offset % bitsInLargeBlock
	byteOffset := offset / bitsInLargeBlock
	// fmt.Printf("\nbitOffset: %d, byteOffset: %d\n", bitOffset, byteOffset)
	// fmt.Printf("m.m: %x\n", m.m)

	readBitCount := bitsInLargeBlock - bitOffset
	if readBitCount > bitsRemaining {
		readBitCount = bitsRemaining
	}
	initial := m.m[byteOffset]
	mask := ^(fullBlock << readBitCount) << bitOffset
	result = uint64((initial & mask) >> bitOffset)

	bitsRemaining -= readBitCount
	byteOffset++

	if bitsRemaining >= 32 {
		// fmt.Printf("--> %064b; %d; %d\n", result, bitsRemaining, 32)
		initial := m.m[byteOffset]
		result |= (uint64(initial) << (uint64(m.bitsPerEntry) - bitsRemaining))
		bitsRemaining -= 32
		byteOffset++
	}

	if bitsRemaining > 0 {
		initial := m.m[byteOffset]
		// fmt.Printf("--> %064b; %d\n", result, bitsRemaining)
		result |= uint64(initial&(fullBlock>>(32-bitsRemaining))) << (uint64(m.bitsPerEntry) - bitsRemaining)
	}
	// fmt.Printf("--> %064b\n", result)

	return result
}
