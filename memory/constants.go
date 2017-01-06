package memory

const (
	// SmallBlock indicates 8 bit (1 byte) blocks
	SmallBlock BlockSize = iota

	// LargeBlock indicates 32 bit (4 byte) blocks
	LargeBlock

	// NoPadding indicates 64 bit (8 byte) blocks, with no padding
	NoPadding
)

const (
	// allUint32bits = ^uint32(0)
	bitsInLargeBlock = 32
	bitsInNoPadding  = 64
	bitsInSmallBlock = 8
	// bitsInBlock   = uint32(unsafe.Sizeof(uint32(0)) * bitsInSmallBlock)
)
