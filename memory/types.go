package memory

// BlockSize is a numeration of block sizes (8, 32)
type BlockSize int

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

var blockSizes = [...]string{
	"SmallBlock",
	"LargeBlock",
	"ExtraLargeBlock",
	"NoPacking",
}

func (b BlockSize) String() string {
	return blockSizes[b-1]
}
