package memory

// import "fmt"
//
// const fullBlock = ^uint32(0)
//
// // Memories32 is all your memories.
// type Memories32 struct {
// 	bitsPerEntry uint32
// 	m            []largeBlock
// }
//
// // Assign sets a value to the internal memory at the given index
// func (m *Memories32) Assign(index uint64, value uint64) {
// 	bitsRemaining := uint64(m.bitsPerEntry)
// 	offset := bitsRemaining * index
// 	byteOffset := offset / bitsInLargeBlock
// 	bitOffset := offset % bitsInLargeBlock
//
// 	// fmt.Printf("\nAssigning %032b to index %d\n", value, index)
// 	// fmt.Printf("byteOffset: %d, bitOffset: %d, bitsRemaining: %d\n", byteOffset, bitOffset, bitsRemaining)
//
// 	writeBitCount := bitsInLargeBlock - bitOffset
// 	if writeBitCount > bitsRemaining {
// 		writeBitCount = bitsRemaining
// 	}
// 	initial := uint64(m.m[byteOffset])
// 	mask := uint64(fullBlock << writeBitCount)
// 	result := (initial & ^(^mask << bitOffset)) | ((value & ^mask) << bitOffset)
// 	m.m[byteOffset] = largeBlock(result)
//
// 	bitsRemaining -= writeBitCount
//
// 	// fmt.Printf("result at %d: %032b\n", byteOffset, m.m[byteOffset])
//
// 	if bitsRemaining > 0 {
// 		initial := uint64(m.m[byteOffset+1])
// 		mask := uint64(fullBlock << bitsRemaining)
// 		result := (initial & mask) | ((value & (^mask << writeBitCount)) >> writeBitCount)
// 		m.m[byteOffset+1] = largeBlock(result)
//
// 		// fmt.Printf("result at %d: %032b\n", byteOffset+1, m.m[byteOffset+1])
// 	}
// }
//
// // Reads the value at a particular offset
// func (m *Memories32) Read(index uint64) (result uint64) {
// 	bitsRemaining := uint64(m.bitsPerEntry)
// 	offset := bitsRemaining * index
// 	bitOffset := offset % bitsInLargeBlock
// 	byteOffset := offset / bitsInLargeBlock
// 	fmt.Printf("\nbitOffset: %d, byteOffset: %d\n", bitOffset, byteOffset)
// 	fmt.Printf("m.m: %x\n", m.m)
//
// 	readBitCount := bitsInLargeBlock - bitOffset
// 	if readBitCount > bitsRemaining {
// 		readBitCount = bitsRemaining
// 	}
// 	initial := uint64(m.m[byteOffset])
// 	mask := uint64(^(fullBlock << readBitCount)) << bitOffset
// 	result = uint64((initial & mask) >> bitOffset)
//
// 	bitsRemaining -= readBitCount
//
// 	fmt.Printf("--> %032b; %d\n", result, bitsRemaining)
// 	if bitsRemaining > 0 {
// 		initial := uint64(m.m[byteOffset+1])
// 		result |= ((initial & uint64(^(fullBlock << bitsRemaining))) << (uint64(m.bitsPerEntry) - bitsRemaining))
// 		fmt.Printf("--> %032b\n", result)
// 	}
//
// 	return result
// }
