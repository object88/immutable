package immutable

import (
	"fmt"
	"testing"
)

type MyBadKey struct {
	value int
}

func (k MyBadKey) Hash() uint32 {
	if k.value%2 == 0 {
		return 0x0
	}
	return 0xffffffff
}

func (k MyBadKey) String() string {
	return fmt.Sprintf("%d", k.value)
}

func Test_Hashmap_CustomKey_BadHash(t *testing.T) {
	data := map[Key]Value{
		MyBadKey{0}: 0,
		MyBadKey{1}: 1,
		MyBadKey{2}: 2,
		MyBadKey{3}: 3,
		MyBadKey{4}: 4,
		MyBadKey{5}: 5,
		MyBadKey{6}: 6,
		MyBadKey{7}: 7,
		MyBadKey{8}: 8,
		MyBadKey{9}: 9,
	}
	original := NewHashMap(data)
	if original == nil {
		t.Fatal("NewHashMap returned nil\n")
	}
}
