package memory

const (
	// SmallBlock packs values into 8 bit (1 byte) blocks
	SmallBlock BlockSize = iota

	// LargeBlock packs values into 32 bit (4 byte) blocks
	LargeBlock

	// ExtraLargeBlock packs values into 64 bit (8 byte) blocks
	ExtraLargeBlock

	// NoPacking stores values directingly into 64 bit (8 byte) blocks, with no
	// packing
	NoPacking
)

const (
	// allUint32bits = ^uint32(0)
	bitsInExtraLargeBlock = 64
	bitsInLargeBlock      = 32
	bitsInSmallBlock      = 8
	// bitsInBlock   = uint32(unsafe.Sizeof(uint32(0)) * bitsInSmallBlock)
)
