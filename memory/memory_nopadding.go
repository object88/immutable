package memory

// MemoriesNoPadding is all your memories.
type MemoriesNoPadding struct {
	m []extraLargeBlock
}

// Assign sets a value to the internal memory at the given index
func (m *MemoriesNoPadding) Assign(index uint64, value uint64) {
	m.m[index] = extraLargeBlock(value)
}

// Reads the value at a particular offset
func (m *MemoriesNoPadding) Read(index uint64) uint64 {
	return uint64(m.m[index])
}
