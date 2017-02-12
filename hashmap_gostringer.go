package immutable

import (
	"bytes"
	"fmt"
)

// GoString provides a programmatic view into a HashMap.  This may be used,
// for example, with the '%#v' operand to fmt.Printf, fmt.Sprintf, etc.
func (h *HashMap) GoString() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Size: %d\n[\n", h.size))
	for k, v := range h.buckets {
		b := v
		if b == nil {
			buffer.WriteString(fmt.Sprintf("  bucket #%d: nil\n", k))
			continue
		}
		buffer.WriteString(fmt.Sprintf("  bucket #%d: {\n    entryCount: %d\n    entries: [\n", k, b.entryCount))
		for b != nil {
			for i := 0; i < int(b.entryCount); i++ {
				key := h.options.KeyConfig.Read(b.keys, i)
				value := h.options.ValueConfig.Read(b.values, i)

				ks := h.options.KeyConfig.Format(key)
				vs := h.options.ValueConfig.Format(value)
				// k, _ := b.keys.Hydrate(i).(core.Element)
				// v := b.values.Hydrate(i)
				buffer.WriteString(fmt.Sprintf("      [0x%016x,%s] -> %s\n", b.hobs.Read(uint64(i)), ks, vs))
			}

			b = b.overflow
		}
		buffer.WriteString("    ]\n  },\n")
	}
	buffer.WriteString("]\n")
	return buffer.String()
}
