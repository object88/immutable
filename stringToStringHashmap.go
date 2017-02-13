package immutable

// import (
// 	"fmt"
// 	"unsafe"
//
// 	"github.com/object88/immutable/core"
// 	"github.com/object88/immutable/handlers/strings"
// )
//
// type StringToStringHashmap struct {
// 	h *HashMap
// }
//
// func NewStringToStringHashmap(contents map[string]string) *StringToStringHashmap {
// 	opts := core.DefaultHashMapOptions()
// 	strings.WithStringKeyMetadata(opts)
// 	strings.WithStringValueMetadata(opts)
// 	// for _, fn := range options {
// 	// 	fn(opts)
// 	// }
//
// 	hash := CreateEmptyHashmap(len(contents), opts)
//
// 	for k, v := range contents {
// 		key, value := k, v
// 		hash.internalSet(unsafe.Pointer(&key), unsafe.Pointer(&value))
// 	}
// 	return &StringToStringHashmap{hash}
// }
//
// func (stsh *StringToStringHashmap) Get(key string) string {
// 	k := unsafe.Pointer(&key)
// 	r, ok, err := stsh.h.Get(k)
// 	if err != nil {
// 		panic("NOPE")
// 	}
// 	if !ok {
// 		fmt.Printf("!ok\n")
// 	}
// 	fmt.Printf("value: %#v\n", r)
// 	v := *(*string)(r)
// 	return v
// }
//
// func (stsh *StringToStringHashmap) Insert(key string, value string) {
//
// }
//
// func (stsh *StringToStringHashmap) Size() int {
// 	return stsh.h.Size()
// }
//
// func (stsh *StringToStringHashmap) String() string {
// 	return stsh.h.String()
// }
