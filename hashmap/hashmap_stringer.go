package hashmap

import (
	"bytes"
	"fmt"

	"github.com/object88/immutable"
)

func (h *HashMap) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Count: %d\n", h.count))
	buffer.WriteString(fmt.Sprintf("[\n"))
	h.ForEach(func(k immutable.Key, v immutable.Value) {
		buffer.WriteString(fmt.Sprintf("  %s: %s\n", k, v))
	})
	buffer.WriteString(fmt.Sprintf("]\n"))
	return buffer.String()
}
