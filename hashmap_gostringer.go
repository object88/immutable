package immutable

import (
	"bytes"
	"fmt"
)

// GoString provides a programmatic view into a HashMap.  This may be used,
// for example, with the '%#v' operand to fmt.Printf, fmt.Sprintf, etc.
func (h *HashMap) GoString() string {
	if h == nil {
		return "(nil)"
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Size: %d\n", h.size))
	buffer.WriteString("[\n")
	for k, v := range h.buckets {
		b := v
		if b == nil {
			buffer.WriteString(fmt.Sprintf("  bucket #%d: nil\n", k))
			continue
		}
		buffer.WriteString(fmt.Sprintf("  bucket #%d: {\n", k))
		buffer.WriteString(fmt.Sprintf("    entryCount: %d\n", b.entryCount))
		buffer.WriteString("    entries: [\n")
		for b != nil {
			for i := uint64(0); i < uint64(b.entryCount); i++ {
				buffer.WriteString(fmt.Sprintf("      [0x%016x,%s] -> %s\n", b.hobs.Read(i), b.entries[i].key, b.entries[i].value))
			}

			b = b.overflow
		}
		buffer.WriteString("    ]\n")
		buffer.WriteString("  },\n")
	}
	buffer.WriteString("]\n")
	return buffer.String()
}
