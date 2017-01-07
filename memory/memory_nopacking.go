package memory

// MemoriesNoPacking is all your memories.
type MemoriesNoPacking struct {
	m []uint64
}

// Assign sets a value to the internal memory at the given index
func (m *MemoriesNoPacking) Assign(index uint64, value uint64) {
	m.m[index] = value
}

// Reads the value at a particular offset
func (m *MemoriesNoPacking) Read(index uint64) uint64 {
	return m.m[index]
}
