package immutable

import (
	"bytes"
	"fmt"

	"github.com/object88/immutable/core"
)

func (h *HashMap) String() string {
	if h == nil {
		return "(nil)"
	}

	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Size: %d\n", h.size))
	buffer.WriteString(fmt.Sprintf("[\n"))
	h.ForEach(func(k core.Element, v core.Element) {
		buffer.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
	})
	buffer.WriteString(fmt.Sprintf("]\n"))
	return buffer.String()
}
