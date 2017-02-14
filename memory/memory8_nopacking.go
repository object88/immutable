package memory

// Memories8NoPacking is all your memories.
type Memories8NoPacking struct {
	m []uint8
}

// Assign sets a value to the internal memory at the given index
func (m *Memories8NoPacking) Assign(index uint64, value uint64) {
	m.m[index] = uint8(value & 0xff)
}

// Reads the value at a particular offset
func (m *Memories8NoPacking) Read(index uint64) uint64 {
	return uint64(m.m[index])
}
