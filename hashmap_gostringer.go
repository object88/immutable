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
			for i := 0; i < int(b.entryCount); i++ {
				k, _ := b.keys.Hydrate(i).(Element)
				v := b.values.Hydrate(i)
				buffer.WriteString(fmt.Sprintf("      [0x%016x,%s] -> %s\n", b.hobs.Read(uint64(i)), k, v))
			}

			b = b.overflow
		}
		buffer.WriteString("    ]\n")
		buffer.WriteString("  },\n")
	}
	buffer.WriteString("]\n")
	return buffer.String()
}
