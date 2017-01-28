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

func (h *StringSubBucket) Hydrate(index int) Value {
	u := h.m[index]
	s := *(*string)(u)
	v := StringValue(s)
	return v
}

func (h *StringSubBucket) Dehydrate(index int, v Value) {
	v1 := v.(StringValue)
	u := unsafe.Pointer(&v1)
	h.m[index] = u
}
