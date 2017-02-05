package immutable_test

import (
	"testing"

	"github.com/object88/immutable"
	"github.com/object88/immutable/core"
	"github.com/object88/immutable/handlers/booleans"
	"github.com/object88/immutable/handlers/integers"
)

func Test_Hashmap_ForEach(t *testing.T) {
	data := map[core.Element]core.Element{
		integers.IntElement(0): booleans.BoolElement(false),
		integers.IntElement(1): booleans.BoolElement(false),
		integers.IntElement(2): booleans.BoolElement(false),
		integers.IntElement(3): booleans.BoolElement(false),
	}
	original := immutable.NewHashMap(data, integers.WithIntKeyMetadata, booleans.WithBoolValueMetadata)
	if original == nil {
		t.Fatal("Failed to create hashmap")
	}
	original.ForEach(func(k core.Element, v core.Element) {
		if bool(v.(booleans.BoolElement)) {
			t.Fatalf("At %s, already visited\n", k)
		}
		data[k] = booleans.BoolElement(true)
	})

	for k, v := range data {
		if !bool(v.(booleans.BoolElement)) {
			t.Fatalf("At %s, not visited\n", k)
		}
	}
}
