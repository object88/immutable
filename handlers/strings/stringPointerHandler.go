package strings

// import (
// 	"sync"
//
// 	"github.com/object88/immutable/core"
// )
//
// var pconfig *core.HandlerConfig
// var ponce sync.Once
//
// // WithStringPointerKeyMetadata establishes the hydrator and dehydrator methods
// // for working with integer keys.
// func WithStringPointerKeyMetadata(o *core.HashmapOptions) {
// 	o.KeyConfig = createOneStringPointerHandler()
// }
//
// func WithStringPointerValueMetadata(o *core.HashmapOptions) {
// 	o.ValueConfig = createOneStringPointerHandler()
// }
//
// func createOneStringPointerHandler() *core.HandlerConfig {
// 	ponce.Do(func() {
// 		pconfig = &core.HandlerConfig{
// 			Compare: func(a, b core.Element) (match bool) {
// 				return string(a.(StringElement)) == string(b.(StringElement))
// 			},
// 			CreateBucket: func(count int) core.SubBucket {
// 				m := make([]*string, count)
// 				return &StringPointerSubBucket{m}
// 			},
// 		}
// 	})
// 	return pconfig
// }
//
// type StringPointerSubBucket struct {
// 	m []*string
// }
//
// func (h *StringPointerSubBucket) Dehydrate(index int, value core.Element) {
// 	if value == nil {
// 		h.m[index] = nil
// 		return
// 	}
// 	s := value.(StringElement)
// 	h.m[index] = (*string)(&s)
// }
//
// func (h *StringPointerSubBucket) Hydrate(index int) core.Element {
// 	s := h.m[index]
// 	if s == nil {
// 		return nil
// 	}
// 	return StringElement(*s)
// }
