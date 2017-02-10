package core

import "unsafe"

type HandlerConfig interface {
	CompareTo(memory unsafe.Pointer, index int, other unsafe.Pointer) (match bool)
	CreateBucket(count int) unsafe.Pointer
	Hash(element unsafe.Pointer, seed uint32) uint64
	Read(memory unsafe.Pointer, index int) (result unsafe.Pointer)
	Format(value unsafe.Pointer) string
	Write(memory unsafe.Pointer, index int, value unsafe.Pointer)
}

type KeyValuePair struct {
	Key   unsafe.Pointer
	Value unsafe.Pointer
}
