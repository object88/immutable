package immutable

import (
	"errors"
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

// Filter returns a subset of the collection, based on the predicate supplied
func (hm *IntToStringHashmap) Filter(predicate func(key int, value string) (bool, error)) (*IntToStringHashmap, error) {
	if hm == nil {
		return nil, errors.New("Pointer receiver is nil")
	}
	newHashmap, err := hm.h.Filter(func(kp unsafe.Pointer, vp unsafe.Pointer) (bool, error) {
		key, value := *(*int)(kp), *(*string)(vp)
		return predicate(key, value)
	})
	if err != nil {
		return nil, err
	}
	return &IntToStringHashmap{newHashmap}, nil
}

// ForEach iterates over each key-value pair in this collection
func (hm *IntToStringHashmap) ForEach(predicate func(key int, value string)) {
	if hm == nil {
		return
	}
	hm.h.ForEach(func(kp unsafe.Pointer, vp unsafe.Pointer) {
		key, value := *(*int)(kp), *(*string)(vp)
		predicate(key, value)
	})
}

func (hm *IntToStringHashmap) Get(key int) (value string, ok bool, err error) {
	k := unsafe.Pointer(&key)
	r, ok, err := hm.h.Get(k)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	v := *(*string)(r)
	return v, true, nil
}

func (hm *IntToStringHashmap) Insert(key int, value string) (result *IntToStringHashmap, err error) {
	if hm == nil {
		return nil, errors.New("Pointer receiver is nil")
	}
	kp, vp := unsafe.Pointer(&key), unsafe.Pointer(&value)
	newHashmap, err := hm.h.Insert(kp, vp)
	if err != nil {
		return nil, err
	}
	if newHashmap == hm.h {
		return hm, nil
	}
	return &IntToStringHashmap{newHashmap}, nil
}

func (hm *IntToStringHashmap) Map(predicate func(key int, value string) (result string, err error)) (*IntToStringHashmap, error) {
	if hm == nil {
		return nil, errors.New("Pointer receiver is nil")
	}
	newHashmap, err := hm.h.Map(func(kp, vp unsafe.Pointer) (newv unsafe.Pointer, err error) {
		key, value := *(*int)(kp), *(*string)(vp)
		newS, err := predicate(key, value)
		if err != nil {
			return nil, err
		}
		return unsafe.Pointer(&newS), nil
	})
	if err != nil {
		return nil, err
	}
	return &IntToStringHashmap{newHashmap}, nil
}

func (hm *IntToStringHashmap) Reduce(accumulator interface{}, predicate func(accumulator interface{}, key int, value string) (interface{}, error)) (result interface{}, err error) {
	if hm == nil {
		return nil, errors.New("Pointer receiver is nil")
	}
	acc := accumulator
	_, err = hm.h.Reduce(func(ap unsafe.Pointer, kp, vp unsafe.Pointer) (unsafe.Pointer, error) {
		key, value := *(*int)(kp), *(*string)(vp)
		acc, err = predicate(acc, key, value)
		return nil, err
	}, nil)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (hm *IntToStringHashmap) Size() int {
	if hm == nil {
		return 0
	}
	return hm.h.Size()
}

func (hm *IntToStringHashmap) String() string {
	return hm.h.String()
}
