package immutable

import (
	"bytes"
	"fmt"
	"unsafe"
)

func (h *HashMap) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Size: ")
	buffer.WriteString(fmt.Sprintf("%d", h.size))
	buffer.WriteString("\n[\n")
	h.ForEach(func(k unsafe.Pointer, v unsafe.Pointer) {
		ks := h.options.KeyConfig.Format(k)
		vs := h.options.ValueConfig.Format(v)
		buffer.WriteString(fmt.Sprintf("  %s: %s\n", ks, vs))
	})
	buffer.WriteString("]\n")
	return buffer.String()
}
