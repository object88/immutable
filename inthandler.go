package immutable

import "github.com/object88/immutable/memory"

type IntSubBucket struct {
	m memory.Memories
}

func NewIntHandler() BucketGenerator {
	return NewIntSubBucket
}

func NewIntSubBucket(count int) SubBucket {
	m := memory.AllocateMemories(memory.NoPacking, uint32(0), uint32(count))
	return &IntSubBucket{m}
}

func (h *IntSubBucket) Hydrate(index int) Element {
	u := h.m.Read(uint64(index))
	v := IntElement(u)
	return v
}

func (h *IntSubBucket) Dehydrate(index int, v Element) {
	u := v.(IntElement)
	// u := *(*int)(unsafe.Pointer(&v))
	h.m.Assign(uint64(index), uint64(u))
}
