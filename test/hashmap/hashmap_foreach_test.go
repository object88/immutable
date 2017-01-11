package immutable

import (
	"testing"

	"github.com/object88/immutable"
)

func Test_Hashmap_ForEach(t *testing.T) {
	data := map[immutable.Key]immutable.Value{
		immutable.IntKey(0): false,
		immutable.IntKey(1): false,
		immutable.IntKey(2): false,
		immutable.IntKey(3): false,
	}
	original := immutable.NewHashMap(data, nil)
	if original == nil {
		t.Fatal("Failed to create hashmap")
	}
	original.ForEach(func(k immutable.Key, v immutable.Value) {
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
