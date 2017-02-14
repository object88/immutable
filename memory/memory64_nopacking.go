package memory

// Memories64NoPacking is all your memories.
type Memories64NoPacking struct {
	m []uint64
}

// Assign sets a value to the internal memory at the given index
func (m *Memories64NoPacking) Assign(index uint64, value uint64) {
	m.m[index] = value
}

// Reads the value at a particular offset
func (m *Memories64NoPacking) Read(index uint64) uint64 {
	return m.m[index]
}
