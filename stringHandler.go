package immutable

import "unsafe"

type StringSubBucket struct {
	m []unsafe.Pointer
}

func NewStringHandler() BucketGenerator {
	return NewStringSubBucket
}

func NewStringSubBucket(count int) SubBucket {
	m := make([]unsafe.Pointer, count)
	return &StringSubBucket{m}
}

func (h *StringSubBucket) Hydrate(index int) Element {
	u := h.m[index]
	if u == nil {
		return nil
	}
	s := *(*string)(u)
	v := StringElement(s)
	return v
}

func (h *StringSubBucket) Dehydrate(index int, v Element) {
	if v == nil {
		h.m[index] = nil
		return
	}
	v1 := v.(StringElement)
	u := unsafe.Pointer(&v1)
	h.m[index] = u
}
