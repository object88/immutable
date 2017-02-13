package immutable

import (
	"bytes"
	"fmt"

	"github.com/object88/immutable/core"
)

// GoString provides a programmatic view into a HashMap.  This may be used,
// for example, with the '%#v' operand to fmt.Printf, fmt.Sprintf, etc.
func (h *HashMap) GoString(config *core.HashmapConfig) string {
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
				key := config.KeyConfig.Read(b.keys, i)
				value := config.ValueConfig.Read(b.values, i)

				ks := config.KeyConfig.Format(key)
				vs := config.ValueConfig.Format(value)
				buffer.WriteString(fmt.Sprintf("      [0x%016x,%s] -> %s\n", b.hobs.Read(uint64(i)), ks, vs))
			}

			b = b.overflow
		}
		buffer.WriteString("    ]\n  },\n")
	}
	buffer.WriteString("]\n")
	return buffer.String()
}
