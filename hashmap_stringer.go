package immutable

import (
	"bytes"
	"fmt"
)

func (h *HashMap) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Size: %d\n", h.size))
	buffer.WriteString(fmt.Sprintf("[\n"))
	h.ForEach(func(k Key, v Value) {
		buffer.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
	})
	buffer.WriteString(fmt.Sprintf("]\n"))
	return buffer.String()
}
