package memory

import "fmt"

// Memories8 is all your memories.
type Memories8 struct {
	bitsPerEntry uint32
	m            []smallBlock
}

// bitsPerEntry = 30 -> 0x1E (ends at 0x3C)
// 0               1               2               3
// 0123456789abcedf0123456789abcedf0123456789abcedf0123456789abcedf
//                              [______________________________]
// Reading:
// Take the first 9 bits.
// While remaining count > 8, shift 8 bits, OR 8 new bits.
// Shift 6 bits, OR 6 new bits
// Writing:
// Get leftmost 9 bits, write to p + (16-9)
// While remaining count > 8, write to p + n
// Remainging bits write to p + n + 1, shifted 16 - remaining

// Assign sets a value to the byte array at the given index
func (m *Memories8) Assign(index uint32, value uint32) {
	bitsRemaining := m.bitsPerEntry
	offset := m.bitsPerEntry * index
	byteOffset := offset / bitsInSmallBlock
	bitOffset := offset % bitsInSmallBlock

	fmt.Printf("\nAssigning %08b to index %d\n", value, index)

	nextBitsRemaining, nextByteOffset, nextReadBitOffset := writeFirstBits8(m.m, bitsRemaining, byteOffset, bitOffset, value)
	for nextBitsRemaining > 0 {
		nextBitsRemaining, nextByteOffset, nextReadBitOffset = writeNextBits32(m.m, nextBitsRemaining, nextByteOffset, value, nextReadBitOffset)
	}
}

func writeFirstBits8(memory []smallBlock, bitsRemaining, byteOffset, bitOffset, value uint32) (nextBitsRemaining, nextByteOffset, nextReadBitOffset uint32) {
	// i := memory[byteOffset]
	fmt.Printf("byteOffset: %d, bitOffset: %d, bitsRemaining: %d\n", byteOffset, bitOffset, bitsRemaining)
	writeBitCount := bitsInSmallBlock - bitOffset
	if writeBitCount > bitsRemaining {
		writeBitCount = bitsRemaining
	}
	// mask32 := generateMask32(writeBitCount, 0)
	mask32 := byte(^(0xff << writeBitCount))
	v := byte(value) & mask32
	fmt.Printf("value: %016b, mask: %08b, v: %08b\n", value, mask32, v)

	// memory[byteOffset] = writeBits(i, byte(v), byte(writeBitCount), byte(bitOffset))
	// mask := generateMask(byte(writeBitCount), byte(bitOffset))
	mask := byte(^(0xff << byte(writeBitCount))) << byte(bitOffset)
	memory[byteOffset] = smallBlock((byte(memory[byteOffset]) & ^mask) | v<<byte(bitOffset))
	fmt.Printf("result at %d: %08b\n", byteOffset, memory[byteOffset])
	return bitsRemaining - writeBitCount, byteOffset + 1, writeBitCount
}

func writeNextBits32(memory []smallBlock, bitsRemaining, byteOffset, value, readBitOffset uint32) (nextBitsRemaining, nextByteOffset, nextReadBitOffset uint32) {
	// i := memory[byteOffset]
	fmt.Printf("byteOffset: %d, bitsRemaining: %d, readBitOffset: %d\n", byteOffset, bitsRemaining, readBitOffset)
	writeBitCount := uint32(bitsInSmallBlock)
	if writeBitCount > bitsRemaining {
		writeBitCount = bitsRemaining
	}
	// mask32 := generateMask32(writeBitCount, readBitOffset)
	mask32 := uint32(^(0xff << writeBitCount)) << readBitOffset
	v := byte(value&mask32) >> readBitOffset
	fmt.Printf("value: %016b, mask: %08b, v: %08b\n", value, mask32, v)
	// memory[byteOffset] = writeBits(i, byte(v), byte(writeBitCount), 0)
	// mask := generateMask(width, shift)
	mask := byte(^(0xff << writeBitCount))
	memory[byteOffset] = smallBlock((byte(memory[byteOffset]) & ^mask) | v)

	fmt.Printf("result at %d: %08b\n", byteOffset, memory[byteOffset])
	return bitsRemaining - writeBitCount, byteOffset + 1, readBitOffset + writeBitCount
}

// Reads the value at a particular offset
func (m *Memories8) Read(index uint32) uint32 {
	bitsRemaining := m.bitsPerEntry
	offset := m.bitsPerEntry * index
	bitOffset := offset % bitsInSmallBlock
	byteOffset := offset / bitsInSmallBlock
	// fmt.Printf("\nbitOffset: %d, byteOffset: %d\n", bitOffset, byteOffset)
	// fmt.Printf("m.m: %x\n", m.m)

	result, nextBitsRemaining, nextByteOffset := readFirstBits(m.m, bitsRemaining, byteOffset, bitOffset)
	// fmt.Printf("--> 0x%x; %d, %d\n", result, nextBitsRemaining, nextByteOffset)
	for nextBitsRemaining > 0 {
		result, nextBitsRemaining, nextByteOffset = readNextBits(result, m.m, nextBitsRemaining, nextByteOffset)
		// fmt.Printf("--> 0x%x; %d, %d\n", result, nextBitsRemaining, nextByteOffset)
	}

	return result
}

func readFirstBits(memory []smallBlock, bitsRemaining, byteOffset, bitOffset uint32) (result, nextBitsRemaining, nextByteOffset uint32) {
	// fmt.Printf("Reading first byte: %d, %d\n", byteOffset, bitOffset)
	i := byte(memory[byteOffset])
	readBitCount := bitsInSmallBlock - bitOffset
	if readBitCount > bitsRemaining {
		readBitCount = bitsRemaining
	}
	b := readBits(i, generateMask(byte(readBitCount), byte(bitOffset)))
	result = uint32(b) >> bitOffset

	// fmt.Printf("bitsRemaining: %d, readBitCount: %d\n", bitsRemaining, readBitCount)
	// fmt.Printf("Returning 0x%x, %d, %d\n", result, bitsRemaining-readBitCount, byteOffset+1)
	return result, bitsRemaining - readBitCount, byteOffset + 1
}

func readNextBits(initial uint32, memory []smallBlock, bitsRemaining, byteOffset uint32) (result, nextBitsRemaining, nextByteoffset uint32) {
	// fmt.Printf("Reading next byte: %d\n", byteOffset)
	i := byte(memory[byteOffset])
	readBitCount := uint32(bitsInSmallBlock)
	if readBitCount > bitsRemaining {
		readBitCount = bitsRemaining
	}
	b := readBits(i, generateMask(byte(readBitCount), 0))
	result = (initial << readBitCount) | uint32(b)

	// fmt.Printf("Returning 0x%x, %d, %d\n", result, bitsRemaining-readBitCount, byteOffset+1)
	return result, bitsRemaining - readBitCount, byteOffset + 1
}

func readBits(original, mask byte) byte {
	return original & mask
}

func writeBits(original, newValue, width, shift byte) byte {
	mask := generateMask(width, shift)
	result := (original & ^mask) | newValue<<shift
	return result
}

func generateMask32(width, shift uint32) uint32 {
	return (^(0xffffffff << width)) << shift
}

func generateMask(width, shift byte) byte {
	return (^(0xff << width)) << shift
}
