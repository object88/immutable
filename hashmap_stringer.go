package immutable

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/object88/immutable/core"
)

func (h *HashMap) String(config *core.HashmapConfig) string {
	var buffer bytes.Buffer
	buffer.WriteString("Size: ")
	buffer.WriteString(fmt.Sprintf("%d", h.size))
	buffer.WriteString("\n[\n")
	h.ForEach(config, func(k unsafe.Pointer, v unsafe.Pointer) {
		ks := config.KeyConfig.Format(k)
		vs := config.ValueConfig.Format(v)
		buffer.WriteString(fmt.Sprintf("  %s: %s\n", ks, vs))
	})
	buffer.WriteString("]\n")
	return buffer.String()
}
