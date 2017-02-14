package immutable

import (
	"errors"
	"unsafe"

	"github.com/object88/immutable/core"
	"github.com/object88/immutable/handlers/integers"
	"github.com/object88/immutable/handlers/strings"
)

type IntToStringHashmap struct {
	h *core.InternalHashmap
	c *core.HashmapConfig
}

func NewIntToStringHashmap(contents map[int]string, options ...core.HashmapOption) *IntToStringHashmap {
	opts := core.DefaultHashmapOptions()
	for _, fn := range options {
		fn(opts)
	}

	c := &core.HashmapConfig{
		KeyConfig:   integers.GetHandler(),
		Options:     opts,
		ValueConfig: strings.GetHandler(),
	}

	hash := core.CreateEmptyInternalHashmap(opts.PackedBucket, len(contents))

	for k, v := range contents {
		key, value := k, v
		hash.InternalSet(c, unsafe.Pointer(&key), unsafe.Pointer(&value))
	}
	return &IntToStringHashmap{hash, c}
}

// Filter returns a subset of the collection, based on the predicate supplied
func (hm *IntToStringHashmap) Filter(predicate func(key int, value string) (bool, error)) (*IntToStringHashmap, error) {
	if hm == nil {
		return nil, errors.New("Pointer receiver is nil")
	}
	newHashmap, err := hm.h.Filter(hm.c, func(kp, vp unsafe.Pointer) (bool, error) {
		key, value := *(*int)(kp), *(*string)(vp)
		return predicate(key, value)
	})
	if err != nil {
		return nil, err
	}
	return &IntToStringHashmap{newHashmap, hm.c}, nil
}

// ForEach iterates over each key-value pair in this collection
func (hm *IntToStringHashmap) ForEach(predicate func(key int, value string)) {
	if hm == nil {
		return
	}
	hm.h.ForEach(hm.c, func(kp, vp unsafe.Pointer) {
		key, value := *(*int)(kp), *(*string)(vp)
		predicate(key, value)
	})
}

func (hm *IntToStringHashmap) Get(key int) (value string, ok bool, err error) {
	if hm == nil {
		return "", false, errors.New("Pointer receiver is nil")
	}
	k := unsafe.Pointer(&key)
	r, ok, err := hm.h.Get(hm.c, k)
	if err != nil {
		return "", false, err
	}
	if !ok {
		return "", false, nil
	}
	v := *(*string)(r)
	return v, true, nil
}

func (hm *IntToStringHashmap) GetKeys() (results []int, err error) {
	if hm == nil {
		return nil, errors.New("Pointer receiver is nil")
	}
	var a []unsafe.Pointer
	a, err = hm.h.GetKeys(hm.c)
	if err != nil {
		return nil, err
	}
	results = make([]int, len(a))
	for k, v := range a {
		results[k] = *(*int)(v)
	}
	return results, nil
}

func (hm *IntToStringHashmap) GoString() string {
	if hm == nil {
		return "(nil)"
	}
	return hm.h.GoString(hm.c)
}

func (hm *IntToStringHashmap) Insert(key int, value string) (result *IntToStringHashmap, err error) {
	if hm == nil {
		return nil, errors.New("Pointer receiver is nil")
	}
	kp, vp := unsafe.Pointer(&key), unsafe.Pointer(&value)
	newHashmap, err := hm.h.Insert(hm.c, kp, vp)
	if err != nil {
		return nil, err
	}
	if newHashmap == hm.h {
		return hm, nil
	}
	return &IntToStringHashmap{newHashmap, hm.c}, nil
}

func (hm *IntToStringHashmap) Map(predicate func(key int, value string) (result string, err error)) (*IntToStringHashmap, error) {
	if hm == nil {
		return nil, errors.New("Pointer receiver is nil")
	}
	newHashmap, err := hm.h.Map(hm.c, func(kp, vp unsafe.Pointer) (newv unsafe.Pointer, err error) {
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
	return &IntToStringHashmap{newHashmap, hm.c}, nil
}

func (hm *IntToStringHashmap) Reduce(accumulator interface{}, predicate func(accumulator interface{}, key int, value string) (interface{}, error)) (result interface{}, err error) {
	if hm == nil {
		return nil, errors.New("Pointer receiver is nil")
	}
	acc := accumulator
	_, err = hm.h.Reduce(hm.c, func(ap, kp, vp unsafe.Pointer) (unsafe.Pointer, error) {
		key, value := *(*int)(kp), *(*string)(vp)
		acc, err = predicate(acc, key, value)
		return nil, err
	}, nil)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

// Remove returns a copy of the provided InternalHashmap with the specified element
// removed.
func (hm *IntToStringHashmap) Remove(key int) (*IntToStringHashmap, error) {
	if hm == nil {
		return nil, errors.New("Pointer receiver is nil")
	}

	kp := unsafe.Pointer(&key)
	newHashmap, err := hm.h.Remove(hm.c, kp)
	if err != nil {
		return nil, err
	}
	if newHashmap == hm.h {
		return hm, nil
	}
	return &IntToStringHashmap{newHashmap, hm.c}, nil
}

func (hm *IntToStringHashmap) Size() int {
	if hm == nil {
		return 0
	}
	return hm.h.Size()
}

func (hm *IntToStringHashmap) String() string {
	if hm == nil {
		return "(nil)"
	}
	return hm.h.String(hm.c)
}
