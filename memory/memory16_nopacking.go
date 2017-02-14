package memory

// Memories16NoPacking is all your memories.
type Memories16NoPacking struct {
	m []uint16
}

// Assign sets a value to the internal memory at the given index
func (m *Memories16NoPacking) Assign(index uint64, value uint64) {
	m.m[index] = uint16(value & 0xffff)
}

// Reads the value at a particular offset
func (m *Memories16NoPacking) Read(index uint64) uint64 {
	return uint64(m.m[index])
}
