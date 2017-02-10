package immutable

import (
	"fmt"
	"unsafe"

	"github.com/object88/immutable/core"
	"github.com/object88/immutable/handlers/integers"
	"github.com/object88/immutable/handlers/strings"
)

type IntToStringHashmap struct {
	h *HashMap
}

func NewIntToStringHashmap(contents map[int]string) *IntToStringHashmap {
	opts := core.DefaultHashMapOptions()
	integers.WithIntKeyMetadata(opts)
	strings.WithStringValueMetadata(opts)
	// for _, fn := range options {
	// 	fn(opts)
	// }

	hash := createHashMap(len(contents), opts)

	for k, v := range contents {
		key, value := k, v
		hash.internalSet(unsafe.Pointer(&key), unsafe.Pointer(&value))
	}
	return &IntToStringHashmap{hash}
}

func (stsh *IntToStringHashmap) Get(key int) string {
	k := unsafe.Pointer(&key)
	r, ok, err := stsh.h.Get(k)
	if err != nil {
		panic("NOPE")
	}
	if !ok {
		fmt.Printf("!ok\n")
	}
	v := *(*string)(r)
	return v
}

func (stsh *IntToStringHashmap) Size() int {
	return stsh.h.Size()
}

func (stsh *IntToStringHashmap) String() string {
	return stsh.h.String()
}
