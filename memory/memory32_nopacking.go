package memory

// Memories32NoPacking is all your memories.
type Memories32NoPacking struct {
	m []uint32
}

// Assign sets a value to the internal memory at the given index
func (m *Memories32NoPacking) Assign(index uint64, value uint64) {
	m.m[index] = uint32(value & 0xffffffff)
}

// Reads the value at a particular offset
func (m *Memories32NoPacking) Read(index uint64) uint64 {
	return uint64(m.m[index])
}
