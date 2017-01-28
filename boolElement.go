package immutable

// BoolElement is a boolean.
type BoolElement bool

// StringElement is a string.
type StringElement string

type BoolSubBucket struct {
	m []bool
}

func NewBoolHandler() BucketGenerator {
	return NewBoolSubBucket
}

func NewBoolSubBucket(count int) SubBucket {
	m := make([]bool, count)
	return &BoolSubBucket{m}
}

func (sb *BoolSubBucket) Hydrate(index int) Element {
	b := sb.m[index]
	e := BoolElement(b)
	return e
}

func (sb *BoolSubBucket) Dehydrate(index int, e Element) {
	be := e.(BoolElement)
	b := bool(be)
	sb.m[index] = b
}

func (b BoolElement) Hash(seed uint32) uint64 {
	return 0
}

func (b BoolElement) String() string {
	if bool(b) {
		return "true"
	}
	return "false"
}
