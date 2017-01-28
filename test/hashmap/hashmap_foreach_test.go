package immutable

import (
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap_ForEach(t *testing.T) {
	data := map[immutable.Element]immutable.Element{
		immutable.IntElement(0): false,
		immutable.IntElement(1): false,
		immutable.IntElement(2): false,
		immutable.IntElement(3): false,
	}
	original := immutable.NewHashMap(data)
	if original == nil {
		t.Fatal("Failed to create hashmap")
	}
	original.ForEach(func(k immutable.Element, v immutable.Element) {
		if v.(bool) {
			t.Fatalf("At %s, already visited\n", k)
		}
		data[k] = true
	})

	for k, v := range data {
		if !v.(bool) {
			t.Fatalf("At %s, not visited\n", k)
		}
	}
}
