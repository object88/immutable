package memory

const (
	// allUint32bits = ^uint32(0)
	bitsInExtraLargeBlock = 64
	bitsInLargeBlock      = 32
	bitsInSmallBlock      = 8
	// bitsInBlock   = uint32(unsafe.Sizeof(uint32(0)) * bitsInSmallBlock)
)
