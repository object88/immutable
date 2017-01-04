package memory

import "unsafe"

const (
	// SmallBlock indicates 8 bit (1 byte) blocks
	SmallBlock BlockSize = 8

	// LargeBlock indicates 32 bit (4 byte) blocks
	LargeBlock = 32
)

const (
	// allUint32bits = ^uint32(0)
	bitsInLargeBlock = uint32(unsafe.Sizeof(uint32(0))) * bitsInSmallBlock
	bitsInSmallBlock = 8
	// bitsInBlock   = uint32(unsafe.Sizeof(uint32(0)) * bitsInSmallBlock)
)
