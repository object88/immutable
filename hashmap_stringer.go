package immutable

import (
	"bytes"
	"fmt"
)

func (h *HashMap) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Count: %d\n", h.count))
	buffer.WriteString(fmt.Sprintf("Size: %d\n", h.size))
	buffer.WriteString(fmt.Sprintf("[\n"))
	for k, v := range h.buckets {
		buffer.WriteString(fmt.Sprintf("\t%d: {\n", k))
		buffer.WriteString(fmt.Sprintf("\t\tentryCount: %d\n", v.entryCount))
		buffer.WriteString("\t\tentries: [\n")
		for i := uint32(0); i < uint32(v.entryCount); i++ {
			buffer.WriteString(fmt.Sprintf("\t\t\t[0x%08x,%s] -> %s\n", v.hobs.Read(i), v.entries[i].key, v.entries[i].value))
		}
		buffer.WriteString("\t\t]\n")
		buffer.WriteString("\t},\n")
	}
	buffer.WriteString(fmt.Sprintf("]\n"))
	return buffer.String()
}
