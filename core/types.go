package core

import "fmt"

type BucketGenerator func(count int) SubBucket

// Element may be a key or a value
type Element interface {
	fmt.Stringer
	Hash(seed uint32) uint64
}

type HandlerConfig struct {
	Compare      func(a, b Element) (match bool)
	CreateBucket func(count int) SubBucket
}

type KeyValuePair struct {
	Key   Element
	Value Element
}

type SubBucket interface {
	Hydrate(index int) (e Element)
	Dehydrate(index int, e Element)
}
