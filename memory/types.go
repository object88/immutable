package memory

// BlockSize is a numeration of block sizes (8, 32)
type BlockSize int

const (
	// SmallBlock packs values into 8 bit (1 byte) blocks
	SmallBlock BlockSize = iota

	// SmallBlockNoPacking stores values directingly into 8 bit (1 byte)
	// blocks, with no packing
	SmallBlockNoPacking

	// MediumBlock packs values into 16 bit (2 byte) blocks
	MediumBlock

	// MediumBlockNoPacking stores values directingly into 16 bit (2 byte)
	// blocks, with no packing
	MediumBlockNoPacking

	// LargeBlock packs values into 32 bit (4 byte) blocks
	LargeBlock

	// LargeBlockNoPacking stores values directingly into 32 bit (4 byte)
	// blocks, with no packing
	LargeBlockNoPacking

	// ExtraLargeBlock packs values into 64 bit (8 byte) blocks
	ExtraLargeBlock

	// ExtraLargeBlockNoPacking stores values directingly into 64 bit (8 byte)
	// blocks, with no packing
	ExtraLargeBlockNoPacking
)

var blockSizes = [...]string{
	"SmallBlock",
	"SmallBlockNoPacking",
	"MediumBlock",
	"MediumBlockNoPacking",
	"LargeBlock",
	"LargeBlockNoPacking",
	"ExtraLargeBlock",
	"ExtraLargeBlockNoPacking",
}

func (b BlockSize) String() string {
	return blockSizes[b-1]
}
