package memory

// Memories8 is all your memories.
type Memories8 struct {
	bitsPerEntry uint32
	m            []byte
}

const fullByteBlock = uint8(0xff)

// Assign sets a value to the internal memory at the given index
func (m *Memories8) Assign(index uint64, value uint64) {
	bitsRemaining := m.bitsPerEntry
	offset := bitsRemaining * uint32(index)
	byteOffset := offset / bitsInSmallBlock
	bitOffset := offset % bitsInSmallBlock

	// fmt.Printf("\nAssigning %064b to index %d\n", value, index)

	writeBitCount := bitsInSmallBlock - bitOffset
	if writeBitCount > bitsRemaining {
		writeBitCount = bitsRemaining
	}
	// fmt.Printf("byteOffset: %d, bitOffset: %d, bitsRemaining: %d, writeBitCount: %d\n", byteOffset, bitOffset, bitsRemaining, writeBitCount)
	initial := m.m[byteOffset]
	mask := ^(fullExtraLargeBlock << writeBitCount)
	result := uint8(value&mask)<<bitOffset | initial&^((^(fullByteBlock << writeBitCount))<<bitOffset)
	m.m[byteOffset] = result

	// fmt.Printf("result at %d: %08b ->  %08b\n", byteOffset, initial, result)

	bitsRemaining -= writeBitCount
	byteOffset++

	for bitsRemaining >= 8 {
		o := m.bitsPerEntry - bitsRemaining
		result := uint8((value & (fullExtraLargeBlock << o)) >> o)
		m.m[byteOffset] = result
		// fmt.Printf("result at %d: xxxxxxxx ->  %08b\n", byteOffset, result)

		bitsRemaining -= 8
		byteOffset++
	}

	if bitsRemaining > 0 {
		writeBitCount = m.bitsPerEntry - bitsRemaining
		initial := m.m[byteOffset]
		mask := fullExtraLargeBlock << bitsRemaining
		result := (initial & (fullByteBlock << bitsRemaining)) | uint8((value&((^mask)<<writeBitCount))>>writeBitCount)
		m.m[byteOffset] = result

		// fmt.Printf("result at %d: %08b ->  %08b\n", byteOffset, initial, result)
	}
}

// Reads the value at a particular offset
func (m *Memories8) Read(index uint64) (result uint64) {
	bitsRemaining := uint64(m.bitsPerEntry)
	offset := bitsRemaining * index
	bitOffset := offset % bitsInSmallBlock
	byteOffset := offset / bitsInSmallBlock
	// fmt.Printf("\nbitOffset: %d, byteOffset: %d\n", bitOffset, byteOffset)
	// fmt.Printf("m.m: %x\n", m.m)

	readBitCount := bitsInSmallBlock - bitOffset
	if readBitCount > bitsRemaining {
		readBitCount = bitsRemaining
	}
	initial := m.m[byteOffset]
	mask := ^(fullByteBlock << readBitCount) << bitOffset
	result = uint64((initial & mask) >> bitOffset)

	bitsRemaining -= readBitCount
	byteOffset++

	for bitsRemaining >= 8 {
		// fmt.Printf("--> %064b; %d; %d\n", result, bitsRemaining, 8)
		initial := m.m[byteOffset]
		result |= (uint64(initial) << (uint64(m.bitsPerEntry) - bitsRemaining))
		bitsRemaining -= 8
		byteOffset++
	}

	if bitsRemaining > 0 {
		initial := m.m[byteOffset]
		// fmt.Printf("--> %064b; %d\n", result, bitsRemaining)
		result |= uint64(initial&(fullByteBlock>>(8-bitsRemaining))) << (uint64(m.bitsPerEntry) - bitsRemaining)
	}
	// fmt.Printf("--> %064b\n", result)

	return result
}
